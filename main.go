package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

// WebData holds the structure of the data we want to scrape
type WebData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	// Instantiate the collector
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	// Define a list of URLs to crawl
	urls := []string{
		"https://en.wikipedia.org/wiki/Robotics",
		"https://en.wikipedia.org/wiki/Robot",
		"https://en.wikipedia.org/wiki/Reinforcement_learning",
		"https://en.wikipedia.org/wiki/Robot_Operating_System",
		"https://en.wikipedia.org/wiki/Intelligent_agent",
		"https://en.wikipedia.org/wiki/Software_agent",
		"https://en.wikipedia.org/wiki/Robotic_process_automation",
		"https://en.wikipedia.org/wiki/Chatbot",
		"https://en.wikipedia.org/wiki/Applications_of_artificial_intelligence",
		"https://en.wikipedia.org/wiki/Android_(robot)",
		// ... other URLs
	}

	// Create a WebData slice to hold our results
	var results []WebData

	// Callback for when content is found
	c.OnHTML("#content", func(e *colly.HTMLElement) {
		title := e.ChildText("h1#firstHeading")
		description := e.ChildText("div.mw-parser-output > p:not([class])")

		// Filtering out empty titles and descriptions
		if title != "" && description != "" {
			if strings.Contains(description, "robotics") || strings.Contains(description, "intelligent systems") {
				results = append(results, WebData{
					Title:       title,
					Description: strings.TrimSpace(description),
				})
			}
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Iterate over the slice of URLs
	for _, url := range urls {
		c.Visit(url)
	}

	// Wait until all threads are finished
	c.Wait()

	// Write results to a .json file
	file, err := os.Create("results.json")
	if err != nil {
		fmt.Println("Could not create file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(results)
	if err != nil {
		fmt.Println("Could not encode results to JSON:", err)
		return
	}

	fmt.Println("Scraping finished, check results.json for output")
}
