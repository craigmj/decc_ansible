package decc

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type GemFile struct {
	*Gem
	Filename string
}

var gemFileRegexp = regexp.MustCompile(`^([^-]+)-(\S+)\.gem`)

func GemFiles(dirname string) ([]*GemFile,error) {
	gems := make([]*GemFile, 0)
	dir, err := os.Open(dirname)
	if nil!=err {
		return nil, err
	}
	defer dir.Close()
	files, err := dir.Readdir(0)
	if nil!=err {
		return nil, err
	}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".gem") && strings.HasPrefix(f.Name(), "decc_") {
			matches := gemFileRegexp.FindStringSubmatch(f.Name())
			if nil==matches {
				fmt.Fprintln(os.Stderr, "FAILED TO MATCH:", f.Name())
				continue
			}
			gems = append(gems, &GemFile{ Gem:  &Gem{Name:matches[1], Version:matches[2]}, Filename: f.Name()} )
		}
	}
	return gems, nil
}

func WhichGem(dirname string) (*GemFile, error) {
	gems, err := GemFiles(dirname)
	if nil!=err {
		return nil, err
	}
	var g *GemFile
	for _, c := range gems {
		if nil==g || g.Gem.EarlierThan(c.Gem) {
			g=c
		}
	}
	return g, nil
}

