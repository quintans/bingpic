package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/quintans/bingpic/internal/files"
)

const (
	bingURL = "https://www.bing.com/HPImageArchive.aspx?format=js&idx=%d&n=%d"
	maxRows = 8
)

// Get BingXML file which contains the URL of the Bing Photo of the day
// idx = Number days previous the present day. 0 means current day, 1 means yesterday, etc
// n = Number of images previous the day given by idx
// mkt denotes your location. e.g. en-US means United States. Put in your country code

type BingResponse struct {
	Images []struct {
		URL     string `json:"url"`
		UrlBase string `json:"urlbase"`
	} `json:"images"`
}

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
	if n > 16 {
		fmt.Println("Number of images must be less or equal than 16.")
		return
	}

	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}

	calls := n / maxRows
	for call := 0; call < calls; call++ {
		err := downloadPage(dest, call*maxRows, maxRows)
		if err != nil {
			fmt.Println("Error downloading page:", err)
			return
		}
	}

	remainder := n % maxRows
	if remainder > 0 {
		err := downloadPage(dest, calls*maxRows+2, remainder)
		if err != nil {
			fmt.Println("Error downloading remainder page:", err)
			return
		}
	}

	err = files.RollOver(dest, roll)
	if err != nil {
		fmt.Println("Error rolling over:", err)
		return
	}
}

func downloadPage(dest string, idx int, n int) error {
	url := fmt.Sprintf(bingURL, idx, n)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("fetching Bing image: %w", err)
	}
	defer resp.Body.Close()

	var bingResp BingResponse
	if err := json.NewDecoder(resp.Body).Decode(&bingResp); err != nil {
		return fmt.Errorf("Error decoding JSON response: %w", err)
	}

	if len(bingResp.Images) == 0 {
		fmt.Println("No images found in the response")
		return nil
	}

	for _, image := range bingResp.Images {
		filename := getFilename(image.UrlBase)
		filename = filepath.Join(dest, filename)
		if files.Exists(filename) {
			fmt.Printf("(Skip) Picture already downloaded: %s\n", filename)
			continue
		}

		err := files.DownloadImage("https://www.bing.com"+image.URL, filename)
		if err != nil {
			return fmt.Errorf("downloading page: %w", err)
		}

		fmt.Printf("Picture downloaded successfully!: %s\n", filename)
	}

	return nil
}

func getFilename(url string) string {
	splits := strings.Split(url, "=")
	name := splits[1]
	idx := strings.Index(name, ".")
	if idx != -1 {
		name = name[idx+1:]
	}

	idx = strings.LastIndex(name, "_")
	if idx != -1 {
		name = name[:idx]
	}

	return name + ".jpg"
}
