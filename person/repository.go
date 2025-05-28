package person

import (
	"fmt"
	"go/ef-mob-api/db"
	"log"

	"gorm.io/gorm"
)

// PersonRepository обеспечивает взаимодействие с базой данных для сущности Person
type PersonRepository struct {
	Database *db.Db
}

// NewPersonRepository создает новый экземпляр PersonRepository
func NewPersonRepository(database *db.Db) *PersonRepository {
	return &PersonRepository{
		Database: database,
	}
}

// Create создает новую запись о человеке в базе данных
// @Summary Создать запись
// @Description Создает новую запись о человеке в базе данных
// @Param person body Person true "Данные человека"
// @Success 200 {object} Person
// @Failure 500 {string} string "Ошибка сервера"
func (repo *PersonRepository) Create(person *Person) (*Person, error) {
	result := repo.Database.DB.Create(person)
	if result.Error != nil {
		return nil, result.Error
	}
	return person, nil
}

// GetById получает запись о человеке по ID
// @Summary Получить по ID
// @Description Возвращает запись о человеке по указанному ID
// @Param id path int true "ID человека"
// @Success 200 {object} Person
// @Failure 404 {string} string  "Запись не найдена"
// @Failure 500 {string} string "Ошибка сервера"
func (repo *PersonRepository) GetById(id uint) (*Person, error) {
	var person Person
	result := repo.Database.DB.First(&person, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &person, nil
}

// Update обновляет существующую запись о человеке
// @Summary Обновить запись
// @Description Обновляет данные о человеке в базе данных
// @Param person body Person true "Обновленные данные"
// @Success 200 {object} Person
// @Failure 500 {string} string  "Ошибка сервера"
func (repo *PersonRepository) Update(person *Person) (*Person, error) {
	result := repo.Database.DB.Save(person)
	if result.Error != nil {
		return nil, result.Error
	}
	return person, nil
}

// Delete удаляет запись о человеке по ID
// @Summary Удалить запись
// @Description Удаляет запись о человеке по указанному ID
// @Param id path int true "ID человека"
// @Success 200 {string} string "Успешное удаление"
// @Failure 500 {string} string  "Ошибка сервера"
func (repo *PersonRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&Person{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetWithFilters возвращает отфильтрованный список людей с пагинацией
// @Summary Получить с фильтрами
// @Description Возвращает список людей с возможностью фильтрации и пагинации
// @Param name query []string false "Фильтр по имени (можно несколько через запятую)" collectionFormat(multi)
// @Param surname query []string false "Фильтр по фамилии (можно несколько через запятую)" collectionFormat(multi)
// @Param gender query []string false "Фильтр по полу (можно несколько через запятую)" collectionFormat(multi)
// @Param nationality query []string false "Фильтр по национальности (можно несколько через запятую)" collectionFormat(multi)
// @Param limit query int false "Лимит записей"
// @Param offset query int false "Смещение"
// @Success 200 {object} GetWithParamResponse
// @Failure 500 {string} string "Ошибка сервера"
func (repo *PersonRepository) GetWithFilters(filters map[string]any, limit, offset int) ([]Person, int64, error) {
	var persons []Person
	var total int64
	dbQuery := repo.Database.DB.Model(&Person{})
	dbQuery = applyFilter(dbQuery, "name", filters["name"])
	dbQuery = applyFilter(dbQuery, "surname", filters["surname"])
	dbQuery = applyFilter(dbQuery, "gender", filters["gender"])
	dbQuery = applyFilter(dbQuery, "narionality", filters["nationality"])
	dbQuery = applyFilter(dbQuery, "age_from", filters["age_from"])
	dbQuery = applyFilter(dbQuery, "age_to", filters["age_to"])
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}
	if err := dbQuery.Find(&persons).Error; err != nil {
		return nil, 0, err
	}

	return persons, total, nil
}

// applyFilter применяет фильтры к запросу
func applyFilter(dbQuery *gorm.DB, field string, value any) *gorm.DB {
	if value == nil {
		return dbQuery
	}
	fmt.Println(value)
	switch v := value.(type) {
	case []string:
		if len(v) > 1 {
			dbQuery = dbQuery.Where(field+" IN (?)", v)
		} else if len(v) == 1 {
			dbQuery = dbQuery.Where(field+" = ?", v[0])
		}
	case int:
		if field == "age_from" {
			dbQuery = dbQuery.Where("age >= ?", v)
		} else {
			dbQuery = dbQuery.Where("age <= ?", v)
		}
	default:
		log.Printf("unsupported type for field %s: %T", field, v)
	}
	return dbQuery

}
