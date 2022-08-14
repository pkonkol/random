package as

import (
	"io/ioutil"
	"net/http"
)

const (
	BGPVIEW_API_URL = "https://api.bgpview.io"
	ASRANK_API_URL  = "https://api.asrank.caida.org/v2/restful"
)

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
