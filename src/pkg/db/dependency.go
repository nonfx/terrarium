package db

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Dependency struct {
	Model

	ID uuid.UUID `gorm:"primaryKey"`
	// Taxonomy    []string
	Title       string
	Description string
	Engine      string
	Version     string
	Inputs      map[string]interface{} `gorm:"type:jsonb"`
	Outputs     map[string]interface{} `gorm:"type:jsonb"`
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateComponent(e *Dependency) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"id"})
}

// MarshalInputOutput serializes the input and output properties to JSON format.
func (c *Dependency) MarshalInputOutput() error {
	inputsJSON, err := json.Marshal(c.Inputs)
	if err != nil {
		return err
	}

	outputsJSON, err := json.Marshal(c.Outputs)
	if err != nil {
		return err
	}

	c.Inputs = nil
	c.Outputs = nil

	if err := json.Unmarshal(inputsJSON, &c.Inputs); err != nil {
		return err
	}

	if err := json.Unmarshal(outputsJSON, &c.Outputs); err != nil {
		return err
	}

	return nil
}
