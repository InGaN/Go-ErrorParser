package main

import (
	"bufio"
	"fmt"
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

func checkContains(tags []string) bool {
	var val bool = false
	var amount int = 0;
	
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input := scanner.Text()
		fmt.Println(input)
		for index, element := range tags {
			if s.Contains(input, element) {
				val = true
				fmt.Printf("idx: %d", index)
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

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	
	t := time.Now()

	fmt.Println("== Starting data checker ==")
	fmt.Printf("%d-%02d-%02d %d:%d:%d\n\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	

	tags := []string{"ERROR", "Error", "error", "Warning"}

	if checkContains(tags) {
		fmt.Println("ERROR FOUND")
	} else {
		fmt.Println("NO ERRORS FOUND")
	}

	// TXT file parse concept
}
