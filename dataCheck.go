package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	s "strings"
	"time"
)

func checkFile(e error) {
	if e != nil {
		panic(e)
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func checkContainsScanner(tags []string) bool {
	defer timeTrack(time.Now(), "Scanner")
	var val bool = false
	var amount int = 0

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input := scanner.Text()
		for _, element := range tags {
			if s.Contains(input, element) {
				val = true
				//fmt.Printf("tag: %s | %s \n\n", tags[index], input)
				amount++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\namount of errors: %d\n", amount)

	return val
}

func checkContainsReadString(tags []string) bool {
	defer timeTrack(time.Now(), "ReadString")
	var val bool = false
	var amount int = 0

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
		//return &stats, err
	}

	r := bufio.NewReader(file)
	for {
		recordRaw, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		for _, element := range tags {
			if s.Contains(recordRaw, element) {
				val = true
				//fmt.Printf("tag: %s | %s \n\n", tags[index], recordRaw)
				amount++
			}
		}
	}
	fmt.Printf("\namount of errors: %d\n", amount)
	file.Close()

	return val
}

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	t := time.Now()

	fmt.Println("== Starting data checker ==")
	fmt.Printf("%d-%02d-%02d %d:%d:%d\n\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	tags := []string{"ERROR", "Error", "error", "Warning"}

	if checkContainsReadString(tags) {
		fmt.Println("ERRORS FOUND")
	} else {
		fmt.Println("NO ERRORS FOUND")
	}
	if checkContainsScanner(tags) {
		fmt.Println("ERRORS FOUND")
	} else {
		fmt.Println("NO ERRORS FOUND")
	}
	if checkContainsReadString(tags) {
		fmt.Println("ERRORS FOUND")
	} else {
		fmt.Println("NO ERRORS FOUND")
	}
	if checkContainsScanner(tags) {
		fmt.Println("ERRORS FOUND")
	} else {
		fmt.Println("NO ERRORS FOUND")
	}
	if checkContainsReadString(tags) {
		fmt.Println("ERRORS FOUND")
	} else {
		fmt.Println("NO ERRORS FOUND")
	}
	// TXT file parse concept
}
