package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/QubitProducts/bamboo/configuration"
	eb "github.com/QubitProducts/bamboo/services/event_bus"
)

type EventSubscriptionAPI struct {
	Conf     *configuration.Configuration
	EventBus *eb.EventBus
}

func (sub *EventSubscriptionAPI) Callback(w http.ResponseWriter, r *http.Request) {
	payload, _ := ioutil.ReadAll(r.Body)
	marathonEvent, ok := convertJsonToMarathonEvent(payload)

	if !ok {
		log.Printf("Unable to decode JSON Marathon event request: %s \n", string(payload))
	}
	sub.EventBus.Publish(*marathonEvent)
	io.WriteString(w, "Got it!")
}

func convertJsonToMarathonEvent(payload []byte) (*eb.MarathonEvent, bool) {
	var content interface{}
	err := json.Unmarshal(payload, &content)
	if err != nil {
		log.Printf("An error occurred while decoding Marathon event request: %s\n (payload=%s)", err, string(payload))
		return nil, false
	}
	contentMap := content.(map[string]interface{})
	return eb.RestoreMarathonEvent(contentMap)
}
