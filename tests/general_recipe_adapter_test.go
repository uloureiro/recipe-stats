package tests

import (
	"recipe-stats/adapters"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	filePath := "./testdata/test_calculation_fixtures_single.json"
	adptr := adapters.NewGeneralRecipeAdapter()
	result, err := adptr.Unmarshal(filePath)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestUnmarshalParsing(t *testing.T) {
	filePath := "./testdata/test_calculation_fixtures_single.json"
	adptr := adapters.NewGeneralRecipeAdapter()

	result, err := adptr.Unmarshal(filePath)

	assert.NoError(t, err)
	assert.IsType(t, adapters.GeneralRecipe{}, (*result)[0])
}
