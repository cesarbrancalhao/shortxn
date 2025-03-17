package postgres

import (
	"Shortxn/internal/domain"
	"database/sql"

	_ "github.com/lib/pq"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(connStr string) (*URLRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &URLRepository{db: db}, nil
}

func (r *URLRepository) Store(url *domain.URL) error {
	query := `
		INSERT INTO urls (id, long_url, short_url, created_at, clicks)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, url.ID, url.LongURL, url.ShortURL, url.CreatedAt, url.Clicks)
	return err
}

func (r *URLRepository) GetByID(id string) (*domain.URL, error) {
	var url domain.URL
	query := `SELECT id, long_url, short_url, created_at, clicks FROM urls WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&url.ID, &url.LongURL, &url.ShortURL, &url.CreatedAt, &url.Clicks)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *URLRepository) GetByLongURL(longURL string) (*domain.URL, error) {
	var url domain.URL
	query := `SELECT id, long_url, short_url, created_at, clicks FROM urls WHERE long_url = $1`

	err := r.db.QueryRow(query, longURL).Scan(&url.ID, &url.LongURL, &url.ShortURL, &url.CreatedAt, &url.Clicks)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *URLRepository) IncrementClicks(id string) error {
	query := `UPDATE urls SET clicks = clicks + 1 WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
