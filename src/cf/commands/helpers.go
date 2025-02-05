package commands

import (
	term "cf/terminal"
	"errors"
	"fmt"
	"github.com/cloudfoundry/loggregatorlib/logmessage"
	"strconv"
	"strings"
	"time"
)

const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
)

func byteSize(bytes int) string {
	unit := ""
	value := float32(bytes)

	switch {
	case bytes >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	}

	stringValue := fmt.Sprintf("%.1f", value)
	stringValue = strings.TrimRight(stringValue, ".0")
	return fmt.Sprintf("%s%s", stringValue, unit)
}

func bytesFromString(s string) (bytes int, err error) {
	unit := string(s[len(s)-1])
	stringValue := s[0 : len(s)-1]

	value, err := strconv.ParseInt(stringValue, 10, 0)
	if err != nil {
		return
	}

	var byteValue int64

	switch unit {
	case "T":
		byteValue = value * TERABYTE
	case "G":
		byteValue = value * GIGABYTE
	case "M":
		byteValue = value * MEGABYTE
	case "K":
		byteValue = value * KILOBYTE
	}

	if byteValue == 0 {
		err = errors.New("Could not parse byte string")
	}

	bytes = int(byteValue)
	return
}

func coloredState(state string) (colored string) {
	switch state {
	case "started", "running":
		colored = term.SuccessColor("running")
	case "stopped":
		colored = term.StoppedColor("stopped")
	case "flapping":
		colored = term.WarningColor("flapping")
	case "starting":
		colored = term.AdvisoryColor("starting")
	default:
		colored = term.FailureColor(state)
	}

	return
}

func logMessageOutput(appName string, lm logmessage.LogMessage) string {
	sourceTypeNames := map[logmessage.LogMessage_SourceType]string{
		logmessage.LogMessage_CLOUD_CONTROLLER: "API",
		logmessage.LogMessage_ROUTER:           "Router",
		logmessage.LogMessage_UAA:              "UAA",
		logmessage.LogMessage_DEA:              "Executor",
		logmessage.LogMessage_WARDEN_CONTAINER: "App",
	}

	sourceType, _ := sourceTypeNames[*lm.SourceType]
	sourceId := "?"
	if lm.SourceId != nil {
		sourceId = *lm.SourceId
	}
	msg := lm.GetMessage()

	t := time.Unix(0, *lm.Timestamp)
	timeString := t.Format("Jan 2 15:04:05")

	channel := ""
	if lm.MessageType != nil && *lm.MessageType == logmessage.LogMessage_ERR {
		channel = "STDERR "
	}

	return fmt.Sprintf("%s %s %s/%s %s%s", timeString, appName, sourceType, sourceId, channel, msg)
}
