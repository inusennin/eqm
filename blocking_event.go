package eqm

import "fmt"

// Operates similarly to a normal EventQueue, but does not have a Before, After, or Default 
// and simply waits on the channel for messages. The intent is to use this to reduce the CPU footprint
// when the other behaviors aren't needed.
type BlockingEventQueue struct {
    Name    string
    Events  EventChannel
    Map     EventMap
}

func (beq BlockingEventQueue) StartQueue() {
    for {
        e := <-beq.Events

        if _, ok := beq.Map[e.Action]; ok {
            e.Reply <-beq.Map[e.Action](e)
        }
    }
}

func (beq BlockingEventQueue) SendEvent(e Event) EventReply {
    if beq.Events == nil {
        return NewEventReply("", fmt.Errorf("eventqueue: Queue of name [%s] not initialized.", beq.Name))
    }

    if e.Reply == nil {
        e.Reply = make(chan EventReply)
    }
    beq.Events <-e
    return <-e.Reply
}

func (beq *BlockingEventQueue) SetEventsBuffer(buffSize int) {
    old := *beq

    old.Events = make(chan Event, buffSize)

    *beq = old
}

func (beq *BlockingEventQueue) MapEvent(key string, eventFunction func(Event) EventReply) {
    old := *beq

    // We should never get this if we're using NewEventQueue, but check for safety
    if old.Map == nil {
        old.Map = make(map[string]func(Event) EventReply)
    }

    old.Map[key] = eventFunction

    *beq = old
}

func NewBlockingEventQueue(name string) BlockingEventQueue {
    return BlockingEventQueue{
        Name:   name,
        Events: make(chan Event),
        Map:    make(map[string]func(Event) EventReply),
    }
}
