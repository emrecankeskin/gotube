package main

import (
	"gotube/gotube"
	"os"
)

func main() {
	a, _ := gotube.GetStreams("https://www.youtube.com/watch?v=go6nQOKakeQ")
	//PrintMemUsage()
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
