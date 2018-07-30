package log

// Level of severity.
type Level int

// Verbose .
type Verbose bool

const (
	debugLevel Level = iota
	infoLevel
	warnLevel
	errorLevel
	fatalLevel
)

var levelNames = [...]string{
	debugLevel: "DEBUG",
	infoLevel:  "INFO",
	warnLevel:  "WARN",
	errorLevel: "ERROR",
	fatalLevel: "FATAL",
}

// String implementation.
func (l Level) String() string {
	return levelNames[l]
}
