package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"user-crud-go/config"
	"user-crud-go/models"
	"user-crud-go/rabbitmq"
	"user-crud-go/request"
	"user-crud-go/response"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

/* ---------------- CREATE USER ---------------- */

func CreateUser(c echo.Context) error {
	req := new(request.CreateUserRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	// Structured validation
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	var existingUser models.User
	if err := config.DB.
		Where("email = ? AND deleted_at IS NULL", req.Email).
		First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusConflict, response.APIResponse{
			Status:  "error",
			Message: "Email already in use",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Status:  "error",
			Message: "Failed to hash password",
		})
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hash),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Status:  "error",
			Message: "Failed to create user",
		})
	}

	// USER_CREATED event
	event := rabbitmq.UserEvent{
		Event:     "USER_CREATED",
		Version:   "1.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data: rabbitmq.UserData{
			UserID: user.ID,
			Name:   user.Name,
			Email:  user.Email,
		},
	}

	if err := rabbitmq.PublishEvent("user.created", event); err != nil {
		log.Println("Failed to publish USER_CREATED event:", err)
	}

	return c.JSON(http.StatusCreated, response.APIResponse{
		Status:  "success",
		Message: "User created successfully",
		Data: response.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

/* ---------------- GET ALL USERS ---------------- */

func GetUsers(c echo.Context) error {
	page := 1
	perPage := 10

	name := c.QueryParam("name")
	email := c.QueryParam("email")

	if p := c.QueryParam("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if pp := c.QueryParam("per_page"); pp != "" {
		fmt.Sscanf(pp, "%d", &perPage)
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	var users []models.User
	var total int64

	query := config.DB.Model(&models.User{}).
		Where("deleted_at IS NULL")

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if email != "" {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Status:  "error",
			Message: "Failed to count users",
		})
	}

	if err := query.
		Limit(perPage).
		Offset(offset).
		Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Status:  "error",
			Message: "Failed to fetch users",
		})
	}

	userResponses := make([]response.UserResponse, 0)
	for _, u := range users {
		userResponses = append(userResponses, response.UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	lastPage := int((total + int64(perPage) - 1) / int64(perPage))

	// ✅ UPDATED RESPONSE STRUCTURE (PR FIX)
	return c.JSON(http.StatusOK, response.APIResponse{
		Status:  "success",
		Message: "Users fetched successfully",
		Data: map[string]interface{}{
			"pagination": response.Pagination{
				CurrentPage:  page,
				PerPage:      perPage,
				LastPage:     lastPage,
				TotalResults: int(total),
			},
			"users": userResponses,
		},
	})
}

/* ---------------- GET SINGLE USER ---------------- */

func GetUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	if err := config.DB.
		Where("deleted_at IS NULL").
		First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, response.APIResponse{
			Status:  "error",
			Message: "User not found",
		})
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Status:  "success",
		Message: "User fetched successfully",
		Data: response.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

/* ---------------- UPDATE USER ---------------- */

func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, response.APIResponse{
			Status:  "error",
			Message: "User not found",
		})
	}

	req := new(request.UpdateUserRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	var existingUser models.User
	if err := config.DB.
		Where("email = ? AND id != ? AND deleted_at IS NULL", req.Email, user.ID).
		First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusConflict, response.APIResponse{
			Status:  "error",
			Message: "Email already in use",
		})
	}

	user.Name = req.Name
	user.Email = req.Email

	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.APIResponse{
				Status:  "error",
				Message: "Failed to hash password",
			})
		}
		user.Password = string(hashed)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Status:  "error",
			Message: "Failed to update user",
		})
	}

	// USER_UPDATED event
	event := rabbitmq.UserEvent{
		Event:     "USER_UPDATED",
		Version:   "1.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data: rabbitmq.UserData{
			UserID: user.ID,
			Name:   user.Name,
			Email:  user.Email,
		},
	}

	if err := rabbitmq.PublishEvent("user.updated", event); err != nil {
		log.Println("Failed to publish USER_UPDATED event:", err)
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Status:  "success",
		Message: "User updated successfully",
		Data: response.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

/* ---------------- DELETE USER ---------------- */

func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, response.APIResponse{
			Status:  "error",
			Message: "User not found",
		})
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Status:  "error",
			Message: "Failed to delete user",
		})
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Status:  "success",
		Message: "User deleted successfully",
	})
}
