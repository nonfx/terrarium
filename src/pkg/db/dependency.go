package db

import (
	"github.com/google/uuid"
)

type Dependency struct {
	Model

	Title   string                 `json:"id" gorm:"uniqueIndex:dependency_unique"`
	Inputs  map[string]interface{} `json:"inputs" gorm:"type:json"`
	Outputs map[string]interface{} `json:"outputs" gorm:"type:json"`
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set the existing record ID in the given object
func (db *gDB) CreateDependencyInterface(e *Dependency) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"id"})
}

// MarshalInputOutput serializes the input and output properties to JSON format.
func (c *Dependency) MarshalInputOutput() error {
	// No need to unmarshal, since Inputs and Outputs are already of type json.RawMessage.
	return nil
}
