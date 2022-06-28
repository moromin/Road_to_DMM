package validate

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationErr struct {
	Errors []string
}

func (v *ValidationErr) Error() string {
	return strings.Join(v.Errors, ", ")
}

func Validate(v *validator.Validate, obj interface{}) error {
	err := v.Struct(obj)
	if err != nil {
		verr := &ValidationErr{}
		if _, ok := err.(*validator.ValidationErrors); ok {
			log.Println(err)
			return nil
		}

		for _, e := range err.(validator.ValidationErrors) {
			verr.Errors = append(verr.Errors, fmt.Sprintf("%s is %s", e.StructField(), e.Tag()))
		}
		return verr
	}
	return nil
}
