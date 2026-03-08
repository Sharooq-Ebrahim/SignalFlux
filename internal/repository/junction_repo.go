package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/srq/signalflux/internal/domain"
)

type junctionRepo struct {
	db *sql.DB
}

func NewJunctionRepo(db *sql.DB) domain.JunctionRepository {
	return &junctionRepo{db: db}
}

func (r *junctionRepo) Create(ctx context.Context, j domain.Junction) (domain.Junction, error) {

	err := r.db.QueryRowContext(ctx,
		"INSERT INTO junctions (id, type, location, created_at) VALUES ($1,$2,$3,$4) RETURNING id",
		j.ID.String(), j.Type, j.Location, j.CreatedAt,
	).Scan(&j.ID)

	if err != nil {
		return domain.Junction{}, err
	}

	return j, nil

}

func (r *junctionRepo) List(ctx context.Context, page, limit int) ([]domain.Junction, int, error) {

	var junctions []domain.Junction

	rows, err := r.db.QueryContext(ctx, "SELECT id, type, location, created_at FROM Junctions")

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var j domain.Junction
		err := rows.Scan(&j.ID, &j.Type, &j.Location, &j.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		junctions = append(junctions, j)
	}

	return junctions, 0, nil

}

func (r *junctionRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Junction, error) {

	var j domain.Junction

	err := r.db.QueryRowContext(ctx, "SELECT id,type,location,created_at FROM junctions WHERE id=$1", id).Scan(&j.ID, &j.Type, &j.Location, &j.CreatedAt)

	if err != nil {
		return domain.Junction{}, err
	}

	return j, nil
}

func (r *junctionRepo) Delete(ctx context.Context, id uuid.UUID) error {

	_, err := r.db.ExecContext(ctx, "DELETE FROM junctions WHERE id=$1", id)

	if err != nil {
		return err
	}
	return nil
}
