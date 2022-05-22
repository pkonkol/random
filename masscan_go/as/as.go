package as

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	BGPVIEW_API_URL = "https://api.bgpview.io"
	ASRANK_API_URL  = "https://api.asrank.caida.org/v2/restful"
	MONGODB_USER    = "test"
	MONGODB_PASS    = "test"
	MONGODB_IP      = "localhost"
	MONGODB_DB      = "masscan_go"
)

func Pubtest() string {
	return "works"
}

// Return AS'es for a given country from the smallest that have registered
// active prefixes and are not scanned yet.
// If an AS is marked as interesting manually then return it and and it's
// unscanned downstream peers.
func GetASFromSmallest(country string) {
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
}

func GetASFromClosest(lat float64, lon float64) {
	// TODO
}

func ScanIDMongo() {

}

func getAsGeocoding() {
	// TODO Later most likely with Google api or bing api as only these
	// didn't suck for example address from TASK
	// Using batch requests should bring down the price
}

func getDetails(as string) (asnDetails, asnPeers, asnPrefixes) {
	url := BGPVIEW_API_URL + "/asn/" + as
	bodyBytes := makeApiCall(url) // test for first element
	var details asnDetails
	json.Unmarshal(bodyBytes, &details)

	// should i alsoget upstreams and downstreams specifically or doesn't matter?
	url = BGPVIEW_API_URL + "/asn/" + as + "/peers"
	bodyBytes = makeApiCall(url)
	var peers asnPeers
	json.Unmarshal(bodyBytes, &peers)
	// calculate prefix count and addresses sum
	url = BGPVIEW_API_URL + "/asn/" + as + "/prefixes"
	bodyBytes = makeApiCall(url)
	var prefixes asnPrefixes
	json.Unmarshal(bodyBytes, &prefixes)
	// calculate prefix count and addresses sum

	return details, peers, prefixes
}

func getWhoisDetails(as string) WhoisDetails {
	asName := fmt.Sprintf("as%s", as)
	raw, err := exec.Command("whois", asName).Output()
	if err != nil {
		panic(err)
	}
	output := string(raw)

	// a := regexp.MustCompile(`^\s+$`).Split(output, -1)
	a := strings.Split(output, "\n\n")
	org := make(map[string]string)
	autnum := make(map[string]string)
	var persons []map[string]string
	for _, s := range a {
		switch {
		case regexp.MustCompile("^aut-num:").MatchString(s):
			fmt.Println("autnum")
			for _, s2 := range strings.Split(s, "\n") {
				s2_content := regexp.MustCompile(`^[^\s]+:\s+`).ReplaceAllString(s2, "")
				switch {
				case regexp.MustCompile("^remarks").MatchString(s2):
					autnum["remarks"] = autnum["remarks"] + ", " + s2_content
				case regexp.MustCompile("^created").MatchString(s2):
					autnum["created"] = s2_content
				case regexp.MustCompile("^last-modified").MatchString(s2):
					autnum["last-modified"] = s2_content
				}
			}
		case regexp.MustCompile("^organisation:").MatchString(s):
			fmt.Println("org")
			for _, s2 := range strings.Split(s, "\n") {
				s2_content := regexp.MustCompile(`^[^\s]+:\s+`).ReplaceAllString(s2, "")
				switch {
				case regexp.MustCompile("^address").MatchString(s2):
					org["address"] = org["address"] + ", " + s2_content
				case regexp.MustCompile("^org-name").MatchString(s2):
					org["org-name"] = s2_content
				case regexp.MustCompile("^organisation").MatchString(s2):
					org["organisation"] = s2_content
				case regexp.MustCompile("^phone").MatchString(s2):
					org["phone"] = s2_content
				case regexp.MustCompile("^fax-no").MatchString(s2):
					org["fax-no"] = s2_content
				}
			}
		case regexp.MustCompile("^person:").MatchString(s):
			fmt.Println("person")
			person := make(map[string]string)
			for _, s2 := range strings.Split(s, "\n") {
				s2_content := regexp.MustCompile(`^[^\s]+:\s+`).ReplaceAllString(s2, "")
				switch {
				case regexp.MustCompile("^address").MatchString(s2):
					person["address"] = person["address"] + ", " + s2_content
				case regexp.MustCompile("^person").MatchString(s2):
					person["person"] = s2_content
				case regexp.MustCompile("^org-name").MatchString(s2):
					person["org-name"] = s2_content
				case regexp.MustCompile("^phone").MatchString(s2):
					person["phone"] = s2_content
				case regexp.MustCompile("^fax-no").MatchString(s2):
					person["fax-no"] = s2_content
				}
			}
			persons = append(persons, person)
		}
	}
	fmt.Println(org)
	fmt.Println("")
	fmt.Println(autnum)
	fmt.Println(persons)

	return WhoisDetails{org, autnum, persons}
}

