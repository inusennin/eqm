package eqm

import (
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestEventQueueManager(t *testing.T) {
    testMan := NewEventQueueManager()

    testQueue := NewEventQueue("TestQueue")
    testQueue.MapEvent("TestAction", func(e Event) EventReply {
        return EventReply{e.Action, nil}
    })
    testQueue.SetEventsBuffer(2) // Set the buffer so we can call this and not have to wait forever for results

    testEvent := NewEvent("TestQueue", "TestAction")

    _, err := testMan.ProcessEvent(testEvent)
    if assert.Error(t, err) {
        assert.True(t, strings.Contains(err.Error(), "eventqueue: No queue of name"))
    }

    testMan["DeadQueue"] = EventQueue{}
    _, err = testMan.ProcessEvent(NewEvent("DeadQueue", "NoAction"))
    if assert.Error(t, err) {
        assert.True(t, strings.Contains(err.Error(), "not initialized"))
    }

    testMan.AddQueue(testQueue)
    res, err := testMan.ProcessEvent(testEvent)
    if assert.NoError(t, err) {
        assert.Equal(t, "TestAction", res.(string))
    }
}
