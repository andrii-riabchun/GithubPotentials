package main

import (
	"fmt"
	potentials "github.com/artisresistance/githubpotentials"
	"os"
	"time"
)

const configPath = "config.json"
const outputPath = "data.json"

var conf config
var startTime time.Time

func init() {
	startTime = time.Now()

	file, err := os.Open(configPath)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	conf, err = loadConfig(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("config loaded")
}

func main() {
	updatedFrom := time.Now().AddDate(0, 0, -1)

	client := potentials.New(conf.Token, updatedFrom)

	errCount := 0
	onError := func(err error) {
		errCount++
		fmt.Fprintln(os.Stderr, err.Error())
	}
	criteria := potentials.CombinedCriteria

	fmt.Println("sorting by:", criteria.String())

	it := client.SearchIterator(onError)

	repos := client.CountStats(it, onError).
		FilterZeroStats(criteria).
		Dump(onError).
		Sort(criteria)

	fmt.Println("done!")
	fmt.Println("error count:", errCount)
	rates, err := client.GetAPIRates()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("rates:", rates)

	fmt.Println("total fetched:", len(repos))
	fmt.Println("selecting top", conf.OutCount)

	trimmed := trim(repos, conf.OutCount)

	out := potentialsResult{
		Updated:  time.Now().Unix(),
		Errors:   errCount,
		Fetched:  len(trimmed),
		SortedBy: criteria.String(),
		Items:    trimmed,
	}

	err = writeToFile(out)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Execution time: %s", elapsed)
}

func writeToFile(result potentialsResult) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = result.Write(file)
	if err != nil {
		return err
	}
	err = file.Sync()
	return err
}

func trim(repos []potentials.Repository, count int) []potentials.Repository {
	bound := count
	if len(repos) < count {
		bound = len(repos)
	}

	return repos[:bound]
}
