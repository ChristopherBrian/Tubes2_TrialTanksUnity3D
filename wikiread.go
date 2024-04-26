package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var temp *template.Template

func init() {
	temp = template.Must(template.ParseGlob("html/*.html"))
}

func main() {
	type listLink struct {
		Link string
	}

	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "mainpage.html", nil)
	})
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		startpoint := r.FormValue("fstart")
		endpoint := r.FormValue("fend")

		webstartpoint := "wikipedia.org/" + startpoint
		webendpoint := "wikipedia.org/" + endpoint

		listOfLink := []listLink{}

		st := colly.NewCollector(colly.AllowedDomains("wikipedia.org"))
		ed := colly.NewCollector(colly.AllowedDomains("wikipedia.org"))

		st.OnHTML("article", func(h *colly.HTMLElement) {
			metaTags := h.DOM.ParentsUntil("~").Find("meta")
			metaTags.Each(func(_ int, s *goquery.Selection) {
				x := listLink{
					Link: h.ChildAttr("a", "href"),
				}
				listOfLink = append(listOfLink, x)
			})
		})

		ed.OnHTML("article", func(h *colly.HTMLElement) {
			metaTags := h.DOM.ParentsUntil("~").Find("meta")
			metaTags.Each(func(_ int, s *goquery.Selection) {

			})
		})

		st.OnHTML("a[href]", func(h *colly.HTMLElement) {
			link := h.Attr("href")
			st.Visit(h.Request.AbsoluteURL(link))
		})

		ed.OnHTML("a[href]", func(h *colly.HTMLElement) {
			link := h.Attr("href")
			st.Visit(h.Request.AbsoluteURL(link))
		})

		st.Limit(&colly.LimitRule{
			DomainGlob:  endpoint,
			RandomDelay: 1 * time.Second,
		})

		ed.Limit(&colly.LimitRule{
			DomainGlob:  startpoint,
			RandomDelay: 1 * time.Second,
		})

		st.Visit(webstartpoint)
		ed.Visit(webendpoint)
	})
	http.ListenAndServe(":9999", nil)
}
