package eqm

import "fmt"

type EventQueueInterface interface {
    StartQueue()
    SendEvent(e Event) EventReply
}

type EventQueue struct {
    Name    string
    Events  EventChannel
    Before  func()         
    Map     EventMap
    Default func()
    After   func()
}

func (eq EventQueue) SendEvent(e Event) EventReply {
    if eq.Events == nil {
        return NewEventReply("", fmt.Errorf("eventqueue: Queue of name [%s] not initialized.", eq.Name))
    }

    if e.Reply == nil {
        e.Reply = make(chan EventReply)
    }
    eq.Events <-e
    return <-e.Reply
}

func (eq *EventQueue) SetEventsBuffer(buffSize int) {
    old := *eq

    old.Events = make(chan Event, buffSize)

    *eq = old
}

func (eq *EventQueue) MapEvent(key string, eventFunction func(Event) EventReply) {
    old := *eq

    // We should never get this if we're using NewEventQueue, but check for safety
    if old.Map == nil {
        old.Map = make(map[string]func(Event) EventReply)
    }

    old.Map[key] = eventFunction

    *eq = old
}

func (eq EventQueue) StartQueue() {
    for {
        eq.processLoop()
    }
}

func (eq EventQueue) processLoop() {
    if eq.Before != nil {
        eq.Before()
    }

    select {
    case event := <-eq.Events:
        if _, ok := eq.Map[event.Action]; ok {
            event.Reply <-eq.Map[event.Action](event)
        }
    default:
        if eq.Default != nil {
            eq.Default()
        }
    }

    if eq.After != nil {
        eq.After()
    }
}

func NewEventQueue(name string) EventQueue {
    return EventQueue{
        Name:   name,
        Events: make(chan Event),
        Map:    make(map[string]func(Event) EventReply),
    }
}
