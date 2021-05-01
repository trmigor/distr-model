package errors

// ErrorCode represents errors
type ErrorCode int

const (
	// OK is an error code for success result.
	OK ErrorCode = iota

	// NotExpectedType is an error code for not expected type of argument.
	NotExpectedType

	// ReadAttemptOutOfBounds is an error code for read attempt with invalid index.
	ReadAttemptOutOfBounds

	// ObjectIsNil is an error code for use attempt of a nil object.
	ObjectIsNil

	// ResourceNotFound is an error code for use attempt of not found resource.
	ResourceNotFound

	// ResourceInUse is an error code for use attempt of already used resource.
	ResourceInUse

	// PrematureEndOfStream is an error code for premature end of a stream.
	PrematureEndOfStream

	// SizeTooBig is an error code for too big requested size.
	SizeTooBig

	// YetNotImplemented is an error code for yet not implemented features.
	YetNotImplemented

	// DuplicateItems is an error code for item duplicate attempt.
	DuplicateItems

	// ItemNotFound is an error code for use attempt of non-existent item.
	ItemNotFound

	// QueueIsEmpty is an error code for use attempt of empty queue.
	QueueIsEmpty

	// SocketError is an error code for socket errors.
	SocketError

	// ConnectionFailed is an error code for failed connection.
	ConnectionFailed

	// ConnectionInUse is an error code for use attempt of already used connection
	ConnectionInUse

	// TimeOut is an error code for lost messages.
	TimeOut
)
