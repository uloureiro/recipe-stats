package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGet(t *testing.T) {
	filePath := "./testdata/test_calculation_fixtures_double.json"
	rnsk := LoadRecipeNameSliceKeeperHelper(filePath)

	filteredRecipes, ok := rnsk.Get("Pork")

	assert.True(t, ok)
	assert.NotNil(t, filteredRecipes)
	assert.Equal(t, len(filteredRecipes), 2)
}

func TestGetSome(t *testing.T) {
	filePath := "./testdata/test_calculation_fixtures_full.json"
	rnsk := LoadRecipeNameSliceKeeperHelper(filePath)

	filteredRecipes := rnsk.GetSome([]string{"Cheese", "Stovetop"})

	assert.Equal(t, filteredRecipes[0].Recipe, "Grilled Cheese and Veggie Jumble")
	assert.Equal(t, filteredRecipes[1].Recipe, "Stovetop Mac 'N' Cheese")
	assert.Equal(t, len(filteredRecipes), 2)
}

func TestGetSomeDuplicatedQuery(t *testing.T) {
	filePath := "./testdata/test_calculation_fixtures_full.json"
	rnsk := LoadRecipeNameSliceKeeperHelper(filePath)

	filteredRecipes := rnsk.GetSome([]string{"Cheese", "Cheese"})

	assert.Equal(t, filteredRecipes[0].Recipe, "Grilled Cheese and Veggie Jumble")
	assert.Equal(t, filteredRecipes[1].Recipe, "Stovetop Mac 'N' Cheese")
	assert.Equal(t, len(filteredRecipes), 2)
}
