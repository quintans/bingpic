package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/quintans/bingpic/internal/files"
)

const bingoURL = "https://bing.gifposter.com"

func main() {
	var n int
	flag.IntVar(&n, "n", 1, "Number of images previous today.")

	var roll int
	flag.IntVar(&roll, "roll", 100, "Number of files to keep in the folder.")

	var dest string
	flag.StringVar(&dest, "dest", ".", "Folder to save the images.")

	flag.Parse()

	if n < 1 {
		fmt.Println("Number of images must be greater than 0.")
		return
	}
	if n > 100 {
		fmt.Println("Number of images must be less or equal than 100.")
		return
	}

	if n > roll {
		fmt.Println("Number of images must be less or equal than roll.")
		return
	}

	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}

	doc, err := loadPage(bingoURL)
	if err != nil {
		fmt.Println("Error loading page:", err)
		return
	}

	// Find and print the title
	var url string
	doc.Find(`.dayimg`).
		Find(`[itemprop="url"]`).Each(func(i int, s *goquery.Selection) {
		attr, exists := s.Attr("href")
		if exists {
			url = attr
		}
	})

	if url == "" {
		fmt.Println("First URL not found")
		return
	}

	err = scrapeImage(bingoURL+url, dest, n)
	if err != nil {
		fmt.Println("Error scraping image:", err)
		return
	}

	err = files.RollOver(dest, roll)
	if err != nil {
		fmt.Println("Error rolling over:", err)
		return
	}
}

func scrapeImage(url, dest string, count int) error {
	doc, err := loadPage(url)
	if err != nil {
		return fmt.Errorf("loading page: %w", err)
	}

	// Find and print the title
	var img string
	doc.Find(`#bing_wallpaper`).Each(func(i int, s *goquery.Selection) {
		attr, exists := s.Attr("src")
		if exists {
			img = attr
		}
	})

	if img == "" {
		return fmt.Errorf("image not found: missing element with id 'bing_wallpaper'")
	}

	// img = strings.ReplaceAll(img, "_1920x1080", "_uhd")
	splits := strings.Split(img, "/")
	filename := filepath.Join(dest, splits[len(splits)-1])
	if !files.Exists(filename) {
		err := files.DownloadImage(img, filename)
		if err != nil {
			return fmt.Errorf("downloading image: %w", err)
		}
		fmt.Printf("Picture downloaded successfully!: %s\n", filename)
	} else {
		fmt.Printf("(Skip) Picture already downloaded: %s\n", filename)
	}

	count--
	if count == 0 {
		return nil
	}

	var next string
	doc.Find(`.icon.next`).Each(func(i int, s *goquery.Selection) {
		attr, exists := s.Attr("href")
		if exists {
			next = attr
		}
	})

	if next == "" {
		fmt.Println("No next URL found")
		return nil
	}

	return scrapeImage(bingoURL+next, dest, count)
}

func loadPage(url string) (*goquery.Document, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
