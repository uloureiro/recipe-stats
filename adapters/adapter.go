package adapters

import (
	"recipe-stats/models"
)

// This is a baseline to keep in mind if the program could scale to other input
// formats

type Adapter interface {
	unwrap(file []byte) ([]AdapterMember, error)
	Unmarshal(filePath string) ([]models.Recipe, error)
}

type AdapterMember interface {
	ToRecipe() []models.Recipe
	ToDelivery() []models.Delivery
}
