package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/srq/signalflux/internal/domain"
)



type junctionService struct {
	junctionRepo domain.JunctionRepository
	signalRepo   domain.SignalRepository
}

func NewJunctionService(jr domain.JunctionRepository, sr domain.SignalRepository) domain.JunctionService {
	return junctionService{junctionRepo: jr, signalRepo: sr}
}

func (js junctionService) Create(ctx context.Context, j domain.Junction) (domain.Junction, error) {

	isValidType := false

	for _, jtype := range domain.ValidJunctionTypes {

		if j.Type == jtype {
			isValidType = true

		}
	}

	if !isValidType {
		return domain.Junction{}, errors.New("Invalid Input")
	}

	if j.Location == "" {
		return domain.Junction{}, errors.New("Invalid Input")
	}

	j.ID = uuid.New()

	return js.junctionRepo.Create(ctx, j)

}

func (js junctionService) List(ctx context.Context, page, limit int) ([]domain.Junction, int, error) {

	return js.junctionRepo.List(ctx, page, limit)

}

func (js junctionService) GetByID(ctx context.Context, id uuid.UUID) (domain.Junction, error) {
	j, err := js.junctionRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Junction{}, err
	}
	return j, err
}

func (js junctionService) Delete(ctx context.Context, id uuid.UUID) error {

	err := js.junctionRepo.Delete(ctx, id)

	if err != nil {
		return errors.New("Failed to delete junction")
	}
	return nil

}

func (js junctionService) GetSignals(ctx context.Context, id uuid.UUID) ([]domain.Signal, error) {

	j, err := js.junctionRepo.GetByID(ctx, id)

	log.Println(j)
	if err != nil {
		return nil, err
	}
	return js.signalRepo.GetByJunction(ctx, id)

}

func (js junctionService) UpdateSignal(ctx context.Context, id uuid.UUID, dir string, secs int) error {

	validDir := false
	for _, d := range domain.ValidDirections {
		if dir == d {
			validDir = true
			break
		}
	}
	if !validDir {
		return errors.New("Invalid Input")
	}

	if secs <= 0 || secs > 300 {
		return errors.New("Invalid Input")
	}

	return js.signalRepo.UpdateDuration(ctx, id, dir, secs)

}
