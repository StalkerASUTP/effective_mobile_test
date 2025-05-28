package person

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	Name        string
	Surname     string
	Patronymics string
	Age         int
	Gender      string
	Nationality string
}

func NewPerson(name, surname, patronymics string) *Person {
	newPerson := &Person{
		Name:        name,
		Surname:     surname,
		Patronymics: patronymics,
	}
	return newPerson
}

// func (person *Person) AddAge() {
// 	person.Age =
// }

// func (person *Person) AddGender() {
// 	person.Gender = "male"
// }

// func (person *Person) AddNation() {
// 	person.Nationality = "tatar"
// }

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
