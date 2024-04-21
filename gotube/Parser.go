package gotube

import (
	"errors"
	"strings"
)

// TODO implement them to accept pointers as texts because copying big string is expensive

type jsonArray struct {
	jsonObject *[]byte
}

type jsonObject struct {
	object *[]byte
}

type _html struct {
	text *[]byte
}
type jsFile struct {
	content *[]byte
	key     string
}

//Returns index of ending point and key's last index. If (keyPtr = len(key)-1) then function is not found key
// in this particular text.

// []byte{'"', 's', 't', 'r', 'e', 'a', 'm', 'i', 'n', 'g', 'D', 'a', 't', 'a', '"'}

// function gets only the object with in { }
func (j jsonObject) getJsonObject(str string) string {
	var cursor int
	var start int
	var length = len(*j.object)
	var bracketCount = 0
	var key = []byte(str)
	var keyPtr = 0

	cursor = 0
	for cursor < length && keyPtr < len(key) {
		if (*j.object)[cursor] == key[keyPtr] {
			keyPtr++
		} else {
			keyPtr = 0
		}
		cursor++
	}

	//object is not found
	if cursor == length-1 {
		return ""
	}

	for string((*j.object)[cursor]) != "{" {
		cursor++
	}
	start = cursor
	cursor++
	bracketCount++
	for bracketCount > 0 {
		if string((*j.object)[cursor]) == "{" {
			bracketCount += 1
		} else if string((*j.object)[cursor]) == "}" {
			bracketCount -= 1
		}
		cursor++
	}
	return string((*j.object)[start:cursor])
}

//[]byte{'"', 'f', 'o', 'r', 'm', 'a', 't', 's', '"'}
/*
	@param str = json array key
	Parsing the json array from json object and returns it
*/
func (j jsonArray) getJsonArray(str string) string {
	var cursor int
	var start int
	var length = len(*j.jsonObject)
	var key = []byte(str)
	var keyPtr = 0

	cursor = 0
	for cursor < length && keyPtr < len(key) {
		if (*j.jsonObject)[cursor] == key[keyPtr] {
			keyPtr++
		} else {
			keyPtr = 0
		}
		cursor++
	}

	if cursor == length-1 {
		return ""
	}

	for string((*j.jsonObject)[cursor]) != "[" {
		cursor++
	}
	start = cursor
	for string((*j.jsonObject)[cursor]) != "]" {
		cursor++
	}
	return string((*j.jsonObject)[start : cursor+1])
}

// TODO create function that gets all the objects of jsonArray array to a string to string map

// only finds for one object's value don't use with arrays
// i can deprecate it because it gets only one value i have to parse all the values from json
func getObjectValue(val string, str string) string {
	var cursor int
	var length = len(val)
	var key = []byte(str)
	var start int
	var keyPtr = 0
	for cursor < length && keyPtr < length {
		if (val)[cursor] == key[keyPtr] {
			keyPtr++
		} else {
			keyPtr = 0
		}
		cursor++
	}
	for (val)[cursor] != '"' {
		cursor++
	}
	cursor++
	start = cursor
	for (val)[cursor] != '"' {
		cursor++
	}
	return (val)[start:cursor]
}
func (j jsonObject) getValue(obj *[]byte, s int, e int, str string) string {
	var cursor int
	var start int
	var length = s - e + 1
	var key = []byte(str)
	var keyPtr = 0

	cursor = s
	for cursor < length && keyPtr < len(key) {
		if (*obj)[cursor] == key[keyPtr] {
			keyPtr++
		} else {
			keyPtr = 0
		}
		cursor++
	}

	if cursor == length-1 {
		return ""
	}

	for string((*obj)[cursor]) != ":" {
		cursor++
	}
	start = cursor
	for string((*obj)[cursor]) != "," {
		cursor++
	}

	return string(*obj)[start+2 : cursor-2]

}

/*
Returns player response object from html file
*/
func (html _html) getPlayerResponse() string {

	var cursor int
	var start int
	var length = len(*html.text)
	var key = []byte("var ytInitialPlayerResponse")
	var keyPtr = 0
	var bracketCount = 0

	cursor = 0
	for cursor < length && keyPtr < len(key) {
		if (*html.text)[cursor] == key[keyPtr] {
			keyPtr++
		} else {
			keyPtr = 0
		}
		cursor++
	}

	if cursor == length-1 {
		return ""
	}

	for string((*html.text)[cursor]) != "{" {
		cursor++
	}
	start = cursor
	bracketCount++
	cursor++
	for bracketCount > 0 {

		if string((*html.text)[cursor]) == "{" {
			bracketCount += 1
		} else if string((*html.text)[cursor]) == "}" {
			bracketCount -= 1
		}

		cursor++
	}
	return string((*html.text)[start:cursor])
}

