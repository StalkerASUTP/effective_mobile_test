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

// Create godoc
// @Summary Создать запись о человеке
// @Description Создает новую запись, обогащая данные из внешних API (возраст, пол, национальность)
// @Tags person
// @Accept json
// @Produce json
// @Param input body CreatePersonRequest true "Данные о человеке"
// @Success 201 {object} Person
// @Failure 400 {string} string "Неверные данные или ошибка при обращении к внешнему API"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /person [post]
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, createdPerson, http.StatusCreated)
	}
}

// Update godoc
// @Summary Обновить данные человека
// @Description Обновляет данные существующей записи
// @Tags person
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Param input body UpdatePersonRequest true "Обновляемые данные"
// @Success 200 {object} Person
// @Failure 400 {string} string "Неверные данные или ошибка при обращении к внешнему API"
// @Failure 404 {string} string "Не найдено"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /person/{id} [patch]
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		res.Json(w, updatedPerson, http.StatusOK)
	}
}

// Delete godoc
// @Summary Удалить запись о человеке
// @Description Удаляет запись по указанному ID
// @Tags person
// @Produce json
// @Param id path int true "ID человека"
// @Success 200
// @Failure 400 {string} string "Неверные данные или ошибка при обращении к внешнему API"
// @Failure 404 {string} string "Не найдено"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /person/{id} [delete]
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
// GetWithParams godoc
// @Summary Получить список людей с фильтрацией
// @Description Возвращает список людей с возможностью фильтрации и пагинации
// @Tags person
// @Produce json
// @Param name query []string false "Фильтр по имени" collectionFormat(multi)
// @Param surname query []string false "Фильтр по фамилии" collectionFormat(multi)
// @Param gender query []string false "Фильтр по полу" collectionFormat(multi)
// @Param nationality query []string false "Фильтр по национальности" collectionFormat(multi)
// @Param age_from query int false "Минимальный возраст"
// @Param age_to query int false "Максимальный возраст"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Лимит записей (по умолчанию 20, максимум 100)"
// @Success 200 {object} GetWithParamResponse
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /person [get]
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
		filters := map[string]any{
			"name":        getParams(query, "name"),
			"surname":     getParams(query, "surname"),
			"gender":      getParams(query, "gender"),
			"nationality": getParams(query, "nationality"),
		}
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
