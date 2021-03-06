package tests

import (
	"recipe-stats/keepers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeToIndex(t *testing.T) {
	type expected struct {
		time  string
		value int
	}
	expectations := []expected{
		expected{time: "12AM", value: 0},
		expected{time: "1AM", value: 1},
		expected{time: "2AM", value: 2},
		expected{time: "3AM", value: 3},
		expected{time: "4AM", value: 4},
		expected{time: "5AM", value: 5},
		expected{time: "6AM", value: 6},
		expected{time: "7AM", value: 7},
		expected{time: "8AM", value: 8},
		expected{time: "9AM", value: 9},
		expected{time: "10AM", value: 10},
		expected{time: "11AM", value: 11},
		expected{time: "12PM", value: 12},
		expected{time: "1PM", value: 13},
		expected{time: "2PM", value: 14},
		expected{time: "3PM", value: 15},
		expected{time: "4PM", value: 16},
		expected{time: "5PM", value: 17},
		expected{time: "6PM", value: 18},
		expected{time: "7PM", value: 19},
		expected{time: "8PM", value: 20},
		expected{time: "9PM", value: 21},
		expected{time: "10PM", value: 22},
		expected{time: "11PM", value: 23},
	}

	for i := 0; i < len(expectations); i++ {
		value := keepers.TimeToIndex(expectations[i].time)
		assert.Equal(t, expectations[i].value, value)
	}
}

func TestCountByInterval(t *testing.T) {
	filePath := "./testdata/test_calculation_fixtures_10122.json"
	dk := LoadDeliveryKeeperHelper(filePath)

	filtered := dk.CountByInterval("10122", "10AM", "2PM")

	assert.Equal(t, 2919, filtered)
}

func TestCountCountByIntervalSingle(t *testing.T) {
	filePath := "./testdata/test_calculation_fixtures_double.json"
	dk := LoadDeliveryKeeperHelper(filePath)

	count := dk.CountByInterval("10145", "8AM", "3PM")

	assert.Equal(t, 1, count)
}

func TestCountByIntervalUpperLimit(t *testing.T) {
	filePath := "./testdata/test_calculation_fixtures_upper_limit_delivery.json"
	dk := LoadDeliveryKeeperHelper(filePath)

	count := dk.CountByInterval("10145", "5PM", "11PM")

	assert.Equal(t, 1, count)
}

func TestCountByIntervalLowerLimit(t *testing.T) {
	filePath := "./testdata/test_calculation_fixtures_lower_limit_delivery.json"
	dk := LoadDeliveryKeeperHelper(filePath)

	filtered := dk.CountByInterval("10174", "12AM", "8AM")

	assert.Equal(t, 1, filtered)
}

func TestBusiestPostcode(t *testing.T) {
	filePath := "./testdata/test_calculation_fixtures_busiest_postalcode.json"
	dk := LoadDeliveryKeeperHelper(filePath)

	result := dk.BusiestPostcode

	assert.NotNil(t, result)
	assert.Equal(t, "10129", result.Code)
	assert.Equal(t, 6, result.Count)
}
