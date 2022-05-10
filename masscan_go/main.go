package main

import (
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

func main() {
	fmt.Println(as.Pubtest())
	os.MkdirAll(filepath.Join(".", "tmp"), os.ModePerm)
	ipRange := "91.209.116.0/24"
	portRange := "21,22,23,25,80,443,3389,110,445,139,3306"

	var massOut chan string
	var prefixChan chan string
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
	select {
	case prefix := <-prefixChan:
		for i := 0; i < TOTAL_SHARDS; i++ {
			// channels[i] = make(chan string)
			go func() {
				outputPath := runMasscan(ipRange, portRange, 1)
				if outputPath != "" {
					massOut <- outputPath
				}
			}()
		}
	case path := <-massOut:
		entries := parseJson(path)
		saveToMongo(entries)
	default:
		fmt.Println("in default")
	}

	// entries := parseJson(out_path)
	// saveToMongo(entries)
}

// func manageOutput() {
// 	entries := parseJson()
// 	saveToMongo(entries)
// }
func getUncheckedAS() {

}

func runNmap(host string, ports string) string {
	// wait := "0"
	// shards := fmt.Sprintf("%d/%d", shard, TOTAL_SHARDS)

	outputPath := fmt.Sprintf("tmp/nmap_{}_{}_{}.xml", host, ports, strconv.Itoa(rand.Int())) // TODO: randomize
	nmapPath, err := exec.LookPath("masscan")
	if err != nil {
		panic(err)
	}

	fmt.Println(nmapPath, host, "-Pn", "-R", "-sV", "-O", "--top-ports", "100", "-p", ports, "-oX", outputPath)
	cmd := exec.Command(nmapPath, host)
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
}

func parseJson(path string) []scanEntry {
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

	entries := make([]scanEntry, len(rawEntries))
	var country string
	var rDns string
	for i, e := range rawEntries {
		country = geoIpLookup(e.Ip)
		rDns = reverseDns(e.Ip)
		t, _ := strconv.Atoi(e.Timestamp)
		entries[i] = scanEntry{
			ip:        e.Ip,
			timestamp: t,
			port:      int(e.Ports[0]["port"].(float64)),
			status:    e.Ports[0]["status"].(string),
			protocol:  e.Ports[0]["proto"].(string),
			country:   country,
			rDNS:      rDns}
	}
	fmt.Printf("%+v\n", entries)
	return entries
}

func geoIpLookup(ip string) string {
	path, err := exec.LookPath("geoiplookup")
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(path, ip)
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimLeft(string(output), "GeoIP Country Edition: ")
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

func runMasscan(ipRange string, portRange string, shard int) string {
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

func saveToMongo(entries []scanEntry) error {
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
	collection := client.Database("masscan_go").Collection("scan_results")
	opts := options.Update().SetUpsert(true)
	for _, e := range entries {
		collection.UpdateOne(ctx, bson.M{"ip": e.ip},
			bson.M{"$addToSet": bson.M{
				"ports": bson.D{
					{"timestamp", e.timestamp},
					{"port", e.port},
					{"status", e.status},
					{"protocol", e.protocol}}}}, opts)
	}
	return nil
}

type scanEntry struct {
	ip        string
	as        int
	timestamp int
	port      int
	status    string
	protocol  string
	country   string
	rDNS      []string
}
