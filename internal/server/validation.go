package server

import (
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"strings"
)

func (a *Server) setupValidation() error {
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)

	var ok bool
	a.translator, ok = uni.GetTranslator("en")
	if !ok {
		return fmt.Errorf("error retrieving locale for translations")

	}
	a.validator = validator.New(validator.WithRequiredStructEnabled())
	err := en2.RegisterDefaultTranslations(a.validator, a.translator)
	if err != nil {
		return fmt.Errorf("%q: %w", "error registering translations", err)
	}
	a.validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
	return nil

}

func (a *Server) createValidationError(ve validator.ValidationErrors) map[string]any {
	e := map[string]any{}
	for _, err := range ve {
		e[err.Field()] = err.Translate(a.translator)
	}

	return e

}
