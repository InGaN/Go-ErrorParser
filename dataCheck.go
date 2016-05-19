package main

import (
	"fmt"
	"io/ioutil"
	s "strings"
	"time"
)

func checkFile(e error) {
	if e != nil {
		panic(e)
	}
}

func checkContains(input string, tags []string) bool {
	fmt.Printf("Amount of tags: %d\n", len(tags))
	var val bool = false

	for index, element := range tags {
		fmt.Printf("idx: %d - %s - ", index, element)
		if s.Contains(input, element) {
			fmt.Print("Error found!\n")
			val = true;
		} else {
			fmt.Print("none\n")
		}
	}
	return val
}

func main() {
	t := time.Now()

	fmt.Println("== Starting data checker ==")
	fmt.Printf("%d-%02d-%02d %d:%d:%d\n\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	dat, err := ioutil.ReadFile("test.txt")
	checkFile(err)
	fmt.Printf("%s\n\n", string(dat))

	tags := []string{"ERROR", "Error", "error", "Warning"}

	if checkContains(string(dat), tags) {
		fmt.Println("ERROR FOUND")
	} else {
		fmt.Println("NO ERRORS FOUND")
	}

	// TXT file parse concept
}
