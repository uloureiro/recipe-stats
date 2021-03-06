package adapters

import (
	"io/ioutil"
	"recipe-stats/models"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

// GeneralRecipeAdapter is the base structure for this adapter
type GeneralRecipeAdapter struct {
	rawRecipes *[]GeneralRecipe
}

// GeneralRecipe is the struct that maps to the input JSON
type GeneralRecipe struct {
	Recipe   string          `json:"recipe"`
	Postcode string          `json:"postcode"`
	Delivery GeneralDelivery `json:"delivery"`
}

// GeneralDelivery is the struct that maps to the delivery data from the input JSON
type GeneralDelivery struct {
	From int
	To   int
}

// jsoniter is an optimized library to encode/decode JSON
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// NewRecipeAdapter provides a usable instance of GeneralRecipeAdapter
func NewGeneralRecipeAdapter() GeneralRecipeAdapter {
	return GeneralRecipeAdapter{}
}

// unwrap does the Unmarshal of the file data into a collection of GeneralRecipe
func (a *GeneralRecipeAdapter) unwrap(file []byte) (*[]GeneralRecipe, error) {
	rawRecipes := new([]GeneralRecipe)

	err := json.Unmarshal(file, &rawRecipes)
	if err != nil {
		return nil, err
	}

	return rawRecipes, nil
}

// Unmarshal gets a file path as input argument, reads the file and call the
// unwrap method to transform that data into the adapter's base struct
func (a *GeneralRecipeAdapter) Unmarshal(filePath string) (*[]GeneralRecipe, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	unwrappedRecipes, err := a.unwrap(file)
	if err != nil {
		return nil, err
	}

	a.rawRecipes = unwrappedRecipes

	return unwrappedRecipes, nil
}

// ToRecipe provides the transforming logic from GeneralRecipe to models.Recipe
func (r *GeneralRecipe) ToRecipe() models.Recipe {
	return models.Recipe{
		Recipe: r.Recipe,
	}
}

// ToDelivery provides the transforming logic from GeneralRecipe to models.Delivery
func (r *GeneralRecipe) ToDelivery() models.Delivery {
	return models.Delivery{
		Postcode: r.Postcode,
		From:     r.Delivery.From,
		To:       r.Delivery.To,
	}
}

// UnmarshalJSON is the custom Unmarshaler that runs whe decoding GeneralDelivery
func (d *GeneralDelivery) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	d.From, d.To = d.GetDeliveryTimes(v.(string))

	return nil
}

// weekdays is a fixed collection of weekdays and their characters lengths
var weekdays = map[string]int{
	"Sa": 8,
	"Su": 6,
	"Mo": 6,
	"Tu": 7,
	"We": 9,
	"Th": 8,
	"Fr": 6,
}

// GetDeliveryTimes is a transformation method that gets a full delivery string
// from GeneralDelivery, such as "Wednesday 8AM - 2PM" and transforms it into two 24h
// hour numbers that represents the respectives from and to times. In the
// example given, the expected result would be 8 and 14.
// For "12AM" it transforms to 0, for "12PM" it transforms to 12.
// Byte manipulation was used in order to achieve performance since Regex
// demonstrated to be very slow.
func (d *GeneralDelivery) GetDeliveryTimes(value string) (int, int) {
	byteR := []byte(value)

	padding := weekdays[string(byteR[0:2])] + 1 // + 1 to eliminate trailing space
	timePortion := byteR[padding:]

	// start time
	var startNumber int
	var startIndicator string
	// the second position is a number, the hour has 2 digits
	if _, err := strconv.Atoi(string(timePortion[1])); err == nil {
		startNumber, _ = strconv.Atoi(string(timePortion[:2]))
		startIndicator = string(timePortion[3:5])
	} else { // the hour has 1 digit
		startNumber, _ = strconv.Atoi(string(timePortion[:1]))
		startIndicator = string(timePortion[1:3])
	}

	if startIndicator == "PM" { // it is noon
		startNumber += 12
	} else { // it is morning
		if startNumber == 12 { // 12 AM
			startNumber = 0
		}
	}

	//end time
	var endNumber int
	var endIndicator string
	// the first position is not a number, the hour has 1 digit
	if _, err := strconv.Atoi(string(timePortion[len(timePortion)-4:][0])); err != nil {
		endNumber, _ = strconv.Atoi(string(timePortion[len(timePortion)-3 : len(timePortion)-2]))
		endIndicator = string(timePortion[len(timePortion)-2:])
	} else { // the hour has 2 digit
		endNumber, _ = strconv.Atoi(string(timePortion[len(timePortion)-4 : len(timePortion)-2]))
		endIndicator = string(timePortion[len(timePortion)-2:])
	}

	if endIndicator == "PM" { // it is noon
		endNumber += 12
	} else { // it is morning
		if endNumber == 12 { // 12 AM
			endNumber = 0
		}
	}

	return startNumber, endNumber
}
