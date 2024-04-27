package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func fetch(page string) ([]string, error) {
	response, error := http.Get("https://en.wikipedia.org/wiki/" + strings.ReplaceAll(page, " ", "_"))
	if error != nil {
		return nil, error
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch page: %s", response.Status)
	}

	doc, error := goquery.NewDocumentFromReader(response.Body)
	if error != nil {
		return nil, error
	}

	var links []string
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists && strings.HasPrefix(link, "/wiki/") {
			links = append(links, strings.TrimPrefix(link, "/wiki/"))
		}
	})

	return links, nil
}

// breadthFirstSearch returns a list of shortest paths from the source to target pages by running a
// bi-directional breadth-first search.
func BFS(sourcePage, targetPage string) [][]string {
	// If the source and target pages are identical, return the trivial path.
	if sourcePage == targetPage {
		return [][]string{{sourcePage}}
	}

	paths := [][]string{}

	// The unvisited dictionaries are a mapping from page name to a list of that page's parents' names.
	// Empty string signifies that the source and target pages have no parent.
	unvisitedForward := map[string][]string{sourcePage: {""}}
	unvisitedBackward := map[string][]string{targetPage: {""}}

	// The visited dictionaries are a mapping from page name to a list of that page's parents' names.
	visitedForward := map[string][]string{}
	visitedBackward := map[string][]string{}

	// Continue the breadth-first search until a path has been found or either of the unvisited lists
	// are empty.
	for len(paths) == 0 && (len(unvisitedForward) != 0 && len(unvisitedBackward) != 0) {
		// Run the next iteration of the breadth-first search in whichever direction has the smaller number
		// of links at the next level.
		forwardDepth := len(visitedForward)
		backwardDepth := len(visitedBackward)

		if forwardDepth < backwardDepth {
			// FORWARD BREADTH-FIRST SEARCH
			// Fetch the links from the currently unvisited forward pages.
			outgoingLinks, err := fetch(sourcePage)
			if err != nil {
				log.Printf("error fetching links for %s: %v", sourcePage, err)
				continue
			}

			// Mark all of the unvisited forward pages as visited.
			for page, parents := range unvisitedForward {
				visitedForward[page] = parents
			}

			// Clear the unvisited forward dictionary.
			unvisitedForward = map[string][]string{}

			for _, targetPage := range outgoingLinks {
				// If the target page is in neither visited forward nor unvisited forward, add it to
				// unvisited forward.
				if _, ok := visitedForward[targetPage]; !ok {
					unvisitedForward[targetPage] = []string{sourcePage}
				}

				// If the target page is in unvisited forward, add the source page as another one of its
				// parents.
				if parents, ok := unvisitedForward[targetPage]; ok {
					unvisitedForward[targetPage] = append(parents, sourcePage)
				}
			}
		} else {
			// BACKWARD BREADTH-FIRST SEARCH
			// Fetch the links to the currently unvisited backward pages.
			incomingLinks, err := fetch(targetPage)
			if err != nil {
				log.Printf("error fetching links for %s: %v", targetPage, err)
				continue
			}

			// Mark all of the unvisited backward pages as visited.
			for page, parents := range unvisitedBackward {
				visitedBackward[page] = parents
			}

			// Clear the unvisited backward dictionary.
			unvisitedBackward = map[string][]string{}

			for _, sourcePage := range incomingLinks {
				// If the source page is in neither visited backward nor unvisited backward, add it to
				// unvisited backward.
				if _, ok := visitedBackward[sourcePage]; !ok {
					unvisitedBackward[sourcePage] = []string{targetPage}
				}

				// If the source page is in unvisited backward, add the target page as another one of its
				// parents.
				if parents, ok := unvisitedBackward[sourcePage]; ok {
					unvisitedBackward[sourcePage] = append(parents, targetPage)
				}
			}
		}

		// CHECK FOR PATH COMPLETION
		// The search is complete if any of the pages are in both unvisited backward and unvisited, so
		// find the resulting paths.
		for page := range unvisitedForward {
			if _, ok := unvisitedBackward[page]; ok {
				pathsFromSource := getPaths(unvisitedForward[page], visitedForward)
				pathsFromTarget := getPaths(unvisitedBackward[page], visitedBackward)

				for _, pathFromSource := range pathsFromSource {
					for _, pathFromTarget := range pathsFromTarget {
						currentPath := append(pathFromSource, page)
						currentPath = append(currentPath, pathFromTarget...)

						// Add the current path to the list of paths if it's not already included
						if !containsPath(paths, currentPath) {
							paths = append(paths, currentPath)
						}
					}
				}
			}
		}
	}

	return paths
}

// getPaths returns a list of paths which go from the provided pages to either the source or target pages.
func getPaths(pageNames []string, visitedDict map[string][]string) [][]string {
	paths := [][]string{}

	for _, pageName := range pageNames {
		if pageName == "" {
			// If the current page name is empty, it is either the source or target page, so return an empty path.
			return [][]string{{}}
		} else {
			// Otherwise, recursively get the paths for the current page's children and append them to paths.
			currentPaths := getPaths(visitedDict[pageName], visitedDict)
			for _, currentPath := range currentPaths {
				newPath := append([]string{}, currentPath...)
				newPath = append(newPath, pageName)
				paths = append(paths, newPath)
			}
		}
	}

	return paths
}

// containsPath checks if the list of paths contains the specified path.
func containsPath(paths [][]string, path []string) bool {
	for _, p := range paths {
		if len(p) == len(path) {
			equal := true
			for i := range p {
				if p[i] != path[i] {
					equal = false
					break
				}
			}
			if equal {
				return true
			}
		}
	}
	return false
}

// func test() {
// 	// Define the source and target page names
// 	sourcePage := "Apple"
// 	targetPage := "Red"

// 	// Call the breadth-first search algorithm to find the shortest paths
// 	paths := BFS(sourcePage, targetPage)

// 	// Print the resulting paths
// 	fmt.Println("Shortest paths:")
// 	for _, path := range paths {
// 		fmt.Println(path)
// 	}
// }
