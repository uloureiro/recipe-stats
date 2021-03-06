package cmd

import (
	"fmt"
	"os"
	"recipe-stats/keepers"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/AlecAivazis/survey/v2"
)

var (
	builtInFilePath        string
	customFilePath         string
	reuseDataset           bool
	recipeKeeper           *keepers.RecipeKeeper
	recipeNameSlicesKeeper *keepers.RecipeNameSlicesKeeper
	deliveryKeeper         *keepers.DeliveryKeeper
	keepersError           error
)

// interactiveFlow is the entrypoint for the interactive execution. It prints a logo
// along with a basic help message. The steps that follows the interactive flow
// are pretty self explanatory.
func interactiveFlow(filePathFromConfig string) {
	builtInFilePath = filePathFromConfig

	fmt.Println(`	
    ___          _          ______       __    
   / _ \___ ____(_)__  ___ / __/ /____ _/ /____
  / , _/ -_) __/ / _ \/ -_)\ \/ __/ _ )/ __(_-<
 /_/|_|\__/\__/_/  __/\__/___/\__/\_,_/\__/___/
               /_/`)

	fmt.Println("\nPlease follow the instructions bellow to setup your recipe stats query.")
	fmt.Println()

	runFlow()
}

// runFlow is the actual place where the interactive flow is controlled.
// It contains all the questions and rules for branching the flow to build-up
// the parameters for the calculation.
func runFlow() {
	wg := *new(sync.WaitGroup)

	var (
		fileOption     string
		customFile     bool
		defaultFile    bool
		filePath       string
		options        []string
		recipesNames   string
		postcode       string
		from           string
		to             string
		askRecipeNames bool
		askPostcode    bool
		recipeCount    bool
		runAgain       bool
	)

	if !reuseDataset {
		// force GC to free up memory to load large chunks again
		debug.FreeOSMemory()

		_ = survey.AskOne(fileOptionQuestion, &fileOption, survey.WithValidator(survey.Required))

		if fileOption == "Inform custom file path" {
			customFile = true
		} else if fileOption == "Use built-in file" {
			defaultFile = true
		}

		if defaultFile {
			filePath = builtInFilePath
		} else if customFile {
			_ = survey.AskOne(customFileQuestion, &customFilePath, survey.WithValidator(survey.Required))
			filePath = customFilePath
		}

		// loading everything while options are selected
		wg.Add(1)
		go func() {
			defer wg.Done()
			recipeKeeper, recipeNameSlicesKeeper, deliveryKeeper, keepersError = loadKeepers(filePath, false)
		}()
	}

	_ = survey.AskOne(optionsQuestion, &options, survey.WithValidator(survey.Required))

	for _, option := range options {
		if option == "Count unique recipes" {
			recipeCount = true
		}
		if option == "Search by recipe name" {
			askRecipeNames = true
		}
		if option == "Search by postcode and time window" {
			askPostcode = true
		}
	}

	if askRecipeNames {
		_ = survey.AskOne(recipeNameQuestion, &recipesNames, survey.WithValidator(survey.Required))
	}

	if askPostcode {
		_ = survey.AskOne(postcodeQuestion, &postcode, survey.WithValidator(survey.Required))
		_ = survey.AskOne(fromTimeQuestion, &from, survey.WithValidator(survey.Required))
		_ = survey.AskOne(toTimeQuestion, &to, survey.WithValidator(survey.Required))
	}

	wg.Wait()

	if keepersError != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong, please try again. Error:%s\n", keepersError.Error())

		var tryAgain bool
		_ = survey.AskOne(tryAgainQuestion, &tryAgain)
		if !tryAgain {
			return
		}
		runFlow()
	}

	runFromInteractive(recipeKeeper, recipeNameSlicesKeeper, deliveryKeeper, recipeCount, strings.Split(recipesNames, ","), postcode, from, to)

	_ = survey.AskOne(runAgainQuestion, &runAgain)
	if runAgain {
		_ = survey.AskOne(reuseDatasetQuestion, &reuseDataset)
		runFlow()
	} else {
		return
	}
}

var fileOptionQuestion = &survey.Select{
	Message: "Which file do you want to use?",
	Options: []string{
		"Inform custom file path",
		"Use built-in file",
	},
}

var customFileQuestion = &survey.Input{
	Message: "Inform the custom file path to use:",
}

var optionsQuestion = &survey.MultiSelect{
	Message: "Select which options you want to use in order to check recipes stats:",
	Options: []string{
		"Count unique recipes",
		"Search by recipe name",
		"Search by postcode and time window",
	},
}

var recipeNameQuestion = &survey.Input{
	Message: "Inform a comma separated list of recipes name to search:",
}

var postcodeQuestion = &survey.Input{
	Message: "Inform the desired postcode to search:",
}

var fromTimeQuestion = &survey.Input{
	Message: "From what time to start searching? (Example: 9AM)",
}

var toTimeQuestion = &survey.Input{
	Message: "Up to what time to end searching? (Example: 2PM)",
}

var tryAgainQuestion = &survey.Confirm{
	Message: "Do you want to try again?",
	Default: false,
}

var runAgainQuestion = &survey.Confirm{
	Message: "Do you want to run again?",
	Default: false,
}

var reuseDatasetQuestion = &survey.Confirm{
	Message: "Use the same dataset?",
	Default: true,
}
