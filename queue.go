package eqm

type EventReply struct {
    Result      interface{}     // The result from processing the Event
    Error       error           // Any errors while processing the corresponding event should be put here
}

func NewEventReply(result interface{}, err error) EventReply {
    return EventReply{result, err}
}

type Event struct {
    QueueName   string              // The name the EventQueue is registered with
    Action      string              // A key defining the function the EventQueuue should call
    Args        []interface{}       // The list of arguments (if any) the EventController needs to take action
    Reply       chan EventReply     // Channel to put the EventReply on
}

func NewEvent(queue, action string, args ...interface{}) Event {
    return Event{
        QueueName:  queue,
        Action:     action,
        Args:       args,
        Reply:      make(chan EventReply),
    }
}


type EventChannel chan Event

type EventMap map[string]func(Event) EventReply

type EventQueue struct {
    Name    string
    Events  EventChannel
    Before  func()         
    Map     EventMap
    Default func()
    After   func()
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

func (eq EventQueue) StartThread() {
    for {
        if eq.Default != nil {
            eq.processLoop()
        } else {
            eq.processNonDefaultLoop()
        }
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

func (eq EventQueue) processNonDefaultLoop() {
    if eq.Before != nil {
        eq.Before()
    }

    select {
    case event := <-eq.Events:
        if _, ok := eq.Map[event.Action]; ok {
            event.Reply <-eq.Map[event.Action](event)
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
