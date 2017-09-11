package abstract

type Logger interface {
	Info(content ...interface{})

	Infof(format string, content ...interface{})

	Warn(content ...interface{})

	Warnf(format string, content ...interface{})

	Notice(content ...interface{})

	Noticef(format string, content ...interface{})

	Debug(content ...interface{})

	Debugf(format string, content ...interface{})

	Error(content ...interface{})

	Errorf(format string, content ...interface{})
}

//EmptyLog is a stub Object for Logger interface
type EmptyLog struct {
}

func (*EmptyLog) Info(content ...interface{}) {

}

func (*EmptyLog) Infof(format string, content ...interface{}) {

}

func (*EmptyLog) Warn(content ...interface{}) {

}

func (*EmptyLog) Warnf(format string, content ...interface{}) {

}

func (*EmptyLog) Notice(content ...interface{}) {

}

func (*EmptyLog) Noticef(format string, content ...interface{}) {

}

func (*EmptyLog) Debug(content ...interface{}) {

}

func (*EmptyLog) Debugf(format string, content ...interface{}) {

}

func (*EmptyLog) Error(content ...interface{}) {

}

func (*EmptyLog) Errorf(format string, content ...interface{}) {

}
