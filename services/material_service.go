package services

import (
	"github.com/KishiEdward/backend-bengkel.git/models"
	"github.com/KishiEdward/backend-bengkel.git/repositories"
)

type MaterialService struct {
	materialRepo *repositories.MaterialRepository
}

func NewMaterialService() *MaterialService {
	return &MaterialService{
		materialRepo:
			repositories.NewMaterialRepository(),
	}
}

// =========================
// CREATE MATERIAL
// =========================
func (s *MaterialService) CreateMaterial(
	material *models.Material,
) error {

	return s.materialRepo.
		Create(material)
}

// =========================
// GET ALL MATERIAL
// =========================
func (s *MaterialService) GetAll() (
	[]models.Material,
	error,
) {

	return s.materialRepo.
		FindAll()
}