package db

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TFProvider struct {
	Model

	Name string `gorm:"unique"`
}

func (tfp *TFProvider) GetCondition() entity {
	return &TFProvider{
		Name: tfp.Name,
	}
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTFProvider(e *TFProvider) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"name"})
}

func (db *gDB) GetTFProvider(e *TFProvider, where *TFProvider) error {
	return get(db.g(), e, where)
}

func (db *gDB) GetOrCreateTFProvider(e *TFProvider) (id uuid.UUID, isNew bool, err error) {
	err = db.g().First(e, e.GetCondition()).Error
	if err != nil && !IsNotFoundError(err) {
		return uuid.Nil, true, err
	}
	if e.ID != uuid.Nil {
		return e.ID, false, nil
	}
	e.GenerateID()
	err = db.g().Create(e).Error
	return e.ID, true, err
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
