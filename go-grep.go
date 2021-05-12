package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"strconv"

	"github.com/morgulbrut/color256"
)

var root, query string
var found = 1
var wg sync.WaitGroup

type result struct {
	root  string
	path  string
	lines []string
}

func readFile(wg *sync.WaitGroup, path string) {
	defer wg.Done()

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return
	}
	var res result

	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i++ {
		if strings.Contains(strings.ToUpper(scanner.Text()), strings.ToUpper(query)) {
			found = 0
			res.root = root
			res.path = path
			res.lines = append(res.lines, fmt.Sprintf("%s: %s", color256.Yellow(strconv.Itoa(i)), scanner.Text()))
		}
	}
	if len(res.lines) > 0 {
		color256.PrintHiGreen(res.path)
		for _, l := range res.lines {
			fmt.Printf("\t%s\n", l)
		}
	}

}

func main() {
	flag.Parse()
	query = flag.Arg(0)
	root = flag.Arg(1)

	if query == "" || root == "" {
		os.Exit(2)
	}

	filepath.Walk(root, func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			wg.Add(1)
			go readFile(&wg, path)
		}
		return nil
	})
	wg.Wait()
	defer os.Exit(found)
}
