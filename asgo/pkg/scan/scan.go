package scan

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkonkol/random/asgo/pkg/as"
)

// const TOTAL_SHARDS int = 1
// const SHARD int = 1
const RATE string = "10000"
const SCANNED_COUNTRY string = "PL"
const SCAN_METHOD string = "as-smallest"
const PORT_RANGE string = "21,22,23,25,80,443,3389,110,445,139,3306"

// Sharding would allow to distribute the scans among multiple machines and then
// merge the databases. This works for masscan, Nmap would have to scan only
// adresses scanned by masscan locally, then it would be sharded too.
const SHARD int = 1
const TOTAL_SHARDS int = 1

func Run() {
	// take paths of masscan output files to parse
	var massOut, nmapOut chan string = make(chan string), make(chan string)
	// take prefixes to scan
	var massIn, nmapIn chan string = make(chan string), make(chan string)
	// var nmapOut chan string = make(chan string)
	// var nmapIn chan string = make(chan string)

	go func() {
		prefix, _ := GetNextMass(SCAN_METHOD)
		massIn <- prefix
	}()
	go func() {
		ipRange, _ := GetNextNmap()
		nmapIn <- ipRange
	}()
	// var path string
	for {
		select {
		case prefix := <-massIn:
			fmt.Printf("Received prefix %s, starting masscan", prefix)
			go func() {
				outputPath := RunMass(prefix, PORT_RANGE, SHARD, TOTAL_SHARDS)
				if outputPath != "" {
					massOut <- outputPath
				}
			}()
		case path := <-massOut:
			fmt.Printf("Masscan returned output to %s, saving to DB", path)
			entries := parseMassJson(path)
			go db.MassToMongo(entries)
			go func() {
				prefix, _ := GetNextMass(SCAN_METHOD)
				massIn <- prefix
			}()
		case ipRange := <-nmapIn:
			fmt.Printf("Received ipRange %s\n, starting nmap", ipRange)
			scan.RunNmap(ipRange, PORT_RANGE)
		case path := <-nmapOut:
			fmt.Printf("Nmap returned output to %s, saving to DB", path)
			entries := parseNmapXML(path)
			go db.NmapToMongo(entries)
			go func() {
				ipRange, _ := GetNextNmap()
				nmapIn <- ipRange
			}()
		}
	}
}

// Returns next available prefix for a ASN until are ASNs are scanned
// then moves on to the next ASN
func GetNextMass(method string) (string, error) {
	var prefix string
	for {
		prefix = as.GetNextPrefix(SCANNED_COUNTRY, method)
		if isPrefixScanned(prefix) {
			as.MarkPrefixScanned(prefix)
			continue
		}
		break
	}
	//  errors.New("next Masscan target not available")
	return prefix, nil
}

// Parse through mass results in the DB and get up to n (100)
// starting from the ones with lowest timestamps.
func GetNextNmap() (string, error) {
	//  errors.New("next Nmap target not available")
	return "", nil
}

func RunNmap(hosts string, ports string) string {
	outputPath := fmt.Sprintf("tmp/nmap_%s_%s_%d.xml", hosts, ports, rand.Int()) // TODO: randomize
	nmapPath, err := exec.LookPath("masscan")
	if err != nil {
		panic(err)
	}

	fmt.Println(nmapPath, hosts, "-Pn", "-R", "-sV", "-O", "--top-ports", "100", "-p", ports, "-oX", outputPath)
	cmd := exec.Command(nmapPath, hosts)
	var buf_out, buf_err bytes.Buffer
	cmd.Stdout = &buf_out
	cmd.Stderr = &buf_err

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	err = cmd.Wait()
	fmt.Println(err)
	fmt.Println(buf_out.String())
	fmt.Println(buf_err.String())

	return outputPath
}

func RunMass(ipRange string, portRange string, shard int, totalShards int) string {
	wait := "0"
	shards := fmt.Sprintf("%d/%d", shard, totalShards)
	outputPath := fmt.Sprintf("tmp/masscan_%s_%d.json", ipRange, rand.Int()) // TODO: randomize
	masscanPath, err := exec.LookPath("masscan")
	fmt.Println(masscanPath)
	if err != nil {
		panic(err)
	}

	fmt.Println(masscanPath, ipRange, "-p", portRange, "--rate", RATE,
		"--wait", wait, "--shards", shards, "-oJ", outputPath)
	cmd := exec.Command(masscanPath, ipRange, "-p", portRange, "--rate", RATE,
		"--wait", wait, "--shards", shards, "-oJ", outputPath)
	var (
		buf_out, buf_err bytes.Buffer
	)
	cmd.Stdout = &buf_out
	cmd.Stderr = &buf_err

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(buf_out.String())
	fmt.Println(buf_err.String())

	return buf_out.String()
}

func parseNmapXML(path string) []nmapEntry {
	return []nmapEntry{{}}
}

func parseMassJson(path string) []massEntry {
	scanOutput, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	scanOutput = []byte(strings.Replace(string(scanOutput), "},\n]", "}\n]", 1))
	var rawEntries []struct {
		Ip        string                   `json:"ip"`
		Timestamp string                   `json:"timestamp"`
		Ports     []map[string]interface{} `json:"ports"`
	}

	err = json.Unmarshal(scanOutput, &rawEntries)
	if err != nil {
		panic(err)
	}

	entries := make([]massEntry, len(rawEntries))
	var rDns []string
	for i, e := range rawEntries {
		rDns = reverseDns(e.Ip)
		t, _ := strconv.Atoi(e.Timestamp)
		entries[i] = massEntry{
			ip:        e.Ip,
			timestamp: t,
			port:      int(e.Ports[0]["port"].(float64)),
			status:    e.Ports[0]["status"].(string),
			protocol:  e.Ports[0]["proto"].(string),
			rDNS:      rDns}
	}
	fmt.Printf("%+v\n", entries)
	return entries
}
