package main

import (
	"fmt"
	"os"
	"strings"

	"decc"
)

var helpText = `deccgem v0.0.1
USAGE: deccgem CMD

COMMANDS:
find		Return the latest DECC gem installed
gemfile 	Write the twenty-fifty Gemfile to Stdout
whichgem 	Find the latest .gem file in current directory 
`

func main() {
	cmd := "help"
	if 2 <= len(os.Args) {
		cmd = strings.ToLower(os.Args[1])
		if 3 == len(os.Args) {
			decc.SetGemCmd(os.Args[2])
		}
	}

	switch cmd {
	case "find":
		fmt.Println(decc.TheDeccGem())
	case "gemfile":
		decc.WriteGemFile(os.Stdout)
	case "help":
		fmt.Println(helpText)
	case "whichgem":
		wdir, err := os.Getwd()
		if nil != err {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		g, err := decc.WhichGem(wdir)
		if nil != err {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		if nil == g {
			fmt.Fprintln(os.Stderr, "No .gem files found")
			os.Exit(1)
		}
		fmt.Println(g.Filename)
	case "parse":
		decc.GemLister()
	}
}
