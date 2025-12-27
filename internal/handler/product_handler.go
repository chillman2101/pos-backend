package handler

import (
	"math"
	"pos-backend/internal/dto"
	"pos-backend/internal/service"
	"pos-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Get filter parameters
	search := c.Query("search")
	categoryID := c.Query("category_id")

	var products []*dto.ProductResponse
	var total int64
	var err error

	// Use filtered query if search or category_id is provided
	if search != "" || categoryID != "" {
		products, total, err = h.productService.GetAllWithFilter(search, categoryID, page, limit)
	} else {
		products, total, err = h.productService.GetAll(page, limit)
	}

	if err != nil {
		response.InternalServerError(c, "Failed to get products", err.Error())
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response.SuccessWithPagination(c, "Products retrieved successfully", products, response.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalRows:  total,
		TotalPages: totalPages,
	})
}

func (h *ProductHandler) GetByCategory(c *gin.Context) {
	categoryID := c.Param("category_id")

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, total, err := h.productService.GetByCategory(categoryID, page, limit)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response.SuccessWithPagination(c, "Products retrieved successfully", products, response.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalRows:  total,
		TotalPages: totalPages,
	})
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetByID(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Product retrieved successfully", product)
}

func (h *ProductHandler) GetBySKU(c *gin.Context) {
	sku := c.Param("sku")

	product, err := h.productService.GetBySKU(sku)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Product retrieved successfully", product)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	product, err := h.productService.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Product created successfully", product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	product, err := h.productService.Update(id, &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Product updated successfully", product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.productService.Delete(id); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Product deleted successfully", nil)
}
