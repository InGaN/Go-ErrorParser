package main

import (
	"flag"
	"fmt"
	tagparser "github.com/InGaN/Go-TagParser"
	//"github.com/InGaN/Go-TagParser/slack"
)

var (	
	flagHelp1	 = flag.Bool("h", false, "help")
	flagHelp2	 = flag.Bool("help", false, "help")
	flagSearchR	 = flag.Bool("r", false, "Recursive search")
	flagFileDir	 = flag.String("f", "logs", "Directory containing files to parse")
	flagTags	 = flag.String("t", "tags.txt", "file containing tags")	
	flagPrint	 = flag.Bool("p", false, "Print in terminal")
	//flagSlack	 = flag.Bool("slack", false, "Send error messages to Slack")
)

func showHelp() {
	fmt.Println("This Go program parses plain text files containing certain tags.")
	fmt.Println("Paramaters that can be used:")
	fmt.Println("-f=<path>\t\tPath to folder to search in")
	fmt.Println("-h, -help\t\tShow help")
	fmt.Println("-p\t\t\tPrint messages in terminal")
	fmt.Println("-r\t\t\tRecursive search")
	fmt.Println("-t=<path>\t\tPath to csv file containing tags")
}

func main() {	
	flag.Parse()
	if(*flagHelp1 || *flagHelp2) {
		showHelp()
	} else {
		//response := tagparser.Parse(false, false, false, "logs", "tags.txt")
		response := tagparser.Parse(*flagSearchR, *flagPrint, *flagFileDir, *flagTags)
		fmt.Println(string(response))
		
		//if(*flagSlack) {
		//	slack.SendJSONtoSlack(fmt.Sprintf("tag: %s | %s \n", tags[index], input))
		//}
	}
}