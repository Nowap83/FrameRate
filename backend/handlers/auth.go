package handlers

import (
    "net/http"

    "github.com/Nowap83/FrameRate/backend/models"
    "github.com/Nowap83/FrameRate/backend/utils"  
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type AuthHandler struct {
    DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
    return &AuthHandler{DB: db}
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

    // create user
    user := models.User{
    		Username:		req.Username,
    		Email:			req.Email,
				IsVerified:	false,
    }

    if err := user.HashPassword(req.Password); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
        return
    }

    if err := h.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // genere et renvoie token pour co directe
    tokenString, err := utils.GenerateToken(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusCreated, AuthResponse{
        Token: tokenString,
        User:  user.ToResponse(),
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
