package vmio

const (
	FATAL = iota
	ERROR
	WARNING
	INFO
	DEBUG
	TRACE
	ALL
)

//VMLogger is the logger for the VM
type VMLogger interface {
	Log(level uint16, s string)
}
