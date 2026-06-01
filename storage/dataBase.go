package storage

import (
	"context"
	"fmt"
	"iot/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, connString string) (*PostgresStorage, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать пул соединений: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("база данных недоступна: %w", err)
	}

	return &PostgresStorage{pool: pool}, nil
}
func (s *PostgresStorage) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}
func (s *PostgresStorage) GetShelves(ctx context.Context) ([]models.Shelf, error) {
	var shelves []models.Shelf

	sqlQuery := `
		SELECT s.shelf_id, s.current_weight, s.status,
		s.updated_at, p.product_id, p.name, p.weight_grams, p.min_weight_grams
		FROM shelves s
		LEFT JOIN products p on s.product_id = p.product_id
		ORDER BY s.shelf_id ASC;
	`
	rows, err := s.pool.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var shelf models.Shelf

		var pID *int
		var pName *string
		var pWeight *float64
		var pMinWeight *float64
		err := rows.Scan(
			&shelf.ID, &shelf.CurrentWeight, &shelf.Status, &shelf.Updated_at,
			&pID, &pName, &pWeight, &pMinWeight,
		)
		if err != nil {
			return nil, err
		}
		if pID != nil {
			shelf.Product = &models.Product{
				ID:        *pID,
				Name:      *pName,
				Weight:    *pWeight,
				MinWeight: *pMinWeight,
			}
		}
		shelves = append(shelves, shelf)
	}
	return shelves, nil
}
func (s *PostgresStorage) GetShelfByID(ctx context.Context, shelfID int) (*models.Shelf, error) {
	query := `
	SELECT s.shelf_id, s.current_weight, s.status,
		s.updated_at, p.product_id, p.name, p.weight_grams, p.min_weight_grams
		FROM shelves s
		LEFT JOIN products p on s.product_id = product_id
		WHERE s.shelf_id = $1;
	`
	row := s.pool.QueryRow(ctx, query, shelfID)
	var shelf models.Shelf
	var pID *int
	var pName *string
	var pWeight *float64
	var pMinWeight *float64
	err := row.Scan(
		&shelf.ID, &shelf.CurrentWeight, &shelf.Status, &shelf.Updated_at,
		&pID, &pName, &pWeight, &pMinWeight,
	)
	if err != nil {
		return nil, err
	}
	if pID != nil {
		shelf.Product = &models.Product{
			ID:        *pID,
			Name:      *pName,
			Weight:    *pWeight,
			MinWeight: *pMinWeight,
		}
	}
	return &shelf, nil
}
func (s *PostgresStorage) UpdateWeight(ctx context.Context, shelfID int, weight float64, status string) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	query := `
	UPDATE shelves
	SET current_weight = $1, status = $2, updated_at = CURRENT_TIMESTAMP
	WHERE shelf_id = $3;
	`
	if _, err := tx.Exec(ctx, query, weight, status, shelfID); err != nil {
		return err
	}

	logQuery := `
	INSERT INTO weight_logs(shelf_id, raw_weight, created_at)
	VALUES($1, $2, CURRENT_TIMESTAMP);
	`
	if _, err := tx.Exec(ctx, logQuery, shelfID, weight); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PostgresStorage) UpdateShelfProduct(ctx context.Context, shelfID int, productID *int) error {
	query := `
	UPDATE shelves
	SET product_id = $1, updated_at = CURRENT_TIMESTAMP
	WHERE shelf_id = $2;
	`
	if _, err := s.pool.Exec(ctx, query, productID, shelfID); err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) AddShelf(ctx context.Context, shelf *models.Shelf) error {
	query := `
		INSERT INTO shelves(product_id, current_weight, status, updated_at)
		VALUES($1, $2, $3, CURRENT_TIMESTAMP)
		RETURNING shelf_id;
	`
	var productID *int
	if shelf.Product != nil {
		productID = &shelf.Product.ID
	}
	err := s.pool.QueryRow(ctx, query, productID, shelf.CurrentWeight, shelf.Status).Scan(&shelf.ID)
	if err != nil {
		return err
	}
	return nil
}
func (s *PostgresStorage) DeleteShelf(ctx context.Context, shelfID int) error {
	query := `
	DELETE FROM shelves
	WHERE shelf_id = $1;
	`
	if _, err := s.pool.Exec(ctx, query, shelfID); err != nil {
		return err
	}
	return nil
}
