package domain

import "context"

type Repository interface {
	Save(ctx context.Context, a *Alert) error
	Update(ctx context.Context, a *Alert) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*Alert, error)
	List(ctx context.Context) ([]Alert, error)
}
