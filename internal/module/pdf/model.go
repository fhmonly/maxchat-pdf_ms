package pdf

import "time"

type PDFFile struct {
	ID           int64      `json:"id"`
	Filename     string     `json:"filename"`
	OriginalName string     `json:"original_name,omitempty"`
	Filepath     string     `json:"filepath"`
	Size         int64      `json:"size,omitempty"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type GenerateRequest struct {
	Title           string `json:"title"`
	InstitutionName string `json:"institution_name"`
	Address         string `json:"address"`
	Phone           string `json:"phone"`
	LogoURL         string `json:"logo_url,omitempty"`
	Content         string `json:"content"`
}
