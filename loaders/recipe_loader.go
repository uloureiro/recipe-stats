package loaders

import (
	"fmt"
	"os"
	"recipe-stats/adapters"
	"recipe-stats/keepers"
	"sync"
	"time"
)

func loadRecipesFromGeneralRecipe(recipes []adapters.GeneralRecipe, wg *sync.WaitGroup, verbose bool) (*keepers.RecipeKeeper, error) {
	defer wg.Done()

	start := time.Now()
	if verbose {
		fmt.Fprintln(os.Stderr, "Loading recipes...")
	}

	recipeKeeper := keepers.NewRecipeKeeper()

	for i := 0; i < len(recipes); i++ {
		recipe := recipes[i].ToRecipe()
		err := recipeKeeper.Add(recipe)

		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "It was impossible to load data into the calculator. The error was: %s\n", err.Error())
			}
			return nil, err
		}
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Loading recipes took %s\n", time.Since(start))
	}

	return &recipeKeeper, nil
}
