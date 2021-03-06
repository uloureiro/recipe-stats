package loaders

import (
	"fmt"
	"os"
	"recipe-stats/keepers"
	"recipe-stats/models"
	"sync"
	"time"
)

func loadRecipeNameSlicesFromRecipes(recipes map[string]models.Recipe, wg *sync.WaitGroup, verbose bool) *keepers.RecipeNameSlicesKeeper {
	defer wg.Done()

	start := time.Now()
	if verbose {
		fmt.Fprintln(os.Stderr, "Building recipes map...")
	}

	rnsk := keepers.NewRecipeNameSlicesKeeper()

	rnsk.Load(recipes)

	if verbose {
		fmt.Fprintf(os.Stderr, "Building recipes map took %s\n", time.Since(start))
	}

	return &rnsk
}
