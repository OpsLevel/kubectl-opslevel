package opslevel

import (
	"github.com/go-playground/validator/v10"
	"github.com/opslevel/moredefaults"
)

var structValidator = validator.New(validator.WithRequiredStructEnabled())

// Validates resource's `validate:""` struct tags
func IsResourceValid[T any](opslevelResource T) error {
	return structValidator.Struct(opslevelResource)
}

// Apply resource's `default:""` struct tags
func SetDefaultsFor[T any](opslevelResource *T) {
	validator.New(validator.WithRequiredStructEnabled())
	if err := moredefaults.Set(opslevelResource); err != nil {
		panic(err)
	}
}

// Apply resource's `example:""` struct tags
func SetExamplesFor[T any](opslevelResource *T) {
	validator.New(validator.WithRequiredStructEnabled())
	if err := moredefaults.Set(opslevelResource, "example"); err != nil {
		panic(err)
	}
}

// Make new OpsLevel resource with defaults set
func NewExampleOf[T any]() T {
	var opslevelResource T
	SetExamplesFor[T](&opslevelResource)
	return opslevelResource
}