// Finds web player's link and returns
func (html _html) getBaseJs() string {
	/*
		Finds the web player object from plain text _html
		and returns the concatenated link
	*/
	var key = []byte("\"jsUrl\"")
	var keyPtr = 0
	var keyLen = len(key)
	var cursor = 0
	var start int
	var json jsonObject

	json = jsonObject{html.text}
	var jsonObj = json.getJsonObject("\"WEB_PLAYER_CONTEXT_CONFIG_ID_KEVLAR_WATCH\"")
	var objLen = len(jsonObj)

	for cursor < objLen && keyPtr < keyLen {
		if jsonObj[cursor] == key[keyPtr] {
			keyPtr++
		} else {
			keyPtr = 0
		}
		cursor++
	}
	if cursor == objLen-1 {
		return ""
	}

	for string(jsonObj[cursor]) != "\"" {
		cursor++
	}
	start = cursor
	cursor++
	for string(jsonObj[cursor]) != "," {
		cursor++
	}
	return "https://www.youtube.com" + jsonObj[start+1:cursor-1]
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
func (js jsFile) getEncryptFunction() string {
	var key = []byte(js.key)
	var cursor = 0
	var keyPtr = 0
	var contentLen = len(*js.content)
	var keyLen = len(key)
	var start = 0

	for cursor < contentLen && keyPtr < keyLen {
		if (key[keyPtr]) == ((*js.content)[cursor]) {
			keyPtr++
		} else {
			keyPtr = 0
		}
		cursor++
	}

	if cursor == contentLen-1 {
		return ""
	}
	start = cursor
	cursor++
	for string((*js.content)[cursor]) != "}" {
		cursor++
	}

	return string((*js.content)[start:cursor])
}

/* Returns the methods that used for encrypting url */
func (js jsFile) getEncryptMethods(encFunc string) ([]string, error) {
	var size = len(encFunc)
	var cursor = 0
	var start = 0
	var encMethods = make([]string, 0)

	for cursor < size {
		if string(encFunc[cursor]) == ";" {
			encMethods = append(encMethods, encFunc[start:cursor+1])
			start = cursor + 1
		}
		cursor++
	}

	if len(encMethods) < 1 {
		return nil, errors.New("NO OPERATIONS")
	}

	return encMethods, nil
}

/*
Deciphers the encrypt methods and maps
[method] : [splice,reverse,swap]
mapping the function with iterating over encryptor string
*/
func (js jsFile) mapEncryptMethods(encryptor string) (map[string]string, error) {
	var cursor = 0
	var start = 0
	var funcStr = make([]string, 0)
	var ops = make(map[string]string)

	for cursor < len(encryptor) {
		if encryptor[cursor] == '}' {
			funcStr = append(funcStr, encryptor[start:cursor+1])
			cursor += 2
			start = cursor
		}
		cursor++
	}
	for _, v := range funcStr {
		if strings.Contains(v, "splice") {
			l := 0
			r := 0
			for v[l] == '\n' || v[l] == ' ' {
				l++
			}
			for v[r] != ':' {
				r++
			}
			ops[v[l:r]] = "splice"
		} else if strings.Contains(v, "reverse") {
			l := 0
			r := 0
			for v[l] == '\n' || v[l] == ' ' {
				l++
			}
			for v[r] != ':' {
				r++
			}
			ops[v[l:r]] = "reverse"
		} else {
			l := 0
			r := 0
			for v[l] == '\n' || v[l] == ' ' {
				l++
			}
			for v[r] != ':' {
				r++
			}
			ops[v[l:r]] = "swap"
		}
	}
	if len(ops) < 1 {
		return nil, errors.New("CANT PARSE OPERATIONS")
	}

	return ops, nil
}

/*
Helper class for mapEncryptMethods to get encryptor function's implementation
returns string of encryptor function
*/
func (js jsFile) getEncryptClassName(str string) string {
	var cursor = 0
	for cursor < len(str) {
		if string(str[cursor]) == "." {
			break
		}
		cursor++
	}

	return "var " + str[:cursor] + "={"
}

/*
Finds the encryptor function that uses encrypt methods
*/
//func (js jsFile) getEncryptBase(encClassName string) (string, error)
func (js jsFile) getEncryptBase(encClassName string) (int, int, error) {
	var cursor = 0
	var keyPtr = 0
	var start int
	var bracketCount = 1
	for cursor < len(*js.content) && keyPtr < len(encClassName) {
		if string((*js.content)[cursor]) == string(encClassName[keyPtr]) {
			keyPtr++
		} else {
			keyPtr = 0
		}
		cursor++
	}
	if cursor == len(*js.content)-1 {
		return 0, 0, errors.New("CAN NOT FIND ENCRYPT BASE FUNCTION")
	}
	start = cursor
	for bracketCount != 0 {
		if string((*js.content)[cursor]) == "}" {
			bracketCount--
		} else if string((*js.content)[cursor]) == "{" {
			bracketCount++
		}
		cursor++
	}
	return start, cursor + 1, nil
}
