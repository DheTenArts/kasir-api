package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"

)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	// ! Inisialisasi  substotal -> jumlah total transaksi keseluruhan
	totalAmount := 0
	// ! Inisialisasi modeling  transactionDetail -> nanti kita insert ke DB
	details := make([]models.TransactionDetail, 0)
	// ? Loop setiap items
	for _, item := range items {
		var productName string
		var productID, price, stock int

		// ? get product dapat price
		err := tx.QueryRow("SELECT id, name, price, stock FROM product WHERE id = $1", item.ProductID).Scan(&productID, &productName, &price, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// ? hitung substotal = cuantity * pricing
		substotal := item.Quantity * price
		// ? Ditambahkan di totalAmount
		totalAmount += substotal

		// ? kurangi jumlah stok yang di transaction
		_, err = tx.Exec("UPDATE product SET stock  =  stock - $1 WHERE id = $2", item.Quantity, productID)
		if err != nil {
			return nil, err
		}

		// ? Itemnya di tambahkan ke transactionDetail
		details = append(details, models.TransactionDetail{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    substotal,
		})

	}

	// ! insert transaction
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// ! insert transaction details
	stmt, err := tx.Prepare("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, d := range details {
		if d.Quantity <= 0 || d.Subtotal <= 0 || d.ProductID == 0 {
			return nil, fmt.Errorf("invalid detail data")
		}
		_, err = stmt.Exec(transactionID, d.ProductID, d.Quantity, d.Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
