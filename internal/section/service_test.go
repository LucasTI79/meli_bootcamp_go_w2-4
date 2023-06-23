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
var sectionNumber = 123
var sectionTemp = 11

func TestCreate(t *testing.T) {
	t.Run("Create valid section", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		sectionPayload := getTestCreateSections()

		expected := *section.MapCreateToDomain(&sectionPayload)
		expected.ID = sectionID

		repositoryMock.On("Exists", mock.Anything, sectionNumber).Return(false)
		repositoryMock.On("Save", mock.Anything, mock.Anything).Return(sectionID, nil)

		result, err := svc.Save(context.TODO(), sectionPayload)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.EqualValues(t, sectionID, result.ID)
	})
	t.Run("Doesn't create section and return error: section number alredy exists", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		sectionPayload := getTestCreateSections()

		repositoryMock.On("Exists", mock.Anything, sectionNumber).Return(true)
		_, err := svc.Save(context.TODO(), sectionPayload)

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrInvalidSectionNumber)
		repositoryMock.AssertNumberOfCalls(t, "Save", 0)
	})
	t.Run("Doesn't create section and return error: saving section", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		sectionPayload := getTestCreateSections()

		repositoryMock.On("Exists", mock.Anything, sectionNumber).Return(false)
		repositoryMock.On("Save", mock.Anything, mock.Anything).Return(0, section.ErrSavingSection)

		_, err := svc.Save(context.TODO(), sectionPayload)

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrSavingSection)
	})
}

func TestRead(t *testing.T) {
	t.Run("Get all products", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		expected := getTestSections()

		repositoryMock.On("GetAll", mock.Anything).Return(expected, nil)
		result, err := svc.GetAll(context.TODO())

		assert.NoError(t, err)
		assert.ElementsMatch(t, expected, result)
	})
	t.Run("Get no products and return error: gettiung sections", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		repositoryMock.On("GetAll", mock.Anything).Return([]domain.Section{}, section.ErrGetSections)
		_, err := svc.GetAll(context.TODO())

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrGetSections)
	})
	t.Run("Get a section by ID", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		expected := getTestSections()

		repositoryMock.On("Get", mock.Anything, sectionID).Return(expected, nil)
		result, err := svc.Get(context.TODO(), sectionID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("Doesn't get a section and return error: not found", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		repositoryMock.On("Get", mock.Anything, sectionID).Return(domain.Section{}, section.ErrNotFound)
		_, err := svc.Get(context.TODO(), sectionID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrNotFound)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete a section by ID", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		repositoryMock.On("Get", mock.Anything, sectionID).Return(domain.Section{}, section.ErrNotFound)
		repositoryMock.On("Delete", mock.Anything, sectionID).Return(nil)
		err := svc.Delete(context.TODO(), sectionID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrNotFound)
	})
	t.Run("Doesn't delete a section and return error: not found", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		repositoryMock.On("Get", mock.Anything, sectionID).Return(domain.Section{}, section.ErrNotFound)
		err := svc.Delete(context.TODO(), sectionID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrNotFound)
		repositoryMock.AssertNumberOfCalls(t, "Delete", 0)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Update a section", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		actualSection := getTestSections()[0]

		updatePayload := section.UpdateSection{
			SectionNumber:      testutil.ToPtr(sectionNumber),
			CurrentTemperature: testutil.ToPtr(sectionTemp),
		}

		expected := getTestSections()[0]
		expected.CurrentTemperature = sectionTemp

		repositoryMock.On("Get", mock.Anything, mock.Anything).Return(actualSection, nil)
		repositoryMock.On("Exists", mock.Anything, mock.Anything).Return(false)
		repositoryMock.On("Update", mock.Anything, mock.Anything).Return(expected, nil)

		result, err := svc.Update(context.TODO(), updatePayload, sectionID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Equal(t, actualSection.ProductTypeID, result.ProductTypeID)
		assert.NotEqual(t, actualSection, result)
	})
	t.Run("Doesn't update a section by id and return error: not found", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		updatePayload := section.UpdateSection{
			SectionNumber:      testutil.ToPtr(sectionNumber),
			CurrentTemperature: testutil.ToPtr(sectionTemp),
		}

		repositoryMock.On("Get", mock.Anything, sectionID).Return(domain.Section{}, section.ErrNotFound)

		_, err := svc.Update(context.TODO(), updatePayload, sectionID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrNotFound)
		repositoryMock.AssertNumberOfCalls(t, "Exists", 0)
		repositoryMock.AssertNumberOfCalls(t, "Update", 0)
	})
	t.Run("Doesn't update a section by sectionNumber and return error: section number alredy exists", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := section.NewService(&repositoryMock)

		updatePayload := section.UpdateSection{
			SectionNumber: testutil.ToPtr(sectionNumber),
		}

		repositoryMock.On("Get", mock.Anything, sectionID).Return(domain.Section{}, nil)
		repositoryMock.On("Exists", mock.Anything, sectionNumber).Return(true)

		_, err := svc.Update(context.TODO(), updatePayload, sectionID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, section.ErrInvalidSectionNumber)
		repositoryMock.AssertNumberOfCalls(t, "Update", 0)
	})
}

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

func (r *RepositoryMock) Exists(ctx context.Context, SectionNumber int) bool {
	args := r.Called(ctx, SectionNumber)
	return args.Get(0).(bool)
}

func (r *RepositoryMock) Save(ctx context.Context, s domain.Section) (int, error) {
	args := r.Called(ctx, s)
	return args.Get(0).(int), args.Error(1)
}

func (r *RepositoryMock) Update(ctx context.Context, s domain.Section) (domain.Section, error) {
	args := r.Called(ctx, s)
	return args.Get(0).(domain.Section), args.Error(1)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}
