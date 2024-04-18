package main

import (
	"gotube/gotube"
	"os"
)

func example(ytLink string) {
	a, _ := gotube.GetStreams(ytLink)
	gotube.PrintMemUsage()
	f, _ := os.Create("signature")
	for _, v := range a {

		_, err := f.WriteString(v.SignatureCipher)
		if err != nil {
			return
		}
		_, err = f.WriteString("\n")
		if err != nil {
			return
		}
	}
}
