package ferry

import "github.com/go-playground/validator/v10"

// Config -- Configuration for ferry
type Config struct {
	Validator *validator.Validate
}
