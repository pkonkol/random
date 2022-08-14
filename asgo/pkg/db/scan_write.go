package db

import (
	"context"
	"time"
)

type MassEntry struct {
	ip        string
	asn       int
	timestamp int
	port      int
	status    string
	protocol  string
	country   string
	rDNS      []string
}

type NmapEntry struct {
	// TODO
	ip    string
	ports []string
}

func MassToMongo(entries []massEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://test:test@localhost:27017/masscan_go"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	scan := client.Database("masscan_go").Collection("scan")
	opts := options.Update().SetUpsert(true)
	for _, e := range entries {
		scan.UpdateOne(ctx, bson.M{"ip": e.ip},
			bson.M{"$addToSet": bson.M{
				"ports": bson.D{
					{"timestamp", e.timestamp},
					{"port", e.port},
					{"status", e.status},
					{"protocol", e.protocol},
					{"asn", e.asn}}}}, opts)
	}
	return nil
}

func NmapToMongo(entries []nmapEntry) {
}
