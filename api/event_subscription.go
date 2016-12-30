package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cloverstd/bamboo/configuration"
	eb "github.com/cloverstd/bamboo/services/event_bus"
)

type EventSubscriptionAPI struct {
	Conf     *configuration.Configuration
	EventBus *eb.EventBus
}

func (sub *EventSubscriptionAPI) Callback(w http.ResponseWriter, r *http.Request) {
	payload, _ := ioutil.ReadAll(r.Body)

	sub.Notify(payload)

	io.WriteString(w, "Got it!")
}

func (sub *EventSubscriptionAPI) Notify(payload []byte) {

	var event eb.MarathonEvent
	err := json.Unmarshal(payload, &event)

	if err != nil {
		log.Printf("Unable to decode JSON Marathon Event request: %s \n", string(payload))
	}

	sub.EventBus.Publish(event)
}
