package domains

import (
	"../utils"
	"net/http"
	"time"
)

type CircuitBraker struct {
	State       string
	Errors      int
	Timeout     int
	ErrorsLimit int
}

var (
	CircuitBreakerGlobal = CircuitBraker{}
)

func (circuitBreaker *CircuitBraker) NewCircuitBreaker(state string, errors int, timeout int, errorsLimit int) {
	circuitBreaker.State = state
	circuitBreaker.Errors = errors
	circuitBreaker.Timeout = timeout
	circuitBreaker.ErrorsLimit = errorsLimit
}

func (circuitBreaker *CircuitBraker) SetState(state string) {
	circuitBreaker.State = state
}

func (circuitBreaker *CircuitBraker) Reset() {
	circuitBreaker.State = "CLOSE"
	circuitBreaker.Errors = 0
	circuitBreaker.Timeout = 15000
	circuitBreaker.ErrorsLimit = 3
}

func GetPing() chan bool {
	ch := make(chan bool)
	time.Sleep(time.Millisecond * 5000)
	go func() {
		urlPing := utils.UrlPing
	FOR:
		for i := 0; i < 5; i++ {
			if _, err := http.Get(urlPing); err == nil {
				ch <- true
				break FOR
			}
		}
		ch <- false
	}()
	return ch
}

func Circuit() (apiErr *utils.ApiError) {

	if CircuitBreakerGlobal.State == "HALF" {
		apiErr = &utils.ApiError{
			Message: "Server is trying connection(HALF)",
			Status:  http.StatusInternalServerError,
		}
		return apiErr
	}

	CircuitBreakerGlobal.SetState("OPEN")
	timer := time.Duration(CircuitBreakerGlobal.Timeout)
	timeout := time.After(time.Millisecond * timer)
	go func() {
	FOR:
		for {
			select {
			case <-timeout:
				CircuitBreakerGlobal.SetState("HALF")
				p := GetPing()
				ping := <-p
				if ping {
					CircuitBreakerGlobal.Reset()
					break FOR
				}
				CircuitBreakerGlobal.SetState("OPEN")
			}
		}
	}()
	if CircuitBreakerGlobal.State == "OPEN" {
		apiErr = &utils.ApiError{
			Message: "Server is bloqued try again in a few seconds",
			Status:  http.StatusInternalServerError,
		}
	}

	return apiErr
}
