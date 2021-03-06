RecipeStats
====
RecipeStats is a CLI calculator for recipes statuses.

It can read an arbitrary JSON file passed as argument and provide some insights about its recipes (although a sample dataset is provided so you can get started right now!).

The input file should be a JSON with objects respecting the following sample structure:

```json5
{
  "postcode": "10145",
  "recipe": "Parmesan-Crusted Pork Tenderloin",
  "delivery": "Wednesday 9AM - 2PM"
}
```

The results are provided as a JSON like this:
```json5
{
    "unique_recipe_count": 15,
    "count_per_recipe": [
        {
            "recipe": "Mediterranean Baked Veggies",
            "count": 1
        },
        {
            "recipe": "Speedy Steak Fajitas",
            "count": 1
        },
        {
            "recipe": "Tex-Mex Tilapia",
            "count": 3
        }
    ],
    "busiest_postcode": {
        "postcode": "10120",
        "delivery_count": 1000
    },
    "count_per_postcode_and_time": {
        "postcode": "10120",
        "from": "11AM",
        "to": "3PM",
        "delivery_count": 500
    },
    "match_by_name": [
        "Mediterranean Baked Veggies", "Speedy Steak Fajitas", "Tex-Mex Tilapia"
    ]
}
```

Setup
---------
Clone this repository to your machine.

You can setup this application both locally or using Docker.

Please be sure to have `make` installed on your machine. If `make` can't be used, please take a look inside `Makefile` to run the commands on your own.

### Setup using Docker

```sh
make setup-docker
```

This will build the image and start the container for you.

Please notice that Docker maps a volume to the `data` folder within the working directory of the application. Any custom file that you want to use should be referenced using that folder. You can find further instructions in the [usage](#usage) section.

### Setup on you local machine

Please be sure to have [Go](https://golang.org/) installed before proceeding. This application was tested only in Mac OS and Linux. If you can't use Mac OS or Linux, please go for Docker.

```sh
make setup-local
```
This will install all the dependencies and build the binary for you.

Usage
---------

You can run this application in interactive or flag mode. 

**Note about using Docker and custom input files**

If you are using Docker and want to use custom files, it is required that you put your custom file into the same folder where the `docker-compose.yml` is located (A.K.A the root of this project). This is due the mapping to the local directory to the `data` folder within Docker so the application can have a channel to communicate to the outside world of the container, accessing custom files.

With that, anytime you want to reference a custom file, you will have to add an extra `data/` to it so the application can lookup within the data folder.

Example:

I put the file `my_custom_data.json` into the root of the project and then I run the application referencing it.

```
recipe-stats -f data/my_custom_data.json
```

The same applies for interactive mode.

### Interactive mode

![demo](https://media.giphy.com/media/icIHIkyLPFkwd5uFHs/giphy.gif)

Using docker, run this command:

```sh
make run-interactive
```

Using local, run this command:

```sh
./recipe-stats -i
```

### Flag mode

If you are using **Docker**, an alias may come in handy in order to make the commands cleaner. You can source an alias for your shell using the script provided.

```sh
source scripts/alias.sh
```

This will live along your current shell session and you can run the application as:

```sh
recipe-stats [flags]
```

If you don't want to use the alias, you should run the application as:

```sh
docker-compose exec recipe-stats recipe-stats [flags]
```

To run it locally, simply execute the binary:
```sh
./recipe-stats [flags]
```

Those are the flags accepted (you can get this message using the `-h` flag).

```sh
Usage:
  recipe-stats [flags]

Flags:
  -f, --file string       The full path of a different input file to analyze (default "sample_data.json")
  -c, --count             Counts the number of unique recipes
  -s, --search strings    Comma separated list of recipe names to find
  -p, --postcode string   Postcode number to lookup. Using that flag will require you to inform the --from and --to flags
      --from string       The starting time for postcode deliveries search. Example: 11AM
      --to string         The ending time for postcode deliveries search. Example: 2PM
  -i, --interactive       Runs the program in interactive mode. Any other flag will be ignored.
  -v, --verbose           Prints profiling and performance messages
  -h, --help              help for recipe-stats
```

Example:

```sh
recipe-stats -f data/my_custom_file.json -c -s Pasta,Cheese -p 10122 --from 9AM --to 2PM
```