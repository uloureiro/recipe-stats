package main

import "recipe-stats/cmd"

/*
Application breakdown
- cmd
	Contains the files related to the CLI.
	- root.go
		It is the main entrypoint of the application and it routes to
		interactiveFlow() or runFromCli() methods.
		Contains the CLI flags configurations, as well as the loading of config.yml.
	- interactive.go
		Contains all the logic to handle the interactive mode.
	- runner.go
		Is the common point of contact from root.go and interactive.go and handles
		the necessary execution steps in the correct order. This is where the pieces
		get bound together, producing the final calculation.
- adapters
	Contains the Adapter and AdapterMember interfaces that aims to provide a
	baseline for input types (aiming on scaling support for input types).
	- general_recipe_adapter.go
		Is the adapter for the fixtures file provided on the requirements. It contains
		all the specific methods to transform the input file into collections of
		Recipes and Deliveries.
- keepers
	Keepers contains the files that holds the collections of pre-processed data of
	Recipes, Deliveries and Recipes Names Slices. Those files  also contains the
	methods can calculate the desired outputs.
	- delivery_keeper.go
		Holds the collections of deliveries and the methods to add deliveries,
		calculate the busiest postcode and searching for deliveries intervals.
	- recipe_keeper.go
		Holds the collection of distinct recipes and the methods to add recipes,
		search by name and count the total of recipes.
	- recipe_names_slices_keeper.go
		Holds the collection of all the possible slices of recipes names and the
		methods to add slices and search recipes by slice.
- loaders
	Building a loaders structure had a sole purpose of helping on managing the
	parallelization of loading tasks to make it perform better.
	general_recipe_loader.go is the central loader for GeneralRecipe and it manages the
	loading order and parallelization of tasks. In the end, it provides instances
	of the required keepers so the runner can execute the calculations.
- models
	Models are the basic common types where the data used throughout the
	application relies on. Every input data gets transformed into one of the
	models so the keepers don't need to change if new input formats should be
	supported.
- reporters
	Reporters contains the structure and encoding methods to generate an output in
	a desired format.
	- json_reporter.go
		Has the structs and encoding methods to create an output within the
		contraints from the problem statement.
*/

func main() {
	cmd.Execute()
}
