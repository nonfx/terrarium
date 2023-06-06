package db

import "github.com/google/uuid"

type TFProvider struct {
	Model

	Name string `gorm:"unique"`
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTFProvider(e *TFProvider) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"name"})
}

func (db *gDB) GetTFProvider(e *TFProvider, where *TFProvider) error {
	return get(db.g(), e, where)
}

func (db *gDB) GetOrCreateTFProvider(e *TFProvider) (isNew bool, err error) {
	result := db.g().Where(TFProvider{Name: e.Name}).FirstOrCreate(e)
	return result.RowsAffected > 0, result.Error
}
