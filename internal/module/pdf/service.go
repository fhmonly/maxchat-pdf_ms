package pdf

import (
	"fmt"
	"io"
	"maxchat/pdf_ms/internal/constants"
	"maxchat/pdf_ms/internal/utils"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/httpimg"
)

type Service struct {
	repo      *Repository
	uploadDir string
	maxSize   int64
}

func NewService(repo *Repository, uploadDir string) *Service {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}
	return &Service{repo: repo, uploadDir: uploadDir, maxSize: 10 << 20}
}

// TASK 1: Generate PDF
func (s *Service) GeneratePDF(req *GenerateRequest) (*PDFFile, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	var opt gofpdf.ImageOptions

	// buat pdf Header
	pdf.SetHeaderFunc(func() {
		// default image = funi
		logoImg := "https://api.duniagames.co.id/api/content/upload/file/5716432721680166042.jpg"
		// logo/image di kiri atas
		if req.LogoURL != "" && utils.HasAllowedExt(req.LogoURL, []string{".jpg", ".jpeg", ".png"}) {
			logoImg = req.LogoURL
		}

		httpimg.Register(pdf, logoImg, "")

		pdf.ImageOptions(logoImg, 10, 10, 20, 0, false, opt, 0, "")

		// nama lembaga di tengah
		pdf.SetFont("Arial", "B", 16)
		pdf.CellFormat(0, 10, req.InstitutionName, "", 1, "C", false, 0, "")

		// alamat dan kontak
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 5, req.Address, "", 1, "C", false, 0, "")
		pdf.CellFormat(0, 5, req.Phone, "", 1, "C", false, 0, "")
		pdf.Ln(5)
	})

	// Footer
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)

		// no halaman
		pageStr := fmt.Sprintf("Page %d of {nb}", pdf.PageNo())

		// tanggal generate
		timestamp := time.Now().Format("2006-01-02 15:04:05")

		footerText := fmt.Sprintf("%s | Generated: %s", pageStr, timestamp)
		pdf.CellFormat(0, 10, footerText, "", 0, "C", false, 0, "")
	})

	// AliasNbPages untuk mengganti {nb} dengan total halaman
	pdf.AliasNbPages("")

	// Tambah halaman pertama
	pdf.AddPage()

	// content

	// title
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, req.Title, "", 1, "L", false, 0, "")
	pdf.Ln(3)

	// tanggal
	pdf.SetFont("Arial", "", 10)
	generateDate := time.Now().Format("02 January 2006")
	pdf.CellFormat(0, 6, fmt.Sprintf("Tanggal: %s", generateDate), "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// isi
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 6, req.Content, "", "L", false)

	// simpan
	// nama file unik: report_YYYYMMDD_uuid.pdf
	filename := fmt.Sprintf("report_%s_%s.pdf",
		time.Now().Format("20060102"),
		uuid.New().String())

	filepath := filepath.Join(s.uploadDir, filename)

	// output pdf ke file
	if err := pdf.OutputFileAndClose(filepath); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	// get file info buat disimpan ke db
	info, err := os.Stat(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// simpan ke db

	file := &PDFFile{
		Filename:     filename,
		OriginalName: req.Title + ".pdf",
		Filepath:     filepath,
		Size:         info.Size(),
		Status:       "CREATED",
		CreatedAt:    time.Now(),
	}

	if err := s.repo.Create(file); err != nil {
		// hapus file jika gagal disimpan ke db
		os.Remove(filepath)
		return nil, fmt.Errorf("failed to save to database: %w", err)
	}

	return file, nil
}

// TASK 2: Upload PDF
func (s *Service) UploadPDF(fileHeader *multipart.FileHeader) (*PDFFile, error) {
	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".pdf") {
		return nil, fmt.Errorf(string(constants.ERROR_INVALID_FILE_EXTENSION))
	}
	if fileHeader.Size > s.maxSize {
		return nil, fmt.Errorf(string(constants.ERROR_FILE_TOO_LARGE))
	}
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	filename := fmt.Sprintf("upload_%s_%s.pdf", time.Now().Format("20060102"), uuid.New().String())
	destPath := filepath.Join(s.uploadDir, filename)
	out, err := os.Create(destPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return nil, err
	}

	stat, _ := out.Stat()
	pdf := &PDFFile{
		Filename:     filename,
		OriginalName: fileHeader.Filename,
		Filepath:     destPath,
		Size:         stat.Size(),
		Status:       "UPLOADED",
	}
	if err := s.repo.Create(pdf); err != nil {
		return nil, err
	}
	return pdf, nil
}

// TASK 3: List PDFs
func (s *Service) ListPDFs(status string, page, limit int) ([]PDFFile, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.List(status, offset, limit)
}

// TASK 4: Soft Delete PDF
func (s *Service) DeletePDF(id int64) (*PDFFile, error) {
	return s.repo.SoftDelete(id)
}
