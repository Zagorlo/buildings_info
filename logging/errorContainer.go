package logging

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"runtime"
	"time"

	"github.com/gdexlab/go-render/render"
)

type ErrorContainer struct {
	err    error
	uuid   uuid.UUID
	time   time.Time
	args   []interface{}
	status int
}

func (ec *ErrorContainer) Err() string {
	if ec != nil && ec.err != nil {
		return ec.err.Error()
	}

	return ""
}

func (ec *ErrorContainer) UUID() string {
	return ec.uuid.String()
}

func (ec *ErrorContainer) Status() int {
	return ec.status
}

func (ec *ErrorContainer) trace() {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	Logging.Error("%s:%d %s, with error: %v with args %v additional info:{uuid: %v ; time: %v}", frame.File, frame.Line, frame.Function, ec.err, render.AsCode(ec.args), ec.uuid, ec.time)
}

func (ec *ErrorContainer) NotNil() bool {
	if ec != nil && ec.err != nil {
		ec.trace()

		return true
	}

	return false
}

func NewErrorContainer(ctx context.Context, err error, status int, args ...interface{}) *ErrorContainer {
	uuidParsed := GetRequestUUID(ctx)

	ec := &ErrorContainer{
		err:    err,
		uuid:   uuidParsed,
		time:   time.Now(),
		args:   args,
		status: status,
	}

	if ec.err != nil {
		ec.trace()
	}

	return ec
}

func NilErrorContainer() *ErrorContainer {
	return &ErrorContainer{
		err:  nil,
		args: nil,
	}
}

var (
	NilErrorContainerVar = NilErrorContainer()
)
