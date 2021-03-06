package keepers

import (
	"recipe-stats/models"
)

// RecipeKeeper is the main struct for the Recipe Keeper, holding the necessary
// data structures to calculate distinct recipes count
type RecipeKeeper struct {
	recipes map[string]models.Recipe
}

// NewRecipeKeeper provides a usable instance of RecipeKeeper.
func NewRecipeKeeper() RecipeKeeper {
	rk := RecipeKeeper{}
	rk.recipes = make(map[string]models.Recipe)

	return rk
}

// Add puts a new recipe on the list of recipes if it does not exists yet.
func (rk *RecipeKeeper) Add(recipe models.Recipe) error {
	if existingRecipe, exists := rk.recipes[recipe.Recipe]; exists {
		existingRecipe.Count += 1
		rk.recipes[recipe.Recipe] = existingRecipe
	} else {
		recipe.Count = 1
		rk.recipes[recipe.Recipe] = recipe
	}

	return nil
}

// Count calculates the number of distinct recipes found.
func (rk *RecipeKeeper) Count() int {
	if rk.recipes == nil {
		return 0
	}

	return len(rk.recipes)
}

// GetMap is a utility to return a copy of the recipes map
func (rk *RecipeKeeper) GetMap() map[string]models.Recipe {
	return rk.recipes
}
