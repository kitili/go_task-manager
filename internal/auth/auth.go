package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"learn-go-capstone/internal/database"
)

// JWTSecret is the secret key for JWT tokens
// In production, this should be loaded from environment variables
var JWTSecret = []byte("your-secret-key-change-in-production")

// Claims represents the JWT claims
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// AuthService handles authentication operations
type AuthService struct {
	repository database.Repository
}

// NewAuthService creates a new authentication service
func NewAuthService(repository database.Repository) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

// RegisterUser registers a new user
func (as *AuthService) RegisterUser(username, email, password string) (*database.User, error) {
	// Validate input
	if username == "" || email == "" || password == "" {
		return nil, errors.New("username, email, and password are required")
	}
	
	// Check if username already exists
	_, err := as.repository.GetUserByUsername(username)
	if err == nil {
		return nil, errors.New("username already exists")
	}
	
	// Check if email already exists
	_, err = as.repository.GetUserByEmail(email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Create user
	user := &database.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		IsActive: true,
	}
	
	err = as.repository.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	return user, nil
}

// LoginUser authenticates a user and returns a JWT token
func (as *AuthService) LoginUser(username, password string) (string, *database.User, error) {
	// Get user by username
	user, err := as.repository.GetUserByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid username or password")
	}
	
	// Check if user is active
	if !user.IsActive {
		return "", nil, errors.New("user account is deactivated")
	}
	
	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid username or password")
	}
	
	// Generate JWT token
	token, err := as.generateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}
	
	return token, user, nil
}

// ValidateToken validates a JWT token and returns the claims
func (as *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	
	return claims, nil
}

// GetUserFromToken gets user information from a JWT token
func (as *AuthService) GetUserFromToken(tokenString string) (*database.User, error) {
	claims, err := as.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	
	user, err := as.repository.GetUser(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	return user, nil
}

// ChangePassword changes a user's password
func (as *AuthService) ChangePassword(userID int, oldPassword, newPassword string) error {
	// Get user
	user, err := as.repository.GetUser(userID)
	if err != nil {
		return err
	}
	
	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("invalid old password")
	}
	
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}
	
	// Update user
	user.Password = string(hashedPassword)
	err = as.repository.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	
	return nil
}

// generateToken generates a JWT token for a user
func (as *AuthService) generateToken(user *database.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// RefreshToken generates a new token for an existing user
func (as *AuthService) RefreshToken(tokenString string) (string, error) {
	claims, err := as.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	
	// Get user to ensure they still exist and are active
	user, err := as.repository.GetUser(claims.UserID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}
	
	if !user.IsActive {
		return "", errors.New("user account is deactivated")
	}
	
	// Generate new token
	return as.generateToken(user)
}

// DeactivateUser deactivates a user account
func (as *AuthService) DeactivateUser(userID int) error {
	user, err := as.repository.GetUser(userID)
	if err != nil {
		return err
	}
	
	user.IsActive = false
	return as.repository.UpdateUser(user)
}

// ActivateUser activates a user account
func (as *AuthService) ActivateUser(userID int) error {
	user, err := as.repository.GetUser(userID)
	if err != nil {
		return err
	}
	
	user.IsActive = true
	return as.repository.UpdateUser(user)
}
