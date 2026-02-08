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
	var totalRevenue, totalTransaksi int

	args := []interface{}{}

	// ! hitung total harian
	queryTotal := "SELECT COALESCE(SUM(total_amount),0), COUNT(*) FROM transactions"
	if tanggalMulai != "" && tanggalAkhir != "" {
		queryTotal += " WHERE created_at::date BETWEEN $1 AND $2"
		args = append(args, tanggalMulai, tanggalAkhir)
	}

	err := repo.db.QueryRow(queryTotal, args...).Scan(&totalRevenue, &totalTransaksi)
	if err != nil {
		return nil, err
	}

	// ! ambil produk terlair
	queryLaris := "SELECT p.name, SUM(td.quantity) AS qty FROM transaction_details td LEFT JOIN transactions t ON td.transaction_id = t.id LEFT JOIN product p ON td.product_id = p.id"
	args2 := []interface{}{}
	if tanggalMulai != "" && tanggalAkhir != "" {
		queryLaris += " WHERE t.created_at::date BETWEEN $1 AND $2"
		args2 = append(args2, tanggalMulai, tanggalAkhir)
	}
	queryLaris += " GROUP BY p.name ORDER BY qty DESC"
	if tanggalMulai != "" && tanggalAkhir != "" {
		queryLaris += " LIMIT 1"
	}
	rows, err := repo.db.Query(queryLaris, args2...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	produkTerlaris := []models.BestProduct{}

	for rows.Next() {
		var name string
		var qty int
		if err := rows.Scan(&name, &qty); err != nil {
			return nil, err
		}
		produkTerlaris = append(produkTerlaris, models.BestProduct{
			Name:    name,
			QtySold: qty,
		})
	}

	return &models.Report{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: produkTerlaris,
	}, nil

}
