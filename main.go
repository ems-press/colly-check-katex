package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

func main() {
	// parse command line arguments
	const usage = `Usage: check-katex [-u|--start-url ...] [-h|--help]

  -u, --start-url Starting point of the crawler (default: https://ems.press/journals)
  -h, --help      prints help information

Will only check URLs deeper than the given start URL. All errors are printed to
stderr, verbose request information is printed to stdout.

Examples:

   go run main.go > /dev/null # only print errors
   go run main.go 2>&1 | tee errors.log # save all errors to a file
`
	startUrl := "https://ems.press/journals"
	flag.StringVar(&startUrl, "start-url", startUrl, "Start point of the crawler")
	flag.StringVar(&startUrl, "u", startUrl, "Start point of the crawler")
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	// code 0: no errors found, code 1: katex errors found, code 255: encountered http errors
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	collector := colly.NewCollector(
		colly.UserAgent("ems.press check-katex"),
		colly.URLFilters(
			// only look at urls deeper than the given start url:
			regexp.MustCompile(startUrl+".*?"),
		),
		colly.Async(),
	)
	collector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 8})
	extensions.Referer(collector)

	// check all links
	collector.OnHTML("a[href]", func(element *colly.HTMLElement) {
		element.Request.Visit(element.Attr("href"))
	})

	// check all rendered formulae for errors
	collector.OnHTML("span.katex", func(element *colly.HTMLElement) {
		katexErrors := element.DOM.Find("span[style*=color]")
		if len(katexErrors.Nodes) > 0 {
			exitCode = 1
			latexSrc := element.DOM.Find("annotation[encoding=\"application/x-tex\"]").Text()
			fmt.Fprintf(os.Stderr, "Error: malformatted latex src $%v$ on URL %s\n", latexSrc, element.Request.URL)
		}
	})

	// print some info about visited pages to stdout
	collector.OnResponse(func(response *colly.Response) {
		fmt.Printf("Checked %s\n", response.Request.URL)
	})

	collector.OnError(func(response *colly.Response, err error) {
		if response.StatusCode == 503 || response.StatusCode == 999 || response.StatusCode == 0 {
			// ignore 503 and 999 and 0 status code to avoid flaky errors
			return
		}

		exitCode = 255
		request := response.Request
		fmt.Fprintf(os.Stderr, "Error: \"%v %s\" while visiting %s; Found on: %s\n", response.StatusCode, err, request.URL, request.Headers.Get("Referer"))
	})

	collector.Visit(startUrl)
	collector.Wait()
}
