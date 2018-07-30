package log

import (
	stdlog "log"
	"os"
)

// StdoutHandler stdout log handler
type StdoutHandler struct {
	out *stdlog.Logger
}

// NewStdout create a stdout log handler
func NewStdout() *StdoutHandler {
	return &StdoutHandler{out: stdlog.New(os.Stdout, "", stdlog.LstdFlags)}
}

// Log stdout loging
func (h *StdoutHandler) Log(lv Level, f D) {
	switch lv {
	case debugLevel:
		h.out.Println(f)
	case infoLevel:
		h.out.Println(f)
	case warnLevel:
		h.out.Println(f)
	case errorLevel:
		h.out.Fatal(f)
	case fatalLevel:
		h.out.Fatal(f)
	default:
		h.out.Panic("unkonw log _level")
	}
}

// Close stdout loging
func (h *StdoutHandler) Close() (err error) {
	return
}
