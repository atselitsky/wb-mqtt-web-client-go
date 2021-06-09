package configs

import "time"

type MQTTmessage struct {
	Path    string
	Message string
	When    time.Time
}
