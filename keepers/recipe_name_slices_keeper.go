package keepers

import (
	"recipe-stats/models"
	"sort"
	"strings"
)

// I first thought of using a balanced tree to store everything from the
// begining because it provides elements ordering out of the box, but I came
// accross a few pitfalls:
// - Searching for a name slice in a tree is slower than in a map
// - Even though searching for each slice can be pretty fast either way, I would
//   have to reorder again the recipes found after joining the results to
//   return (recipes are already ordered in the tree but the search by slice
//   can't guarantee an ordered search/result)
// - It means that I would have to iterate twice over the recipes found:
//   1st when feeding the tree
//   2nd when reading again from the tree
// - I decided then to use Sort since its worst case will be O(n^2) but in
//   in general tends to be O(n log n)
// - I have to make a hard decision here: using just maps and slices may turn
//   dealing with repeated results into something very slow... BUT the usual
//   scenario is that just a tiny number of resulting recipes is found
//   (in fact, the number of unique recipes is very low) thus I will assume the
//   risk of adding one more iteration on the search having in mind that it will
//   not impact performance at all

// RecipeNameSlicesKeeper is the main struct for the Recipes Name Slice
// reference to all the required objects for the logic of controlling deliveries and the
// busiest postcode.
type RecipeNameSlicesKeeper struct {
	recipeNameSlices map[string][]models.Recipe
}

func NewRecipeNameSlicesKeeper() RecipeNameSlicesKeeper {
	rnsk := RecipeNameSlicesKeeper{}
	rnsk.recipeNameSlices = make(map[string][]models.Recipe)

	return rnsk
}

// Load receives a map of Recipes (the key is the recipe name) and breaks it
// into name slices, distributing the words found within the names map along
// the recipe itself
func (rnsk *RecipeNameSlicesKeeper) Load(recipes map[string]models.Recipe) {
	for _, recipe := range recipes {
		recipeNameSlices := strings.Split(recipe.Recipe, " ")

		for _, recipeNameSlice := range recipeNameSlices {
			rnsk.recipeNameSlices[recipeNameSlice] = append(rnsk.recipeNameSlices[recipeNameSlice], recipe)
		}
	}
}

// Get finds a single name slice and returns the recipes related to it. Also
// returns true if find something and false if not.
func (rnsk *RecipeNameSlicesKeeper) Get(recipeNameSlice string) ([]models.Recipe, bool) {
	recipes, found := rnsk.recipeNameSlices[recipeNameSlice]

	return recipes, found
}

// GetSome finds multiple name slices and returns the recipes related to them.
func (rnsk *RecipeNameSlicesKeeper) GetSome(recipeNameSlices []string) []models.Recipe {
	recipesFound := []models.Recipe{}
	for _, slice := range recipeNameSlices {
		if recipes, found := rnsk.Get(slice); found {
			recipesFound = append(recipesFound, recipes...)
		}
	}

	// this sort has worst case O(n^2) but usually O(n log n)
	// https://golang.org/pkg/sort/#Sort
	sort.Sort(ByRecipe(recipesFound))

	// given it is sorted, it is easier to remove duplicates
	singleRecipesFound := []models.Recipe{}
	var previousRecipe *models.Recipe
	for _, recipe := range recipesFound {
		if previousRecipe == nil || recipe.Recipe != previousRecipe.Recipe {
			cachedRecipe := recipe
			previousRecipe = &cachedRecipe
			singleRecipesFound = append(singleRecipesFound, recipe)
		}
	}

	return singleRecipesFound
}

type ByRecipe []models.Recipe

func (a ByRecipe) Len() int           { return len(a) }
func (a ByRecipe) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRecipe) Less(i, j int) bool { return a[i].Recipe < a[j].Recipe }
