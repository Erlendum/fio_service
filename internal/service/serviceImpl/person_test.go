package serviceImpl

import (
	"context"
	"fio_service/internal/models"
	mock_repository "fio_service/internal/repository/mocks"
	"fio_service/internal/service"
	"fio_service/pkg/errors/repositoryErrors"
	"fio_service/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

type personServiceFields struct {
	personRepositoryMock *mock_repository.MockPersonRepository
}

func createPersonServiceFields(controller *gomock.Controller) *personServiceFields {
	fields := new(personServiceFields)

	fields.personRepositoryMock = mock_repository.NewMockPersonRepository(controller)

	return fields
}

func createPersonService(fields *personServiceFields) service.PersonService {
	return NewPersonServiceImplementation(fields.personRepositoryMock, logger.New("/dev/null", ""), nil, 0)
}

var testCreateSuccess = []struct {
	TestName  string
	InputData struct {
		person *models.Person
	}
	Prepare     func(fields *personServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			person *models.Person
		}{person: &models.Person{Name: "Vasya", Surname: "Pupkin"}},
		Prepare: func(fields *personServiceFields) {
			fields.personRepositoryMock.EXPECT().Create(context.Background(), &models.Person{Name: "Vasya", Surname: "Pupkin"}).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

func TestPersonServiceImplementation_Create(t *testing.T) {
	t.Parallel()

	for _, tt := range testCreateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createPersonServiceFields(ctrl)
			tt.Prepare(fields)

			personService := createPersonService(fields)

			err := personService.Create(context.Background(), tt.InputData.person)

			tt.CheckOutput(t, err)
		})
	}

}

var testGetSuccess = []struct {
	TestName  string
	InputData struct {
		id uint64
	}
	Prepare     func(fields *personServiceFields)
	CheckOutput func(t *testing.T, person *models.Person, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *personServiceFields) {
			fields.personRepositoryMock.EXPECT().Get(context.Background(), uint64(1)).Return(&models.Person{Name: "Vasya", Surname: "Pupkin"}, nil)
		},
		CheckOutput: func(t *testing.T, person *models.Person, err error) {
			require.NoError(t, err)
			require.Equal(t, &models.Person{Name: "Vasya", Surname: "Pupkin"}, person)
		},
	},
}

var testGetFailed = []struct {
	TestName  string
	InputData struct {
		id uint64
	}
	Prepare     func(fields *personServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "person does not exists",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *personServiceFields) {
			fields.personRepositoryMock.EXPECT().Get(context.Background(), uint64(1)).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.ObjectDoesNotExists)
		},
	},
}

func TestPersonServiceImplementation_Get(t *testing.T) {
	t.Parallel()

	for _, tt := range testGetSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createPersonServiceFields(ctrl)
			tt.Prepare(fields)

			personService := createPersonService(fields)

			p, err := personService.Get(context.Background(), tt.InputData.id)

			tt.CheckOutput(t, p, err)
		})
	}
	for _, tt := range testGetFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createPersonServiceFields(ctrl)
			tt.Prepare(fields)

			personService := createPersonService(fields)

			_, err := personService.Get(context.Background(), tt.InputData.id)

			tt.CheckOutput(t, err)
		})
	}

}

var testDeleteSuccess = []struct {
	TestName  string
	InputData struct {
		id uint64
	}
	Prepare     func(fields *personServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *personServiceFields) {
			fields.personRepositoryMock.EXPECT().Delete(context.Background(), uint64(1)).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testDeleteFailed = []struct {
	TestName  string
	InputData struct {
		id uint64
	}
	Prepare     func(fields *personServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "person does not exists",
		InputData: struct {
			id uint64
		}{id: 1},
		Prepare: func(fields *personServiceFields) {
			fields.personRepositoryMock.EXPECT().Delete(context.Background(), uint64(1)).Return(repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.ObjectDoesNotExists)
		},
	},
}

func TestPersonServiceImplementation_Delete(t *testing.T) {
	t.Parallel()

	for _, tt := range testDeleteSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createPersonServiceFields(ctrl)
			tt.Prepare(fields)

			personService := createPersonService(fields)

			err := personService.Delete(context.Background(), tt.InputData.id)

			tt.CheckOutput(t, err)
		})
	}
	for _, tt := range testDeleteFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createPersonServiceFields(ctrl)
			tt.Prepare(fields)

			personService := createPersonService(fields)
			err := personService.Delete(context.Background(), tt.InputData.id)

			tt.CheckOutput(t, err)
		})
	}
}

