package files

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

func RollOver(dest string, roll int) error {
	files, err := glob(dest, filepath.Join(dest, "*.jpg"))
	if err != nil {
		return fmt.Errorf("listing files: %w", err)
	}

	if len(files) > roll {
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().After(files[j].ModTime())
		})

		for i := roll; i < len(files); i++ {
			fmt.Println("Removing file:", files[i].Name())
			err := os.Remove(filepath.Join(dest, files[i].Name()))
			if err != nil {
				return fmt.Errorf("removing file: %s :%w", files[i].Name(), err)
			}
		}
	}

	return nil
}

func glob(root, pattern string) ([]os.FileInfo, error) {
	files := make([]os.FileInfo, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		ok, err := filepath.Match(pattern, path)
		if err != nil {
			return fmt.Errorf("matching pattern: %w", err)
		}

		if ok {
			files = append(files, info)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking to %s: %w", pattern, err)
	}

	return files, nil
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func DownloadImage(imageURL string, filename string) error {
	resp, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("downloading image: %s: %w", imageURL, err)
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file '%s': %w", filename, err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("saving image of file '%s': %w", filename, err)
	}

	return nil
}
