package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/kroekle/zipcode-generator/pop"
)

func main() {
	help := flag.Bool("help", false, "Prints Usage")
	zipcodes := flag.Int("zipcodes", 41000, "Approx number of zip codes to generate")
	output_file := flag.String("out", "zips.json", "Name of the file to output")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}
	log.Printf("zipcode: %v outfile: %v", *zipcodes, *output_file)

	var minState, maxState string
	var totalPop, max int
	min := 100000000000
	states := make([]string, 0)

	for k, v := range pop.STATE_POPULATIONS {
		ignore := map[string]bool{"PW": true, "MP": true, "AS": true, "MH": true, "VI": true, "FM": true, "GU": true}
		if ignore[k] {
			continue
		}
		totalPop += v
		if v < min {
			minState = k
			min = v
		}
		if v > max {
			maxState = k
			max = v
		}
		for i := 0; i < rnd100k(v); i++ {
			states = append(states, k)
		}
	}

	zips := make(map[int]string, 0)
	for i := 0; i < *zipcodes; i++ {
		zips[getAvailableZip(zips)] = randState(states)
	}

	log.Printf("total US population: %v; Min: %v(%v); Max %v(%v); size: %v; zips: %v", totalPop, minState, rnd100k(min), maxState, rnd100k(max), len(states), len(zips))
	stateNames := make([]string, 0, len(pop.STATE_POPULATIONS))
	for k, _ := range pop.STATE_POPULATIONS {
		stateNames = append(stateNames, k)
	}
	sort.Strings(stateNames)

	all_zips := make([]Zip, 0, len(zips))
	for _, v := range stateNames {

		total := 0

		for k, s := range zips {
			if s == v {
				total++
				//TODO: only doing this here because I also want counts printed
				all_zips = append(all_zips, Zip{
					State:   v,
					Zipcode: fmt.Sprintf("%05d", k),
				})
			}
		}

		log.Printf("%v: %v\n", v, total)
	}

	json, err := json.Marshal(all_zips)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(*output_file, json, 0644)
}

type Zip struct {
	Zipcode string `json:"zip_code"`
	State   string `json:"state"`
}

func randState(states []string) string {
	rand.Seed(time.Now().UnixMicro())
	return states[rand.Intn(len(states)-1)]
}

func getAvailableZip(zips map[int]string) int {

	rand.Seed(time.Now().UnixMicro())
	z := rand.Intn(99999)
	for ; zips[z] != ""; z = rand.Intn(99999) {
	}
	return z
}

func rnd100k(num int) int {
	return int(math.Round(float64(num) / float64(100000)))
}
