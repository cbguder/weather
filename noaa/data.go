package noaa

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const baseUrl = "https://www.ncei.noaa.gov/pub/data/ghcn/daily/"

var CacheDir string

func openDataFile(path string) (*os.File, error) {
	fpath := filepath.Join(CacheDir, path)

	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		err = downloadDataFile(path)
		if err != nil {
			return nil, err
		}
	}

	return os.Open(fpath)
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

	fpath := filepath.Join(CacheDir, path)

	dir := filepath.Dir(fpath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(fpath)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
