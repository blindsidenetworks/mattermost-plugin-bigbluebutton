/*
Copyright 2018 Blindside Networks

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helpers

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/mattermost"
	"io/ioutil"
	"net/http"
)

var PluginVersion string

//sends a get request to the url given, returns the result as a string
func HttpGet(url string) (string, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "bbb-mm-" + PluginVersion)

	response, err := client.Do(req)
	if err != nil {
		mattermost.API.LogError("HTTP GET ERROR: " + err.Error())
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		mattermost.API.LogError("HTTP GET ERROR: " + err.Error())
		return "", err
	}

	return string(body), nil
}

func GetChecksum(toConvert string) string {
	toByte := []byte(toConvert)
	checkSumString := sha1.Sum(toByte)

	return hex.EncodeToString(checkSumString[:])
}

func ReadXML(response string, data interface{}) error {
	err := xml.Unmarshal([]byte(response), data)
	if nil != err {
		mattermost.API.LogError("XML PARSE ERROR: " + err.Error())
	}
	return err
}
