package eqm

import (
    "fmt"
)

type EventQueueManager map[string]EventQueueInterface

func (eqm EventQueueManager) AddQueue(name string, queue EventQueueInterface) {
    // If we use NewEventQueueManager this should never happen but check for safety
    if eqm == nil {
        eqm = make(map[string]EventQueueInterface)
    }

    go queue.StartQueue()
    eqm[name] = queue
}

func (eqm EventQueueManager) ProcessEvent(e Event) (interface{}, error) {
    if _, ok := eqm[e.QueueName]; !ok {
        err := fmt.Errorf("eventqueue: No queue of name [%s] in manager to send action [%s].",
            e.QueueName, e.Action,
        )
        return "",  err
    }

    result := eqm[e.QueueName].SendEvent(e)
    return result.Result, result.Error
}

func NewEventQueueManager() EventQueueManager {
    return make(map[string]EventQueueInterface)
}
