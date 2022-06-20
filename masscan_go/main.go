package main

import (
	"flag"

	"goas.xd/as"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// const TOTAL_SHARDS int = 1
// const SHARD int = 1
const RATE string = "10000"
const SCANNED_COUNTRY string = "PL"
const SCAN_METHOD string = "as-smallest"

var fullFlag = flag.Bool("full", true, "start full scanning")
var testFlag = flag.String("test", "", "test the tests for flags")

// Sharding would allow to distribute the scans among multiple machines and then
// merge the databases. This works for masscan, Nmap would have to scan only
// adresses scanned by masscan locally, then it would be sharded too.
var shardFlag = flag.Int("shard", 1, "which shard should run on this instance")
var totalShardsFlag = flag.Int("totalShards", 1, "how many shards are deployed in general")

// TODO use just one mongo client for package
// var DB *mongo.Client =

func main() {
	flag.Parse()
	if !*fullFlag {
		fmt.Println(as.Pubtest())
		os.Exit(1)
	}
	os.MkdirAll(filepath.Join(".", "tmp"), os.ModePerm)
	const portRange = "21,22,23,25,80,443,3389,110,445,139,3306"

	// take paths of masscan output files to parse
	var massOut, nmapOut chan string = make(chan string), make(chan string)
	// take prefixes to scan
	var massIn, nmapIn chan string = make(chan string), make(chan string)
	// var nmapOut chan string = make(chan string)
	// var nmapIn chan string = make(chan string)

	// Initialize input channels with data
	// TODO integrate it with the main loop in a clean way
	go func() {
		prefix, _ := getNextMass(SCAN_METHOD)
		massIn <- prefix
	}()
	go func() {
		ipRange, _ := getNextNmap()
		nmapIn <- ipRange
	}()
	// var path string
	for {
		select {
		case prefix := <-massIn:
			fmt.Printf("Received prefix %s, starting masscan", prefix)
			go func() {
				outputPath := runMass(prefix, portRange, *shardFlag, *totalShardsFlag)
				if outputPath != "" {
					massOut <- outputPath
				}
			}()
		case path := <-massOut:
			fmt.Printf("Masscan returned output to %s, saving to DB", path)
			entries := parseMassJson(path)
			go massToMongo(entries)
			go func() {
				prefix, _ := getNextMass(SCAN_METHOD)
				massIn <- prefix
			}()
		case ipRange := <-nmapIn:
			fmt.Printf("Received ipRange %s\n, starting nmap", ipRange)
			runNmap(ipRange, portRange)
		case path := <-nmapOut:
			fmt.Printf("Nmap returned output to %s, saving to DB", path)
			entries := parseNmapXML(path)
			go nmapToMongo(entries)
			go func() {
				ipRange, _ := getNextNmap()
				nmapIn <- ipRange
			}()
		}
	}

}

// Returns next available prefix for a ASN until are ASNs are scanned
// then moves on to the next ASN
func getNextMass(method string) (string, error) {
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

func isPrefixScanned(prefix string) bool {
	// Consult DB to check that
	return false
}

// Parse through mass results in the DB and get up to n (100)
// starting from the ones with lowest timestamps.
func getNextNmap() (string, error) {
	//  errors.New("next Nmap target not available")
	return "", nil
}

func runNmap(hosts string, ports string) string {
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

func parseNmapXML(path string) []nmapEntry {
	return []nmapEntry{{}}
}

func reverseDns(ip string) []string {
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

func runMass(ipRange string, portRange string, shard int, totalShards int) string {
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

func massToMongo(entries []massEntry) error {
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

func nmapToMongo(entries []nmapEntry) {

}

type massEntry struct {
	ip        string
	asn       int
	timestamp int
	port      int
	status    string
	protocol  string
	country   string
	rDNS      []string
}

type nmapEntry struct {
	// TODO
	ip    string
	ports []string
}