var testGetListSuccess = []struct {
	TestName  string
	InputData struct {
	}
	Prepare     func(fields *personServiceFields)
	CheckOutput func(t *testing.T, persons []models.Person, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
		}{},
		Prepare: func(fields *personServiceFields) {
			fields.personRepositoryMock.EXPECT().GetList(context.Background()).Return([]models.Person{{Name: "Vasya", Surname: "Pupkin"}}, nil)
		},
		CheckOutput: func(t *testing.T, persons []models.Person, err error) {
			require.NoError(t, err)
			require.Equal(t, []models.Person{{Name: "Vasya", Surname: "Pupkin"}}, persons)
		},
	},
}

var testGetListFailed = []struct {
	TestName  string
	InputData struct {
	}
	Prepare     func(fields *personServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "person does not exists",
		InputData: struct {
		}{},
		Prepare: func(fields *personServiceFields) {
			fields.personRepositoryMock.EXPECT().GetList(context.Background()).Return(nil, repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.ObjectDoesNotExists)
		},
	},
}

func TestPersonServiceImplementation_GetList(t *testing.T) {
	t.Parallel()

	for _, tt := range testGetListSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createPersonServiceFields(ctrl)
			tt.Prepare(fields)

			personService := createPersonService(fields)

			p, err := personService.GetList(context.Background())

			tt.CheckOutput(t, p, err)
		})
	}
	for _, tt := range testGetListFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createPersonServiceFields(ctrl)
			tt.Prepare(fields)

			personService := createPersonService(fields)

			_, err := personService.GetList(context.Background())

			tt.CheckOutput(t, err)
		})
	}

}

var testUpdateSuccess = []struct {
	TestName  string
	InputData struct {
		id             uint64
		fieldsToUpdate models.PersonFieldsToUpdate
	}
	Prepare     func(fields *personServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "usual test",
		InputData: struct {
			id             uint64
			fieldsToUpdate models.PersonFieldsToUpdate
		}{id: 1, fieldsToUpdate: map[models.PersonField]any{models.PersonFieldName: "Jora"}},
		Prepare: func(fields *personServiceFields) {
			fields.personRepositoryMock.EXPECT().Update(context.Background(), uint64(1), map[models.PersonField]any{models.PersonFieldName: "Jora"}).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

var testUpdateFailed = []struct {
	TestName  string
	InputData struct {
		id             uint64
		fieldsToUpdate models.PersonFieldsToUpdate
	}
	Prepare     func(fields *personServiceFields)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "person does not exists",
		InputData: struct {
			id             uint64
			fieldsToUpdate models.PersonFieldsToUpdate
		}{id: 1, fieldsToUpdate: map[models.PersonField]any{models.PersonFieldName: "Jora"}},
		Prepare: func(fields *personServiceFields) {
			fields.personRepositoryMock.EXPECT().Update(context.Background(), uint64(1), map[models.PersonField]any{models.PersonFieldName: "Jora"}).Return(repositoryErrors.ObjectDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			require.ErrorIs(t, err, repositoryErrors.ObjectDoesNotExists)
		},
	},
}

func TestPersonServiceImplementation_Update(t *testing.T) {
	t.Parallel()

	for _, tt := range testUpdateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createPersonServiceFields(ctrl)
			tt.Prepare(fields)

			personService := createPersonService(fields)

			err := personService.Update(context.Background(), tt.InputData.id, tt.InputData.fieldsToUpdate)

			tt.CheckOutput(t, err)
		})
	}
	for _, tt := range testUpdateFailed {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := createPersonServiceFields(ctrl)
			tt.Prepare(fields)

			personService := createPersonService(fields)
			err := personService.Update(context.Background(), tt.InputData.id, tt.InputData.fieldsToUpdate)

			tt.CheckOutput(t, err)
		})
	}
}
