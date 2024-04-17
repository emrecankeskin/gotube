package gotube

import (
	"strconv"
	"strings"
)

type adaptiveObject struct {
	object     *[]byte
	encMap     map[string]string
	encMethods []string
}

type videoFormat struct {
	SignatureCipher string
	MimeType        string
	Quality         string
	QualityLabel    string
	Bitrate         int
	Width           int
	Height          int
	Fps             int
}

type audioFormat struct {
	SignatureCipher string
	MimeType        string
	Quality         string
	AudioQuality    string
	Bitrate         int
	AudioChannels   int
	LoudnessDb      int
}

/*
TODO parse adaptiveFormat "[tag]": [val],

TODO parse all of the keys of json object then only return wanted ones
TODO figure out how to handle inner json objects
*/
func (j adaptiveObject) getAdaptiveFormat() (map[string]videoFormat, map[string]audioFormat) {

	var video = make(map[string]videoFormat)
	var audio = make(map[string]audioFormat)
	var cursor = 0
	var start = 0
	var count = 0
	var dec = decrypter{j.encMap, j.encMethods}
	// parse adaptive format obj by obj then get values of this obj
	for cursor < len(*j.object) {
		count = 0
		for cursor < len(*j.object) && string((*j.object)[cursor]) != "{" {
			cursor++
		}
		count = 1
		start = cursor
		cursor++
		for count != 0 && cursor < len(*j.object) {
			if string((*j.object)[cursor]) == "{" {
				count++
			} else if string((*j.object)[cursor]) == "}" {
				count--
			}
			cursor++
		}
		cursor++
		(*j.object)[cursor-2] = ','
		j.parseJsonObject(string((*j.object)[start:cursor-1]), &dec, video, audio)

	}
	return video, audio
}

// TODO how to handle str with pointers because it consumes memory
func (j adaptiveObject) parseJsonObject(str string, dec *decrypter, videoMap map[string]videoFormat, audioMap map[string]audioFormat) {
	var video videoFormat
	var audio audioFormat
	var cursor = 0
	var keyStart = 0
	var keyEnd = 0
	var valStart = 0
	var valEnd = 0

	var mp = make(map[string]string)

	for cursor < len(str) {
		for string(str[cursor]) != "\"" {
			cursor++
		}
		cursor++
		keyStart = cursor
		for string(str[cursor]) != "\"" {
			cursor++
		}
		keyEnd = cursor
		cursor++
		for string(str[cursor]) == " " || string(str[cursor]) == ":" {
			cursor++
		}
		if string(str[cursor]) == "{" {
			for string(str[cursor]) != "}" {
				cursor++
			}
			cursor++
		} else {
			for string(str[cursor]) == " " || string(str[cursor]) == ":" {
				cursor++
			}
			valStart = cursor
			for string(str[cursor]) != "," {
				cursor++
			}
			valEnd = cursor
			cursor++
			mp[str[keyStart:keyEnd]] = str[valStart:valEnd]
		}

	}
	if strings.Contains(mp["mimeType"], "video") {
		val, ok := mp["signatureCipher"]
		if ok {
			video.SignatureCipher = dec.decodeUrl(val[1 : len(mp["signatureCipher"])-1])
		} else {
			video.SignatureCipher = mp["url"][1 : len(mp["url"])-1]
		}
		video.MimeType = mp["mimeType"][1 : len(mp["mimeType"])-1]
		video.Quality = mp["quality"][1 : len(mp["quality"])-1]
		video.QualityLabel = mp["qualityLabel"][1 : len(mp["qualityLabel"])-1]
		video.Bitrate, _ = strconv.Atoi(mp["bitrate"])
		video.Width, _ = strconv.Atoi(mp["width"])
		video.Height, _ = strconv.Atoi(mp["height"])
		video.Fps, _ = strconv.Atoi(mp["fps"])
		videoMap[mp["itag"]] = video
	} else {

		val, ok := mp["signatureCipher"]
		if ok {
			audio.SignatureCipher = dec.decodeUrl(val[1 : len(mp["signatureCipher"])-1])
		} else {
			audio.SignatureCipher = mp["url"][1 : len(mp["url"])-1]
		}
		audio.MimeType = mp["mimeType"][1 : len(mp["mimeType"])-1]
		audio.Quality = mp["quality"][1 : len(mp["quality"])-1]
		audio.AudioQuality = mp["audioQuality"][1 : len(mp["audioQuality"])-1]
		audio.Bitrate, _ = strconv.Atoi(mp["bitrate"])
		audio.AudioChannels, _ = strconv.Atoi(mp["audioChannels"])
		audio.LoudnessDb, _ = strconv.Atoi(mp["loudnessDb"])
		audioMap[mp["itag"]] = audio
	}
	mp = nil
}
func (j jsonObject) getVideoTitle() string {
	return ""
}
