package decorator

import (
	"net/http"
	"testing"

	"github.com/papa-stiflera/go-cabbage"
	. "gopkg.in/check.v1"
)

type TestSuite struct {
	dummyClient *dummy
}

type dummy struct{}

func (d *dummy) Do(r *http.Request) (*http.Response, error) {
	return &http.Response{Request: r}, nil
}

type callOrder []string

var (
	_ = Suite(&TestSuite{})
)

func Test(t *testing.T) { TestingT(t) }

func (s *TestSuite) SetUpSuite(c *C) {
	s.dummyClient = &dummy{}
}

func (s *TestSuite) send(r *http.Request, decorator cabbage.DecoratorFunc) (*http.Response, error) {
	return cabbage.Decorate(s.dummyClient, decorator).Do(r)
}

func createDecorator(name string, order *callOrder) cabbage.DecoratorFunc {
	return func(c cabbage.HttpClientInterface) cabbage.HttpClientInterface {
		return cabbage.ClientFunc(func(r *http.Request) (*http.Response, error) {
			*order = append(*order, name)
			return c.Do(r)
		})
	}
}

func (s *TestSuite) TestDecorate(c *C) {
	var order callOrder
	client := cabbage.Decorate(s.dummyClient, createDecorator("A", &order), createDecorator("B", &order), createDecorator("C", &order))

	client.Do(&http.Request{})

	c.Assert(callOrder{"C", "B", "A"}, DeepEquals, order)
}

func (s *TestSuite) TestHeaderDecorator(c *C) {
	r, _ := http.NewRequest("GET", "/", nil)
	res, _ := s.send(r, HeaderDecorator("X-Foo", "123"))

	c.Assert(res.Request.Header.Get("X-Foo"), Equals, "123")
}

func (s *TestSuite) TestHeadersDecorator(c *C) {
	r, _ := http.NewRequest("GET", "/", nil)
	res, _ := s.send(r, HeadersDecorator(map[string]string{"X-Foo": "123", "X-Bar": "456"}))

	c.Assert(res.Request.Header.Get("X-Foo"), Equals, "123")
	c.Assert(res.Request.Header.Get("X-Bar"), Equals, "456")
}

func (s *TestSuite) TestPanicDecorator(c *C) {
	panicTriggerDecorator := func(c cabbage.HttpClientInterface) cabbage.HttpClientInterface {
		return cabbage.ClientFunc(func(r *http.Request) (res *http.Response, err error) {
			panic("oops")
		})
	}
	r, _ := http.NewRequest("GET", "/", nil)
	client := cabbage.Decorate(s.dummyClient, panicTriggerDecorator, RecoverDecorator())
	res, err := client.Do(r)

	c.Assert(res, IsNil)
	c.Assert(err, ErrorMatches, "*oops")
}
