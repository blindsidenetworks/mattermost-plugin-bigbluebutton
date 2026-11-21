package main

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// Embed the SDK interface so unexpected API calls fail instead of being silently ignored.
type dependencyTestAPI struct {
	plugin.API
	command *model.Command
	post    *model.Post
}

func (a *dependencyTestAPI) LoadPluginConfiguration(dest any) error {
	*dest.(*Configuration) = Configuration{BaseURL: testBBBURL, Secret: "secret", AllowRecordings: true}
	return nil
}
func (a *dependencyTestAPI) RegisterCommand(command *model.Command) error {
	a.command = command
	return nil
}
func (a *dependencyTestAPI) GetConfig() *model.Config {
	siteURL := "https://chat.example:8443/mattermost"
	return &model.Config{ServiceSettings: model.ServiceSettings{SiteURL: &siteURL}}
}
func (a *dependencyTestAPI) GetUser(id string) (*model.User, *model.AppError) {
	return &model.User{Id: id, Username: "alice", Email: "alice@example.com"}, nil
}
func (a *dependencyTestAPI) CreatePost(post *model.Post) (*model.Post, *model.AppError) {
	a.post = post
	post.Id = "post-id"
	return post, nil
}
func (a *dependencyTestAPI) GetChannel(id string) (*model.Channel, *model.AppError) {
	return &model.Channel{Id: id, Name: "alice__bob", Type: model.ChannelTypeDirect}, nil
}

func TestSDKPluginLifecycle(t *testing.T) {
	api := &dependencyTestAPI{}
	p := &Plugin{}
	p.SetAPI(api)
	if err := p.OnActivate(); err != nil {
		t.Fatal(err)
	}
	if api.command == nil || api.command.Trigger != "bbb" || api.command.AutocompleteData.RoleID != model.SystemUserRoleId {
		t.Fatalf("slash command registration failed: %+v", api.command)
	}
	recipient, show, err := p.shouldShowCallPopup("alice", "channel")
	if err != nil || !show || recipient != "bob" {
		t.Fatalf("direct call recipient = %q, %v, %v", recipient, show, err)
	}
	// Recording processing is disabled, so shutdown must tolerate no scheduled job.
	if err := p.OnDeactivate(); err != nil {
		t.Fatal(err)
	}
}

func TestSDKMeetingPostAndCallbacks(t *testing.T) {
	for _, allowRecording := range []bool{false, true} {
		t.Run(map[bool]string{false: "unrecorded", true: "recorded"}[allowRecording], func(t *testing.T) {
			api := &dependencyTestAPI{}
			p := &Plugin{}
			p.SetAPI(api)
			p.configuration.Store(&Configuration{BaseURL: testBBBURL, Secret: "secret", AllowRecordings: true})
			meeting := &dataStructs.MeetingRoom{}
			if err := p.PopulateMeeting(meeting, nil, "test meeting", "alice", "channel", allowRecording); err != nil {
				t.Fatal(err)
			}
			if meeting.Meta_bbb_origin_server_name != "chat.example" {
				t.Fatalf("unexpected origin: %q", meeting.Meta_bbb_origin_server_name)
			}
			if !strings.HasPrefix(meeting.Meta_bbb_recording_ready_url, "https://chat.example:8443/mattermost/plugins/bigbluebutton/") {
				t.Fatal(meeting.Meta_bbb_recording_ready_url)
			}
			if meeting.AllowStartStopRecording != allowRecording {
				t.Fatal("recording preference lost")
			}
			if err := p.createStartMeetingPost("alice", "channel", meeting); err != nil {
				t.Fatal(err)
			}
			// Check the wire format rather than the new SDK's Go type names.
			data, err := json.Marshal(api.post.GetProp("attachments"))
			if err != nil {
				t.Fatal(err)
			}
			var attachments []struct {
				Actions []struct {
					Name        string `json:"name"`
					Integration struct {
						URL string `json:"url"`
					} `json:"integration"`
				} `json:"actions"`
				Fields []struct {
					Title string `json:"title"`
					Value string `json:"value"`
				} `json:"fields"`
			}
			if err := json.Unmarshal(data, &attachments); err != nil {
				t.Fatal(err)
			}
			if len(attachments) != 1 || len(attachments[0].Actions) != 2 || attachments[0].Actions[0].Integration.URL != "/plugins/bigbluebutton/joinmeeting" {
				t.Fatalf("meeting actions changed: %s", data)
			}
			expected := "Disabled"
			if allowRecording {
				expected = "Allowed"
			}
			if attachments[0].Fields[1].Value != expected {
				t.Fatalf("recording field changed: %s", data)
			}
		})
	}
}

func TestPublicConfigurationOmitsCredentials(t *testing.T) {
	config := &Configuration{BaseURL: testBBBURL, Secret: "must-not-leak", AllowRecordings: true, AllowExternalUsers: true}
	data, err := json.Marshal(config.AsMap())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "SALT") || strings.Contains(string(data), config.Secret) || strings.Contains(string(data), config.BaseURL) {
		t.Fatalf("configuration exposes connection details: %s", data)
	}
	p := &Plugin{}
	p.configuration.Store(config)
	response := httptest.NewRecorder()
	p.handleGetConfig(response, httptest.NewRequest("POST", "/config", nil))
	if response.Code != 200 || response.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("config response: %v", response.Result())
	}
	var flags map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &flags); err != nil {
		t.Fatal(err)
	}
	if flags["ALLOW_RECORDINGS"] != true || flags["ALLOW_EXTERNAL_USERS"] != true {
		t.Fatalf("public flags missing: %v", flags)
	}
	if _, found := flags["SALT"]; found {
		t.Fatal("secret field exposed")
	}
}

const testBBBURL = "https://bbb.example/api"
