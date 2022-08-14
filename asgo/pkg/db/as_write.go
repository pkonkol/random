package db

import (
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

func AsnDetailsMongo(
	ad asnDetails, apr asnPrefixes, ape asnPeers, wd WhoisDetails,
) error {
	client, ctx, _ := getMongoClient()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("masscan_go").Collection("as")
	opts := options.Update().SetUpsert(true)
	fmt.Println("ASN IS: " + strconv.Itoa(ad.Data.Asn))
	res, err := collection.UpdateOne(ctx,
		bson.M{"as_number": strconv.Itoa(ad.Data.Asn)},
		bson.M{"$set": bson.M{
			"description_short": ad.Data.DescriptionShort,
			"description_full":  ad.Data.DescriptionFull,

			"peers":     ape.Data.Ipv4Peers,
			"ipV6peers": ape.Data.Ipv6Peers,

			"prefixes":     apr.Data.Ipv4Prefixes,
			"ipV6prefixes": apr.Data.Ipv6Prefixes,

			"organisation": bson.M{
				"address":      wd.organisation["address"],
				"org-name":     wd.organisation["org-name"],
				"phone":        wd.organisation["phone"],
				"fax-no":       wd.organisation["fax-no"],
				"organisation": wd.organisation["organisation"],
			},
			"autnum": bson.M{
				"created":       wd.autnum["created"],
				"last-modified": wd.autnum["last-modified"],
				"remarks":       wd.autnum["remarks"],
			},
			"persons": wd.persons}},
		opts)
	if err != nil {
		return err
	}
	fmt.Println(res)
	fmt.Println(err)
	return nil
}
