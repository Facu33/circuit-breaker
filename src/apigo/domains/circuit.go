package domains

import (
	"strconv"
	"time"
)

type State int
type StateService int

const (
	STATE_CLOSE         State        = 1
	STATE_HALF_OPEN     State        = 2
	STATE_OPEN          State        = 3
	SERVICE_AVAILABLE   StateService = 1
	SERVICE_UNAVAILABLE StateService = 0
)

type CounterResult struct {
	TotalRequests uint32
	TotalSucceses uint32
	TotalFailures uint32
	TotalRejects  uint32
}

type Breaker struct {
	options        OptionsConfig
	counter        CounterResult
	events         []chan BreakerEvent
	state          State
	consecFailures uint32
}

type BreakerEvent struct {
	Code    StateService
	Message string
}
type OptionsConfig struct {
	Attempts     uint32
	TimeoutState time.Duration
	LimitFailure uint32
	MaxRequests  uint32
	NameService  string
}

func NewBreaker(options *OptionsConfig) *Breaker {
	if options == nil {
		options = &OptionsConfig{}
	}
	if options.Attempts == 0 {
		options.Attempts = 3
	}
	if options.TimeoutState == 0 {
		options.TimeoutState = 4 * time.Second
	}
	if options.LimitFailure == 0 {
		options.LimitFailure = 5
	}
	if options.MaxRequests == 0 {
		options.MaxRequests = 10000000
	}
	if options.NameService == "" {
		options.NameService = strconv.Itoa(int(time.Now().UnixNano()))
	}
	return &Breaker{options: *options, counter: CounterResult{}, state: STATE_CLOSE, consecFailures: 0}
}

func (cb *Breaker) Subscriber() <-chan BreakerEvent {
	evenReader := make(chan BreakerEvent)
	outputChannel := make(chan BreakerEvent, 100)
	go func() {
		for event := range evenReader {
			select {
			case outputChannel <- event:
			default:
				<-outputChannel
				outputChannel <- event
			}
		}
	}()
	cb.events = append(cb.events, evenReader)
	return outputChannel
}


