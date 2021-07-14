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

package main

import (
	"fmt"
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/helpers"
	"github.com/mattermost/mattermost-plugin-api/cluster"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/mattermost"

	bbbAPI "github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/api"
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/robfig/cron"
)

var PluginVersion string

const closeWindowScript = `<!doctype html>
				<html>
						<head><script>
							window.onload = function load() {window.open('', '_self', ''); window.close();};
						</script></head>
					<body></body>
				</html>`

const (
	jobInterval = 2 * time.Minute
)

type Plugin struct {
	plugin.MattermostPlugin

	c             *cron.Cron
	configuration atomic.Value
	job           *cluster.Job
	handler       http.Handler
}

//OnActivate runs as soon as plugin activates
func (p *Plugin) OnActivate() error {
	mattermost.API = p.API
	if err := p.OnConfigurationChange(); err != nil {
		p.API.LogError(err.Error())
		return err
	}

	config := p.config()
	if err := config.IsValid(); err != nil {
		p.API.LogError(err.Error())
		return err
	}

	bbbAPI.SetAPI(config.BaseURL+"/", config.Secret)

	helpers.PluginVersion = PluginVersion

	if config.ProcessRecordings {
		if err := p.schedule(); err != nil {
			return err
		}
	}

	if err := p.setupStaticFileServer(); err != nil {
		return err
	}

	// register slash command '/bbb' to create a meeting
	return p.API.RegisterCommand(&model.Command{
		Trigger:          "bbb",
		AutoComplete:     true,
		AutoCompleteDesc: "Create a BigBlueButton meeting",
	})
}

func (p *Plugin) schedule() error {
	if p.job != nil {
		if err := p.job.Close(); err != nil {
			return err
		}
	}

	job, err := cluster.Schedule(
		p.API,
		"BigBlueButtonRecordingProcessor",
		cluster.MakeWaitForRoundedInterval(jobInterval),
		p.Loopthroughrecordings,
	)

	if err != nil {
		p.API.LogError(fmt.Sprintf("Unable to schedule job for processing recordings. Error: {%s}", err.Error()))
		return err
	}

	p.job = job
	return nil
}

//following method is to create a meeting from '/bbb' slash command
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	meetingpointer := new(dataStructs.MeetingRoom)

	if err := p.PopulateMeeting(meetingpointer, nil, "", args.UserId, args.ChannelId); err != nil {
		return nil, model.NewAppError("ExecuteCommand", "Please provide a 'Site URL' in Settings > General > Configuration", nil, err.Error(), http.StatusInternalServerError)
	}

	p.createStartMeetingPost(args.UserId, args.ChannelId, meetingpointer)
	if err := p.SaveMeeting(meetingpointer); err != nil {
		return nil, model.NewAppError("ExecuteCommand", "Unable so save meeting", nil, err.Error(), http.StatusInternalServerError)
	}

	return &model.CommandResponse{}, nil

}

//this is the router to handle our server calls
//methods are all in responsehandlers.go
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {

	config := p.config()
	if err := config.IsValid(); err != nil {
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
		return
	}

	path := r.URL.Path
	if path == "/joinmeeting" {
		p.handleJoinMeeting(w, r)
	} else if strings.HasPrefix(path, "/endmeeting") {
		p.handleEndMeeting(w, r)
	} else if path == "/create" {
		p.handleCreateMeeting(w, r)
	} else if strings.HasPrefix(path, "/recordingready") {
		p.handleRecordingReady(w, r)
	} else if path == "/getattendees" {
		p.handleGetAttendeesInfo(w, r)
	} else if path == "/publishrecordings" {
		p.handlePublishRecordings(w, r)
	} else if path == "/deleterecordingsconfirmation" {
		p.handleDeleteRecordingsConfirmation(w, r)
	} else if path == "/deleterecordings" {
		p.handleDeleteRecordings(w, r)
	} else if strings.HasPrefix(path, "/meetingendedcallback") {
		p.handleImmediateEndMeetingCallback(w, r, path)
	} else if path == "/ismeetingrunning" {
		p.handleIsMeetingRunning(w, r)
	} else if path == "/redirect" {
		// html file to automatically close a window
		// nolint:staticcheck
		_, _ = fmt.Fprintf(w, closeWindowScript)
	} else {
		p.handler.ServeHTTP(w, r)
	}
}

func (p *Plugin) setupStaticFileServer() error {
	exe, err := os.Executable()
	if err != nil {
		p.API.LogError("Couldn't find plugin executable path", err, nil)
		return err
	}

	p.handler = http.FileServer(http.Dir(filepath.Dir(exe) + "/../assets"))
	return nil
}

func (p *Plugin) OnDeactivate() error {
	//on deactivate, save meetings details, stop check recordings looper, destroy webhook
	p.c.Stop()
	return nil
}

func main() {
	plugin.ClientMain(&Plugin{})
}
