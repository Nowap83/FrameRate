package handlers

import (
    "net/http"
		"encoding/hex"
		"time"
		"crypto/rand"

    "github.com/Nowap83/FrameRate/backend/models"
    "github.com/Nowap83/FrameRate/backend/utils"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type AuthHandler struct {
    DB	*gorm.DB
		EmailService *utils.EmailService
}

func NewAuthHandler(db *gorm.DB, emailService *utils.EmailService) *AuthHandler {
    return &AuthHandler{
			DB: db,
			EmailService: emailService,
		}
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50,username"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8,strongpassword"`
}

type LoginRequest struct {
    Login    string `json:"login" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    Token string              `json:"token"`
    User  models.UserResponse `json:"user"`
}

func generateVerificationToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // check email existe
    var existingUser models.User
    if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
        return
    }

    // check username existe
    if err := h.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
        return
    }

		// generation token mail
		verificationToken, err := generateVerificationToken()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
        return
    }

    expiresAt := time.Now().Add(24 * time.Hour)

    // create user
    user := models.User{
    		Username:		req.Username,
    		Email:			req.Email,
				IsVerified:	false,
				VerificationToken: &verificationToken,
				TokenExpiresAt: &expiresAt,
    }

    if err := user.HashPassword(req.Password); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
        return
    }

    if err := h.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

		// envoit du mail verif
    if err := h.EmailService.SendVerificationEmail(user.Email, user.Username, verificationToken); err != nil {
				// cas user created mais mail pas send
        c.JSON(http.StatusCreated, gin.H{
            "message": "Account created but verification email failed. Please contact support.",
            "user": gin.H{
                "id":       user.ID,
                "username": user.Username,
                "email":    user.Email,
            },
        })
        return
    }
		
		// reponse si bien send
		c.JSON(http.StatusCreated, gin.H{
        "message": "Registration successful! Please check your email to verify your account.",
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

		// check email OU username
    var user models.User
    if err := h.DB.Where("email = ? OR username = ?", req.Login, req.Login).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if !user.ValidatePassword(req.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
		
		// bloque login si mail pas verif
		if !user.IsVerified {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "Please verify your email before logging in. Check your inbox.",
        })
        return
    }

    tokenString, err := utils.GenerateToken(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, AuthResponse{
        Token: tokenString,
        User:  user.ToResponse(),
    })
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
    token := c.Query("token")

    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
        return
    }

    var user models.User

    // trouve user avec token valide
    if err := h.DB.Where("verification_token = ? AND token_expires_at > ?", token, time.Now()).First(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification token"})
        return
    }

    // cas deja verif
    if user.IsVerified {
        c.JSON(http.StatusOK, gin.H{"message": "Email already verified"})
        return
    }

    // marque comme verif
    user.IsVerified = true
    user.VerificationToken = nil
    user.TokenExpiresAt = nil

    if err := h.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
        return
    }

    // genere token JWT pour co directe
    tokenString, err := utils.GenerateToken(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate login token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Email verified successfully! You are now logged in.",
        "token":   tokenString,
        "user":    user.ToResponse(),
    })
}
