package internal

// NodeError type ensure server to only stop if there are any fatal errors
type NodeError struct {
	error
	IsFatal bool
}

func NewNodeEror(e error, isFatal bool) NodeError { return NodeError{e, isFatal} }
