package serviceImpl

import (
	"context"
	"fio_service/internal/models"
	"fio_service/internal/repository"
	"fio_service/internal/service"
	"fio_service/pkg/cache"
	"fio_service/pkg/errors/repositoryErrors"
	"fio_service/pkg/logger"
	"strconv"
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

func (p *personServiceImplementation) CreateWithEnrichment(ctx context.Context, person *models.Person) error {
	fields := map[string]interface{}{"name": person.Name, "surname": person.Surname}
	if len(person.Name) == 0 || len(person.Surname) == 0 {
		return repositoryErrors.MissingRequiredFields
	}
	err := p.personRepository.Create(ctx, person)
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
