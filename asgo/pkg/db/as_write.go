package db

import (
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AsnDetailsMongo(
	ad AsnDetails, apr AsnPrefixes, ape AsnPeers, wd WhoisDetails,
) error {
	client, ctx := DB.Client, DB.Ctx
	// defer func() {
	// 	if err := client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()
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
				"address":      wd.Organisation["address"],
				"org-name":     wd.Organisation["org-name"],
				"phone":        wd.Organisation["phone"],
				"fax-no":       wd.Organisation["fax-no"],
				"organisation": wd.Organisation["organisation"],
			},
			"autnum": bson.M{
				"created":       wd.Autnum["created"],
				"last-modified": wd.Autnum["last-modified"],
				"remarks":       wd.Autnum["remarks"],
			},
			"persons": wd.Persons}},
		opts)
	if err != nil {
		return err
	}
	fmt.Println(res)
	fmt.Println(err)
	return nil
}

type AsnDetails struct {
	Status        string `json:"status"`
	StatusMessage string `json:"status_message"`
	Data          struct {
		Asn              int         `json:"asn"`
		Name             interface{} `json:"name"`
		DescriptionShort string      `json:"description_short"`
		DescriptionFull  []string    `json:"description_full"`
		CountryCode      string      `json:"country_code"`
		Website          interface{} `json:"website"`
	} `json:"data"`
}

type AsnPeers struct {
	Status        string `json:"status"`
	StatusMessage string `json:"status_message"`
	Data          struct {
		Ipv4Peers []struct {
			Asn         int         `json:"asn"`
			Name        interface{} `json:"name"` // seems to be always null
			Description string      `json:"description"`
			CountryCode string      `json:"country_code"`
		} `json:"ipv4_peers"`
		Ipv6Peers []struct {
			Asn         int         `json:"asn"`
			Name        interface{} `json:"name"`
			Description string      `json:"description"`
			CountryCode string      `json:"country_code"`
		} `json:"ipv6_peers"`
	} `json:"data"`
}

type AsnPrefixes struct {
	Status        string `json:"status"`
	StatusMessage string `json:"status_message"`
	Data          struct {
		Ipv4Prefixes []struct {
			// Prefix string `json:"prefix"`
			IP          string `json:"ip"`
			Cidr        int    `json:"cidr"`
			Name        string `json:"name"`
			Description string `json:"description"`
			CountryCode string `json:"country_code"`
			Parent      struct {
				Prefix string `json:"prefix"`
				IP     string `json:"ip"`
				Cidr   int    `json:"cidr"`
			} `json:"parent"`
		} `json:"ipv4_prefixes"`
		Ipv6Prefixes []struct {
			// Prefix string `json:"prefix"`
			IP          string `json:"ip"`
			Cidr        int    `json:"cidr"`
			Name        string `json:"name"`
			Description string `json:"description"`
			CountryCode string `json:"country_code"`
			Parent      struct {
				Prefix string `json:"prefix"`
				IP     string `json:"ip"`
				Cidr   int    `json:"cidr"`
			} `json:"parent"`
		} `json:"ipv6_prefixes"`
	} `json:"data"`
}

type WhoisDetails struct {
	Organisation map[string]string
	Autnum       map[string]string
	Persons      []map[string]string
}
