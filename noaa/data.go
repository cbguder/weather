package noaa

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

const (
	baseUrl = "https://www.ncei.noaa.gov/pub/data/ghcn/daily/"

	numWorkers = 8
)

var CacheDir string

func openDataFile(path string) (*os.File, error) {
	err := preloadDataFiles([]string{path})
	if err != nil {
		return nil, err
	}

	return os.Open(cachePath(path))
}

func isFileCached(path string) bool {
	fpath := cachePath(path)

	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return false
	}

	return true
}

func cachePath(path string) string {
	return filepath.Join(CacheDir, path)
}

func preloadDataFiles(paths []string) error {
	var toDownload []string
	dirs := make(map[string]struct{})

	for _, path := range paths {
		if !isFileCached(path) {
			toDownload = append(toDownload, path)

			fpath := cachePath(path)

			dir := filepath.Dir(fpath)
			dirs[dir] = struct{}{}
		}
	}

	if len(toDownload) == 0 {
		return nil
	}

	for dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	pb := progressbar.Default(int64(len(toDownload)), "Downloading data files")

	jobs := make(chan string, len(toDownload))
	results := make(chan error, len(toDownload))

	for i := 0; i < numWorkers; i++ {
		go downloadWorker(jobs, results)
	}

	for _, path := range toDownload {
		jobs <- path
	}

	close(jobs)

	var errors int

	for i := 0; i < len(toDownload); i++ {
		err := <-results
		if err != nil {
			errors += 1
		}
		pb.Add(1)
	}

	pb.Finish()

	if errors > 0 {
		return fmt.Errorf("failed to download %d files", errors)
	}

	return nil
}

func downloadWorker(paths <-chan string, results chan<- error) {
	for path := range paths {
		results <- downloadDataFile(path)
	}
}

func downloadDataFile(path string) error {
	loc, err := url.JoinPath(baseUrl, path)
	if err != nil {
		return err
	}

	resp, err := http.Get(loc)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	fpath := cachePath(path)

	f, err := os.Create(fpath)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
