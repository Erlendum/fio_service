package serviceImpl

import (
	"context"
	"encoding/json"
	"fio_service/internal/models"
	"fio_service/internal/repository"
	"fio_service/internal/service"
	"fio_service/pkg/cache"
	"fio_service/pkg/errors/repositoryErrors"
	"fio_service/pkg/logger"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type personServiceImplementation struct {
	personRepository repository.PersonRepository
	logger           *logger.Logger
	cache            cache.Cache
	ttlCache         time.Duration
}

func NewPersonServiceImplementation(personRepository repository.PersonRepository, logger *logger.Logger, cache cache.Cache, ttlCache time.Duration) service.PersonService {
	return &personServiceImplementation{
		personRepository: personRepository,
		logger:           logger,
		cache:            cache,
		ttlCache:         ttlCache,
	}
}

func (p *personServiceImplementation) Create(ctx context.Context, person *models.Person) error {
	fields := map[string]interface{}{"name": person.Name, "surname": person.Surname}
	err := p.personRepository.Create(ctx, person)
	if err != nil {
		p.logger.WithFields(fields).Error("person create failed: " + err.Error())
		return err
	}
	p.logger.WithFields(fields).Info("person create completed")
	return nil
}

type ageResponse struct {
	Count int64  `json:"count"`
	Name  string `json:"name"`
	Age   uint64 `json:"age"`
}

type genderResponse struct {
	Count       int64   `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type country struct {
	Country_id  string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type nationalityResponse struct {
	Count   int64     `json:"count"`
	Name    string    `json:"name"`
	Country []country `json:"country"`
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func (p *personServiceImplementation) CreateWithEnrichment(ctx context.Context, person *models.Person) error {
	fields := map[string]interface{}{"name": person.Name, "surname": person.Surname}
	if len(person.Name) == 0 || len(person.Surname) == 0 {
		return repositoryErrors.MissingRequiredFields
	}

	urlAge := "https://api.agify.io/?name=%s"
	urlGender := "https://api.genderize.io/?name=%s"
	urlNationality := "https://api.nationalize.io/?name=%s"
	urlAge = fmt.Sprintf(urlAge, url.QueryEscape(person.Name))
	urlGender = fmt.Sprintf(urlGender, url.QueryEscape(person.Name))
	urlNationality = fmt.Sprintf(urlNationality, url.QueryEscape(person.Name))

	ageResp := new(ageResponse)
	err := getJson(urlAge, ageResp)
	if err != nil {
		p.logger.WithFields(fields).Error("get age from api failed: " + err.Error())
		return err
	}

	genderResp := new(genderResponse)
	err = getJson(urlGender, genderResp)
	if err != nil {
		p.logger.WithFields(fields).Error("get gender from api failed: " + err.Error())
		return err
	}

	nationalityResp := new(nationalityResponse)
	err = getJson(urlNationality, nationalityResp)
	if err != nil {
		p.logger.WithFields(fields).Error("get nationality from api failed: " + err.Error())
		return err
	}

	person.Age = ageResp.Age
	genderResp.Gender = strings.ToUpper(string(genderResp.Gender[0])) + genderResp.Gender[1:]
	person.Gender = models.PersonGender(genderResp.Gender)
	person.Nationality = nationalityResp.Country[0].Country_id

	err = p.personRepository.Create(ctx, person)
	if err != nil {
		p.logger.WithFields(fields).Error("person create failed: " + err.Error())
		return err
	}
	p.logger.WithFields(fields).Info("person create completed")
	return nil
}

func (p *personServiceImplementation) Delete(ctx context.Context, id uint64) error {
	fields := map[string]interface{}{"id": id}
	err := p.personRepository.Delete(ctx, id)
	if err != nil {
		p.logger.WithFields(fields).Error("person delete failed: " + err.Error())
		return err
	}
	p.logger.WithFields(fields).Info("person delete completed")
	return nil
}

func (p *personServiceImplementation) Update(ctx context.Context, id uint64, fieldsToUpdate models.PersonFieldsToUpdate) error {
	fields := map[string]interface{}{"id": id}
	err := p.personRepository.Update(ctx, id, fieldsToUpdate)
	if err != nil {
		p.logger.WithFields(fields).Error("person update failed: " + err.Error())
		return err
	}
	p.logger.WithFields(fields).Info("person update completed")
	return nil
}

func (p *personServiceImplementation) Get(ctx context.Context, id uint64) (*models.Person, error) {
	fields := map[string]interface{}{"id": id}

	if p.cache != nil {
		cachedPerson, err := p.cache.Get(ctx, "person:"+strconv.Itoa(int(id)))

		if err == nil {
			cachedData, ok := cachedPerson.(models.Person)
			if ok {
				p.logger.WithFields(fields).Info("person update completed")
				return &cachedData, nil
			}
		}
	}

	person, err := p.personRepository.Get(ctx, id)

	if err != nil {
		p.logger.WithFields(fields).Error("person get failed: " + err.Error())
		return person, err
	}

	if p.cache != nil {
		if err := p.cache.Set(ctx, "person:"+strconv.Itoa(int(id)), person, p.ttlCache); err != nil {
			p.logger.WithFields(fields).Error("person caching failed: " + err.Error())
		}
	}
	p.logger.WithFields(fields).Info("person get completed")
	return person, nil
}

func (p *personServiceImplementation) GetList(ctx context.Context) ([]models.Person, error) {
	if p.cache != nil {
		cachedPerson, err := p.cache.Get(ctx, "persons")

		if err == nil {
			cachedData, ok := cachedPerson.([]models.Person)
			if ok {
				return cachedData, nil
			}
		}
	}
	persons, err := p.personRepository.GetList(ctx)

	if err != nil {
		p.logger.Error("person get list failed: " + err.Error())
		return persons, err
	}

	if p.cache != nil {
		if err := p.cache.Set(ctx, "persons", persons, p.ttlCache); err != nil {
			p.logger.Error("person list caching failed: " + err.Error())
		}
	}
	p.logger.Info("person get list completed")
	return persons, nil
}
