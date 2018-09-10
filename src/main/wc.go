package main

import (
	"fmt"
	"mapreduce"
	"os"
	"strings"
	"unicode"
	"strconv"
)

//
// The map function is called once for each file of input. The first
// argument is the name of the input file, and the second is the
// file's complete contents. You should ignore the input file name,
// and look only at the contents argument. The return value is a slice
// of key/value pairs.
//
func shouldSplit(word rune) bool {
	if !unicode.IsLetter(word) {
		return true
	}
	return false
}

func mapF(filename string, contents string) []mapreduce.KeyValue {
	wordCount := make(map[string]int)
	kvs := []mapreduce.KeyValue{}

	words := strings.FieldsFunc(contents, shouldSplit)
	for _, word := range words {
		wordCount[word] = wordCount[word] + 1
	}
	for word, count := range wordCount {
		kvs = append(kvs, mapreduce.KeyValue{Key: word, Value: strconv.Itoa(count)})
	}
	return kvs
	// Your code here (Part II).
}

//
// The reduce function is called once for each key generated by the
// map tasks, with a list of all the values created for that key by
// any map task.
//
func reduceF(key string, values []string) string {
	// Your code here (Part II).
	var sum int
	for _, value := range values {
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println(intvalue, "is not a Itoa-compatible string")
			continue
		}

		sum = sum + intvalue

	}
	return strconv.Itoa(sum)
}

// Can be run in 3 ways:
// 1) Sequential (e.g., go run wc.go master sequential x1.txt .. xN.txt)
// 2) Master (e.g., go run wc.go master localhost:7777 x1.txt .. xN.txt)
// 3) Worker (e.g., go run wc.go worker localhost:7777 localhost:7778 &)
func main() {
	if len(os.Args) < 4 {
		fmt.Printf("%s: see usage comments in file\n", os.Args[0])
	} else if os.Args[1] == "master" {
		var mr *mapreduce.Master
		if os.Args[2] == "sequential" {
			mr = mapreduce.Sequential("wcseq", os.Args[3:], 3, mapF, reduceF)
		} else {
			mr = mapreduce.Distributed("wcseq", os.Args[3:], 3, os.Args[2])
		}
		mr.Wait()
	} else {
		mapreduce.RunWorker(os.Args[2], os.Args[3], mapF, reduceF, 100, nil)
	}
}
