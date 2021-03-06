package tests

import (
	"recipe-stats/keepers"
	"recipe-stats/loaders"
)

func LoadRecipeKeeperHelper(filePath string) keepers.RecipeKeeper {
	rk, _, _, err := loaders.LoadFromGeneralRecipe(filePath, false)

	if err != nil {
		panic(err)
	}

	return *rk
}

func LoadRecipeNameSliceKeeperHelper(filePath string) *keepers.RecipeNameSlicesKeeper {
	_, rnsk, _, err := loaders.LoadFromGeneralRecipe(filePath, false)

	if err != nil {
		panic(err)
	}

	return rnsk
}

func LoadDeliveryKeeperHelper(filePath string) keepers.DeliveryKeeper {
	_, _, dk, err := loaders.LoadFromGeneralRecipe(filePath, false)

	if err != nil {
		panic(err)
	}

	return *dk
}
