package as

import (
	"encoding/json"
	"fmt"

	"github.com/pkonkol/random/asgo/pkg/db"
)

func GenerateDBOverview() {
	fmt.Println("Generating DB overview")
	return
	asrankToMongo()
	asrankOrgMongo()
}

func asrankToMongo() {
	const reqEntries int = 1000
	collection := db.DB.client.Database("masscan_go").Collection("as")
	opts := options.Update().SetUpsert(true)

	var total_count int64
	for first, offset := reqEntries, 0; ; {
		url := fmt.Sprintf("%s/asns/?first=%d&offset=%d", ASRANK_API_URL, first, offset)
		bytes := makeApiCall(url)
		var entry AsRankEntry
		json.Unmarshal(bytes, &entry)
		fmt.Println(entry.Data.Asns.PageInfo.HasNextPage)
		fmt.Println(entry.Data.Asns.Edges[0].Node.Rank)
		fmt.Println(entry.Data.Asns.Edges[len(entry.Data.Asns.Edges)-1].Node.Rank)
		if entry.Data.Asns.PageInfo.HasNextPage != true {
			total_count = int64(entry.Data.Asns.TotalCount)
			break
		}

		for _, e := range entry.Data.Asns.Edges {
			e := e.Node
			_, err := collection.UpdateOne(ctx, bson.M{"as_number": e.Asn},
				bson.M{"$set": bson.M{
					"name":             e.AsnName,
					"country":          e.Country.Iso,
					"rank":             e.Rank,
					"interesting":      false,
					"org_id_caida":     e.Organization.OrgID,
					"lat_caida":        e.Latitude,
					"lon_caida":        e.Longitude,
					"number_asns":      e.Cone.NumberAsns,
					"number_prefixes":  e.Cone.NumberPrefixes,
					"number_addresses": e.Cone.NumberAddresses,
					"customers":        e.AsnDegree.Customer,
					"peers":            e.AsnDegree.Peer,
					"providers":        e.AsnDegree.Provider}},
				opts)
			if err != nil {
				panic(err)
			}
		}
		offset += reqEntries
	}
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	fmt.Println(count, "== t: ", total_count)
	if count != total_count {
		fmt.Println("smth fucked up")
	}

	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func asrankOrgMongo() {
	const reqEntries int = 1000
	collection := db.DB.client.Database("masscan_go").Collection("org")
	opts := options.Update().SetUpsert(true)

	for first, offset := reqEntries, 0; ; {
		url := fmt.Sprintf("%s/organizations/?first=%d&offset=%d", ASRANK_API_URL, first, offset)
		bytes := makeApiCall(url)
		var entry asrankOrgs
		json.Unmarshal(bytes, &entry)
		if entry.Data.Organizations.PageInfo.HasNextPage != true {
			break
		}

		for _, e := range entry.Data.Organizations.Edges {
			var member_asns string
			for _, huj := range e.Node.Members.Asns.Edges {
				member_asns = member_asns + huj.Node.Asn
			}
			e := e.Node
			_, err := collection.UpdateOne(ctx, bson.M{"orgId": e.OrgID},
				bson.M{"$set": bson.M{
					"name":               e.OrgName,
					"country":            e.Country.Iso,
					"rank":               e.Rank,
					"number_asns":        e.Cone.NumberAsns,
					"number_prefixes":    e.Cone.NumberPrefixes,
					"number_addresses":   e.Cone.NumberAddresses,
					"asn_degree_total":   e.AsnDegree.Total,
					"asn_degree_transit": e.AsnDegree.Transit,
					"member_asns":        member_asns}},
				opts)
			if err != nil {
				panic(err)
			}
		}
		offset += reqEntries
	}
	// cur, err := collection.Find(ctx, bson.D{})
	_, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}
