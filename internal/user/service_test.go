package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/tvitcom/local-adverts/internal/entity"
	"github.com/tvitcom/local-adverts/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

var errCRUD = errors.New("error crud")

func TestCreateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     CreateUserRequest
		wantError bool
	}{
		{"success", CreateUserRequest{Name: "test"}, false},
		{"required", CreateUserRequest{Name: ""}, true},
		{"too long", CreateUserRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     UpdateUserRequest
		wantError bool
	}{
		{"success", UpdateUserRequest{Name: "test"}, false},
		{"required", UpdateUserRequest{Name: ""}, true},
		{"too long", UpdateUserRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	logger, _ := log.NewForTest()
	s := NewService(&mockRepository{}, logger)

	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, 0, count)

	// successful creation
	user, err := s.Create(ctx, CreateUserRequest{Name: "test"})
	assert.Nil(t, err)
	assert.NotEmpty(t, user.ID)
	id := user.ID
	assert.Equal(t, "test", user.Name)
	assert.NotEmpty(t, user.CreatedAt)
	assert.NotEmpty(t, user.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// validation error in creation
	_, err = s.Create(ctx, CreateUserRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// unexpected error in creation
	_, err = s.Create(ctx, CreateUserRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	_, _ = s.Create(ctx, CreateUserRequest{Name: "test2"})

	// update
	user, err = s.Update(ctx, id, UpdateUserRequest{Name: "test updated"})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", user.Name)
	_, err = s.Update(ctx, "none", UpdateUserRequest{Name: "test updated"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, id, UpdateUserRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// unexpected error in update
	_, err = s.Update(ctx, id, UpdateUserRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	user, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", user.Name)
	assert.Equal(t, id, user.ID)

	// query
	users, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, len(users))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	user, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, user.ID)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)
}

type mockRepository struct {
	items []entity.User
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.User, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return entity.User{}, sql.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int) ([]entity.User, error) {
	return m.items, nil
}

func (m *mockRepository) Create(ctx context.Context, user entity.User) error {
	if user.Name == "error" {
		return errCRUD
	}
	m.items = append(m.items, user)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, user entity.User) error {
	if user.Name == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.ID == user.ID {
			m.items[i] = user
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
