package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/gofund/shared/jwt"
	"github.com/gofund/shared/metrics"
	"github.com/gofund/shared/models"
	"github.com/gofund/shared/password"
	"github.com/gofund/users-service/internal/dto"
	"github.com/gofund/users-service/internal/repository"
	"github.com/google/uuid"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo     *repository.UserRepository
	sessionRepo  *repository.SessionRepository
	jwtService   *jwt.JWTService
	eventService *EventService
}

// NewAuthService creates a new auth service instance
func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, jwtService *jwt.JWTService, eventService *EventService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		jwtService:   jwtService,
		eventService: eventService,
	}
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		metrics.TrackLoginFailure("invalid_email")
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	valid, err := password.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		metrics.TrackLoginFailure("verification_error")
		return nil, errors.New("authentication failed")
	}
	if !valid {
		metrics.TrackLoginFailure("invalid_password")
		return nil, errors.New("invalid credentials")
	}

	// Generate token pair with user role
	roles := []string{string(user.Role)}
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID.String(), user.Email, roles)
	if err != nil {
		metrics.TrackLoginFailure("token_generation_failed")
		return nil, errors.New("failed to generate tokens")
	}

	// Create session for refresh token tracking
	refreshTokenHash := s.hashToken(tokenPair.RefreshToken)
	session := &models.Session{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
		Metadata: map[string]interface{}{
			"login_time": time.Now(),
			"ip":         "", // Will be set by controller
			"user_agent": "", // Will be set by controller
		},
	}

	if err := s.sessionRepo.CreateSession(session); err != nil {
		metrics.TrackLoginFailure("session_creation_failed")
		return nil, errors.New("failed to create session")
	}

	// Track successful login and metrics
	metrics.TrackLoginSuccess(user.ID.String())
	metrics.TrackSessionCreated(user.ID.String())
	metrics.TrackJWTIssued("access_token")
	metrics.TrackJWTIssued("refresh_token")

	return &dto.AuthResponse{
		User: &dto.UserResponse{
			ID:            user.ID.String(),
			Email:         user.Email,
			Username:      user.Username,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Phone:         user.Phone,
			EmailVerified: user.EmailVerified,
			PhoneVerified: user.PhoneVerified,
			KYCVerified:   user.KYCVerified,
			KYCVerifiedAt: user.KYCVerifiedAt,
			Role:          user.Role,
			CreatedAt:     user.CreatedAt,
		},
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// Register creates a new user account (supports full registration and email-only)
func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if email exists
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err == nil {
		// User exists, update settlement account if provided and not full registration
		if req.Password == "" && req.SettlementBankName != "" {
			user.SettlementBankName = req.SettlementBankName
			user.SettlementAccountNumber = req.SettlementAccountNumber
			user.SettlementAccountName = req.SettlementAccountName
			if err := s.userRepo.UpdateUser(user); err != nil {
				return nil, errors.New("failed to update settlement account")
			}
		}
		
		// If it's a "silent" registration for an existing user, just return success
		if req.Password == "" {
			return &dto.AuthResponse{
				User: mapUserToResponse(user),
			}, nil
		}
		
		return nil, errors.New("email already exists")
	}

	// Prepare user fields
	var hashedPassword string
	hasSetPassword := false
	if req.Password != "" {
		hashed, err := password.HashPassword(req.Password, nil)
		if err != nil {
			return nil, errors.New("failed to process password")
		}
		hashedPassword = hashed
		hasSetPassword = true
	}

	// Infer names if not provided
	firstName := req.FirstName
	lastName := req.LastName
	if firstName == "" {
		emailParts := strings.Split(req.Email, "@")
		usernamePart := emailParts[0]
		nameParts := strings.FieldsFunc(usernamePart, func(r rune) bool {
			return r == '.' || r == '_' || r == '-'
		})
		firstName = "User"
		if len(nameParts) > 0 {
			firstName = strings.Title(strings.ToLower(nameParts[0]))
		}
		if len(nameParts) > 1 && lastName == "" {
			lastName = strings.Title(strings.ToLower(nameParts[1]))
		}
	}

	// Generate a unique username if not provided
	username := req.Username
	if username == "" {
		emailParts := strings.Split(req.Email, "@")
		baseUsername := emailParts[0]
		username = baseUsername
		counter := 1
		for {
			exists, _ := s.userRepo.UsernameExists(username)
			if !exists {
				break
			}
			username = baseUsername + string(rune(counter))
			counter++
		}
	}

	// Create user
	user = &models.User{
		Email:                   req.Email,
		Username:                username,
		PasswordHash:            hashedPassword,
		FirstName:               firstName,
		LastName:                lastName,
		Phone:                   req.Phone,
		HasSetPassword:          hasSetPassword,
		SettlementBankName:      req.SettlementBankName,
		SettlementAccountNumber: req.SettlementAccountNumber,
		SettlementAccountName:   req.SettlementAccountName,
		Role:                    models.UserRoleUser,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, errors.New("registration failed")
	}

	// Publish user signed up event
	if s.eventService != nil {
		if err := s.eventService.PublishUserSignedUp(user); err != nil {
			// Log error but don't fail registration
		}
	}

	// If no password set, don't issue tokens (they need to set password later)
	if !hasSetPassword {
		return &dto.AuthResponse{
			User: mapUserToResponse(user),
		}, nil
	}

	// Generate token pair
	roles := []string{string(user.Role)}
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID.String(), user.Email, roles)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Create session
	refreshTokenHash := s.hashToken(tokenPair.RefreshToken)
	session := &models.Session{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		Metadata: map[string]interface{}{
			"registration_time": time.Now(),
		},
	}

	if err := s.sessionRepo.CreateSession(session); err != nil {
		return nil, errors.New("failed to create session")
	}

	return &dto.AuthResponse{
		User:         mapUserToResponse(user),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}


