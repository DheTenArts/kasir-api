package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(name string) ([]models.Product, error) {
	return s.repo.GetAll(name)
}

func (s *ProductService) Create(data *models.Product) error {
	return s.repo.Create(data)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(product *models.Product) error {
	return s.repo.Update(product)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}

// ? ============================================ Category ============================================

func (s *ProductService) GetAllCategory() ([]models.Category, error) {
	return s.repo.GetAllCategory()
}

func (s *ProductService) CreateCategory(data *models.Category) error {
	return s.repo.CreateCategory(data)
}

func (s *ProductService) GetByIDCategory(id int) (*models.Category, error) {
	return s.repo.GetByIDCategory(id)
}

func (s *ProductService) UpdateCategory(category *models.Category) error {
	return s.repo.UpdateCategory(category)
}

func (s *ProductService) DeleteCategory(id int) error {
	return s.repo.DeleteCategory(id)
}