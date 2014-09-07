package decc

import (
	// "bufio"
	"encoding/json"
	"fmt"
	"log"
	"io"
	"os"
	"os/exec"
	// "regexp"
	"strings"
	"text/template"

	// "lexer"
)

var _ = os.Stdout
var _ = exec.Command

var gemCmd = "gem"

func SetGemCmd(cmd string) {
	gemCmd = cmd
}

type Gem struct {
	Name    string
	Version string
	Lang    string
	Binary  string
}

func (g *Gem) String() string {
	data, _ := json.Marshal(g)
	return string(data)
}

func (g *Gem) EarlierThan(c *Gem) bool {
	if strings.HasPrefix(g.Version, "3") {
		if strings.HasPrefix(c.Version, "3") {
			return g.Version < c.Version
		}
		return true
	} else {
		if strings.HasPrefix(c.Version, "3") {
			return false
		}
		return g.Version < c.Version
	}
}

func GemLine(out chan *Gem) error {
	return nil
}

func GemList() (string, error) {
	cmd := exec.Command(gemCmd, "list")
	output, err := cmd.Output()
	if nil != err {
		fmt.Fprintln(os.Stderr, err.Error())
		return "", err
	}
	return string(output), nil
}

func GemLister() {

	gems, _ := GemList()
	for g := range NewGemReadChan(strings.NewReader(gems)) {
		if nil != g.Err {
			fmt.Println("ERROR:", g.Err)
		} else {
			fmt.Println(g.G)
		}
	}

	// lexer.DumpTokens(strings.NewReader(gems))
}

func Gems(out chan *Gem, filter func(*Gem) bool) error {
	defer close(out)
	gems, _ := GemList()
	// 	gems = `
	// *** LOCAL GEMS ***

	// bundler (1.3.5)
	// decc_2050_model (3.5.1pre, 0.71.20140319pre)
	// ffi (1.8.1)
	// net-http-persistent (2.8)
	// rdoc (3.9.4)
	// thor (0.18.1)
	// `

	for g := range NewGemReadChan(strings.NewReader(gems)) {
		if nil != g.Err {
			fmt.Fprintln(os.Stderr, "ERROR: ", g.Err)
		} else {
			if filter(g.G) {
				// fmt.Println("Gems(..) about to send: ", g.G.Name)
				out <- g.G
				// fmt.Println("Gems(..) sent")
			}
		}
	}
	return nil
}

func DeccGems(out chan *Gem) {
	Gems(out, func(g *Gem) bool {
		return strings.HasPrefix(g.Name, "decc_")
	})
}

func TheDeccGem() *Gem {
	var g *Gem
	C := make(chan *Gem)
	go DeccGems(C)
	for c := range C {
		// fmt.Println("TheDeccGem: have ", g, " and received ", c)
		if nil == g {
			g = c
		} else {
			if g.EarlierThan(c) {
				// fmt.Println(g, " EarlierThan ", c)
				g = c
			}
		}
		// fmt.Println("TDG: g=", g)
	}

	return g
}

func WriteGemFile(out io.Writer) {
	gem := TheDeccGem()
	if nil==gem {
		log.Printf("Failed to get theDeccGem")
		os.Exit(1)
	}
	err := gemfileTemplate.Execute(out, gem)
	if nil!=err {
		log.Printf(err.Error())
		os.Exit(1)
	}
}

var gemfileTemplate = template.Must(template.New("gemfile").Parse(`source 'http://rubygems.org'
# Framework
gem 'sinatra'

# Views
gem 'json'
gem 'sass'
gem 'haml'
gem 'uglifier'#, '>= 1.0.3'
gem 'sprockets'
gem 'coffee-script'

# Model
gem 'ffi'
gem '{{.Name}}', '{{.Version}}'

`))
