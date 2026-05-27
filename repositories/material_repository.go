package repositories

import (
	"github.com/KishiEdward/backend-bengkel.git/config"
	"github.com/KishiEdward/backend-bengkel.git/models"
)

type MaterialRepository struct{}

func NewMaterialRepository() *MaterialRepository {
	return &MaterialRepository{}
}

// =========================
// CREATE MATERIAL
// =========================
func (r *MaterialRepository) Create(
	material *models.Material,
) error {

	return config.DB.
		Create(material).Error
}

// =========================
// GET ALL MATERIAL
// =========================
func (r *MaterialRepository) FindAll() (
	[]models.Material,
	error,
) {

	var materials []models.Material

	err := config.DB.
		Order("nama ASC").
		Find(&materials).Error

	return materials, err
}