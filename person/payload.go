package person

type CreatePersonRequest struct {
	Name        string `json:"name" validate:"required,alphaunicode"`
	Surname     string `json:"surname"  validate:"required,alphaunicode"` // не учитываем двойные фамилии
	Patronymics string `json:"patronymics" validate:"omitempty,alphaunicode"`
}
type UpdatePersonRequest struct {
	Name        *string `json:"name" validate:"omitempty,alphaunicode"`
	Surname     *string `json:"surname" validate:"omitempty,alphaunicode"`
	Patronymics *string `json:"patronymics" validate:"omitempty,alphaunicode"`
	Age         *int    `json:"age" validate:"omitempty,gte=0,lte=120"`
	Gender      *string `json:"gender" validate:"omitempty,oneof=male female"`
	Nationality *string `json:"nationality" validate:"omitempty,iso3166_1_alpha2"`
}
type GetWithParamResponse struct {
	Persons []Person `json:"persons"`
	Count   int64    `json:"count"`
}
type AgifyResponse struct {
	Age int `json:"age"`
}
type GenderizeResponse struct {
	Gender string `json:"gender"`
}
type NationalizeResponse struct {
	Country []Country `json:"country"`
}
type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
