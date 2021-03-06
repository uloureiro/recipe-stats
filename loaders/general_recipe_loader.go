package loaders

import (
	"fmt"
	"os"
	"recipe-stats/adapters"
	"recipe-stats/keepers"
	"sync"
	"time"
)

func LoadFromGeneralRecipe(filePath string, verbose bool) (*keepers.RecipeKeeper, *keepers.RecipeNameSlicesKeeper, *keepers.DeliveryKeeper, error) {
	wg := *new(sync.WaitGroup)

	recipes, err := loadGeneralRecipesFile(filePath, verbose)
	if err != nil {
		return nil, nil, nil, err
	}

	recipeKeeper := new(keepers.RecipeKeeper)
	recipeNameSlicesKeeper := new(keepers.RecipeNameSlicesKeeper)
	wg.Add(2)
	go func() {
		recipeKeeper, _ = loadRecipesFromGeneralRecipe(*recipes, &wg, verbose)
		recipeNameSlicesKeeper = loadRecipeNameSlicesFromRecipes(recipeKeeper.GetMap(), &wg, verbose)
	}()

	deliveryKeeper := new(keepers.DeliveryKeeper)
	wg.Add(1)
	go func() {
		deliveryKeeper = loadDeliveriesFromGeneralRecipe(*recipes, &wg, verbose)
	}()

	wg.Wait()
	return recipeKeeper, recipeNameSlicesKeeper, deliveryKeeper, nil
}

func loadGeneralRecipesFile(filePath string, verbose bool) (*[]adapters.GeneralRecipe, error) {
	if verbose {
		fmt.Fprintln(os.Stderr, "Reading recipes file...")
	}
	start := time.Now()
	recipesAdapter := adapters.NewGeneralRecipeAdapter()
	recipes, err := recipesAdapter.Unmarshal(filePath)
	if err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "It was impossible to parse the input file. The error was: %s\n", err.Error())
		}
		return nil, err
	}
	if verbose {
		fmt.Fprintf(os.Stderr, "Reading recipes file took %s\n", time.Since(start))
	}

	return recipes, nil
}
