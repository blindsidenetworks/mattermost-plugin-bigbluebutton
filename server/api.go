package main

import (
	"net/http"
	"strings"
)

const (
	userIDHeaderName = "Mattermost-User-Id"
)

type HttpHandler func(w http.ResponseWriter, r *http.Request)

func (p *Plugin) GetRouteHandler(path string) HttpHandler {
	switch {
	case path == "/joinmeeting":
		return middlewareRequireAuth(p.handleJoinMeeting)
	case path == "/joinmeeting/external":
		return p.handleJoinMeetingExternalUser
	case strings.HasPrefix(path, "/endmeeting"):
		return middlewareRequireAuth(p.handleEndMeeting)
	case path == "/create":
		return middlewareRequireAuth(p.handleCreateMeeting)
	case strings.HasPrefix(path, "/recordingready"):
		return p.handleRecordingReady
	case path == "/getattendees":
		return middlewareRequireAuth(p.handleGetAttendeesInfo)
	case path == "/publishrecordings":
		return middlewareRequireAuth(p.handlePublishRecordings)
	case path == "/deleterecordingsconfirmation":
		return middlewareRequireAuth(p.handleDeleteRecordingsConfirmation)
	case path == "/deleterecordings":
		return middlewareRequireAuth(p.handleDeleteRecordings)
	case strings.HasPrefix(path, "/meetingendedcallback"):
		return p.handleImmediateEndMeetingCallback
	case path == "/ismeetingrunning":
		return middlewareRequireAuth(p.handleIsMeetingRunning)
	case path == "/config":
		return middlewareRequireAuth(p.handleGetConfig)
	case path == "/redirect":
		return p.handleRedirect
	case path == "/joininvite":
		return middlewareRequireAuth(p.handleJoinInvite)
	default:
		// for static files
		return p.handler.ServeHTTP
	}
}

func middlewareRequireAuth(handler HttpHandler) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(userIDHeaderName)
		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
		}

		handler(w, r)
	}
}
