package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string = "config.yml"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "recipe-stats",
	Short: "Calculates different attributes about an informed collection of recipes.",
	Long: `This program can calculate the following attributes over a collection of recipes from an arbitrary file passed as argument:

- Unique recipes count
- Counting per recipe found (when searching by recipes partial names)
- Busiest postcode
- Deliveries count for searched postcode and time intervals
- Recipes found by partial recipe name

Use the flags described bellow to achieve those results.

Example: recipe-stats -s=Cheese,Grilled -f=path_to_file -c
`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.PersistentFlags().GetString("file")
		recipeCount, _ := cmd.PersistentFlags().GetBool("count")
		namesToSearch, _ := cmd.PersistentFlags().GetStringSlice("search")
		postcodeToSearch, _ := cmd.PersistentFlags().GetString("postcode")
		from, _ := cmd.PersistentFlags().GetString("from")
		to, _ := cmd.PersistentFlags().GetString("to")
		verbose, _ := cmd.PersistentFlags().GetBool("verbose")
		interactive, _ := cmd.PersistentFlags().GetBool("interactive")

		if interactive {
			interactiveFlow(filePath)
		} else {
			runFromCli(filePath, recipeCount, namesToSearch, postcodeToSearch, from, to, verbose)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	initConfig()

	rootCmd.PersistentFlags().StringP("file", "f", viper.GetString("file_path"), "The full path of a different input file to analyze")
	_ = viper.BindPFlag("file_path", rootCmd.PersistentFlags().Lookup("file"))
	rootCmd.PersistentFlags().BoolP("count", "c", false, "Counts the number of unique recipes")
	rootCmd.PersistentFlags().StringSliceP("search", "s", nil, "Comma separated list of recipe names to find")
	rootCmd.PersistentFlags().StringP("postcode", "p", "", "Postcode number to lookup. Using that flag will require you to inform the --from and --to flags")
	rootCmd.PersistentFlags().String("from", "", "The starting time for postcode deliveries search. Example: 11AM")
	rootCmd.PersistentFlags().String("to", "", "The ending time for postcode deliveries search. Example: 2PM")
	rootCmd.PersistentFlags().BoolP("interactive", "i", false, "Runs the program in interactive mode. Any other flag will be ignored.")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Prints profiling and performance messages")

	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".recipe-stats" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".recipe-stats")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file.")
	}
}
