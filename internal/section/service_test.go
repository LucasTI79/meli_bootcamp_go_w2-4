package section_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/section"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var sectionID = 1

// Units tests

func TestRead(t *testing.T) {
	t.Run("Return all sections successfully", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		expected := getTestSections()

		repositoryMock.On("GetAll", mock.Anything).Return(expected, nil)
		result, err := svc.GetAll(context.TODO())

		assert.NoError(t, err)
		assert.ElementsMatch(t, expected, result)
		assert.Len(t, expected, 2)
	})
	t.Run("Does not get any section and returns error: getting sections", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		repositoryMock.On("GetAll", mock.Anything).Return([]domain.Section{}, section.ErrGetSections)
		_, err := svc.GetAll(context.TODO())

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrGetSections)
	})
	t.Run("Return a section by ID successfully", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		expected := getTestSections()[0]

		repositoryMock.On("Get", mock.Anything, mock.Anything).Return(expected, nil)
		result, err := svc.Get(context.TODO(), sectionID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("Does not get any section and returns error: not found", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		repositoryMock.On("Get", mock.Anything, mock.Anything).Return(domain.Section{}, section.ErrNotFound)
		_, err := svc.Get(context.TODO(), sectionID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrNotFound)
	})
}

func TestCreate(t *testing.T) {
	t.Run("Create a section successfully", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		body := getTestCreateSections()

		expected := getTestSections()[0]

		repositoryMock.On("Exists", mock.Anything, mock.Anything).Return(false)
		repositoryMock.On("Save", mock.Anything, mock.Anything).Return(sectionID, nil)

		result, err := svc.Create(context.TODO(), body)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.EqualValues(t, sectionID, result.ID)
	})
	t.Run("Does not create any section and returns error: section number alredy exists", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		body := getTestCreateSections()

		repositoryMock.On("Exists", mock.Anything, mock.Anything).Return(true)
		_, err := svc.Create(context.TODO(), body)

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrInvalidSectionNumber)
		repositoryMock.AssertNumberOfCalls(t, "Save", 0)
	})
	t.Run("Does not create any section and returns error: saving section", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		body := getTestCreateSections()

		repositoryMock.On("Exists", mock.Anything, mock.Anything).Return(false)
		repositoryMock.On("Save", mock.Anything, mock.Anything).Return(0, section.ErrSavingSection)

		_, err := svc.Create(context.TODO(), body)

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrSavingSection)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Update a section successfully", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		actualSection := domain.Section{
			ID:                 1,
			SectionNumber:      123,
			CurrentTemperature: 10,
			MinimumTemperature: 5,
			CurrentCapacity:    15,
			MinimumCapacity:    10,
			MaximumCapacity:    20,
			WarehouseID:        321,
			ProductTypeID:      2,
		}

		updates := section.UpdateSection{
			SectionNumber:      testutil.ToPtr(1234),
			CurrentTemperature: testutil.ToPtr(11),
			MinimumTemperature: testutil.ToPtr(0),
			CurrentCapacity:    testutil.ToPtr(16),
			MinimumCapacity:    testutil.ToPtr(0),
			MaximumCapacity:    testutil.ToPtr(21),
			WarehouseID:        testutil.ToPtr(3210),
			ProductTypeID:      testutil.ToPtr(3),
		}

		expected := domain.Section{
			ID:                 1,
			SectionNumber:      1234,
			CurrentTemperature: 11,
			MinimumTemperature: 0,
			CurrentCapacity:    16,
			MinimumCapacity:    0,
			MaximumCapacity:    21,
			WarehouseID:        3210,
			ProductTypeID:      3,
		}

		repositoryMock.On("Get", mock.Anything, mock.Anything).Return(actualSection, nil)
		repositoryMock.On("Exists", mock.Anything, mock.Anything).Return(false)
		repositoryMock.On("Update", mock.Anything, mock.Anything).Return(nil)

		result, err := svc.Update(context.TODO(), updates, sectionID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.NotEqual(t, actualSection, result)
	})
	t.Run("Does not update any section and returns error: not found", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		body := getUpdateSection()

		repositoryMock.On("Get", mock.Anything, mock.Anything).Return(domain.Section{}, section.ErrNotFound)

		_, err := svc.Update(context.TODO(), body, sectionID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrNotFound)
		repositoryMock.AssertNumberOfCalls(t, "Exists", 0)
		repositoryMock.AssertNumberOfCalls(t, "Update", 0)
	})
	t.Run("Does not update any section and returns error: section number alredy exists", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		body := getUpdateSection()

		repositoryMock.On("Get", mock.Anything, mock.Anything).Return(domain.Section{}, nil)
		repositoryMock.On("Exists", mock.Anything, mock.Anything).Return(true)

		_, err := svc.Update(context.TODO(), body, sectionID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrInvalidSectionNumber)
		repositoryMock.AssertNumberOfCalls(t, "Update", 0)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete a section successfully", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		repositoryMock.On("Get", mock.Anything, mock.Anything).Return(domain.Section{}, nil)
		repositoryMock.On("Delete", mock.Anything, mock.Anything).Return(nil)
		err := svc.Delete(context.TODO(), sectionID)

		assert.NoError(t, err)
	})
	t.Run("Does not delete any section and returns error: not found", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		repositoryMock.On("Get", mock.Anything, mock.Anything).Return(domain.Section{}, section.ErrNotFound)
		err := svc.Delete(context.TODO(), sectionID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrNotFound)
		repositoryMock.AssertNumberOfCalls(t, "Delete", 0)
	})
}

