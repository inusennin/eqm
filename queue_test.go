package eqm

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestEventQueueBefore(t *testing.T) {
    chanResult := make(chan string)
    testQueue := NewEventQueue("TestQueue")
    testQueue.Before = func() {
        chanResult <-"test_result"
    }
    
    go testQueue.processLoop()
    assert.Equal(t, "test_result", <-chanResult)
}

func TestEventQueueMappedEvent(t *testing.T) {
    testQueue := NewEventQueue("TestQueue")

    testQueue.SetEventsBuffer(2) // Set the buffer so we can call this and not have to wait forever for results
    testQueue.MapEvent("TestAction", func(e Event) EventReply {
        result := e.Args[0].(int) + e.Args[1].(int)
        return EventReply{result, nil}
    })
    
    testEvent := NewEvent("TestQueue", "TestAction", 1, 2)
    testQueue.Events <-testEvent
    go testQueue.processLoop()

    res := <-testEvent.Reply
    if assert.NoError(t, res.Error) {
        assert.Equal(t, 3, res.Result.(int))
    }
}

func TestEventQueueDefault(t *testing.T) {
    chanResult := make(chan string)
    testQueue := NewEventQueue("TestQueue")
    testQueue.Default = func() {
        chanResult <-"test_result"
    }
    
    go testQueue.processLoop()
    assert.Equal(t, "test_result", <-chanResult)
}

func TestEventQueueAfter(t *testing.T) {
    chanResult := make(chan string)
    testQueue := NewEventQueue("TestQueue")
    testQueue.After = func() {
        chanResult <-"test_result"
    }
    
    go testQueue.processLoop()
    assert.Equal(t, "test_result", <-chanResult)
}

func TestEventQueueProcessLoop(t *testing.T) {
    chanResult := make(chan string, 3)
    testQueue := NewEventQueue("TestQueue")
        
    testQueue.Before = func() {
        chanResult <-"Before"
    }
    testQueue.Default = func() {
        chanResult <-"Default"
    }
    testQueue.After = func() {
        chanResult <-"After"
    }
    testQueue.MapEvent("TestAction", func(e Event) EventReply {
        return EventReply{"TestAction", nil}
    })
    
    // We're going to frontload an action so we can secure getting the default  at the correct time
    testQueue.SetEventsBuffer(2) // Set the buffer so we can call this and not have to wait forever for results
    testEvent := NewEvent("TestQueue", "TestAction")
    testQueue.Events <-testEvent

    go testQueue.StartThread()
    assert.Equal(t, "Before", <-chanResult)

    res := <-testEvent.Reply
    if assert.NoError(t, res.Error) {
        assert.Equal(t, "TestAction", res.Result.(string))
    }
    assert.Equal(t, "After", <-chanResult)

    assert.Equal(t, "Before", <-chanResult)
    assert.Equal(t, "Default", <-chanResult)
    assert.Equal(t, "After", <-chanResult)
}
