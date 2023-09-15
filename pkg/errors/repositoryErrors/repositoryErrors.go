package repositoryErrors

import (
	"errors"
	"fmt"
)

var (
	InternalRepositoryError = errors.New("internal repository error")

	DoesNotExists       = errors.New("does not exists")
	ObjectDoesNotExists = fmt.Errorf("object %w", DoesNotExists)

	InvalidField = errors.New("invalid fields")
)
