package main

// Import package dan API yang diperlukan
import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Fungsi untuk mencari semua external links yang mungkin menghubungan page asal dan page tujuan
// Mengembalikan daftar page dan pesan error
func getLink(page string) ([]string, error) {
	response, error := http.Get("https://en.wikipedia.org/wiki/" + strings.ReplaceAll(page, " ", "_"))
	// Mengecek page dengan alamat tersebut ada atau tidak
	if error != nil {
		return nil, error
	}
	defer response.Body.Close()
	// Jika gagal mengambil link
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch page: %s", response.Status)
	}
	// Parsing konten HTML dari HTTP response body menggunakan API goquery
	document, error := goquery.NewDocumentFromReader(response.Body)
	// Jika parsing gagal
	if error != nil {
		return nil, error
	}
	// Filtering anchor elements, menyimpan semua page Wikipedia yang merupakan artikel (educated guess)
	var links []string
	document.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists && link != "" && strings.HasPrefix(link, "/wiki/") {
			if !strings.Contains(link, ":") && !strings.Contains(link, "#") && !strings.Contains(link, "{") && !strings.Contains(link, "Main_Page") && !strings.Contains(link, "/File:") && !strings.Contains(link, "/Special") && !strings.Contains(link, "/Template") && !strings.Contains(link, "/Template_page:") && !strings.Contains(link, "/Help:") && !strings.Contains(link, "/Category:") && !strings.Contains(link, "Special:") && !strings.Contains(link, "/Wikipedia:") && !strings.Contains(link, "/Portal:") && !strings.Contains(link, "/Talk:") && !strings.Contains(link, "_(identifier)") {
				links = append(links, strings.TrimPrefix(link, "/wiki/"))
			}
		}
	})
	return links, nil
}

