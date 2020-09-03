package postgres

import (
	"fmt"

	"github.com/christiaan-janssen/goreddit"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CommentStore struct {
	*sqlx.DB
}

func (s *CommentStore) Comment(id uuid.UUID) (goreddit.Comment, error) {
	var c goreddit.Comment
	if err := s.Get(&c, `SELECT * FROM comments WHERE id = $1`, id); err != nil {
		return goreddit.Comment{}, fmt.Errorf("error getting Comment: %w", err)
	}
	return c, nil
}

func (s *CommentStore) CommentsByPost(postId uuid.UUID) ([]goreddit.Comment, error) {
	var cc []goreddit.Comment
	if err := s.Select(&cc, `SELECT * from comments where post_id = $1`, postId); err != nil {
		return []goreddit.Comment{}, fmt.Errorf("error getting Comments: %w", err)
	}
	return cc, nil
}

func (s *CommentStore) CreateComment(c *goreddit.Comment) error {
	if err := s.Get(c, `INSERT INTO Comments VALUES ($1, $2, $3, $4) RETURNING *`,
		c.ID,
		c.PostID,
		c.Content,
		c.Votes); err != nil {
		return fmt.Errorf("error creating Comment: %w", err)
	}
	return nil
}

func (s *CommentStore) UpdateComment(c *goreddit.Comment) error {
	if err := s.Get(c, `UPDATE Comments SET post_id = $1, title = $2, content = $3, votes = $4 WHERE id = $4 RETURNING *`,
		c.PostID,
		c.Content,
		c.Votes,
		c.ID); err != nil {
		return fmt.Errorf("error updating Comment: %w", err)
	}
	return nil
}

func (s *CommentStore) DeleteComment(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM Comments WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting Comment: %w", err)
	}
	return nil
}
