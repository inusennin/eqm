package eqm

import (
    "fmt"
)

type EventQueueManager map[string]EventQueue

func (eqm EventQueueManager) AddQueue(queue EventQueue) {
    // If we use NewEventQueueManager this should never happen but check for safety
    if eqm == nil {
        eqm = make(map[string]EventQueue)
    }

    go queue.StartThread()
    eqm[queue.Name] = queue
}

func (eqm EventQueueManager) ProcessEvent(e Event) (interface{}, error) {
    if _, ok := eqm[e.QueueName]; !ok {
        err := fmt.Errorf("eventqueue: No queue of name [%s] in manager to send action [%s].",
            e.QueueName, e.Action,
        )
        return "",  err
    }

    if eqm[e.QueueName].Events == nil {
        err := fmt.Errorf("eventqueue: Queue of name [%s] not initialized.", e.QueueName)
        return "", err
    }

    // If we didn't define the reply buffer previously, we need to make sure we have it now
    if e.Reply == nil {
        e.Reply = make(chan EventReply)
    }

    eqm[e.QueueName].Events <-e
    result := <-e.Reply
    return result.Result, result.Error
}

func NewEventQueueManager() EventQueueManager {
    return make(map[string]EventQueue)
}
