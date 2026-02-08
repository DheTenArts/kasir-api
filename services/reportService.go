package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetDayReport(tanggalMulai, tanggalAkhir string) (*models.Report, error) {
	return s.repo.GetDayReport(tanggalMulai, tanggalAkhir)
}