// Generate test objects

func getTestSections() []domain.Section {
	return []domain.Section{
		{
			ID:                 1,
			SectionNumber:      123,
			CurrentTemperature: 10,
			MinimumTemperature: 5,
			CurrentCapacity:    15,
			MinimumCapacity:    10,
			MaximumCapacity:    20,
			WarehouseID:        321,
			ProductTypeID:      2,
		},
		{
			ID:                 2,
			SectionNumber:      1234,
			CurrentTemperature: 11,
			MinimumTemperature: 6,
			CurrentCapacity:    16,
			MinimumCapacity:    11,
			MaximumCapacity:    21,
			WarehouseID:        3210,
			ProductTypeID:      3,
		},
	}
}

func getTestCreateSections() section.CreateSection {
	return section.CreateSection{
		SectionNumber:      123,
		CurrentTemperature: 10,
		MinimumTemperature: 5,
		CurrentCapacity:    15,
		MinimumCapacity:    10,
		MaximumCapacity:    20,
		WarehouseID:        321,
		ProductTypeID:      2,
	}
}

func getUpdateSection() section.UpdateSection {
	return section.UpdateSection{
		SectionNumber:      testutil.ToPtr(123),
		CurrentTemperature: testutil.ToPtr(11),
		MinimumTemperature: testutil.ToPtr(6),
		CurrentCapacity:    testutil.ToPtr(16),
		MinimumCapacity:    testutil.ToPtr(11),
		MaximumCapacity:    testutil.ToPtr(21),
		WarehouseID:        testutil.ToPtr(3210),
		ProductTypeID:      testutil.ToPtr(3),
	}
}

// Mock repository functions

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]domain.Section, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Section), args.Error(1)
}

func (r *RepositoryMock) Get(ctx context.Context, id int) (domain.Section, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Section), args.Error(1)
}

func (r *RepositoryMock) Exists(ctx context.Context, sectionNumber int) bool {
	args := r.Called(ctx, sectionNumber)
	return args.Get(0).(bool)
}

func (r *RepositoryMock) Save(ctx context.Context, s domain.Section) (int, error) {
	args := r.Called(ctx, s)
	return args.Get(0).(int), args.Error(1)
}

func (r *RepositoryMock) Update(ctx context.Context, s domain.Section) error {
	args := r.Called(ctx, s)
	return args.Error(0)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}
