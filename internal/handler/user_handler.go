package handler

import (
	"math"
	"pos-backend/internal/dto"
	"pos-backend/internal/service"
	"pos-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, total, err := h.userService.GetAll(page, limit)
	if err != nil {
		response.InternalServerError(c, "Failed to get users", err.Error())
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response.SuccessWithPagination(c, "Users retrieved successfully", users, response.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalRows:  total,
		TotalPages: totalPages,
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetByID(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "User retrieved successfully", user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.Update(id, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "User updated successfully", user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Prevent user from deleting themselves
	userID := c.GetString("user_id")
	if id == userID {
		response.BadRequest(c, "Cannot delete your own account", nil)
		return
	}

	if err := h.userService.Delete(id); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "User deleted successfully", nil)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	id := c.Param("id")

	// Only allow users to change their own password (or admin can change anyone's)
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if id != userID && role != "admin" {
		response.Forbidden(c, "You can only change your own password")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := h.userService.ChangePassword(id, &req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Password changed successfully", nil)
}
