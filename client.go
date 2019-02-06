package cabbage

import "net/http"

// HttpClientInterface sends http.Requests and returns http.Responses or errors in  case of failure.
type HttpClientInterface interface {
	Do(*http.Request) (*http.Response, error)
}

// ClientFunc is a function type that implements the Client interface.
type ClientFunc func(*http.Request) (*http.Response, error)

// Do Client interface support
func (f ClientFunc) Do(r *http.Request) (*http.Response, error) {
	return f(r)
}

// DecoratorFunc wraps a Client with extra behaviour.
type DecoratorFunc func(HttpClientInterface) HttpClientInterface

// Decorate decorates a Client c with all the given Decorators, in order.
func Decorate(c HttpClientInterface, ds ...DecoratorFunc) HttpClientInterface {
	result := c
	for _, decorate := range ds {
		result = decorate(result)
	}
	return result
}
