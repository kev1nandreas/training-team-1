package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"securetask/database"
	"securetask/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"os"
)

// Rate limiter for login attempts
type LoginAttempt struct {
	Count    int
	LastTime time.Time
}

var (
	loginAttempts = make(map[string]*LoginAttempt)
	mu            sync.Mutex
)

func checkRateLimit(email string) error {
	mu.Lock()
	defer mu.Unlock()

	attempt, exists := loginAttempts[email]
	if !exists {
		loginAttempts[email] = &LoginAttempt{Count: 1, LastTime: time.Now()}
		return nil
	}

	// Reset if lockout period has passed
	if time.Since(attempt.LastTime) > LockoutDuration {
		loginAttempts[email] = &LoginAttempt{Count: 1, LastTime: time.Now()}
		return nil
	}

	// Check if max attempts exceeded
	if attempt.Count >= MaxLoginAttempts {
		return fmt.Errorf("too many login attempts. please try again after %d minutes", int(LockoutDuration.Minutes()))
	}

	// Increment count
	attempt.Count++
	attempt.LastTime = time.Now()
	return nil
}

func resetLoginAttempts(email string) {
	mu.Lock()
	defer mu.Unlock()
	delete(loginAttempts, email)
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not found in environment")
	}
	return []byte(secret)
}

func getAdminAPIKey() string {
	key := os.Getenv("ADMIN_API_KEY")
	if key == "" {
		panic("ADMIN_API_KEY not found in environment")
	}
	return key
}

func validatePassword(password string) error {
	if len(password) < 16 {
		return fmt.Errorf("password must be at least 16 characters long")
	}

	hasUpper := false
	hasLower := false
	hasSymbol := false

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if !((char >= '0' && char <= '9') || (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')) {
			hasSymbol = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasSymbol {
		return fmt.Errorf("password must contain at least one symbol")
	}

	return nil
}

func validateEmail(email string) error {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailPattern, email)
	if err != nil || !matched {
		return fmt.Errorf("email must be a valid email format")
	}

	// Check for duplicate email
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err == nil {
		return fmt.Errorf("email already exists")
	}

	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password")
	}
	return string(hashedPassword), nil
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate email format
	if err := validateEmail(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password strength
	if err := validatePassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password using bcrypt
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     sanitizeHTML(req.Name),
		Role:     "user",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// VULNERABILITY #2: No rate limiting - vulnerable to brute force attacks
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	// Check rate limit
	if err := checkRateLimit(req.Email); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Reset rate limit on successful login
	resetLoginAttempts(req.Email)

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  user,
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", uint(claims["user_id"].(float64)))
			c.Set("email", claims["email"].(string))
			c.Set("role", claims["role"].(string))
		}

		c.Next()
	}
}
