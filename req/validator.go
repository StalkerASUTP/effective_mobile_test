package req

import (
	"github.com/go-playground/validator/v10"
)

// глобальный валидатор
var validate *validator.Validate

// инициализация валидатора один раз
func init() {
	validate = validator.New()
	// здесь можно регистрировать кастомные валидаторы, если нужно
	// validate.RegisterValidation("adult", isAdult)
}

// универсальная функция валидации
func IsValid[T any](body T) error {
	return validate.Struct(body)
}
