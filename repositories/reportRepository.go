package repositories

import (
	"database/sql"
	"kasir-api/models"

)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetDayReport(tanggalMulai, tanggalAkhir string) (*models.Report, error) {

	// ? total penghasilan
	var totalRevenue , totalTransaksi, banyakProduk  int

	// ? terlaris
	var terlaris string


	// ! hitung total harian
	err := repo.db.QueryRow("SELECT COALESCE(SUM(total_amount),0), COUNT(*) FROM transactions WHERE created_at::date BETWEEN $1 AND $2", tanggalMulai, tanggalAkhir).Scan(&totalRevenue, &totalTransaksi)
	if err != nil {
		return nil, err
	}

	// ! ambil produk terlair
	err = repo.db.QueryRow("SELECT p.name, SUM(td.quantity) AS qty FROM transaction_details td LEFT JOIN transactions t ON td.transaction_id = t.id LEFT JOIN product p ON td.product_id = p.id WHERE t.created_at::date BETWEEN $1 AND $2 GROUP BY p.name", tanggalMulai, tanggalAkhir).Scan(&terlaris, & banyakProduk)
	if err == sql.ErrNoRows {
		terlaris = ""
		banyakProduk = 0
	} else if err != nil {
		return nil, err
	}

	reports := &models.Report{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: models.BestProduct{
			Name:    terlaris,
			QtySold: banyakProduk,
		},
	}

	return reports, nil
}