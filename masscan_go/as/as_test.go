package as

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestGetAsDetails(t *testing.T) {
	asNumber := "12831" // TASK gdansk

	details, peers, prefixes := getDetails(asNumber)
	fmt.Printf("%v\n%v\n%v\n", details, prefixes, peers)
	if details.Status != "ok" || prefixes.Status != "ok" || peers.Status != "ok" {
		t.Fatalf("Query statuses are not right")
	}
}

func TestGetMongoClient(t *testing.T) {
	client, ctx, cancel := getMongoClient()
	defer cancel()
	db := client.Database("masscan_go")
	fmt.Println(client)
	fmt.Println(db)
	fmt.Println(db.ListCollectionNames(ctx, bson.D{{}}))
	fmt.Println(client.ListDatabaseNames(ctx, nil))
}

func TestAsrankToMongo(t *testing.T) {
	client, ctx, cancel := getMongoClient()
	// fmt.Println(c.CountDocuments(ctx, bson.M{}))
	defer cancel()
	client.Database("masscan_go").Collection("as").Drop(ctx)
	asrankToMongo()
	c := client.Database("masscan_go").Collection("as")
	fmt.Println(c.CountDocuments(ctx, bson.M{}))
	// if count, _ := c.CountDocuments(ctx, nil); count != 2000 {
	// 	t.Fatalf("wrong count of docs Got: %d Want: %d", count, 2000)
	// }
}

func TestAsrankOrgsToMongo(t *testing.T) {
	client, ctx, cancel := getMongoClient()
	// fmt.Println(c.CountDocuments(ctx, bson.M{}))
	defer cancel()
	client.Database("masscan_go").Collection("org").Drop(ctx)
	asrankOrgMongo()
	c := client.Database("masscan_go").Collection("org")
	fmt.Println(c.CountDocuments(ctx, bson.M{}))
}

func TestAsDetailsMongo(t *testing.T) {
	asNumber := "12831" // TASK gdansk
	details, peers, prefixes := getDetails(asNumber)
	whois := getWhoisDetails(asNumber)
	err := asnDetailsMongo(details, prefixes, peers, whois)
	if err != nil {
		t.Fatalf("upload failed")
	}
}

func TestWhoisGet(t *testing.T) {
	asNumber := "12831" // TASK gdansk
	out := getWhoisDetails(asNumber)
	fmt.Printf("%+v\n\n", out)
	fmt.Printf("%+v\n", out.persons)
}

func TestGetInterestingAS(t *testing.T) {
	GetASFromSmallest("PL")
}
