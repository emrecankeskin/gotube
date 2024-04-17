package gotube

import (
	"io"
	"log"
	"net/http"
)

type downloader struct {
}

func (d downloader) getPlainText(link string) *[]byte {
	resp, err := http.Get(link)
	if err != nil {
		return nil
	}
	content, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil
	}
	err = resp.Body.Close()
	if err != nil {
		return nil
	}
	log.Printf("[*] File size is = %d KiB", len(content)/1024)
	return &content
}
