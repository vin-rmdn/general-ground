package postgresql

import (
	"context"
	"fmt"
)

func (p *PostgreSQL) CreateUser(ctx context.Context, user string) error {
	const queryTemplate = `INSERT INTO users (name) VALUES ($1)`

	_, err := p.conn.Exec(ctx, queryTemplate, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
