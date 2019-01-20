package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const IMAGE_HOST = "https://images.craigslist.org"

type ImageDesc struct {
	ImgId string `json:"imgid"`
}

func imageIds(imgListDecl string) (ids []string) {
	parts := strings.Split(imgListDecl, " = ")
	if len(parts) != 2 {
		fmt.Println("warning: could not parse imgListDecl: " + imgListDecl)
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
			fmt.Println("warning: could not parse ids in imgListDecl: " + imgListDecl)
			return
		}
		ids = append(ids, parts[1])
	}
	return
}

func imageUrls(reader io.Reader) (urls []string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "var imgList =") {
			ids := imageIds(line)
			for _, id := range ids {
				url := fmt.Sprintf("%s/%s_1200x900.jpg", IMAGE_HOST, id)
				urls = append(urls, url)
			}
			return
		}
	}
	return
}

func (ci *CraigslistItem) GetImageUrls() (urls []string) {
	resp, err := http.Get(ci.Url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return
	}
	urls = imageUrls(resp.Body)
	return
}
