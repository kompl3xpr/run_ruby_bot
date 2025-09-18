package service

import (
	"fmt"
	"html"
	"log/slog"
	"run_ruby_bot/core"
)

type SourceCode string
type HTMLMessage string

func shorten(msg string) string {
	if len(msg) > core.Cfg().Bot.MessageMaxLength {
		return msg[:core.Cfg().Bot.MessageMaxLength] + "\n..."
	}
	return msg
}

var taskSem chan struct{}
var waitSem chan struct{}

func InitService() {
	slog.Info("initializing task pool and queue...")
	taskSem = make(chan struct{}, core.Cfg().Task.TaskPoolCapacity)
	waitSem = make(chan struct{}, core.Cfg().Task.QueueCapacity)
}

func RunInterpretTask(srcCode SourceCode, ch chan<- HTMLMessage) {
	defer close(ch)

	select {
	case waitSem <- struct{}{}:
		ch <- HTMLMessage("<em>in queue...</em>")
	default:
		ch <- HTMLMessage("<em>too many requests</em>")
		return
	}

	taskSem <- struct{}{}

	ch <- HTMLMessage("<em>running...</em>")
	result := core.Interpret(string(srcCode))
	msg := html.EscapeString(shorten(result.Msg))
	switch result.Status {
	case core.SuccessWithOutput:
		ch <- HTMLMessage(fmt.Sprintf("<code>%v</code>", msg))
	case core.SuccessWithoutOutput:
		ch <- HTMLMessage("<em>no output</em>")
	case core.Timeout:
		ch <- HTMLMessage("<em>timeout</em>")
	case core.CodeError:
		ch <- HTMLMessage(fmt.Sprintf("<pre>%v</pre>", msg))
	case core.InternalError:
		ch <- HTMLMessage(fmt.Sprintf("<em>internal error: %v</em>", msg))
	}

	<-taskSem
	<-waitSem
}