// RefreshToken generates new access token using refresh token
func (s *AuthService) RefreshToken(req *dto.RefreshRequest) (*jwt.TokenPair, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if session exists
	refreshTokenHash := s.hashToken(req.RefreshToken)
	session, err := s.sessionRepo.GetSessionByTokenHash(refreshTokenHash)
	if err != nil {
		return nil, errors.New("invalid session")
	}

	// Check if session is expired
	if session.IsExpired() {
		s.sessionRepo.DeleteSession(refreshTokenHash)
		return nil, errors.New("session expired")
	}

	// Get user
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get user role
	roles := []string{string(user.Role)}

	// Generate new access token (keep same refresh token)
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID.String(), user.Email, roles)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	return &jwt.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken, // Keep the same refresh token
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
	}, nil
}

// ValidateAccessToken validates an access token and returns user info
func (s *AuthService) ValidateAccessToken(tokenString string) (*dto.Claims, error) {
	claims, err := s.jwtService.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &dto.Claims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Roles:  claims.Roles,
	}, nil
}

// Logout invalidates a refresh token session
func (s *AuthService) Logout(refreshToken string) error {
	refreshTokenHash := s.hashToken(refreshToken)
	return s.sessionRepo.DeleteSession(refreshTokenHash)
}

// LogoutAllSessions invalidates all sessions for a user
func (s *AuthService) LogoutAllSessions(userID uuid.UUID) error {
	return s.sessionRepo.DeleteUserSessions(userID)
}


// ForgotPassword initiates password reset process
func (s *AuthService) ForgotPassword(req *dto.ForgotPasswordRequest) error {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		// Return success even if user doesn't exist for security
		return nil
	}

	// Generate reset token
	resetToken := uuid.New().String()
	tokenHash := s.hashToken(resetToken)

	// Create password reset token
	passwordResetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
		Used:      false,
	}

	if err := s.userRepo.CreatePasswordResetToken(passwordResetToken); err != nil {
		return errors.New("failed to create reset token")
	}

	// Publish password reset requested event
	if s.eventService != nil {
		if err := s.eventService.PublishPasswordResetRequested(user, resetToken); err != nil {
			// Log error but don't fail the request
		}
	}

	return nil
}

// ResetPassword resets user password using token
func (s *AuthService) ResetPassword(req *dto.ResetPasswordRequest) error {
	tokenHash := s.hashToken(req.Token)

	// Get password reset token
	resetToken, err := s.userRepo.GetPasswordResetToken(tokenHash)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// Check if token is expired
	if resetToken.IsExpired() {
		return errors.New("token has expired")
	}

	// Get user
	user, err := s.userRepo.GetUserByID(resetToken.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Hash new password
	hashedPassword, err := password.HashPassword(req.NewPassword, nil)
	if err != nil {
		return errors.New("failed to process password")
	}

	// Update user password
	user.PasswordHash = hashedPassword
	if err := s.userRepo.UpdateUser(user); err != nil {
		return errors.New("failed to update password")
	}

	// Mark token as used
	if err := s.userRepo.MarkPasswordResetTokenUsed(tokenHash); err != nil {
		// Log error but don't fail since password was already updated
	}

	// Logout all sessions for security
	s.LogoutAllSessions(user.ID)

	return nil
}

// hashToken creates a SHA256 hash of a token for storage
func (s *AuthService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// mapUserToResponse maps user model to user response DTO
func mapUserToResponse(user *models.User) *dto.UserResponse {
	if user == nil {
		return nil
	}
	return &dto.UserResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		Username:      user.Username,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Phone:         user.Phone,
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
		KYCVerified:   user.KYCVerified,
		KYCVerifiedAt: user.KYCVerifiedAt,
		Role:          user.Role,
		CreatedAt:     user.CreatedAt,
	}
}

// SetPassword handles first-time password setup
func (s *AuthService) SetPassword(req *dto.SetPasswordRequest) error {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.HasSetPassword {
		return errors.New("password already set")
	}

	hashedPassword, err := password.HashPassword(req.Password, nil)
	if err != nil {
		return errors.New("failed to process password")
	}

	user.PasswordHash = hashedPassword
	user.HasSetPassword = true

	return s.userRepo.UpdateUser(user)
}

// UpdateProfile updates user profile
func (s *AuthService) UpdateProfile(userID uuid.UUID, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, errors.New("failed to update profile")
	}

	return mapUserToResponse(user), nil
}

