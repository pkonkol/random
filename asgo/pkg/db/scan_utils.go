package db

import "fmt"

func IsPrefixScanned(prefix string) bool {
	// Consult DB to check that
	return false
}

func GetNextPrefix(country string, method string) string {
	switch method {
	case "as-smallest":
		// TODO set asn in case these statements
		fmt.Println("Returning next prefix by smallest AS")
	case "org-smallest":
		fmt.Println("Unimplemented")
		return ""
	case "interesting":
		fmt.Println("Unimplemented")
		return ""
	}
	asn := getAsnFromSmallest(country)
	if !areDetailsInDB(asn) {
		details, peers, prefixes := getDetails(asn)
		whois := getWhoisDetails(asn)
		asnDetailsMongo(details, prefixes, peers, whois)
	}
	prefix := getFirstUnscannedPrefix(asn)
	return prefix
}

func areDetailsInDB(asn string) bool {
	return false
}

func MarkPrefixScanned(prefix string, asn string) {
}

// Return AS'es for a given country from the smallest that have registered
// active prefixes and are not scanned yet.
// If an AS is marked as interesting manually then return it and and it's
// unscanned downstream peers.
func getAsnFromSmallest(country string) string {
	client, ctx, _ := getMongoClient()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("masscan_go").Collection("as")
	opts := options.Find().SetSort(bson.D{{"number_addresses", 1}})
	cursor, err := collection.Find(ctx,
		bson.M{"country": country, "number_addresses": bson.M{"$gt": 0}},
		opts)
	if err != nil {
		panic(err)
	}

	var results []bson.D
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}
	for _, r := range results[0:10] {
		fmt.Println(r)
	}
	return "" // This functionality is on hold
}

func getAsnFromClosest(lat float64, lon float64) string {
	return "" // Functionality on hold
}

func getFirstUnscannedPrefix(asn string) string {
	return "" // Functionality on hold
}
