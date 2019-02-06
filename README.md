# Cabbage

It`s middleware approach for using http.Client. You can wrap your client with different functionality: 

 - log every request
 - append auth headers
 - use http cache
 - and whatever you want!
 
**Just like a cabbage!**

![](http://2.bp.blogspot.com/-LtmW_ktxtXU/Um28ElCtQnI/AAAAAAAAB04/aaXWbmTdbnE/s1600/cabbage.png)

## Client interface

Internal http package doesn`t have any interface for http clients, so Cabbage provides very simple client interface:
```go
type HttpClientInterface interface {
	Do(*http.Request) (*http.Response, error)
}
```
`http.Client` supports it out of the box.

## Decorators

Decorator is like a middleware:
```go
type DecoratorFunc func(HttpClientInterface) HttpClientInterface
```

Cabbage provides some helpful decorators for you:

- ```HeadersDecorator(values map[string]string)``` Adds headers to requests
- ```HeaderDecorator(name, value string)``` Like headers, but add only one header. 
- ```RecoverDecorator()``` Converts all panics into errors
- ```BaseURLDecorator(baseURL string)``` Replaces scheme and host to baseURL value.

## Usage

```go
client := http.DefaultClient

decoratedClient := cabbage.Decorate(
    client,
    decorator.HeaderDecorator("X-Auth", "123"),
    decorator.RecoverDecorator(), // better to place it last to recover panics from decorators too
)
```

## Create your own decorator

There are two ways of creating new decorators.

You can create some new struct:
```go
struct AwesomeStuffClient {
    client cabbage.Client
}

func(c *AwesomeStuffClient) Do(r *http.Request) (*http.Response, error) {
    // some stuff before call
    res, err := c.client.Do(r)
    // some stuff after call
    
    return res, err
}

func AwesomeStuffDecorator(c cabbage.HttpClientInterface) cabbage.HttpClientInterface {
    return &AwesomeStuffClient{client: c}
}
```

Or you can create just a function with type:
```go 
type ClientFunc func(*http.Request) (*http.Response, error)
```

So the same example will be looks like:
```go
func AwesomeStuffDecorator(c cabbage.HttpClientInterface) cabbage.HttpClientInterface {
	return cabbage.ClientFunc(func(r *http.Request) (*http.Response, error) {
		// some stuff before call
        res, err := c.client.Do(r)
        // some stuff after call
        
        return res, err
	})
}
```

Sometimes it`s required to pass some params in decorator, for details see Headers decorator.
