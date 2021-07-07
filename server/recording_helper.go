package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/server/bigbluebuttonapiwrapper/dataStructs"
	"github.com/thoas/go-funk"
)

const (
	prefixRecording     = "r_"
	prefixRecordingList = "r_list"
)

func (p *Plugin) GetRecordingWaitingList() ([]string, error) {
	data, appErr := p.API.KVGet(prefixRecordingList)
	if appErr != nil {
		p.API.LogError(fmt.Sprintf("Unable to fetch recording waiting list. Wrror: {%s}", appErr.Error()))
		return nil, errors.New(appErr.Error())
	}

	// This handles the case of no data present in KV store.
	// Happens on fresh installation.
	if len(data) == 0 {
		data = []byte("[]")
	}

	var meetingList *[]string
	err := json.Unmarshal(data, &meetingList)
	if err != nil {
		return nil, err
	}

	return *meetingList, nil
}

func (p *Plugin) GetMeetingWaitingForRecording(meetingID string) (*dataStructs.MeetingRoom, error) {
	data, appErr := p.API.KVGet(prefixRecording + meetingID)
	if appErr != nil {
		p.API.LogError(fmt.Sprintf("Unable to fetch list of meeting recording. Error: {%s}", appErr.Error()))
		return nil, errors.New(appErr.Error())
	}

	var meeting *dataStructs.MeetingRoom
	_ = json.Unmarshal(data, &meeting)
	return meeting, nil
}

func (p *Plugin) AddMeetingWaitingForRecording(meeting *dataStructs.MeetingRoom) error {
	if !p.config().ProcessRecordings {
		return nil
	}

	if err := p.saveMeetingForRecording(meeting); err != nil {
		return err
	}

	if err := p.addToRecordingWaitingList(meeting.MeetingID_); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) RemoveMeetingWaitingForRecording(meetingId string) error {
	if err := p.removeFromWaitingForRecording(meetingId); err != nil {
		return err
	}

	if err := p.deleteMeetingForRecording(meetingId); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) saveMeetingForRecording(meeting *dataStructs.MeetingRoom) error {
	data, err := json.Marshal(meeting)
	if err != nil {
		p.API.LogError(fmt.Sprintf("Unable to marshal meeting for storing for recording. Meeting ID: {%s}, error: {%s}", meeting.MeetingID_, err.Error()))
		return err
	}

	if appErr := p.API.KVSet(prefixRecording+meeting.MeetingID_, data); appErr != nil {
		p.API.LogError("Unable to store meeting for storing for recording. Meeting ID: {%s}, error: {%s}", meeting.MeetingID_, appErr.Error())
		return errors.New(appErr.Error())
	}

	return nil
}

func (p *Plugin) deleteMeetingForRecording(meetingID string) error {
	if appErr := p.API.KVDelete(prefixRecording + meetingID); appErr != nil {
		p.API.LogError("Unable to store meeting for storing for recording. Meeting ID: {%s}, error: {%s}", meetingID, appErr.Error())
		return errors.New(appErr.Error())
	}

	return nil
}

func (p *Plugin) addToRecordingWaitingList(meetingId string) error {
	meetingList, err := p.GetRecordingWaitingList()
	if err != nil {
		return err
	}

	meetingList = append(meetingList, meetingId)

	if err := p.saveRecordingList(meetingList); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) removeFromWaitingForRecording(meetingId string) error {
	meetingRecordingList, err := p.GetRecordingWaitingList()
	if err != nil {
		return err
	}

	i := funk.IndexOf(meetingRecordingList, meetingId)
	meetingRecordingList = append(meetingRecordingList[:i], meetingRecordingList[i+1:]...)

	if err := p.saveRecordingList(meetingRecordingList); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) saveRecordingList(meetingList []string) error {
	data, err := json.Marshal(meetingList)
	if err != nil {
		p.API.LogError(fmt.Sprintf("unable to marshal data for saving meeting recording list. Error: {%s}", err.Error()))
		return err
	}

	if appErr := p.API.KVSet(prefixRecordingList, data); appErr != nil {
		p.API.LogError(fmt.Sprintf("Unable to save recording waiting list in KV store. Error: {%s}", appErr.Error()))
		return errors.New(appErr.Error())
	}

	return nil
}
