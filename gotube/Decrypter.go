package gotube

import (
	"errors"
	"html"
	"net/url"
	"strconv"
	"strings"
)

type decrypter struct {
	encMap     map[string]string
	encMethods []string
}

/*
Encryptor function example:

		var FP ={
	        Pq: function(a, b) {
	            a.splice(0, b)
	        },
	        kD: function(a) {
	            a.reverse()
	        },
	        aO: function(a, b) {
	            var c = a[0];
	            a[0] = a[b % a.length];
	            a[b % a.length] = c
	        }
	    };
		var FP ={ is enough for searching
*/
/*
	Main method of decrypter other methods are just helper for this function
*/
func (d decrypter) decodeUrl(link string) string {
	// create array of parameters of the link
	// decode the signature with decodeSignature
	var signature string
	var arr []string
	var decodedLink string
	link = strings.ReplaceAll(link, "\n", "")
	//TODO try to find better solution to decode unicode
	link = strings.ReplaceAll(link, "\\u0026", "&")
	link, _ = d.urlDecode(link)
	arr = strings.Split(link, "&")

	signature = arr[0][2:]
	//fmt.Println("SIGNATURE:", signature)
	decodedLink = reconstructLink(arr, d.decodeSignature(signature))
	//fmt.Println("DECODED LINK: ", decodedLink)
	return decodedLink
}

// TODO create method parser for decode and simulate all the encoding process
// TODO it is decoding signature cipher
func (d decrypter) decodeSignature(signature string) string {
	for _, v := range d.encMethods {
		var key, val = parseMethod(v)
		//fmt.Printf("key : %s val : %s\n", key, val)
		if d.encMap[key] == "splice" {
			converted, _ := strconv.Atoi(val)
			signature = signature[converted:]
		} else if d.encMap[key] == "reverse" {
			signature = reverseString(signature)
		} else if d.encMap[key] == "swap" {
			b, _ := strconv.Atoi(val)
			signature = swapString(signature, b)
		}
	}
	return signature
}
func reconstructLink(arr []string, signature string) string {
	var sb strings.Builder
	var ptr = 3
	sb.WriteString(arr[2][4:])
	for ptr < len(arr) {
		sb.Write([]byte("&" + arr[ptr]))
		ptr++
	}
	sb.Write([]byte("&sig=" + signature))

	return sb.String()
}

func reverseString(str string) string {
	n := 0
	arr := make([]rune, len(str))
	for _, r := range str {
		arr[n] = r
		n++
	}
	arr = arr[0:n]
	// Reverse
	for i := 0; i < n/2; i++ {
		arr[i], arr[n-1-i] = arr[n-1-i], arr[i]
	}
	// Convert back to UTF-8.
	return string(arr)
}
func swapString(str string, b int) string {
	var arr = []byte(str)
	var c = arr[0]
	arr[0] = arr[b%len(arr)]
	arr[b%len(arr)] = c

	return string(arr)
}
func parseMethod(line string) (string, string) {
	var cursor = 0
	var startMethod = 0
	var endMethod = 0
	var startVal = 0
	var endVal = 0

	for cursor < len(line) {
		if string(line[cursor]) == "." {
			startMethod = cursor + 1
		} else if string(line[cursor]) == "(" {
			endMethod = cursor
		} else if string(line[cursor]) == "," {
			startVal = cursor + 1
		} else if string(line[cursor]) == ")" {
			endVal = cursor
		}
		cursor++
	}

	return line[startMethod:endMethod], line[startVal:endVal]

}

// TODO decode unicode chars
func (d decrypter) unicodeDecode(url string) string {
	return html.UnescapeString(url)
}

// TODO url decode learn how to do in golang
func (d decrypter) urlDecode(link string) (string, error) {

	decoded, err := url.QueryUnescape(link)
	if err != nil {
		return "", errors.New("CANT DECODE URL")
	}

	return decoded, nil
}
