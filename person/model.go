package person

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"gorm.io/gorm"
)

// Person представляет модель данных человека
// @Description Модель данных человека с основной информацией
type Person struct {
	// ID записи
	ID uint `gorm:"primarykey" json:"id" example:"1"`
	// Дата создания
	CreatedAt time.Time `json:"created_at" example:"2025-05-29T00:00:00Z"`
	// Дата обновления
	UpdatedAt time.Time `json:"updated_at" example:"2025-05-29T00:00:00Z"`
	// Дата удаления (мягкое удаление)
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggertype:"string" example:"null"`
	// Имя человека
	Name string `json:"name" example:"Иван"`
	// Фамилия человека
	Surname string `json:"surname" example:"Иванов"`
	// Отчество (если есть)
	Patronymics string `json:"patronymics,omitempty" example:"Иванович"`
	// Возраст
	Age int `json:"age" example:"30"`
	// Пол (male/female)
	Gender string `json:"gender" example:"male"`
	// Национальность (код страны)
	Nationality string `json:"nationality" example:"RU"`
}

// NewPerson создает новую структуру Person
func NewPerson(name, surname, patronymics string) *Person {
	return &Person{
		Name:        name,
		Surname:     surname,
		Patronymics: patronymics,
	}
}

func GetAge(name string) (int, error) {
	encodeName := url.QueryEscape(name)
	apiUrl := fmt.Sprintf("https://api.agify.io/?name=%s", encodeName)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var result AgifyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}
	return result.Age, nil
}
func GetGender(name string) (string, error) {
	encodeName := url.QueryEscape(name)
	apiUrl := fmt.Sprintf("https://api.genderize.io/?name=%s", encodeName)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result GenderizeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	return result.Gender, nil
}
func GetNation(name string) (string, error) {
	encodeName := url.QueryEscape(name)
	apiUrl := fmt.Sprintf("https://api.nationalize.io/?name=%s", encodeName)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result NationalizeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	gender, probability := result.Country[0].CountryId, result.Country[0].Probability
	for _, e := range result.Country {
		if e.Probability > probability {
			gender = e.CountryId
			probability = e.Probability
		}
	}
	return gender, nil
}
