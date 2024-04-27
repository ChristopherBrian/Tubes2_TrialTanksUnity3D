package main

import (
	"html/template"
	"net/http"
	"time"
)

var temp *template.Template

func init() {
	temp = template.Must(template.ParseGlob("html/*.html"))
}

type ViewData struct {
	Paths                [][]string
	Paths2               [][]string
	TotalPagesVisited    int
	TotalPagesVisited2   int
	ShortestPathDepth    int
	ShortestPathDepth2   int
	SearchTimeInSeconds  float64
	SearchTimeInSeconds2 float64 // Separate variable for IDS search time
}

func main() {
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Render the mainpage.html template with an empty ViewData
		temp.ExecuteTemplate(w, "mainpage.html", ViewData{})
	})
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		startTime := time.Now() // Record start time

		startpoint := r.FormValue("start")
		endpoint := r.FormValue("end")

		// Define the source and target page names
		sourcePage := startpoint
		targetPage := endpoint

		// Call the breadth-first search algorithm to find the shortest paths
		paths := BFS(sourcePage, targetPage)

		// Prepare the BFS paths for display in HTML
		var formattedPaths [][]string
		for _, path := range paths {
			formattedPaths = append(formattedPaths, path)
		}
		searchTime := time.Since(startTime).Seconds()
		startTime2 := time.Now()
		// Call the iterative deepening search algorithm to find the shortest paths
		paths2 := IDS(sourcePage, targetPage)

		// Prepare the IDS paths for display in HTML
		var formattedPaths2 [][]string
		for _, path2 := range paths2 {
			formattedPaths2 = append(formattedPaths2, []string{path2}) // Convert path to []string and append to formattedPaths2
		}

		// Calculate the total amount of pages visited
		totalPagesVisited := len(paths)
		totalPagesVisited2 := len(paths2)

		// Determine the depth of the shortest path
		shortestPathDepth := len(paths[0]) - 1 // Subtract 1 to exclude the source node
		shortestPathDepth2 := len(paths2[0]) - 1

		// Calculate the time taken for the search operation
		searchTime2 := time.Since(startTime2).Seconds()

		// Combine both sets of paths into a single ViewData along with additional information
		data := ViewData{
			Paths:                formattedPaths,
			Paths2:               formattedPaths2,
			TotalPagesVisited:    totalPagesVisited,
			TotalPagesVisited2:   totalPagesVisited2,
			ShortestPathDepth:    shortestPathDepth,
			ShortestPathDepth2:   shortestPathDepth2,
			SearchTimeInSeconds:  searchTime,
			SearchTimeInSeconds2: searchTime2,
		}

		// Render the mainpage.html template with both sets of paths and additional information
		temp.ExecuteTemplate(w, "mainpage.html", data)
	})

	http.ListenAndServe(":9999", nil)
}
