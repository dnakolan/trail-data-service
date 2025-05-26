package services

import (
	"context"
	"errors"
	"testing"

	"github.com/dnakolan/trail-data-service/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTrailStorage is a mock implementation of storage.TrailStorage
type MockTrailStorage struct {
	mock.Mock
}

func (m *MockTrailStorage) Save(ctx context.Context, trail *models.Trail) error {
	args := m.Called(ctx, trail)
	return args.Error(0)
}

func (m *MockTrailStorage) FindAll(ctx context.Context, filter *models.TrailFilter) ([]*models.Trail, error) {
	args := m.Called(ctx, filter)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Trail), nil
}

func (m *MockTrailStorage) FindById(ctx context.Context, uid string) (*models.Trail, error) {
	args := m.Called(ctx, uid)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Trail), nil
}

func (m *MockTrailStorage) Delete(ctx context.Context, uid string) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *MockTrailStorage) Clear(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestCreateTrail(t *testing.T) {
	mockStorage := new(MockTrailStorage)
	service := NewTrailsService(mockStorage)
	ctx := context.Background()

	trail := models.NewTrail("Test Trail", 45.5231, -122.6765, models.TrailDifficultyMedium, 10.5)

	tests := []struct {
		name        string
		trail       *models.Trail
		setupMock   func()
		expectError bool
		errorMsg    string
	}{
		{
			name:  "successful creation - no duplicates",
			trail: trail,
			setupMock: func() {
				// Mock the duplicate check
				mockStorage.On("FindAll", ctx, mock.MatchedBy(func(f *models.TrailFilter) bool {
					return f.Name != nil && *f.Name == "Test Trail" &&
						f.Lat != nil && *f.Lat == 45.5231 &&
						f.Lon != nil && *f.Lon == -122.6765 &&
						f.RadiusKm != nil
				})).Return([]*models.Trail{}, nil).Once()

				// Mock the save
				mockStorage.On("Save", ctx, trail).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name:  "duplicate trail",
			trail: trail,
			setupMock: func() {
				// Mock finding a duplicate
				mockStorage.On("FindAll", ctx, mock.Anything).Return([]*models.Trail{trail}, nil).Once()
			},
			expectError: true,
			errorMsg:    "trail already exists",
		},
		{
			name:  "storage error during duplicate check",
			trail: trail,
			setupMock: func() {
				mockStorage.On("FindAll", ctx, mock.Anything).Return(nil, errors.New("database error")).Once()
			},
			expectError: true,
			errorMsg:    "database error",
		},
		{
			name:  "storage error during save",
			trail: trail,
			setupMock: func() {
				mockStorage.On("FindAll", ctx, mock.Anything).Return([]*models.Trail{}, nil).Once()
				mockStorage.On("Save", ctx, trail).Return(errors.New("database error")).Once()
			},
			expectError: true,
			errorMsg:    "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := service.CreateTrail(ctx, tt.trail)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Equal(t, tt.errorMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestGetTrail(t *testing.T) {
	mockStorage := new(MockTrailStorage)
	service := NewTrailsService(mockStorage)
	ctx := context.Background()

	uid := uuid.New()
	trail := models.NewTrail("Test Trail", 45.5231, -122.6765, models.TrailDifficultyMedium, 10.5)
	trail.UID = uid

	tests := []struct {
		name        string
		uid         string
		setupMock   func()
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful retrieval",
			uid:  uid.String(),
			setupMock: func() {
				mockStorage.On("FindById", ctx, uid.String()).Return(trail, nil).Once()
			},
			expectError: false,
		},
		{
			name: "trail not found",
			uid:  "non-existent",
			setupMock: func() {
				mockStorage.On("FindById", ctx, "non-existent").Return(nil, errors.New("trail not found")).Once()
			},
			expectError: true,
			errorMsg:    "trail not found",
		},
		{
			name: "storage error",
			uid:  uid.String(),
			setupMock: func() {
				mockStorage.On("FindById", ctx, uid.String()).Return(nil, errors.New("database error")).Once()
			},
			expectError: true,
			errorMsg:    "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			trail, err := service.GetTrail(ctx, tt.uid)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, trail)
				if tt.errorMsg != "" {
					assert.Equal(t, tt.errorMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, trail)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestUpdateTrail(t *testing.T) {
	mockStorage := new(MockTrailStorage)
	service := NewTrailsService(mockStorage)
	ctx := context.Background()

	trail := models.NewTrail("Test Trail", 45.5231, -122.6765, models.TrailDifficultyMedium, 10.5)

	tests := []struct {
		name        string
		trail       *models.Trail
		setupMock   func()
		expectError bool
		errorMsg    string
	}{
		{
			name:  "successful update",
			trail: trail,
			setupMock: func() {
				mockStorage.On("Save", ctx, trail).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name:  "storage error",
			trail: trail,
			setupMock: func() {
				mockStorage.On("Save", ctx, trail).Return(errors.New("database error")).Once()
			},
			expectError: true,
			errorMsg:    "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := service.UpdateTrail(ctx, tt.trail)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Equal(t, tt.errorMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestDeleteTrail(t *testing.T) {
	mockStorage := new(MockTrailStorage)
	service := NewTrailsService(mockStorage)
	ctx := context.Background()

	uid := uuid.New()

	tests := []struct {
		name        string
		uid         string
		setupMock   func()
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful deletion",
			uid:  uid.String(),
			setupMock: func() {
				mockStorage.On("Delete", ctx, uid.String()).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name: "trail not found",
			uid:  "non-existent",
			setupMock: func() {
				mockStorage.On("Delete", ctx, "non-existent").Return(errors.New("trail not found")).Once()
			},
			expectError: true,
			errorMsg:    "trail not found",
		},
		{
			name: "storage error",
			uid:  uid.String(),
			setupMock: func() {
				mockStorage.On("Delete", ctx, uid.String()).Return(errors.New("database error")).Once()
			},
			expectError: true,
			errorMsg:    "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := service.DeleteTrail(ctx, tt.uid)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Equal(t, tt.errorMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestGetAllTrails(t *testing.T) {
	mockStorage := new(MockTrailStorage)
	service := NewTrailsService(mockStorage)
	ctx := context.Background()

	trail1 := models.NewTrail("Trail 1", 45.5231, -122.6765, models.TrailDifficultyMedium, 10.5)
	trail2 := models.NewTrail("Trail 2", 45.5232, -122.6766, models.TrailDifficultyHard, 15.5)
	trails := []*models.Trail{trail1, trail2}

	difficulty := models.TrailDifficultyMedium
	filter := &models.TrailFilter{
		CreateTrailRequest: models.CreateTrailRequest{
			Difficulty: &difficulty,
		},
	}

	tests := []struct {
		name        string
		filter      *models.TrailFilter
		setupMock   func()
		expectError bool
		errorMsg    string
	}{
		{
			name:   "successful retrieval - with filter",
			filter: filter,
			setupMock: func() {
				mockStorage.On("FindAll", ctx, filter).Return([]*models.Trail{trail1}, nil).Once()
			},
			expectError: false,
		},
		{
			name:   "successful retrieval - no filter",
			filter: nil,
			setupMock: func() {
				mockStorage.On("FindAll", ctx, (*models.TrailFilter)(nil)).Return(trails, nil).Once()
			},
			expectError: false,
		},
		{
			name:   "storage error",
			filter: filter,
			setupMock: func() {
				mockStorage.On("FindAll", ctx, filter).Return(nil, errors.New("database error")).Once()
			},
			expectError: true,
			errorMsg:    "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			trails, err := service.GetAllTrails(ctx, tt.filter)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, trails)
				if tt.errorMsg != "" {
					assert.Equal(t, tt.errorMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, trails)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}
