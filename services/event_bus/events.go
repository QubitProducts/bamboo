package event_bus

import (
	"bytes"
	"fmt"
	"sort"
)

var eventKeys = make(map[string]bool)

func init() {
	eventKeys["clientIp"] = true
	eventKeys["uri"] = true
	eventKeys["frameworkId"] = true
	eventKeys["master"] = true
	eventKeys["callbackUrl"] = true
	eventKeys["appId"] = true
	eventKeys["version"] = true
	eventKeys["taskId"] = true
	eventKeys["alive"] = true
	eventKeys["groupId"] = true
	eventKeys["id"] = true
	eventKeys["taskStatus"] = true
	eventKeys["slaveId"] = true
	eventKeys["host"] = true
	eventKeys["ports"] = true
	eventKeys["executorId"] = true
}

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
		if eventKeys[k] {
			keys = append(keys, k)
		}
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
