package album

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/tvitcom/local-adverts/internal/entity"
	"github.com/tvitcom/local-adverts/pkg/log"
	"time"
)

// Agregator encapsulates usecase logic for albums
type Agregator interface {
	Get(ctx context.Context, id string) (Album, error)
	Query(ctx context.Context, offset, limit int) ([]Album, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateAlbumRequest) (Album, error)
	Update(ctx context.Context, id string, input UpdateAlbumRequest) (Album, error)
	Delete(ctx context.Context, id string) (Album, error)
}

// Album represents the data about an album.
type Album struct {
	entity.Album
}

// CreateAlbumRequest represents an album creation request.
type CreateAlbumRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreateAlbumRequest fields
func (m CreateAlbumRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

// UpdateAlbumRequest represents an album update request.
type UpdateAlbumRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreateAlbumRequest fields
func (m UpdateAlbumRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

type agregator struct {
	repo   Repository
	logger log.Logger
}

// NewAgregator creates a new album agregator.
func NewAgregator(repo Repository, logger log.Logger) Agregator {
	return agregator{repo, logger}
}

// Get returns the album with the specified the album ID.
func (ag agregator) Get(ctx context.Context, id string) (Album, error) {
	album, err := ag.repo.Get(ctx, id)
	if err != nil {
		return Album{}, err
	}
	return Album{album}, nil
}

// Create creates a new album.
func (ag agregator) Create(ctx context.Context, req CreateAlbumRequest) (Album, error) {
	
	if err := req.Validate(); err != nil {
		return Album{}, err
	}
	id := entity.GenerateUUID()
	now := time.Now()
	err := ag.repo.Create(ctx, entity.Album{
		ID:        id,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return Album{}, err
	}
	return ag.Get(ctx, id)
}

// Update updates the album with the specified ID.
func (ag agregator) Update(ctx context.Context, id string, req UpdateAlbumRequest) (Album, error) {
	if err := req.Validate(); err != nil {
		return Album{}, err
	}

	album, err := ag.Get(ctx, id)
	if err != nil {
		return album, err
	}
	album.Name = req.Name
	album.UpdatedAt = time.Now()

	if err := ag.repo.Update(ctx, album.Album); err != nil {
		return album, err
	}
	return album, nil
}

// Delete deletes the album with the specified ID.
func (ag agregator) Delete(ctx context.Context, id string) (Album, error) {
	album, err := ag.Get(ctx, id)
	if err != nil {
		return Album{}, err
	}
	if err = ag.repo.Delete(ctx, id); err != nil {
		return Album{}, err
	}
	return album, nil
}

// Count returns the number of albums
func (ag agregator) Count(ctx context.Context) (int, error) {
	return ag.repo.Count(ctx)
}

// Query returns the albums with the specified offset and limit.
func (ag agregator) Query(ctx context.Context, offset, limit int) ([]Album, error) {
	items, err := ag.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := []Album{}
	for _, item := range items {
		result = append(result, Album{item})
	}
	return result, nil
}
