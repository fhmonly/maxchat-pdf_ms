package pdf

import (
	"encoding/json"
	"fmt"
	"maxchat/pdf_ms/internal/constants"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// TASK 1: Generate PDF
func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	file, err := h.service.GeneratePDF(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true, "message": "PDF generated successfully", "data": file,
	})
}

// TASK 2: Upload PDF
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// max upload 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid multipart form",
		})
		return
	}

	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "File not found",
		})
		return
	}

	uploaded, err := h.service.UploadPDF(fileHeader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		// tangani error file terlalu besar / extension
		var statusCode int = http.StatusBadRequest

		message := err.Error()
		if strings.Contains(message, string(constants.ERROR_FILE_TOO_LARGE)) {
			statusCode = http.StatusBadRequest
			message = fmt.Sprintf(
				"File too large, max %dMB",
				h.service.maxSize/(1024*1024),
			)

		} else if strings.Contains(message, string(constants.ERROR_INVALID_FILE_EXTENSION)) {
			statusCode = http.StatusBadRequest
			message = "Invalid file extension, only .pdf allowed"
		}

		resp := map[string]interface{}{
			"success":    false,
			"error_code": err.Error(),
			"message":    message,
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "PDF uploaded successfully",
		"data":    uploaded,
	})
}

// TASK 3: List PDFs
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	files, total, err := h.service.ListPDFs(status, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    files,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// TASK 4: Delete PDF
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	file, err := h.service.DeletePDF(id)
	if err != nil {
		http.Error(w, "file not found or already deleted", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "PDF deleted successfully",
		"data":    file,
	})
}
