package imageloader

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type ImageLoader interface {
	Load(url string) ([]byte, error)
}

type HTTPImageLoader struct {
}

func NewHTTPImageLoader() *HTTPImageLoader {
	return &HTTPImageLoader{}
}

func (h *HTTPImageLoader) Load(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "error performing HTTP GET when loading image from URL: %v", url)
	}
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("got HTTP status %v when loading image from URL: %v", resp.StatusCode, url)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return data, nil
}
