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

package webhook

import (
	"errors"
	"net/url"

	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/helpers"
)

//see documentation:http://docs.bigbluebutton.org/dev/webhooks.html
//webhook was designed to be used for

var BaseUrl string
var secret string

func SetWebhookAPI(url string, secretParam string) {
	BaseUrl = url
	secret = secretParam
}

func CreateHook(wh *dataStructs.WebHook) (string, error) {
	if wh.CallBackURL == "" {
		return "", errors.New("error, must indicate callback url")
	}
	callback := "callbackURL=" + url.QueryEscape(wh.CallBackURL)
	getRaw := "&getRaw=true"
	params := callback + getRaw
	checkSum := helpers.GetChecksum("hooks/create" + params + secret)

	response, err := helpers.HttpGet(BaseUrl + "create?" + params + "&checksum=" + checkSum)
	if err != nil {
		return "", err
	}
	err = helpers.ReadXML(response, &(wh.WebhookResponse))

	if nil != err {
		return "", err
	}
	wh.HookID = wh.WebhookResponse.HookID
	if wh.WebhookResponse.Returncode == "SUCCESS" {
		return "webhook successfully created " + wh.HookID, nil
	} else {
		return "", errors.New(wh.WebhookResponse.Message)
	}
}

func DestroyHook(hookID string) (string, error) {
	hookId := "hookID=" + url.QueryEscape(hookID)
	params := hookId
	checkSum := helpers.GetChecksum("hooks/destroy" + params + secret)

	response, err := helpers.HttpGet(BaseUrl + "destroy?" + params + "&checksum=" + checkSum)
	if err != nil {
		return "", err
	}
	var responseXML dataStructs.DestroyedWebhookResponse
	err = helpers.ReadXML(response, &responseXML)

	if nil != err {
		return "", err
	}
	if responseXML.Returncode == "SUCCESS" {
		return "webhook " + hookID + " destroyed", nil
	}
	return "", errors.New("Can't delete webbook " + hookID + " : " + responseXML.Message)
}
