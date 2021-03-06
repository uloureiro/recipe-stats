package loaders

import (
	"fmt"
	"os"
	"recipe-stats/adapters"
	"recipe-stats/keepers"
	"sync"
	"time"
)

func loadDeliveriesFromGeneralRecipe(recipes []adapters.GeneralRecipe, wg *sync.WaitGroup, verbose bool) *keepers.DeliveryKeeper {
	defer wg.Done()

	start := time.Now()
	if verbose {
		fmt.Fprintln(os.Stderr, "Mapping deliveries...")
	}

	deliveryKeeper := keepers.NewDeliveryKeeper()

	for i := 0; i < len(recipes); i++ {
		delivery := recipes[i].ToDelivery()
		deliveryKeeper.Add(delivery)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Mapping deliveries took %s\n", time.Since(start))
	}

	return &deliveryKeeper
}
