package webhook

import (
	"log"
	"net/url"

	"github.com/mattermost/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
	"github.com/mattermost/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/helpers"
)

//see documentation:http://docs.bigbluebutton.org/dev/webhooks.html
//webhook was designed to be used for

var BASE_URL string
var SALT string

func SetWebhookAPI(url string, salt string) {
	BASE_URL = url
	SALT = salt
}

func CreateHook(wh *dataStructs.WebHook) string {
	if wh.CallBackURL == "" {
		return "Error, must indicate callback url"
	}
	callback := "callbackURL=" + url.QueryEscape(wh.CallBackURL)
	getRaw := "&getRaw=true"
	params := callback + getRaw
	checkSum := helpers.GetChecksum("hooks/create" + params + SALT)

	response := helpers.HttpGet(BASE_URL + "create?" + params + "&checksum=" + checkSum)

	if "ERROR" == response {
		log.Println("ERROR: HTTP ERROR.")
		return "ERROR: HTTP ERROR."
	}
	err := helpers.ReadXML(response, &(wh.WebhookResponse))

	if nil != err {
		log.Println("XML PARSE ERROR: " + err.Error())
		return "ERROR: XML PARSE ERROR."
	}
	wh.HookID = wh.WebhookResponse.HookID
	if wh.WebhookResponse.Returncode == "SUCCESS" {
		return "webhook successfully created " + wh.HookID
	} else {
		return wh.WebhookResponse.Message
	}

	return "aaaaaaahhh"
}

func DestroyHook(hookID string) string {
	hook_id := "hookID=" + url.QueryEscape(hookID)
	params := hook_id
	checkSum := helpers.GetChecksum("hooks/destroy" + params + SALT)

	response := helpers.HttpGet(BASE_URL + "destroy?" + params + "&checksum=" + checkSum)

	if "ERROR" == response {
		log.Println("ERROR: HTTP ERROR.")
		return "ERROR: HTTP ERROR."
	}
	var responseXML dataStructs.DestroyedWebhookResponse
	err := helpers.ReadXML(response, &responseXML)

	if nil != err {
		log.Println("XML PARSE ERROR: " + err.Error())
		return "ERROR: XML PARSE ERROR."
	}
	if responseXML.Returncode == "SUCCESS" {
		return "webhook " + hookID + " destroyed"
	}
	return "Can't delete webbook " + hookID + " : " + responseXML.Message
}
