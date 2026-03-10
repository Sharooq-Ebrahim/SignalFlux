package domain



import (
	"context"

	"github.com/google/uuid"
)

// TODO: Define JunctionService interface
type JunctionService interface {
	Create(ctx context.Context, j Junction) (Junction, error)
	List(ctx context.Context, page, limit int) ([]Junction, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (Junction, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetSignals(ctx context.Context, junctionID uuid.UUID) ([]Signal, error)
	UpdateSignal(ctx context.Context, junctionID uuid.UUID, dir string, secs int) error
}
