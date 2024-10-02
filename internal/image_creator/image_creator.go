package image_creator

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
)

type Quote struct {
	Author string   `json:"author"`
	Quote  string   `json:"quote"`
	Tags   []string `json:"tags"`
}

type UnsplashPhoto struct {
	ID     string `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URLs   struct {
		Regular string `json:"regular"`
	} `json:"urls"`
}

func getRandomPhoto(accessKey string) (*UnsplashPhoto, error) {
	url := fmt.Sprintf("https://api.unsplash.com/photos/random?client_id=%s", accessKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка запроса: %s", resp.Status)
	}

	var photo UnsplashPhoto
	if err := json.NewDecoder(resp.Body).Decode(&photo); err != nil {
		return nil, err
	}

	return &photo, nil
}

func GenerateImageForQuote(quote Quote, outputPath string, accessKey string) error {
	const W = 1024
	const H = 1024

	photo, err := getRandomPhoto(accessKey)
	if err != nil {
		return err
	}

	resp, err := http.Get(photo.URLs.Regular)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return err
	}

	dc := gg.NewContext(W, H)
	dc.DrawImage(img, 0, 0)

	fontPath, err := filepath.Abs("fonts/Cheri.ttf")
	if err != nil {
		return err
	}

	if err := dc.LoadFontFace(fontPath, 48); err != nil {
		return err
	}

	dc.SetRGB(1, 1, 1)
	dc.DrawStringWrapped(quote.Quote, W/2, H/2, 0.5, 0.5, W-50, 1.5, gg.AlignCenter)

	if err := dc.LoadFontFace(fontPath, 36); err != nil {
		return err
	}

	dc.SetRGB(0.8, 0.8, 0.8)
	dc.DrawStringAnchored("- "+quote.Author, 50, H-50, 0, 0)

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, dc.Image())
	if err != nil {
		return err
	}

	return nil
}

func LoadQuotes(filePath string) ([]Quote, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var quotes []Quote
	err = json.Unmarshal(file, &quotes)
	if err != nil {
		return nil, err
	}

	return quotes, nil
}
