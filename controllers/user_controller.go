package controllers

import (
	"net/http"
	"strconv"
	"tusk/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

// Request structs untuk input yang aman
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Response structs untuk output yang aman (tanpa password)
type UserResponse struct {
	Id        int    `json:"id"`
	Role      string `json:"role"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (u *UserController) Login(c *gin.Context) {
	var loginReq LoginRequest

	// Bind dan validasi input
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Cari user berdasarkan email
	errDB := u.DB.Where("email = ?", loginReq.Email).First(&user).Error
	if errDB != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or Password is Wrong"})
		return
	}

	// Verifikasi password
	errHash := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(loginReq.Password),
	)
	if errHash != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or Password is Wrong"})
		return
	}

	// Return user data tanpa password
	userResponse := UserResponse{
		Id:        user.Id,
		Role:      user.Role,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    userResponse,
	})
}

func (u *UserController) CreateAccount(c *gin.Context) {
	var createReq CreateUserRequest

	// Bind dan validasi input
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah email sudah ada
	var existingUser models.User
	if u.DB.Where("email = ?", createReq.Email).First(&existingUser).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(createReq.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Buat user baru
	newUser := models.User{
		Name:     createReq.Name,
		Email:    createReq.Email,
		Password: string(hashedPasswordBytes),
		Role:     "Employee",
	}

	errDB := u.DB.Create(&newUser).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	// Return response tanpa password
	userResponse := UserResponse{
		Id:        newUser.Id,
		Role:      newUser.Role,
		Name:      newUser.Name,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: newUser.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    userResponse,
	})
}

func (u *UserController) Delete(c *gin.Context) {
	idParam := c.Param("id")

	// Validasi ID
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Cek apakah user ada
	var user models.User
	if u.DB.First(&user, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Hapus user
	errDB := u.DB.Delete(&user).Error
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"deletedUser": gin.H{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func (u *UserController) GetEmployee(c *gin.Context) {
	var users []models.User

	errDB := u.DB.Select("id, name, email, role, created_at, updated_at").
		Where("role = ?", "Employee").
		Find(&users).Error

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	// Convert ke response format
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			Id:        user.Id,
			Role:      user.Role,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Employees retrieved successfully",
		"count":     len(userResponses),
		"employees": userResponses,
	})
}
