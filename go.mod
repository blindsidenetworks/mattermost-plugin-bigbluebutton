module github.com/blindsidenetworks/mattermost-plugin-bigbluebutton

go 1.16

replace github.com/mattermost/viper v1.0.3-0.20181112161711-f99c30686b86 => github.com/mattermost/viper v1.0.3-0.20181112161711-f99c30686b86

require (
	github.com/mattermost/mattermost-plugin-api v0.0.16
	github.com/mattermost/mattermost-server/v5 v5.39.0
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron v1.2.0
	github.com/segmentio/ksuid v1.0.3
	github.com/thoas/go-funk v0.8.0
)
