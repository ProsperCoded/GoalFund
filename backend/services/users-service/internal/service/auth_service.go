package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gofund/shared/jwt"
	"github.com/gofund/shared/models"
	"github.com/gofund/shared/password"
	"github.com/gofund/users-service/internal/repository"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	jwtService  *jwt.JWTService
}

// NewAuthService creates a new auth service instance
func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, jwtService *jwt.JWTService) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtService:  jwtService,
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
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Phone         string    `json:"phone"`
	EmailVerified bool      `json:"email_verified"`
	PhoneVerified bool      `json:"phone_verified"`
	Roles         []string  `json:"roles"`
	CreatedAt     time.Time `json:"created_at"`
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

	// Get user roles
	roles, err := s.userRepo.GetUserRoles(user.ID)
	if err != nil {
		roles = []string{} // Default to empty roles if error
	}

	// Generate token pair
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
			Roles:         roles,
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
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, errors.New("registration failed")
	}

	// Assign default user role
	defaultRole, err := s.userRepo.GetRoleByName("user")
	if err == nil {
		s.userRepo.AssignRole(user.ID, defaultRole.ID)
	}

	// Get user roles
	roles, err := s.userRepo.GetUserRoles(user.ID)
	if err != nil {
		roles = []string{"user"} // Default role
	}

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
			Roles:         roles,
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

	// Get user roles
	roles, err := s.userRepo.GetUserRoles(user.ID)
	if err != nil {
		roles = []string{} // Default to empty roles if error
	}

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

// hashToken creates a SHA256 hash of a token for storage
func (s *AuthService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}