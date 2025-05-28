package person

// CreatePersonRequest представляет структуру запроса для создания Person
// @Description Запрос на создание новой записи о человеке
type CreatePersonRequest struct {
	// Имя (только буквы, обязательное поле)
	Name        string `json:"name" validate:"required,alphaunicode" example:"Иван"`
	// Фамилия (только буквы, обязательное поле)
	Surname     string `json:"surname" validate:"required,alphaunicode" example:"Иванов"`
	// Отчество (необязательное, только буквы)
	Patronymics string `json:"patronymics,omitempty" validate:"alphaunicode" example:"Иванович"`
}

// UpdatePersonRequest представляет структуру запроса для обновления Person
// @Description Запрос на обновление данных о человеке (все поля необязательные)
type UpdatePersonRequest struct {
	// Имя (только буквы)
	Name        *string `json:"name,omitempty" validate:"alphaunicode" example:"Петр"`
	// Фамилия (только буквы)
	Surname     *string `json:"surname,omitempty" validate:"alphaunicode" example:"Петров"`
	// Отчество (только буквы)
	Patronymics *string `json:"patronymics,omitempty" validate:"alphaunicode" example:"Петрович"`
	// Возраст (от 0 до 120)
	Age         *int    `json:"age,omitempty" validate:"gte=0,lte=120" example:"30"`
	// Пол (male/female)
	Gender      *string `json:"gender,omitempty" validate:"oneof=male female" example:"male"`
	// Национальность (код страны ISO 3166-1 alpha-2)
	Nationality *string `json:"nationality,omitempty" validate:"iso3166_1_alpha2" example:"RU"`
}

// GetWithParamResponse представляет структуру ответа с пагинацией
// @Description Ответ со списком людей и общим количеством
type GetWithParamResponse struct {
	// Список людей
	Persons []Person `json:"persons"`
	// Общее количество записей
	Count   int64    `json:"count" example:"100"`
}

// AgifyResponse представляет структуру ответа от Agify API
type AgifyResponse struct {
	// Предполагаемый возраст
	Age int `json:"age"`
}

// GenderizeResponse представляет структуру ответа от Genderize API
type GenderizeResponse struct {
	// Предполагаемый пол (male/female)
	Gender string `json:"gender"`
}

// NationalizeResponse представляет структуру ответа от Nationalize API
type NationalizeResponse struct {
	// Список стран с вероятностями
	Country []Country `json:"country"`
}

// Country представляет информацию о стране из Nationalize API
type Country struct {
	// Код страны (ISO 3166-1 alpha-2)
	CountryId   string  `json:"country_id"`
	// Вероятность соответствия
	Probability float64 `json:"probability"`
}