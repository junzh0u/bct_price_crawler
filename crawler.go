package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Product struct {
	Page  string
	Name  string
	Price float64
	Time  string
}

func main() {
	now := time.Now().Format(time.RFC3339)
	products := []Product{}

	c := colly.NewCollector(
		colly.AllowedDomains("bridgecitytools.com"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML("span.gf_product-price.money", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Request.URL.Path, "/products/") {
			price, err := strconv.ParseFloat(strings.ReplaceAll(e.Text, "$", ""), 64)
			if err == nil {
				product := Product{
					Page: e.Request.URL.String(),
					Name: strings.TrimSpace(strings.ReplaceAll(
						e.DOM.ParentsFiltered(":root").Find("title").Text(),
						"– Bridge City Tool Works",
						"")),
					Price: price,
					Time:  now,
				}
				products = append(products, product)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.Visit("https://bridgecitytools.com/")

	b, err := json.Marshal(products)
	if err == nil {
		fmt.Println(string(b))
	}
}
