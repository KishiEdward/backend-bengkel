package handlers

import (
	"net/http"

	"github.com/KishiEdward/backend-bengkel.git/models"
	"github.com/KishiEdward/backend-bengkel.git/services"
	"github.com/gin-gonic/gin"
)

type MaterialHandler struct {
	materialService *services.MaterialService
}

func NewMaterialHandler() *MaterialHandler {
	return &MaterialHandler{
		materialService:
			services.NewMaterialService(),
	}
}

// =========================
// CREATE MATERIAL
// =========================
func (h *MaterialHandler) Create(
	c *gin.Context,
) {

	var material models.Material

	if err := c.ShouldBindJSON(
		&material,
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

	if err := h.materialService.
		CreateMaterial(&material); err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": "Gagal membuat material",
			},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"success": true,
			"data": material,
		},
	)
}

// =========================
// GET ALL MATERIAL
// =========================
func (h *MaterialHandler) GetAll(
	c *gin.Context,
) {

	materials, err :=
		h.materialService.GetAll()

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": "Gagal mengambil material",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data": materials,
		},
	)
}