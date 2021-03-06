package keepers

import (
	"recipe-stats/models"
	"strconv"
)

// postcode is the struct that internally holds all the deliveries within a time
// range for a given postcode. It also keeps in account the count of deliveries
// and the and the postcode number for the sake of searching later.
type postcode struct {
	Code string
	// I am using an slice initialized with  all the possible time ranges (12AM to
	// 12PM) so I don't need to sort later and I can also access time ranges
	// instantly
	Deliveries      [24][24][]models.Delivery
	DeliveriesCount int
}

// BusiestPostCode is the struct that serves as the base for the accounting of
// the busiest post code.
type BusiestPostcode struct {
	Code  string
	Count int
}

// DeliveryKeeper is the main struct for the Delivery Keeper and holds reference
// to all the required objects for the logic of controlling deliveries and the
// busiest postcode.
type DeliveryKeeper struct {
	postcodes       map[string]*postcode
	BusiestPostcode BusiestPostcode
}

// NewDeliveryKeeper provides a usable instance of DeliveryKeeper.
func NewDeliveryKeeper() DeliveryKeeper {
	return DeliveryKeeper{
		postcodes:       map[string]*postcode{},
		BusiestPostcode: BusiestPostcode{},
	}
}

// Add puts a new delivery on the list of deliveries taking its time range and
// postcode into account. It also updates the busiest postcode if applicable.
func (dk *DeliveryKeeper) Add(delivery models.Delivery) {
	// Mapping postcodes
	foundPostcode, found := dk.postcodes[delivery.Postcode]
	if found {
		foundPostcode.DeliveriesCount++
	} else {
		foundPostcode = &postcode{
			Code:            delivery.Postcode,
			DeliveriesCount: 1,
		}
		dk.postcodes[delivery.Postcode] = foundPostcode
	}

	// Mapping deliveries within time range
	foundPostcode.Deliveries[delivery.From][delivery.To] =
		append(foundPostcode.Deliveries[delivery.From][delivery.To],
			delivery)

	if dk.BusiestPostcode.Count < foundPostcode.DeliveriesCount {
		dk.BusiestPostcode.Code = foundPostcode.Code
		dk.BusiestPostcode.Count = foundPostcode.DeliveriesCount
	}
}

// CountByInterval takes a postcode and an start and end times in 12h format,
// finds and counts all the deliveries for that postcode within the time range.
// If nothing was found or the parameters are empty, it returns 0.
func (dk *DeliveryKeeper) CountByInterval(postcode string, start string, end string) int {
	if postcode == "" || start == "" || end == "" {
		return 0
	}
	rangeBottom := TimeToIndex(start)
	rangeTop := TimeToIndex(end) + 1

	foundPostcode, found := dk.postcodes[postcode]
	if !found {
		return 0
	}

	// start filtering by start time
	fromStartTimeDeliveries := foundPostcode.Deliveries[rangeBottom:rangeTop]

	var count int
	for _, endTimes := range fromStartTimeDeliveries {
		for _, deliveries := range endTimes[rangeBottom:rangeTop] {
			if len(deliveries) > 0 {
				count += len(deliveries)
			}
		}
	}

	return count
}

// TimeToIndex is a transform function that transforms a 12h time string in a
// 24h time number. Example: 5PM turns into 17, 12AM turns into 0.
// Byte manipulation was chosen in favor of performance.
func TimeToIndex(time string) int {
	byteDr := []byte(time)
	suffix := string(byteDr[len(byteDr)-2:])
	prefix := string(byteDr[0 : len(byteDr)-2])
	value, _ := strconv.Atoi(prefix)
	if suffix == "PM" {
		if value == 12 {
			return value // 12PM
		} else {
			return value + 12
		}
	} else {
		if value == 12 {
			return 0 // 12AM
		} else {
			return value
		}
	}
}
