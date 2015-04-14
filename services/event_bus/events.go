package event_bus

import (
	"bytes"
	"fmt"
	"sort"
)

type ZookeeperEvent struct {
	Source    string
	EventType string
}

type ServiceEvent struct {
	EventType string
}

type MarathonEvent struct {
	// EventType can be
	// api_post_event, status_update_event, subscribe_event
	EventType string
	Timestamp string
	plaintext string
}

func (me *MarathonEvent) Plaintext() string {
	return me.plaintext
}

func RestoreMarathonEvent(contentMap map[string]interface{}) (*MarathonEvent, bool) {
	if contentMap == nil {
		return nil, false
	}
	eventType, ok := hasStringValue(contentMap, "eventType")
	if !ok {
		return nil, false
	}
	timestamp, ok := hasStringValue(contentMap, "timestamp")
	if !ok {
		return nil, false
	}
	keys := make([]string, 0)
	for k, _ := range contentMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	marathonEvent := MarathonEvent{eventType, timestamp, ""}
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for _, key := range keys {
		buffer.WriteString(fmt.Sprintf("%s: %v, ", key, contentMap[key]))
	}
	buffer.Truncate(buffer.Len() - 2)
	buffer.WriteString("}")
	marathonEvent.plaintext = buffer.String()
	return &marathonEvent, true
}

func hasStringValue(contentMap map[string]interface{}, key string) (string, bool) {
	var value interface{}
	var str string
	var ok bool
	if value, ok = contentMap[key]; ok {
		str, ok = value.(string)
	}
	return str, ok
}
