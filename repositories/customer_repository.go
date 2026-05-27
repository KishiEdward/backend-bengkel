package repositories

import (
	"github.com/KishiEdward/backend-bengkel.git/config"
	"github.com/KishiEdward/backend-bengkel.git/models"
)

type CustomerRepository struct{}

func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{}
}

// =========================
// CREATE CUSTOMER
// =========================
func (r *CustomerRepository) Create(
	customer *models.Customer,
) error {

	return config.DB.
		Create(customer).Error
}

// =========================
// GET ALL CUSTOMER
// =========================
func (r *CustomerRepository) FindAll() (
	[]models.Customer,
	error,
) {

	var customers []models.Customer

	err := config.DB.
		Order("nama ASC").
		Find(&customers).Error

	return customers, err
}