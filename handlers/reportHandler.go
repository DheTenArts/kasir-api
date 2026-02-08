package handlers

import (
	"encoding/json"
	"kasir-api/services"
	"net/http"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) HandleReports(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetDayReport(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h ReportHandler) GetDayReport(w http.ResponseWriter, r *http.Request) {
	tanggalMulai := r.URL.Query().Get("tanggal_mulai")
	tanggalAkhir := r.URL.Query().Get("tanggal_akhir")
	// fmt.Print(name)
	reports, err := h.service.GetDayReport(tanggalMulai, tanggalAkhir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}