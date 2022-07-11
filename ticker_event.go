package eqm

import (
    "fmt"
    "time"
)

// TickerEvents are blocking queues that use a ticker action in lieu of a default path for channels.
// To conserve memory/CPU, before and after are ommitted.
type TickerEventQueue struct {
    Name        string
    Events      EventChannel
    Map         EventMap
    Ticker      *time.Ticker
    TickerEvent func()
}

func (teq TickerEventQueue) StartQueue() {
    // We don't want to throw panics due to potential llack of user definitions
    if teq.Ticker.C == nil {
        teq.Ticker = time.NewTicker(time.Minute)
    }

    if teq.TickerEvent == nil {
        teq.TickerEvent = func() {
            // No function implemented
        }
    }

    for {
        select {
        case e := <-teq.Events:
            if _, ok := teq.Map[e.Action]; ok {
                e.Reply <-teq.Map[e.Action](e)
            }
        case _ = <-teq.Ticker.C:
            teq.TickerEvent()
        }
    }
}

func (teq *TickerEventQueue) SetEventsBuffer(buffSize int) {
    old := *teq

    old.Events = make(chan Event, buffSize)

    *teq = old
}

func (teq *TickerEventQueue) SetTickerDuration(dur time.Duration) {
    old := *teq

    old.Ticker = time.NewTicker(dur)

    *teq = old
}

func (teq TickerEventQueue) SendEvent(e Event) EventReply {
    if teq.Events == nil {
        return NewEventReply("", fmt.Errorf("eventqueue: Queue of name [%s] not initialized.", teq.Name))
    }

    if e.Reply == nil {
        e.Reply = make(chan EventReply)
    }
    teq.Events <-e
    return <-e.Reply
}

func (teq *TickerEventQueue) MapEvent(key string, eventFunction func(Event) EventReply) {
    old := *teq

    // We should never get this if we're using NewEventQueue, but check for safety
    if old.Map == nil {
        old.Map = make(map[string]func(Event) EventReply)
    }

    old.Map[key] = eventFunction

    *teq = old
}
