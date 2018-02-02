package event

const (
	// EvtCodeGOPATHNotFound indicates an event: not found $GOPATH env variable
	EvtCodeGOPATHNotFound = iota
	// EvtCodeGOROOTNotFound indicates an event: not found $GOROOT env variable
	EvtCodeGOROOTNotFound
	// EvtCodeGocodeNotFound indicates an event: not found gocode
	EvtCodeGocodeNotFound
	// EvtCodeIDEStubNotFound indicates an event: not found gotools
	EvtCodeIDEStubNotFound
	// EvtCodeServerInternalError indicates an event: server internal error
	EvtCodeServerInternalError
)

const maxQueueLength = 10

var EventQueue = make(chan *Event, maxQueueLength)

type Event struct {
	Code int         `json:"code"` // event code
	Sid  string      `json:"sid"`  // IDE session id related
	Data interface{} `json:"data"` // event data
}

type UserEventQueue struct {
	Sid      string      // IDE session id related
	Queue    chan *Event // queue
	Handlers []Handler   // event handlers
}

type Handler interface {
	Handle(event *Event)
}

func (ev *Event) Copy() *Event {
	return &Event{
		Code: ev.Code,
		Sid:  ev.Sid,
		Data: ev.Data,
	}
}
