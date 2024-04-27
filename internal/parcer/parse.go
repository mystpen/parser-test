package parcer

import (
	"log"
	"net/http"

	"github.com/mystpen/parser-test/internal/model"
	"golang.org/x/net/html"
)

func Parse(resp *http.Response) (*[]model.InfluencerInfo, error) {
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
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
	// recursive call
	processAll(doc)

	return &Influencers, nil
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
