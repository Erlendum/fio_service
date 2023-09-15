package serviceImpl

import (
	"context"
	"fio_service/internal/models"
	"fio_service/internal/repository"
	"fio_service/internal/service"
	"fio_service/pkg/cache"
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

func (p *personServiceImplementation) Create(ctx context.Context, user *models.Person) error {
	return p.personRepository.Create(ctx, user)
}

func (p *personServiceImplementation) Delete(ctx context.Context, id uint64) error {
	return p.personRepository.Delete(ctx, id)
}

func (p *personServiceImplementation) Update(ctx context.Context, id uint64, fieldsToUpdate models.PersonFieldsToUpdate) error {
	return p.personRepository.Update(ctx, id, fieldsToUpdate)
}

func (p *personServiceImplementation) Get(ctx context.Context, id uint64) (*models.Person, error) {
	cachedPerson, err := p.cache.Get(ctx, "person:"+strconv.Itoa(int(id)))

	if err == nil {
		cachedData, ok := cachedPerson.(models.Person)
		if ok {
			return &cachedData, nil
		}
	}

	person, err := p.personRepository.Get(ctx, id)

	if err != nil {
		return person, err
	}

	if err := p.cache.Set(ctx, "person:"+strconv.Itoa(int(id)), person, p.ttlCache); err != nil {
		p.logger.Printf("error caching person with ID %d:", id)
	}

	return person, nil
}

func (p *personServiceImplementation) GetList(ctx context.Context) ([]models.Person, error) {
	cachedPerson, err := p.cache.Get(ctx, "persons")

	if err == nil {
		cachedData, ok := cachedPerson.([]models.Person)
		if ok {
			return cachedData, nil
		}
	}

	persons, err := p.personRepository.GetList(ctx)

	if err != nil {
		return persons, err
	}

	if err := p.cache.Set(ctx, "persons", persons, p.ttlCache); err != nil {
		p.logger.Print("error caching persons")
	}

	return persons, nil
}
