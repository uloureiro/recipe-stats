package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	filePath := "./testdata/test_calculation_fixtures_full.json"
	rk := LoadRecipeKeeperHelper(filePath)

	assert.Equal(t, 26, rk.Count())
}
