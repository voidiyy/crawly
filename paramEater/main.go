package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"paramEater/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gookit/color"
)

var (
	domains       []string
	depth         int
	outputFile    string
	valid         []string
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: program <domain1,domain2,...> <depth> <output path>")
	}

	domains = strings.Split(os.Args[1], ",")
	depth, _ = strconv.Atoi(os.Args[2])
	outputFile = os.Args[3]

	color.Info.Println("Starting crawl with the following parameters:")
	color.Info.Println("Domains:", domains)
	color.Info.Println("Depth:", depth)
	color.Info.Println("Output file:", outputFile)

	fmt.Println("please pres enter...")
	fmt.Scanln()

	crawlDomains()

	color.Info.Println("Crawl completed. Writing results to file...")
	color.Info.Println("Results written successfully.")
}

func loadLines(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading file:", err)
	}
	return lines
}

func crawlDomains() {
	for _, domain := range domains {
		c := colly.NewCollector(
			colly.Async(true),
			colly.MaxDepth(depth),
		)

		c.WithTransport(&http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			DisableKeepAlives:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		})

		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("User-Agent", utils.RandomAgent())
			r.Headers.Set("Accept", "*/*")
			color.Info.Printf("Visiting %s\n", r.URL.String())
		})

		c.OnResponse(func(r *colly.Response) {
			color.Debug.Printf("Received response from %s\n", r.Request.URL.String())
		})

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			valid = append(valid, link)
			_ = e.Request.Visit(link)
		})

		color.Info.Printf("Starting crawl for domain: %s\n", domain)
		err := c.Visit("https://" + domain)
		if err != nil {
			color.Error.Printf("Visit error: %s\n", err)
		}

		c.Wait()
	}
	err := writeToFile(outputFile, valid)
	if err != nil {
		log.Fatal(err)
	}
}

func writeToFile(filename string, data []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	for _, line := range data {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	fmt.Printf("%v URLs successfully written to %v\n", len(data), filename)
	return nil
}
