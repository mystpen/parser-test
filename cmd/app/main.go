package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mystpen/parser-test/config"
	"github.com/mystpen/parser-test/internal/model"

	"golang.org/x/net/html"
)

var ErrNoResponce = errors.New("no responce from API")

func main() {
	// URL to make the HTTP request to
	// Make the GET request
	resp, err := http.Get(config.Config.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error:", "no responce from API")
		return
	}

	// // Read the response body
	// bytes, _ := io.ReadAll(resp.Body)

	// // Print the body as a string
	// fmt.Println("HTML:\n\n", string(bytes))

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var Influencers []model.InfluencerInfo

	var processAll func(*html.Node)
	processAll = func(n *html.Node) {
		if n.Type == html.ElementNode {
			
			if n.Data == "span" && HaveAttr(n.Attr, map[string]string{"data-v-b11c405a": ""}) {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						Influencers = append(Influencers, model.InfluencerInfo{})
						Influencers[len(Influencers)-1].Rank = c.Data
						
					}
				}
				return
			}

			if n.Data == "img" && HaveAttr(n.Attr, map[string]string{"data-v-c9cd5c3e": "", "src": ""}) {
				for _, a := range n.Attr {
					if a.Key == "src" {
						ImageURL := a.Val
						Influencers[len(Influencers)-1].AvatarImage = ImageURL

						
					}
				}
				return
			}

			if n.Data == "div" && HaveAttr(n.Attr, map[string]string{"class": "contributor__name-content"}) {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						Influencers[len(Influencers)-1].Account = c.Data
						
					}
				}
				return
			}
			if n.Data == "div" && HaveAttr(n.Attr, map[string]string{"class": "contributor__title"}) {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						Influencers[len(Influencers)-1].Name = c.Data
						
					}
				}
				return
			}
			if n.Data == "div" && HaveAttr(n.Attr, map[string]string{"class": "tag__content ellipsis"}) {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						Influencers[len(Influencers)-1].Category = append(Influencers[len(Influencers)-1].Category, c.Data)
					}
				}
				return
			}

			if n.Data == "div" && HaveAttr(n.Attr, map[string]string{"class": "row-cell subscribers"}) {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						Influencers[len(Influencers)-1].Subscribers = c.Data
						
					}
				}
				return
			}
			if n.Data == "div" && HaveAttr(n.Attr, map[string]string{"class": "row-cell audience"}) {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						Influencers[len(Influencers)-1].Country = c.Data
						
					}
				}
				return
			}

			if n.Data == "div" && HaveAttr(n.Attr, map[string]string{"class": "row-cell authentic"}) {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						Influencers[len(Influencers)-1].EngAuth = c.Data
						
					}
				}
				return
			}

			if n.Data == "div" && HaveAttr(n.Attr, map[string]string{"class": "row-cell engagement"}) {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						Influencers[len(Influencers)-1].EngAvg = c.Data
						
					}
				}
				return
			}

		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processAll(c)
		}
	}
	// make a recursive call to your function
	processAll(doc)
	fmt.Println(Influencers[len(Influencers)-1])

	///////////////////////////////////////////////////////////////////

	file, err := os.Create("influencers.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Rank", "Account", "Name", "Avatar Image", "Category", "Subscribers", "Country", "Eng. (Auth.)", "Eng. (Avg.)"}
	if err := writer.Write(header); err != nil {
		fmt.Println("Error writing header:", err)
		return
	}

	// Write data rows
	for _, info := range Influencers {
		record := []string{
			info.Rank,
			info.Account,
			info.Name,
			info.AvatarImage,
			strings.Join(info.Category, ", "),
			info.Subscribers,
			info.Country,
			info.EngAuth,
			info.EngAvg,
		}
		if err := writer.Write(record); err != nil {
			fmt.Println("Error writing record:", err)
			return
		}
	}
}

func HaveAttr(attr []html.Attribute, m map[string]string) bool {
	count := 0

	for _, a := range attr {
		if i, ok := m[a.Key]; (i == "" || i == a.Val) && ok {
			count++
			continue
		}
	}
	return count == len(m)
}

