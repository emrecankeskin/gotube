package gotube

import (
	"fmt"
	"log"
	"runtime"
)

// TODO can i create a youtube object that contains videoFormat audioFormat
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\nHeap Objects = %v Mib", m.HeapAlloc/1024/1024)
	fmt.Printf("\nAlloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\nTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\nSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\nNumGC = %v\n", m.NumGC)
}

func GetStreams(urlLink string) (map[string]videoFormat, map[string]audioFormat) {
	var video = make(map[string]videoFormat)
	var audio = make(map[string]audioFormat)
	var err error
	var down = downloader{}
	var className string
	var adaptiveVal []byte
	var data []byte
	var encMap map[string]string
	var encMethods []string

	var htmlText = _html{down.getPlainText(urlLink)}
	var jsonObj = jsonObject{htmlText.text}
	var baseJsLink = htmlText.getBaseJs()
	var js = jsFile{content: down.getPlainText(baseJsLink), key: "a=a.split(\"\");"}

	methodArr, _ := js.getEncryptMethods(js.getEncryptFunction())
	className = js.getEncryptClassName(methodArr[0])
	start, end, _ := js.getEncryptBase(className)
	encryptor := (*js.content)[start:end]
	encFunc := js.getEncryptFunction()
	encMethods, _ = js.getEncryptMethods(encFunc)
	encMap, err = js.mapEncryptMethods(string(encryptor))

	if err != nil {
		log.Println(err)
		return nil, nil
	}

	data = []byte((jsonObj.getJsonObject("streamingData")))
	var jsonArr = jsonArray{&data}
	adaptiveVal = []byte(jsonArr.getJsonArray("adaptiveFormats"))
	adaptive := adaptiveObject{&adaptiveVal, encMap, encMethods}
	video, audio = adaptive.getAdaptiveFormat()
	var details = jsonObj.getJsonObject("videoDetails")
	fmt.Println(getObjectValue(details, "title"))
	return video, audio
}
