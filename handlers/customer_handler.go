package handlers

import (
	"net/http"

	"github.com/KishiEdward/backend-bengkel.git/models"
	"github.com/KishiEdward/backend-bengkel.git/services"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerService *services.CustomerService
}

func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{
		customerService:
			services.NewCustomerService(),
	}
}

// =========================
// CREATE CUSTOMER
// =========================
func (h *CustomerHandler) Create(
	c *gin.Context,
) {

	var customer models.Customer

	if err := c.ShouldBindJSON(
		&customer,
	); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)
		return
	}

	if err := h.customerService.
		CreateCustomer(&customer); err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": "Gagal membuat customer",
			},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"success": true,
			"data": customer,
		},
	)
}

// =========================
// GET ALL CUSTOMER
// =========================
func (h *CustomerHandler) GetAll(
	c *gin.Context,
) {

	customers, err :=
		h.customerService.GetAll()

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": "Gagal mengambil customer",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data": customers,
		},
	)
}