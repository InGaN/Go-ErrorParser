# Go - Tag Parser (warning, contains some issues)

Module written to find tags or keywords in plain text files and JSON files (.txt and .json) and return output in JSON.

Tags can be passed as comma separated text file or as an array argument.

* Search directory for text files
* Recursive search option
* Concurrent searching
* Output can be printed in terminal
* Output returned as JSON byte[] object

## Issues
* only counts a tag per line in file
* channels sometimes lock

## Installation
```
go get github.com/InGaN/Go-TagParser
```

## Basic Usage
```
import (
	"fmt"
	tagparser "github.com/InGaN/Go-TagParser"
)

func main() {	
	// if an array with tags is used, the tags.txt file is omitted
	arr := []string{"Warning", "ERROR"}
	response, err := tagparser.Parse(false, false, "logs2", "tags.txt", arr)
	if(err != nil) {
		fmt.Print(err)
	}
	fmt.Println(string(response))	
}
```

## tags.txt
```
ERROR, Error, error, Warning, WARNING
```