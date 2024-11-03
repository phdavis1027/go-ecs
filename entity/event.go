package entity

type EventType int
const (
	// Event types
	EventTypeCollide EventType = iota
)

type EventPriority int
const (
	// Priority levels
	PriorityLow EventPriority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

type Event struct {
	Type   EventType
	Sender string
	Data   interface{}
}

type Mailboxes struct {
	low      chan Event
	normal   chan Event
	high     chan Event
	critical chan Event
}

func NewMailboxesWithCapacity(cap int) *Mailboxes {
	return &Mailboxes{
		low:      make(chan Event, cap),
		normal:   make(chan Event, cap),
		high:     make(chan Event, cap),
		critical: make(chan Event, cap),
	}
}

func (self *Mailboxes) SendToPriority(priority EventPriority, event Event) {
	switch priority {
	case PriorityLow:
		self.low <- event
	case PriorityNormal:
		self.normal <- event
	case PriorityHigh:
		self.high <- event
	case PriorityCritical:
		self.critical <- event
	}
}

func (self *Mailboxes) ReceiveFromPriority(priority EventPriority) Event {
	switch priority {
	case PriorityLow:
		return <-self.low
	case PriorityNormal:
		return <-self.normal
	case PriorityHigh:
		return <-self.high
	case PriorityCritical:
		return <-self.critical
	}
	return Event{}
}

func (self *Mailboxes) NextEvent() Event {
	select {
	case event := <-self.critical:
		return event
	case event := <-self.high:
		return event
	case event := <-self.normal:
		return event
	case event := <-self.low:
		return event
	}
}

func (self *Mailboxes) Close() {
	close(self.low)
	close(self.normal)
	close(self.high)
	close(self.critical)
}
