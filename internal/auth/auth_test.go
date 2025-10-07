package auth

import (
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"learn-go-capstone/internal/database"
)

func TestAuthService(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_auth.db"
	defer os.Remove(tempDB)
	
	// Create database configuration
	config := &database.Config{
		Driver: "sqlite3",
		DSN:    tempDB,
	}
	
	// Connect to database
	db, err := database.Connect(config)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)
	
	// Run migrations
	migrationManager := database.NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	
	// Create repository and auth service
	repository := database.NewSQLiteRepository(db)
	authService := NewAuthService(repository)
	
	// Test user registration
	user, err := authService.RegisterUser("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	
	if user.ID == 0 {
		t.Error("Expected user ID to be set after registration")
	}
	
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}
	
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
	
	if !user.IsActive {
		t.Error("Expected user to be active after registration")
	}
	
	// Test duplicate username registration
	_, err = authService.RegisterUser("testuser", "test2@example.com", "password123")
	if err == nil {
		t.Error("Expected error for duplicate username")
	}
	
	// Test duplicate email registration
	_, err = authService.RegisterUser("testuser2", "test@example.com", "password123")
	if err == nil {
		t.Error("Expected error for duplicate email")
	}
	
	// Test user login
	token, loggedInUser, err := authService.LoginUser("testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}
	
	if token == "" {
		t.Error("Expected token to be generated")
	}
	
	if loggedInUser.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, loggedInUser.ID)
	}
	
	// Test invalid login
	_, _, err = authService.LoginUser("testuser", "wrongpassword")
	if err == nil {
		t.Error("Expected error for invalid password")
	}
	
	_, _, err = authService.LoginUser("nonexistent", "password123")
	if err == nil {
		t.Error("Expected error for nonexistent user")
	}
	
	// Test token validation
	claims, err := authService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}
	
	if claims.UserID != user.ID {
		t.Errorf("Expected user ID %d in claims, got %d", user.ID, claims.UserID)
	}
	
	if claims.Username != user.Username {
		t.Errorf("Expected username %s in claims, got %s", user.Username, claims.Username)
	}
	
	// Test getting user from token
	userFromToken, err := authService.GetUserFromToken(token)
	if err != nil {
		t.Fatalf("Failed to get user from token: %v", err)
	}
	
	if userFromToken.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, userFromToken.ID)
	}
	
	// Test password change
	err = authService.ChangePassword(user.ID, "password123", "newpassword123")
	if err != nil {
		t.Fatalf("Failed to change password: %v", err)
	}
	
	// Test login with new password
	_, _, err = authService.LoginUser("testuser", "newpassword123")
	if err != nil {
		t.Fatalf("Failed to login with new password: %v", err)
	}
	
	// Test login with old password (should fail)
	_, _, err = authService.LoginUser("testuser", "password123")
	if err == nil {
		t.Error("Expected error when logging in with old password")
	}
	
	// Test invalid old password
	err = authService.ChangePassword(user.ID, "wrongpassword", "anotherpassword")
	if err == nil {
		t.Error("Expected error for invalid old password")
	}
	
	// Test token refresh
	newToken, err := authService.RefreshToken(token)
	if err != nil {
		t.Fatalf("Failed to refresh token: %v", err)
	}
	
	if newToken == "" {
		t.Error("Expected new token to be generated")
	}
	
	// Test user deactivation
	err = authService.DeactivateUser(user.ID)
	if err != nil {
		t.Fatalf("Failed to deactivate user: %v", err)
	}
	
	// Test login with deactivated user
	_, _, err = authService.LoginUser("testuser", "newpassword123")
	if err == nil {
		t.Error("Expected error when logging in with deactivated user")
	}
	
	// Test user activation
	err = authService.ActivateUser(user.ID)
	if err != nil {
		t.Fatalf("Failed to activate user: %v", err)
	}
	
	// Test login with reactivated user
	_, _, err = authService.LoginUser("testuser", "newpassword123")
	if err != nil {
		t.Fatalf("Failed to login with reactivated user: %v", err)
	}
}

func TestJWTTokenExpiration(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_jwt_expiration.db"
	defer os.Remove(tempDB)
	
	// Create database configuration
	config := &database.Config{
		Driver: "sqlite3",
		DSN:    tempDB,
	}
	
	// Connect to database
	db, err := database.Connect(config)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)
	
	// Run migrations
	migrationManager := database.NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	
	// Create repository and auth service
	repository := database.NewSQLiteRepository(db)
	authService := NewAuthService(repository)
	
	// Register and login user
	_, err = authService.RegisterUser("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	
	token, _, err := authService.LoginUser("testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}
	
	// Test valid token
	_, err = authService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate valid token: %v", err)
	}
	
	// Test invalid token
	_, err = authService.ValidateToken("invalid.token.here")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
	
	// Test empty token
	_, err = authService.ValidateToken("")
	if err == nil {
		t.Error("Expected error for empty token")
	}
}

func TestPasswordHashing(t *testing.T) {
	password := "testpassword123"
	
	// Test password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	if string(hashedPassword) == password {
		t.Error("Hashed password should not be the same as original password")
	}
	
	// Test password verification
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		t.Fatalf("Failed to verify password: %v", err)
	}
	
	// Test wrong password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrongpassword"))
	if err == nil {
		t.Error("Expected error for wrong password")
	}
}
