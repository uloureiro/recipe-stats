package cmd

import (
	"fmt"
	"os"
	"recipe-stats/keepers"
	"recipe-stats/loaders"
	"recipe-stats/reporters"
	"time"
)

//  Runner contains all the methods focused on triggering the data load from
//  Loaders, getting the response back, and doing the dirty work of retrieving the
//  results from the Keepers and agglutinating them into the Reporter.
//  Even though the Reporter is the one who knows how to build up the dataset into
//  the required JSON format, Runner is the place where the result is printed.

// runFromCli is the entrypoint for the CLI execution. It is called when the flag
// `--interactive` is not set
func runFromCli(filePath string, recipeCount bool, namesToSearch []string, postcodeToSearch string, from string, to string, verbose bool) {
	totalStart := time.Now()

	recipeKeeper, recipeNameSlicesKeeper, deliveryKeeper, err := loadKeepers(filePath, verbose)
	if err != nil {
		jsonOutput := reporters.JSONReporter{}
		formattedOutput, _ := jsonOutput.Marshal()
		fmt.Fprintln(os.Stdout, formattedOutput)
		return
	}

	calculate(recipeKeeper, recipeNameSlicesKeeper, deliveryKeeper, recipeCount, namesToSearch, postcodeToSearch, from, to, verbose)

	if verbose {
		fmt.Fprintf(os.Stderr, "Total execution took %s\n", time.Since(totalStart))
	}
}

// runFromInteractive is the entrypoint for the interactive
func runFromInteractive(recipeKeeper *keepers.RecipeKeeper, recipeNameSlicesKeeper *keepers.RecipeNameSlicesKeeper, deliveryKeeper *keepers.DeliveryKeeper, recipeCount bool, namesToSearch []string, postcodeToSearch string, from string, to string) {
	calculate(recipeKeeper, recipeNameSlicesKeeper, deliveryKeeper, recipeCount, namesToSearch, postcodeToSearch, from, to, false)
}

func loadKeepers(filePath string, verbose bool) (*keepers.RecipeKeeper, *keepers.RecipeNameSlicesKeeper, *keepers.DeliveryKeeper, error) {
	recipeKeeper, recipeNameSlicesKeeper, deliveryKeeper, err := loaders.LoadFromGeneralRecipe(filePath, verbose)
	if err != nil {
		return nil, nil, nil, err
	}

	return recipeKeeper, recipeNameSlicesKeeper, deliveryKeeper, nil
}

// It supports the verbose option
func calculate(recipeKeeper *keepers.RecipeKeeper, recipeNameSlicesKeeper *keepers.RecipeNameSlicesKeeper, deliveryKeeper *keepers.DeliveryKeeper, recipeCount bool, namesToSearch []string, postcodeToSearch string, from string, to string, verbose bool) {
	start := time.Now()
	jsonOutput := reporters.JSONReporter{}

	if verbose {
		fmt.Fprintln(os.Stderr, "Calculating...")
	}

	if recipeCount {
		jsonOutput.UniqueRecipeCount = recipeKeeper.Count()
	}

	recipesFound := recipeNameSlicesKeeper.GetSome(namesToSearch)
	recipesFoundNames := []string{}
	recipesFoundCounts := []reporters.CountPerRecipe{}
	for _, recipe := range recipesFound {
		recipesFoundNames = append(recipesFoundNames, recipe.Recipe)
		recipesFoundCounts = append(recipesFoundCounts, reporters.CountPerRecipe{Recipe: recipe.Recipe, Count: recipe.Count})
	}
	jsonOutput.MatchByName = recipesFoundNames
	jsonOutput.CountPerRecipe = recipesFoundCounts

	jsonOutput.BusiestPostCode = &reporters.BusiestPostCode{
		Postcode:      deliveryKeeper.BusiestPostcode.Code,
		DeliveryCount: deliveryKeeper.BusiestPostcode.Count,
	}

	countByPostcode := deliveryKeeper.CountByInterval(postcodeToSearch, from, to)
	if countByPostcode > 0 {
		jsonOutput.CountPerPostcodeAndTime = &reporters.CountPerPostcodeAndTime{
			From:          from,
			To:            to,
			Postcode:      postcodeToSearch,
			DeliveryCount: countByPostcode,
		}
	}

	formattedOutput, err := jsonOutput.Marshal()

	if err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "It was impossible to format the output. The error was:%s\n", err.Error())
		}
	}

	fmt.Fprintln(os.Stdout, formattedOutput)

	if verbose {
		fmt.Fprintf(os.Stderr, "Calculating took %s\n", time.Since(start))
	}
}
