package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// ! ambil semua isi product dengan relasi kategori
func (repo *ProductRepository) GetAll(name string) ([]models.Product, error) {
	query := "SELECT p.id, p.name, p.price, p.stock, COALESCE(p.category_id, 0), COALESCE(c.id, 0), COALESCE(c.name, '') FROM product p LEFT JOIN category_product c ON p.category_id = c.id"
	
	var args [] interface{}
	if name != "" {
		query += " WHERE p.name ILIKE $1 ORDER BY p.id ASC"
		args = append(args, "%" + name + "%")
	}
	
	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		var c models.Category
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		p.Category = c
		products = append(products, p)
	}
	return products, nil
}

// ! masukkan item produk (termasuk category_id)
func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO product (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

// ! ambil isi produk berdasarkan id (dengan join kategori)
func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	var p models.Product
	var c models.Category
	query := "SELECT p.id, p.name, p.price, p.stock, COALESCE(p.category_id, 0), COALESCE(c.id, 0), COALESCE(c.name, '') FROM product p LEFT JOIN category_product c ON p.category_id = c.id WHERE p.id = $1"
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &c.ID, &c.Name)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	p.Category = c
	return &p, nil
}

// ! update isi produk berdasarkan id
func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE product SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return nil
}

// ! delete isi produk berdasarkan id
func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM product WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return err
}

// ? ============================================ Category ============================================

// ! ambil semua isi category_produk
func (repo *ProductRepository) GetAllCategory() ([]models.Category, error) {
	query := "SELECT id, name FROM category_product"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kategori := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		kategori = append(kategori, c)
	}

	return kategori, nil
}

// ! masukkan item category
func (repo *ProductRepository) CreateCategory(kategori *models.Category) error {
	query := "INSERT INTO category_product (name) VALUES ($1) RETURNING id"
	err := repo.db.QueryRow(query, kategori.Name).Scan(&kategori.ID)
	return err
}

// ! ambil isi category berdasarkan id
func (repo *ProductRepository) GetByIDCategory(id int) (*models.Category, error) {
	query := "SELECT id, name FROM category_product WHERE id = $1"

	var c models.Category
	err := repo.db.QueryRow(query, id).Scan(&c.ID, &c.Name)
	if err == sql.ErrNoRows {
		return nil, errors.New("category tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// ! update isi category berdasarkan id
func (repo *ProductRepository) UpdateCategory(category *models.Category) error {
	query := "UPDATE category_product SET name = $1 WHERE id = $2"
	result, err := repo.db.Exec(query, category.Name, category.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("category tidak ditemukan")
	}

	return nil
}

// ! delete isi category berdasarkan id
func (repo *ProductRepository) DeleteCategory(id int) error {
	query := "DELETE FROM category_product WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("category tidak ditemukan")
	}

	return err
}
