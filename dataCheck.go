package main

import (
	"bufio"
	"bytes"	
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	s "strings"
	"time"
)

var i = 0
var pTotalAmount *int = &i
var x = 0
var pTotalFilesTagged *int = &x
var a = 0
var pTotalFiles *int = &a
var b int64 = 0 //this all seems redundant
var pTotalFileSize *int64 = &b
var buffer bytes.Buffer
var pTags *[]string

var (	
	helpParam1 = flag.Bool("h", false, "help")
	helpParam2 = flag.Bool("help", false, "help")
	searchR = flag.Bool("r", false, "Recursive search")
	echoMsg = flag.Bool("e", false, "Echo messages containing tags")
	fileDir = flag.String("f", ".", "Directory containing files to parse")
	tagParam = flag.String("t", "tags.txt", "file containing tags")
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func visit(path string, f os.FileInfo, err error) error {
	//fmt.Printf("Visited: %s\n", path)
	arr := (s.Split(path,"."))
	mimeTxt := mime.TypeByExtension(".txt")
	if(len(arr)>1) {
		if(mime.TypeByExtension("."+arr[len(arr)-1]) == mimeTxt) {
			checkContainsScanner(path, *pTags)
		}
	}  
  return nil
} 

func checkFiles(path string) {	
	mimeTxt := mime.TypeByExtension(".txt")		
	if(*searchR) {
		err := filepath.Walk(path, visit)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		files, _ := ioutil.ReadDir(path)
		for _, file := range files {
			arr := (s.Split(file.Name(),"."))
			if(len(arr)>1) {
				if(mime.TypeByExtension("."+arr[len(arr)-1]) == mimeTxt) {
					*pTotalFiles++
					//fmt.Printf("Type: %s\n", file.Name())
					*pTotalFileSize = *pTotalFileSize + file.Size()
					checkContainsScanner(fmt.Sprintf("%s\\%s",path, file.Name()), *pTags)
				}				
			}
		}	
	}
	
}

func checkContainsScanner(path string, tags []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	changeFlag := *pTotalAmount
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input := scanner.Text()
		parseLine(input, tags)
	}	
	if(changeFlag < *pTotalAmount) {
		buffer.WriteString(fmt.Sprintf("%s (%d)\n", file.Name(), (*pTotalAmount-changeFlag)))
		*pTotalFilesTagged++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func checkContainsReadString(path string, tags []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	changeFlag := *pTotalAmount
	
	r := bufio.NewReader(file)
	for {
		recordRaw, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		parseLine(recordRaw, tags)
	}			
	if(changeFlag < *pTotalAmount) {
		buffer.WriteString(fmt.Sprintf("%s (%d)\n", file.Name(), (*pTotalAmount-changeFlag)))
		*pTotalFilesTagged++
	}
	file.Close()
}

func parseLine(input string, tags []string) {
	amount := 0
	for index, element := range tags {
		if s.Contains(input, element) {
			if(*echoMsg) {
				fmt.Printf("tag: %s | %s \n", tags[index], input)
			}
			amount++
		}
	}
	*pTotalAmount = *pTotalAmount + amount
}

func parseTagFile(path string) []string{
	extension := (s.Split(path,"."))
	if(mime.TypeByExtension("."+extension[len(extension)-1]) == mime.TypeByExtension(".txt")) {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		arr := s.Split(string(file), ",")		
		return arr
	}	
	return nil
}

func showHelp() {
	fmt.Println("This Go program parses plain text files containing certain tags.")
	fmt.Println("Paramaters that can be used:")
	fmt.Println("-e\t\t\techo messages containing tags")
	fmt.Println("-f=<path>\t\tpath to folder to search in")
	fmt.Println("-h, -help\t\tshow help")
	fmt.Println("-r\t\t\trecursive search")
	fmt.Println("-t=<path>\t\tpath to csv file containing tags")
}

func getByteSize(value int64) string {
	sizes := []string{"Byte", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	
	if(value < 0) {
		return "-"
	}
	i := 0
	for (value/1024 >= 1) {
		value = value / 1024
		i++
	}
	return fmt.Sprintf("%d%s", value, sizes[i])
}

func main() {
	flag.Parse()
	if(*helpParam1 || *helpParam2) {
		showHelp()
	} else {
		defer timeTrack(time.Now(), "Operation")
		if *fileDir == "" {
			fmt.Fprintln(os.Stderr, "require a folder")
			flag.Usage()
			os.Exit(1)
		}	
		fmt.Println("\nStarting parse...")		
		
		tags := parseTagFile(*tagParam)
		pTags = &tags
		
		checkFiles(*fileDir)
		
		t := time.Now()
		fmt.Println("\n=== FILES ===")
		fmt.Printf(buffer.String())
		fmt.Println("\n=== RESULTS ===")
		fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
		fmt.Println("Directory: ", *fileDir)
		fmt.Println("Tags file: ", *tagParam)
		fmt.Printf("Tags: %v\n", parseTagFile(*tagParam))
		fmt.Println("Recursive: ", *searchR)	
		fmt.Printf("Total files scanned: %d\n", *pTotalFiles)			
		fmt.Printf("total file size: %s\n", getByteSize(*pTotalFileSize))
		fmt.Printf("total amount of tags found: %d\n", *pTotalAmount)
		fmt.Printf("in %d files\n", *pTotalFilesTagged)					
	}
}
