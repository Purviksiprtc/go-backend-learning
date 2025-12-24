package handlers

import (
	"net/http"

	"user-crud-go/config"
	"user-crud-go/models"
	"user-crud-go/request"
	"user-crud-go/response"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c echo.Context) error {
	req := new(request.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	// Email uniqueness check
	var existingUser models.User
	if err := config.DB.
		Where("email = ? AND deleted_at IS NULL", req.Email).
		First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusConflict, "Email already in use")
	}

	// ✅ FIX: handle bcrypt error
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to hash password")
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hash),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to create user")
	}

	resp := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return c.JSON(http.StatusCreated, resp)
}

func GetUsers(c echo.Context) error {
	var users []models.User
	var resp []response.UserResponse

	if err := config.DB.
		Where("deleted_at IS NULL").
		Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to fetch users")
	}

	for _, u := range users {
		resp = append(resp, response.UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func GetUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	if err := config.DB.
		Where("deleted_at IS NULL").
		First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	resp := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return c.JSON(http.StatusOK, resp)
}

func UpdateUser(c echo.Context) error {
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	req := new(request.UpdateUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	// ✅ Email uniqueness check (excluding self)
	var existingUser models.User
	if err := config.DB.
		Where("email = ? AND id != ? AND deleted_at IS NULL", req.Email, user.ID).
		First(&existingUser).Error; err == nil {

		return c.JSON(http.StatusConflict, "Email already in use")
	}

	user.Name = req.Name
	user.Email = req.Email

	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to update user")
	}

	resp := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return c.JSON(http.StatusOK, resp)
}

func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to delete user")
	}

	return c.JSON(http.StatusOK, "User soft deleted")
}
