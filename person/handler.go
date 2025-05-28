package person

import (
	"go/ef-mob-api/configs"
	"go/ef-mob-api/req"
	"go/ef-mob-api/res"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type PersonHandlerDeps struct {
	PersonRepository *PersonRepository
	Config           *configs.Config
}
type PersonHandler struct {
	PersonRepository *PersonRepository
}

func NewPersonHandler(router *http.ServeMux, deps PersonHandlerDeps) {
	handler := &PersonHandler{PersonRepository: deps.PersonRepository}

	router.Handle("POST /person", handler.Create())
	router.Handle("PATCH /person/{id}", handler.Update())
	router.Handle("DELETE /person/{id}", handler.Delete())
	router.Handle("GET /person", handler.GetWithParams())
}

func (handler *PersonHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[CreatePersonRequest](&w, r)
		if err != nil {
			return
		}
		person := NewPerson(body.Name, body.Surname, body.Patronymics)
		person.Age, err = GetAge(person.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		person.Gender, err = GetGender(person.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		person.Nationality, err = GetNation(person.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		createdPerson, err := handler.PersonRepository.Create(person)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, createdPerson, http.StatusCreated)
	}
}
func (handler *PersonHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}
		person, err := handler.PersonRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		body, err := req.HandleBody[UpdatePersonRequest](&w, r)
		if err != nil {
			return
		}
		updatePersonFields(person, body)
		updatedPerson, err := handler.PersonRepository.Update(person)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}
		res.Json(w, updatedPerson, http.StatusOK)
	}
}
func (handler *PersonHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}
		_, err = handler.PersonRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = handler.PersonRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		res.Json(w, nil, http.StatusOK)

	}
}

func (handler *PersonHandler) GetWithParams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		page, _ := strconv.Atoi(query.Get("page"))
		if page <= 0 {
			page = 1
		}

		limit, _ := strconv.Atoi(query.Get("limit"))
		if limit <= 0 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit
		filters := make(map[string]any)
		filters["name"] = getParams(query, "name")
		filters["surname"] = getParams(query, "surname")
		filters["gender"] = getParams(query, "gender")
		filters["nationality"] = getParams(query, "nationality")
		if ageFrom := query.Get("age_from"); ageFrom != "" {
			if age, err := strconv.Atoi(ageFrom); err == nil {
				filters["age_from"] = age
			}
		}
		if ageTo := query.Get("age_to"); ageTo != "" {
			if age, err := strconv.Atoi(ageTo); err == nil {
				filters["age_to"] = age
			}
		}
		persons, total, err := handler.PersonRepository.GetWithFilters(filters, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := GetWithParamResponse{
			Persons: persons,
			Count:   total,
		}
		res.Json(w, response, http.StatusOK)
	}
}

func updatePersonFields(person *Person, update *UpdatePersonRequest) {
	if update.Name != nil {
		person.Name = *update.Name
	}
	if update.Surname != nil {
		person.Surname = *update.Surname
	}
	if update.Patronymics != nil {
		person.Patronymics = *update.Patronymics
	}
	if update.Age != nil {
		person.Age = *update.Age
	}
	if update.Gender != nil {
		person.Gender = *update.Gender
	}
	if update.Nationality != nil {
		person.Nationality = *update.Nationality
	}
}

func getParams(query url.Values, paramName string) []string {
	values := query[paramName]
	if len(values) == 0 {
		return nil
	}
	var result []string
	for _, v := range values {
		parts := strings.Split(v, ",")
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
	}
	sort.Strings(result)
	result = result[:uniqueSlice(result)]
	return result
}

func uniqueSlice(slice []string) int {
	i := 0
	mappa := make(map[string]bool)
	for _, v := range slice {
		if _, dupl := mappa[v]; !dupl {
			mappa[v] = true
			slice[i] = v
			i++
		}
	}
	slice = slice[:i]
	return len(slice)
}