// Fungsi untuk mencari semua path yang bisa diambil dari page tertentu ke page asal atau page tujuan
// Mengembalikan daftar path
func getPaths(pageNames []string, visitedDict map[string][]string) [][]string {
	paths := [][]string{}

	for _, pageName := range pageNames {
		if pageName == "" {
			// Jika page yang sedang dicek kosong, page tersebut adalah page asal atau page tujuan, kembalikan path kosong
			return [][]string{{}}
		} else {
			// Jika page tidak kosong, cari path secara rekursif untuk child node dari page yang sedang dicek, append ke path
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

// Fungsi untuk mengecek apakah daftar path mengandung path tertentu
// Mengembalikan boolean
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

// Fungsi untuk melakukan Breadth First Search secara bi-direksional
// Mengembalikan daftar path terpendek untuk mencapai page tujuan dari page asal
func BFS(source, target string) [][]string {
	// Jika page asal dan tujuan sama, kembalikan page tersebut
	if source == target {
		return [][]string{{source}}
	}
	paths := [][]string{}
	// Dictionary yang berisi mapping nama page ke daftar nama page parent dari page tersebut
	// Berisi page yang belum dikunjungi dalam proses pencarian
	unvisitedForward := map[string][]string{source: {""}}
	unvisitedBackward := map[string][]string{target: {""}}
	// Sama seperti di atas, tetapi berisi page yang sudah dikunjungi
	visitedForward := map[string][]string{}
	visitedBackward := map[string][]string{}

	// Dilakukan proses pathfinding sampai ditemukan semua path terpendek atau semua page telah dikunjungi
	for len(paths) == 0 && (len(unvisitedForward) != 0 && len(unvisitedBackward) != 0) {
		// Melakukan iterasi berikutnya pada arah dengan jumlah link terkecil di tingkat berikutnya
		forwardDepth := len(visitedForward)
		backwardDepth := len(visitedBackward)

		if forwardDepth < backwardDepth {
			// BFS arah maju (dari page asal ke page tujuan)
			// Mengambil link dari page yang belum dikunjungi dengan arah maju
			outgoingLinks, error := getLink(source)
			if error != nil {
				log.Printf("error fetching links for %s: %v", source, error)
				continue
			}
			// Tandai semua page yang belum dikunjungi sebagai telah dikunjungi
			for page, parents := range unvisitedForward {
				visitedForward[page] = parents
			}
			// Kosongkan dictionary page yang belum dikunjungi
			unvisitedForward = map[string][]string{}
			for _, targetPage := range outgoingLinks {
				// Jika page tujuan tidak ada pada dictionary telah dikunjungi atau belum dikunjungi, masukkan ke belum dikunjungi
				if _, ok := visitedForward[targetPage]; !ok {
					unvisitedForward[targetPage] = []string{source}
				}
				// Jika page tujuan ada pada dictionary belum dikunjungi, tambahkan page asal sebagai salah satu parent nodenya
				if parents, ok := unvisitedForward[targetPage]; ok {
					unvisitedForward[targetPage] = append(parents, source)
				}
			}
		} else {
			// BFS arah mundur (dari page tujuan ke page asal)
			// Mengambil link dari page yang belum dikunjungi dengan arah mundur

			incomingLinks, error := getLink(target)
			if error != nil {
				log.Printf("error fetching links for %s: %v", target, error)
				continue
			}
			// Tandai semua page yang belum dikunjungi sebagai telah dikunjungi
			for page, parents := range unvisitedBackward {
				visitedBackward[page] = parents
			}
			// Kosongkan dictionary page yang belum dikunjungi
			unvisitedBackward = map[string][]string{}

			for _, sourcePage := range incomingLinks {
				// Jika page tujuan tidak ada pada dictionary telah dikunjungi atau belum dikunjungi, masukkan ke belum dikunjungi
				if _, ok := visitedBackward[sourcePage]; !ok {
					unvisitedBackward[sourcePage] = []string{target}
				}
				// Jika page tujuan ada pada dictionary belum dikunjungi, tambahkan page asal sebagai salah satu parent nodenya
				if parents, ok := unvisitedBackward[sourcePage]; ok {
					unvisitedBackward[sourcePage] = append(parents, target)
				}
			}
		}
		// Lakukan pengecekan penyelesaian pencarian
		// Pencarian selesai jika salah satu page terdapat di belum dikunjungi arah maju dan juga belum dikunjungi arah mundur
		// Selanjutnya, cari semua path terpendek
		for page := range unvisitedForward {
			if _, ok := unvisitedBackward[page]; ok {
				pathsFromSource := getPaths(unvisitedForward[page], visitedForward)
				pathsFromTarget := getPaths(unvisitedBackward[page], visitedBackward)

				for _, pathFromSource := range pathsFromSource {
					for _, pathFromTarget := range pathsFromTarget {
						currentPath := append(pathFromSource, page)
						currentPath = append(currentPath, pathFromTarget...)
						// Tambahkan path yang sedang dicek ke daftar path jika belum ada
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

// Fungsi untuk melakukan Iterative Deepening Search
// Depth dibatasi sampai 10 degree
func IDS(source, target string) []string {
	// Iterasi untuk setiap depth
	for depth := 1; depth <= 10; depth++ {
		visited := make(map[string]bool)
		path := make([]string, 0)
		// Cari path ke target page
		found := DLS(source, target, depth, visited, &path)
		// Jika path valid ditemukan
		if found {
			return path
		}
	}
	// Jika tidak
	return nil
}

// Fungsi untuk melakukan Depth Limited Search di setiap iterasi IDS
func DLS(page, target string, depth int, visited map[string]bool, path *[]string) bool {
	// Jika limit depth dicapai dan page yang sedang dicek bukan page target, kembalikan false
	if depth == 0 && page != target {
		return false
	}
	// Jika page yang sedang dicek merupakan page target, tambahkan dalam path, kembalikan true
	if page == target {
		*path = append(*path, page)
		return true
	}
	// Jika page yang sedang dicek sudah dikunjungi, kembalikan false
	if _, ok := visited[page]; ok {
		return false
	}
	// Tandai page yang sedang dicek sebagai sudah dikunjungi, ambil page lain yang terhubung melalui link
	visited[page] = true
	links, error := getLink(page)
	// Validasi
	if error != nil {
		log.Printf("error fetching links for %s: %v", page, error)
		return false
	}
	// Pemanggilan secara rekursif setiap page yang terhubung, dengan depth dikurangi setiap iterasi
	for _, link := range links {
		// Jika nilai true dikembalikan, berarti page target ditemukan, tambahkan dalam path
		if DLS(link, target, depth-1, visited, path) {
			*path = append(*path, page)
			return true
		}
	}
	return false
}
