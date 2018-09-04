package main

var fMsg *FlashBag
var fError *FlashBag

func init() {
	fMsg = &FlashBag{}
	fError = &FlashBag{}
}

type FlashBag struct {
	Value string
}

func (f *FlashBag) String() string {
	latest := f.Value
	f.Value = ""
	return latest
}

func flashMsg(msg string) {
	fMsg.Value = msg
}

func flashError(error string) {
	fError.Value = error
}
