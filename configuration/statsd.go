package configuration

import (
	"github.com/peterbourgon/g2s"
	"log"
	"strings"
	"time"
)

type StatsD struct {
	Enabled bool
	Host    string
	Prefix  string

	Client g2s.Statter
}

func (s *StatsD) CreateClient() {
	if s.Enabled && s.Client == nil {
		log.Println("StatsD is enabled")
		client, err := g2s.Dial("udp", s.Host)
		if err != nil {
			log.Fatalf("Cannot connect to statsd server %v: %v ", s.Host, err)
		}
		s.Client = client
	}

}

func (s *StatsD) Increment(sampleRate float32, bucket string, n int) {
	if s.Client != nil {
		s.Client.Counter(sampleRate, fullBucket(s.Prefix, bucket), n)
	}
}

func (s *StatsD) Timing(sampleRate float32, bucket string, d time.Duration) {
	if s.Client != nil {
		s.Client.Timing(sampleRate, fullBucket(s.Prefix, bucket), d)
	}
}

func (s *StatsD) Gauge(sampleRate float32, bucket string, value string) {
	if s.Client != nil {
		s.Client.Gauge(sampleRate, fullBucket(s.Prefix, bucket), value)
	}
}

func fullBucket(prefix string, bucket string) string {
	if strings.HasSuffix(prefix, ".") {
		return prefix + bucket
	} else {
		return prefix + "." + bucket
	}
}
