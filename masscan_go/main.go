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

const TOTAL_SHARDS int = 1
const RATE string = "10000"

var fullFlag = flag.Bool("full", true, "start full scanning")

// var DB *mongo.Client =

func main() {
	fmt.Println(as.Pubtest())
	os.Exit(1)
	os.MkdirAll(filepath.Join(".", "tmp"), os.ModePerm)
	// ipRange := "91.209.116.0/24"
	const portRange = "21,22,23,25,80,443,3389,110,445,139,3306"

	// take paths of masscan output files to parse
	var massOut chan string = make(chan string, TOTAL_SHARDS)
	// take prefixes to scan
	var massIn chan string = make(chan string)
	var nmapIn chan string = make(chan string)
	var nmapOut chan string = make(chan string)
	// for i := 0; i < TOTAL_SHARDS; i++ {
	// 	// channels[i] = make(chan string)
	// 	go func() {
	// 		outputPath := runMasscan(ipRange, portRange, 1)
	// 		if outputPath != "" {
	// 			massOut <- outputPath
	// 		}
	// 	}()
	// }

	// var path string
	for {
		select {
		case prefix := <-massIn:
			for i := 0; i < TOTAL_SHARDS; i++ {
				// channels[i] = make(chan string)
				go func() {
					outputPath := runMass(ipRange, portRange, 1)
					if outputPath != "" {
						massOut <- outputPath
					}
				}()
			}
		case path := <-massOut:
			entries := parseJson(path)
			saveToMongo(entries)
		case ipRange := <-nmapIn:
			runNmap(ipRange, portRange)
		case <-nmapOut:
			fmt.Println("another scan done nmap")
		default:
			fmt.Println("in default")
		}
	}

}

func getNextMass() {

}

// Parse through mass results in the DB and get up to n (100)
// starting from the ones with lowest timestamps.
func getNextNmap() string {

}

func runNmap(hosts string, ports string) string {
	// wait := "0"
	// shards := fmt.Sprintf("%d/%d", shard, TOTAL_SHARDS)

	outputPath := fmt.Sprintf("tmp/nmap_{}_{}_{}.xml", hosts, ports, strconv.Itoa(rand.Int())) // TODO: randomize
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

func parseNmapJson() {

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

func runMass(ipRange string, portRange string, shard int) string {
	wait := "0"
	shards := fmt.Sprintf("%d/%d", shard, TOTAL_SHARDS)
	outputPath := fmt.Sprintf("tmp/masscan_{}_{}.json", ipRange, strconv.Itoa(rand.Int())) // TODO: randomize
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
	var country string
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

func nmapToMongo() {

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
