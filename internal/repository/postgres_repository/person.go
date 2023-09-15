package postgres_repository

import (
	"context"
	"database/sql"
	"errors"
	"fio_service/internal/models"
	"fio_service/internal/repository"
	"fio_service/pkg/errors/repositoryErrors"
	"fio_service/pkg/queries"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type PersonPostgres struct {
	Id          uint64              `db:"id"`
	Name        string              `db:"name"`
	Surname     string              `db:"surname"`
	Patronymic  string              `db:"patronymic"`
	Gender      models.PersonGender `db:"gender"`
	Age         uint64              `db:"age"`
	Nationality string              `db:"nationality"`
}

var personFieldToDBField = map[models.PersonField]string{
	models.PersonFieldName:        "name",
	models.PersonFieldSurname:     "surname",
	models.PersonFieldPatronymic:  "patronymic",
	models.PersonFieldAge:         "age",
	models.PersonFieldGender:      "gender",
	models.PersonFieldNationality: "nationality",
}

type PersonPostgresRepository struct {
	db *sqlx.DB
}

func NewPersonPostgresRepository(db *sqlx.DB) repository.PersonRepository {
	return &PersonPostgresRepository{db: db}
}

func (p *PersonPostgresRepository) Create(ctx context.Context, user *models.Person) error {
	query := `insert into service.persons (name, surname, patronymic, age, gender, nationality) values
											 ($1, $2, $3, $4, $5, $6);`
	_, err := p.db.ExecContext(ctx, query, user.Name, user.Surname, user.Patronymic, user.Age,
		user.Gender, user.Nationality)
	if err != nil {
		return err
	}
	return nil
}

func (p *PersonPostgresRepository) Delete(ctx context.Context, id uint64) error {
	query := `delete from service.persons where id = $1`
	res, err := p.db.ExecContext(ctx, query, id)
	count, _ := res.RowsAffected()
	if count == 0 || errors.Is(err, sql.ErrNoRows) {
		return repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return err
	}
	return nil
}

func (p *PersonPostgresRepository) Update(ctx context.Context, id uint64, fieldsToUpdate models.PersonFieldsToUpdate) error {
	if len(fieldsToUpdate) == 0 {
		return nil
	}
	updateFields := make(map[string]any, len(fieldsToUpdate))
	for key, value := range fieldsToUpdate {
		field, err := personFieldToDBField[key]
		if !err {
			return repositoryErrors.InvalidField
		}
		updateFields[field] = value
	}

	query, fields := queries.CreateSQLUpdateQuery("service.persons", updateFields)

	fields = append(fields, id)
	query += ` where id = $` + strconv.Itoa(len(fields)) + ";"

	res, err := p.db.ExecContext(ctx, query, fields...)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 || errors.Is(err, sql.ErrNoRows) {
		return repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return err
	}
	return nil
}

func (p *PersonPostgresRepository) Get(ctx context.Context, id uint64) (*models.Person, error) {
	query := `select * from service.persons where id = $1`
	personPostgres := &PersonPostgres{}

	err := p.db.GetContext(ctx, personPostgres, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}
	person := &models.Person{}
	err = copier.Copy(person, personPostgres)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (p *PersonPostgresRepository) GetList(ctx context.Context) ([]models.Person, error) {
	query := `select * from service.persons order by id;`

	var personsPostgres []PersonPostgres
	var persons []models.Person
	err := p.db.SelectContext(ctx, &personsPostgres, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}

	for i := range personsPostgres {
		person := &models.Person{}
		err = copier.Copy(person, &personsPostgres[i])
		if err != nil {
			return nil, err
		}
		persons = append(persons, *person)
	}
	return persons, nil
}
