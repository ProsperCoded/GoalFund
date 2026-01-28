package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/gofund/shared/jwt"
	"github.com/gofund/shared/models"
	"github.com/gofund/shared/password"
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

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50"`
	Phone     string `json:"phone"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         *UserResponse    `json:"user"`
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	TokenType    string           `json:"token_type"`
	ExpiresIn    int64            `json:"expires_in"`
}

// UserResponse represents user data in response
type UserResponse struct {
	ID            string           `json:"id"`
	Email         string           `json:"email"`
	Username      string           `json:"username"`
	FirstName     string           `json:"first_name"`
	LastName      string           `json:"last_name"`
	Phone         string           `json:"phone"`
	EmailVerified bool             `json:"email_verified"`
	PhoneVerified bool             `json:"phone_verified"`
	Role          models.UserRole  `json:"role"`
	CreatedAt     time.Time        `json:"created_at"`
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" binding:"omitempty,min=2,max=50"`
	Phone     string `json:"phone" binding:"omitempty"`
}

// ForgotPasswordRequest represents a forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents a reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	valid, err := password.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, errors.New("authentication failed")
	}
	if !valid {
		return nil, errors.New("invalid credentials")
	}

	// Generate token pair with user role
	roles := []string{string(user.Role)}
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID.String(), user.Email, roles)
	if err != nil {
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
		return nil, errors.New("failed to create session")
	}

	return &AuthResponse{
		User: &UserResponse{
			ID:            user.ID.String(),
			Email:         user.Email,
			Username:      user.Username,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Phone:         user.Phone,
			EmailVerified: user.EmailVerified,
			PhoneVerified: user.PhoneVerified,
			Role:          user.Role,
			CreatedAt:     user.CreatedAt,
		},
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// Register creates a new user account
func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if email exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, errors.New("registration failed")
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Check if username exists
	exists, err = s.userRepo.UsernameExists(req.Username)
	if err != nil {
		return nil, errors.New("registration failed")
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := password.HashPassword(req.Password, nil)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Role:         models.UserRoleUser, // Default role
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, errors.New("registration failed")
	}

	// Publish user signed up event
	if s.eventService != nil {
		if err := s.eventService.PublishUserSignedUp(user); err != nil {
			// Log error but don't fail registration
			// In production, you might want to use a proper logger
		}
	}

	// Generate token pair with user role
	roles := []string{string(user.Role)}

	// Generate token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID.String(), user.Email, roles)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Create session
	refreshTokenHash := s.hashToken(tokenPair.RefreshToken)
	session := &models.Session{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
		Metadata: map[string]interface{}{
			"registration_time": time.Now(),
			"ip":                "", // Will be set by controller
			"user_agent":        "", // Will be set by controller
		},
	}

	if err := s.sessionRepo.CreateSession(session); err != nil {
		return nil, errors.New("failed to create session")
	}

	return &AuthResponse{
		User: &UserResponse{
			ID:            user.ID.String(),
			Email:         user.Email,
			Username:      user.Username,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Phone:         user.Phone,
			EmailVerified: user.EmailVerified,
			PhoneVerified: user.PhoneVerified,
			Role:          user.Role,
			CreatedAt:     user.CreatedAt,
		},
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// RefreshToken generates new access token using refresh token
func (s *AuthService) RefreshToken(req *RefreshRequest) (*jwt.TokenPair, error) {
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
func (s *AuthService) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := s.jwtService.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &Claims{
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

// Claims represents user claims for validation
type Claims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
}

// GetProfile retrieves user profile by ID
func (s *AuthService) GetProfile(userID string) (*UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &UserResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		Username:      user.Username,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Phone:         user.Phone,
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
		Role:          user.Role,
		CreatedAt:     user.CreatedAt,
	}, nil
}

// UpdateProfile updates user profile
func (s *AuthService) UpdateProfile(userID string, req *UpdateProfileRequest) (*UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
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

	return &UserResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		Username:      user.Username,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Phone:         user.Phone,
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
		Role:          user.Role,
		CreatedAt:     user.CreatedAt,
	}, nil
}

// ForgotPassword initiates password reset process
func (s *AuthService) ForgotPassword(req *ForgotPasswordRequest) error {
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
func (s *AuthService) ResetPassword(req *ResetPasswordRequest) error {
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