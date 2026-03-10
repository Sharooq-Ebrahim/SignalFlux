package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/srq/signalflux/internal/domain"
)

type signalRepo struct {
	db *sql.DB
}

func NewSignalRepo(db *sql.DB) domain.SignalRepository {
	return &signalRepo{db: db}
}

func (r signalRepo) GetByJunction(ctx context.Context, junctionID uuid.UUID) ([]domain.Signal, error) {

	row, err := r.db.QueryContext(ctx, "SELECT id,junction_id,direction,state,duration_seconds,updated_at FROM signals WHERE junction_id = $1 ORDER BY direction ASC", junctionID)

	if err != nil {
		return nil, err
	}

	var signals []domain.Signal
	for row.Next() {

		var signal domain.Signal

		err := row.Scan(&signal.ID, &signal.JunctionID, &signal.Direction, &signal.State, &signal.DurationSeconds, &signal.UpdatedAt)

		if err != nil {
			return nil, err
		}
		signals = append(signals, signal)
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return signals, nil

}

func (r signalRepo) UpdateDuration(ctx context.Context, junctionID uuid.UUID, dir string, secs int) error {

	row, err := r.db.ExecContext(ctx, "UPDATE signals SET duration =$1 updated_at =NOW() WHERE junction_id = $2 AND direction =$3  ", secs, junctionID, dir)

	if err != nil {
		return err
	}

	rowsAffected, err := row.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Resource not found")
	}

	return nil
}

func (r signalRepo) UpdateDurations(ctx context.Context, junctionID uuid.UUID, durations map[string]int) error {

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	for dir, secs := range durations {

		result, err := tx.ExecContext(
			ctx,
			`UPDATE signals 
			 SET duration = $1, updated_at = NOW()
			 WHERE junction_id = $2 AND direction = $3`,
			secs, junctionID, dir,

		
		)

		if err != nil {
			tx.Rollback()
			return err
		}

		rowAffected, err := result.RowsAffected()

		if err != nil {
			tx.Rollback()
			return err
		}

		if rowAffected == 0 {
			tx.Rollback()
			return errors.New("Resource not found")
		}

		tx.Commit()

	}

	return nil
}
