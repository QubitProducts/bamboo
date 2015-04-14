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
	if me.plaintext == "" {
		me.plaintext = generatePlaintext(
			map[string]interface{}{
				"eventType": me.EventType,
				"timestamp": me.Timestamp,
			})
	}
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
	marathonEvent := MarathonEvent{eventType, timestamp, ""}
	marathonEvent.plaintext = generatePlaintext(contentMap)
	return &marathonEvent, true
}

func generatePlaintext(m map[string]interface{}) string {
	keys := make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for _, key := range keys {
		buffer.WriteString(fmt.Sprintf("%s: %v, ", key, m[key]))
	}
	buffer.Truncate(buffer.Len() - 2)
	buffer.WriteString("}")
	return buffer.String()
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
