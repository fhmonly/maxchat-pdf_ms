package pdf

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create new PDF record
func (r *Repository) Create(file *PDFFile) error {
	query := `INSERT INTO pdf_files (filename, original_name, filepath, size, status)
	          VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.Exec(query, file.Filename, file.OriginalName, file.Filepath, file.Size, file.Status)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	file.ID = id
	file.CreatedAt = time.Now()
	return nil
}

// List files with optional status, pagination
func (r *Repository) List(status string, offset, limit int) ([]PDFFile, int, error) {
	baseQuery := "SELECT id, filename, original_name, filepath, size, status, created_at FROM pdf_files"
	countQuery := "SELECT COUNT(*) FROM pdf_files"
	args := []interface{}{}
	where := ""
	if status != "" {
		where = " WHERE status = ?"
		args = append(args, status)
	}
	var total int
	if err := r.db.QueryRow(countQuery+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	query := baseQuery + where + " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var files []PDFFile
	for rows.Next() {
		var f PDFFile
		if err := rows.Scan(&f.ID, &f.Filename, &f.OriginalName, &f.Filepath, &f.Size, &f.Status, &f.CreatedAt); err != nil {
			return nil, 0, err
		}
		files = append(files, f)
	}
	return files, total, nil
}

// Soft delete file by ID
func (r *Repository) SoftDelete(id int64) (*PDFFile, error) {
	var f PDFFile
	err := r.db.QueryRow("SELECT id, filename, status FROM pdf_files WHERE id = ?", id).
		Scan(&f.ID, &f.Filename, &f.Status)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if f.Status == "DELETED" {
		return nil, sql.ErrNoRows
	}

	now := time.Now()
	_, err = r.db.Exec("UPDATE pdf_files SET status='DELETED', deleted_at=? WHERE id=?", now, id)
	if err != nil {
		return nil, err
	}
	f.Status = "DELETED"
	f.DeletedAt = &now
	return &f, nil
}
