package main

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