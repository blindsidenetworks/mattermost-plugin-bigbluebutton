package helpers

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

//sends a get request to the url given, returns the result as a string
func HttpGet(url string) string {
	response, err := http.Get(url)

	if err != nil {
		log.Println("HTTP GET ERROR: " + err.Error())
		return "ERROR"
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if nil != err {
		log.Println("HTTP GET ERROR: " + err.Error())
		return "ERROR"
	}

	return string(body)
}

func GetChecksum(toConvert string) string {
	toByte := []byte(toConvert)
	checkSumString := sha1.Sum(toByte)

	return hex.EncodeToString(checkSumString[:])
}

func ReadXML(response string, data interface{}) error {
	err := xml.Unmarshal([]byte(response),data)
	if nil != err {
		log.Println("XML PARSE ERROR: " + err.Error())
	}
	return err
}
