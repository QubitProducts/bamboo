package api


import (
	"github.com/QubitProducts/bamboo/configuration"
	eb "github.com/QubitProducts/bamboo/services/event_bus"
	"net/http"
	"io"
	"log"
	"encoding/json"
	"io/ioutil"
)

type EventSubscriptionAPI struct {
	Conf *configuration.Configuration
	EventBus *eb.EventBus
}

func (sub *EventSubscriptionAPI) Callback(w http.ResponseWriter, r *http.Request) {
	var event eb.MarathonEvent

	payload, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(payload, &event)

	if err != nil {
		log.Printf("Unable to decode JSON Marathon Event request: %s \n", string(payload))
	}

	sub.EventBus.Publish(event)
	io.WriteString(w, "Got it!")
}
