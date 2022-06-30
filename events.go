package eqm

type EventChannel chan Event
type EventMap map[string]func(Event) EventReply

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
