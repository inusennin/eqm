package eqm

import "fmt"

type SingleEventQueue struct {
    Name    string
    Events  EventChannel
    Call    func(Event) EventReply
}

func (seq SingleEventQueue) StartQueue() {
    for {
        e := <-seq.Events
        e.Reply<-seq.Call(e)
    }
} 

func  (seq SingleEventQueue) SendEvent(e Event) EventReply {
    if seq.Call == nil {
        return NewEventReply("", fmt.Errorf("eventqueue: Call for SingleEventQueue [%s] not initialized.", seq.Name))
    }

    if seq.Events == nil {
        return NewEventReply("", fmt.Errorf("eventqueue: Queue of name [%s] not initialized.", seq.Name))
    }

    if e.Reply == nil {
        e.Reply = make(chan EventReply)
    }

    seq.Events <-e
    return <-e.Reply
}

func (seq *SingleEventQueue) SetEventsBuffer(buffSize int) {
    old := *seq

    old.Events = make(chan Event, buffSize)

    *seq = old
}

func NewSingleEventQueue(name string) SingleEventQueue {
    return SingleEventQueue {
        Name:   name,
        Events: make(EventChannel),
    }
}
