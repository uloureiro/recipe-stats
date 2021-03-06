package reporters

import "encoding/json"

// CountPerRecipe is the building block of the counting per recipe output
type CountPerRecipe struct {
	Recipe string `json:"recipe,omitempty"`
	Count  int    `json:"count,omitempty"`
}

// BusiestPostCode is the building block of the busiest postcode output
type BusiestPostCode struct {
	Postcode      string `json:"postcode,omitempty"`
	DeliveryCount int    `json:"delivery_count,omitempty"`
}

// CountPerPostcodeAndTime is the building block of the count per postcode
//output
type CountPerPostcodeAndTime struct {
	Postcode      string `json:"postcode"`
	From          string `json:"from"`
	To            string `json:"to"`
	DeliveryCount int    `json:"delivery_count"`
}

// JSONReporter is the main struct for this reporter, holding all the building
// blocks for this type of output.
type JSONReporter struct {
	UniqueRecipeCount       int                      `json:"unique_recipe_count,omitempty"`
	CountPerRecipe          []CountPerRecipe         `json:"count_per_recipe,omitempty"`
	BusiestPostCode         *BusiestPostCode         `json:"busiest_postcode,omitempty"`
	CountPerPostcodeAndTime *CountPerPostcodeAndTime `json:"count_per_postcode_and_time,omitempty"`
	MatchByName             []string                 `json:"match_by_name,omitempty"`
}

// Marshal is the encodinf function for JSONReporter and creates a formatted
// JSON string.
func (jr *JSONReporter) Marshal() (string, error) {
	marshaled, err := json.MarshalIndent(jr, "", "  ")

	if err != nil {
		return "", err
	}

	return string(marshaled), nil
}
