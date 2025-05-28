package person

import (
	"fmt"
	"go/ef-mob-api/db"
	"log"

	"gorm.io/gorm"
)

type PersonRepository struct {
	Database *db.Db
}

func NewPersonRepository(database *db.Db) *PersonRepository {
	return &PersonRepository{
		Database: database,
	}
}
func (repo *PersonRepository) Create(person *Person) (*Person, error) {
	result := repo.Database.DB.Create(person)
	if result.Error != nil {
		return nil, result.Error
	}
	return person, nil
}
func (repo *PersonRepository) GetById(id uint) (*Person, error) {
	var person Person
	result := repo.Database.DB.First(&person, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &person, nil
}
func (repo *PersonRepository) Update(person *Person) (*Person, error) {
	result := repo.Database.DB.Save(person)
	if result.Error != nil {
		return nil, result.Error
	}
	return person, nil
}
func (repo *PersonRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&Person{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
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
