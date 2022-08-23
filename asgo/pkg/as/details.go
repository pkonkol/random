package as

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkonkol/random/asgo/pkg/db"
)

func ReverseDns(ip string) []string {
	path, err := exec.LookPath("dig")
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(path, ip)
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	var rDnsList []string
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)
		rDnsList = append(rDnsList, fields[4]+" "+fields[3])
	}
	return rDnsList
}

func getDetails(as string) (db.AsnDetails, db.AsnPeers, db.AsnPrefixes) {
	url := BGPVIEW_API_URL + "/asn/" + as
	bodyBytes := makeApiCall(url) // test for first element
	var details db.AsnDetails
	json.Unmarshal(bodyBytes, &details)

	// should i alsoget upstreams and downstreams specifically or doesn't matter?
	url = BGPVIEW_API_URL + "/asn/" + as + "/peers"
	bodyBytes = makeApiCall(url)
	var peers db.AsnPeers
	json.Unmarshal(bodyBytes, &peers)
	// calculate prefix count and addresses sum
	url = BGPVIEW_API_URL + "/asn/" + as + "/prefixes"
	bodyBytes = makeApiCall(url)
	var prefixes db.AsnPrefixes
	json.Unmarshal(bodyBytes, &prefixes)
	// calculate prefix count and addresses sum

	return details, peers, prefixes
}

func getWhoisDetails(as string) db.WhoisDetails {
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

	return db.WhoisDetails{org, autnum, persons}
}
