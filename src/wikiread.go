package src

import (
	"html/template"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var temp *template.Template

func init(){
	temp = template.Must(template.ParseGlob("html/*.html"))
}

func handleFunc(w http.ResponseWriter, r *http.Request){
	temp.ExecuteTemplate(w, "mainpage.html",nil)
}

func processor(w http.ResponseWriter, r *http.Request){
	if r.Method != "Post"{
		http.Redirect(w,r,"/",http.StatusSeeOther)
		return
	}

	startpoint := r.FormValue("fstart")
	endpoint := r.FormValue("fend")
	
	webstartpoint := "wikipedia.org/" + startpoint
	webendpoint := "wikipedia.org/" + endpoint

	st := colly.NewCollector(colly.AllowedDomains(webstartpoint))
	ed := colly.NewCollector(colly.AllowedDomains(webendpoint))

	st.OnHTML("article", func(h *colly.HTMLElement) {
		metaTags := h.DOM.ParentsUntil("~").Find("meta")
		metaTags.Each(func(_ int, s *goquery.Selection) {

		})
	})

	ed.OnHTML("article", func(h *colly.HTMLElement) {
		metaTags := h.DOM.ParentsUntil("~").Find("meta")
		metaTags.Each(func(_ int, s *goquery.Selection) {

		})
	})

	st.OnHTML("a[href]",func(h *colly.HTMLElement) {
		link := h.Attr("href")
		st.Visit(h.Request.AbsoluteURL(link))
	})

	ed.OnHTML("a[href]",func(h *colly.HTMLElement) {
		link := h.Attr("href")
		st.Visit(h.Request.AbsoluteURL(link))
	})

	st.Limit(&colly.LimitRule{
		DomainGlob: "*",
		RandomDelay: 1 * time.Second,
	})

	ed.Limit(&colly.LimitRule{
		DomainGlob: "*",
		RandomDelay: 1 * time.Second,
	})

	st.Visit(webstartpoint)
	ed.Visit(webendpoint)
}

func main(){
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/",http.StripPrefix("/css", fs))
	http.HandleFunc("/",handleFunc)
	http.HandleFunc("/process",processor)
	http.ListenAndServe(":9999",nil)
}