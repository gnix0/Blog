package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID             int64    `json:"id"`
	Content        string   `json:"content"`
	Title          string   `json:"title"`
	UserID         int64    `json:"user_id"`
	Tags           []string `json:"tags"`
	CoverImageURL  string   `json:"cover_image_url,omitempty"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags, cover_image_url)
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at
	`

	return s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
		nullableString(post.CoverImageURL),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
	SELECT id, title, content, user_id, tags, COALESCE(cover_image_url, ''), created_at, updated_at
	FROM posts WHERE id = $1
	`

	post := &Post{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CoverImageURL,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostStore) List(ctx context.Context) ([]Post, error) {
	query := `
	SELECT id, title, content, user_id, tags, COALESCE(cover_image_url, ''), created_at, updated_at
	FROM posts ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.UserID,
			pq.Array(&p.Tags),
			&p.CoverImageURL,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts
	SET title = $1, content = $2, tags = $3, cover_image_url = $4, updated_at = $5
	WHERE id = $6
	RETURNING updated_at
	`

	return s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		nullableString(post.CoverImageURL),
		time.Now().UTC(),
		post.ID,
	).Scan(&post.UpdatedAt)
}

func (s *PostStore) Delete(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// nullableString returns nil for an empty string so optional TEXT columns store NULL.
func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