func asnDetailsMongo(
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

func asrankToMongo() {
	const reqEntries int = 1000
	client, ctx, _ := getMongoClient()
	collection := client.Database("masscan_go").Collection("as")
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
	// cur, err := collection.Find(ctx, bson.D{})
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
	client, ctx, _ := getMongoClient()
	collection := client.Database("masscan_go").Collection("org")
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

func getMongoClient() (*mongo.Client, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Printf("%#v", cancel)
	fmt.Printf("mongodb://%s:%s@%s:27017/%s\n",
		MONGODB_USER, MONGODB_PASS, MONGODB_IP, MONGODB_DB)
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:27017/%s",
			MONGODB_USER, MONGODB_PASS, MONGODB_IP, MONGODB_DB)))
	if err != nil {
		panic(err)
	}
	return client, ctx, cancel
}

func makeApiCall(url string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return bodyBytes
}

type WhoisDetails struct {
	organisation map[string]string
	autnum       map[string]string
	persons      []map[string]string
}

type asnDetails struct {
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

type asnPrefixes struct {
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

type asnPeers struct {
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

type AsRankEntry struct {
	Data struct {
		Asns struct {
			TotalCount int `json:"totalCount"`
			PageInfo   struct {
				HasNextPage bool `json:"hasNextPage"`
			} `json:"pageInfo"`
			Edges []struct {
				Node struct {
					Rank         int         `json:"rank"`
					Asn          string      `json:"asn"`
					AsnName      string      `json:"asnName"`
					Source       string      `json:"source"`
					Seen         bool        `json:"seen"`
					Ixp          interface{} `json:"ixp"`
					Longitude    float64     `json:"longitude"`
					Latitude     float64     `json:"latitude"`
					Organization struct {
						OrgID string `json:"orgId"`
					} `json:"organization"`
					Cone struct {
						NumberAsns      int `json:"numberAsns"`
						NumberPrefixes  int `json:"numberPrefixes"`
						NumberAddresses int `json:"numberAddresses"`
					} `json:"cone"`
					Country struct {
						Iso string `json:"iso"`
					} `json:"country"`
					AsnDegree struct {
						Total    int `json:"total"`
						Customer int `json:"customer"`
						Peer     int `json:"peer"`
						Provider int `json:"provider"`
					} `json:"asnDegree"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"asns"`
	} `json:"data"`
}

type asOrgEntry struct {
	Data struct {
		Organization struct {
			Rank    int    `json:"rank"`
			OrgID   string `json:"orgId"`
			OrgName string `json:"orgName"`
			Seen    bool   `json:"seen"`
			Cone    struct {
				NumberAsns      int `json:"numberAsns"`
				NumberPrefixes  int `json:"numberPrefixes"`
				NumberAddresses int `json:"numberAddresses"`
			} `json:"cone"`
			Country struct {
				Iso string `json:"iso"`
			} `json:"country"`
		} `json:"organization"`
	} `json:"data"`
}

type asrankOrgs struct {
	Data struct {
		Organizations struct {
			PageInfo struct {
				HasNextPage bool `json:"hasNextPage"`
			} `json:"pageInfo"`
			Edges []struct {
				Node struct {
					Rank      int    `json:"rank"`
					OrgID     string `json:"orgId"`
					OrgName   string `json:"orgName"`
					AsnDegree struct {
						Total   int `json:"total"`
						Transit int `json:"transit"`
					} `json:"asnDegree"`
					Cone struct {
						NumberAsns      int `json:"numberAsns"`
						NumberPrefixes  int `json:"numberPrefixes"`
						NumberAddresses int `json:"numberAddresses"`
					} `json:"cone"`
					Country struct {
						Iso string `json:"iso"`
					} `json:"country"`
					Members struct {
						Asns struct {
							Edges []struct {
								Node struct {
									Asn string `json:"asn"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"asns"`
					} `json:"members"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"organizations"`
	} `json:"data"`
}
