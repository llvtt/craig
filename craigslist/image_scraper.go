package craigslist

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
	"io"
	"net/http"
	"strings"
)

const IMAGE_HOST = "https://images.craigslist.org"

type ImageScraper interface {
	GetImageUrls(ci *types.CraigslistItem) (urls []string, err error)
}

type imageScraper struct {
	logger log.Logger
}

type ImageDesc struct {
	ImgId string `json:"imgid"`
}

func NewImageScraper(logger log.Logger) ImageScraper {
	return &imageScraper{
		logger: logger,
	}
}

func (s *imageScraper) imageIds(imgListDecl string) (ids []string) {
	parts := strings.Split(imgListDecl, " = ")
	if len(parts) != 2 {
		level.Warn(s.logger).Log("msg", "could not parse imgListDecl: " + imgListDecl)
		return
	}
	value := strings.Trim(parts[1], "; ")
	var descs []ImageDesc
	if err := json.Unmarshal([]byte(value), &descs); err != nil {
		panic(err)
	}
	for _, desc := range descs {
		// Remove stuff before the first colon
		parts := strings.SplitN(desc.ImgId, ":", 2)
		if len(parts) != 2 {
			level.Warn(s.logger).Log("msg", "could not parse ids in imgListDecl: " + imgListDecl)
			return
		}
		ids = append(ids, parts[1])
	}
	return
}

func (s *imageScraper) imageUrls(reader io.Reader) (urls []string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "var imgList =") {
			ids := s.imageIds(line)
			for _, id := range ids {
				url := fmt.Sprintf("%s/%s_1200x900.jpg", IMAGE_HOST, id)
				urls = append(urls, url)
			}
			return
		}
	}
	return
}

func (s *imageScraper) GetImageUrls(ci *types.CraigslistItem) ([]string, error) {
	resp, err := http.Get(ci.Url)
	if err != nil {
		return nil, utils.WrapError("could not get image from url: "+ci.Url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, nil
	}
	urls := s.imageUrls(resp.Body)
	return urls, nil
}
