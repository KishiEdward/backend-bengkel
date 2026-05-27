package services

import (
	"github.com/KishiEdward/backend-bengkel.git/models"
	"github.com/KishiEdward/backend-bengkel.git/repositories"
)

type CustomerService struct {
	customerRepo *repositories.CustomerRepository
}

func NewCustomerService() *CustomerService {
	return &CustomerService{
		customerRepo:
			repositories.NewCustomerRepository(),
	}
}

// =========================
// CREATE CUSTOMER
// =========================
func (s *CustomerService) CreateCustomer(
	customer *models.Customer,
) error {

	return s.customerRepo.
		Create(customer)
}

// =========================
// GET ALL CUSTOMER
// =========================
func (s *CustomerService) GetAll() (
	[]models.Customer,
	error,
) {

	return s.customerRepo.
		FindAll()
}