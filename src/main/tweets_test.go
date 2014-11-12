package main

import (
	"fmt"
	"testing"
)

func BenchmarkFunction() {
	encodeJson(findLocationForTweets(parseJson(searchTweets(string("ibm")))))
}

func main() {
	br := testing.Benchmark(BenchmarkFunction)
	fmt.Println(br)
}
