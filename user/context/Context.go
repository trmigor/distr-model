package context

// Context represents user-defined context for a process.
type Context interface{}

// Common is a common context without any fields.
type Common struct{}

// SetX is a context for setting one value algorithm.
type SetX struct {
	X int
}

// Contexts is a map of all contexts available to processes.
var Contexts = map[string]Context{
	"Common": Common{},
	"SetX":   SetX{},
}
