package center

import (
	"context"

	"github.com/HugoKovac/rdvisit/db/sqlc"
	"github.com/HugoKovac/rdvisit/internal/domain"
	"github.com/HugoKovac/rdvisit/pkg/errors"
)

type IRepository interface {
	GetCenters(ctx context.Context) ([]*domain.Center, error)
}

// ==================================
//	Params
// ==================================

type CreateParams struct {
	Email        string
	FirstName    string
	LastName     string
	PasswordHash string
}

//==================================
//			Implementation
//==================================

type sqlcRepository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) IRepository {
	return &sqlcRepository{q: q}
}

func (r *sqlcRepository) GetCenters(ctx context.Context) ([]*domain.Center, error) {
	centers, err := r.q.GetAllCenters(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	rtn := make([]*domain.Center, 0, len(centers))
	for _, center := range centers {
		rtn = append(rtn, &domain.Center{
			ID:           center.ID,
			Name:         center.Name,
			Street:       center.Street,
			StreetNumber: center.StreetNumber,
			City:         center.City,
			PostalCode:   center.PostalCode,
			Region:       center.Region,
			Country:      center.Country,
			CreatedAt:    center.CreatedAt,
			UpdatedAt:    center.UpdatedAt,
		})
	}

	return rtn, nil
}
