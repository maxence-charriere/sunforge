package vi

// List of possible resource states.
const (
	// Undefined indicates that the status has not been set.
	Undefined Status = iota

	// Loading indicates that the resource is being fetched or initialized.
	Loading

	// Creating indicates that the resource is currently being created.
	Creating

	// Updating indicates that the resource is currently being updated.
	Updating

	// Deleting indicates that the resource is currently being deleted.
	Deleting

	// Synced indicates that the resource is aligned with its intended or desired state.
	Synced

	// Deleted indicates that the resource has been removed and is no longer active.
	Deleted
)

// Status represents the current lifecycle state of a resource or operation.
//
// It is typically used to describe whether a value is being loaded, created,
// updated, deleted, or has reached a terminal state such as synced or deleted.
type Status int

// String returns a human-readable representation of the status.
func (s Status) String() string {
	switch s {
	case Loading:
		return "Loading"

	case Creating:
		return "Creating"

	case Updating:
		return "Updating"

	case Deleting:
		return "Deleting"

	case Synced:
		return "Synced"

	case Deleted:
		return "Deleted"

	default:
		return "Undefined"
	}
}

// State represents a generic container that wraps a value of type T together
// with its current status and any associated error.
//
// It is commonly used to propagate computation or request results where a
// value, status flag, and error need to be tracked together.
//
// The type parameter T defines the type of the contained value.
type State[T any] struct {
	// Value holds the wrapped value of type T.
	Value T `json:",omitempty"`

	// Status represents the current state of the value or operation.
	Status Status `json:",omitempty"`

	// StatusLabel provides additional, human-readable context about the current
	// Status. It can be used to describe what the state is doing or why it is in
	// this status (for example, the type or phase of an update).
	StatusLabel string `json:",omitempty"`

	// Err contains an error associated with the state, if any.
	Err error `json:"-"`
}

// StatusString returns a human-readable representation of the state's status.
//
// If StatusLabel is set, it is returned to provide additional context.
// Otherwise, it falls back to the string representation of Status.
func (s State[T]) StatusString() string {
	if s.StatusLabel != "" {
		return s.StatusLabel
	}
	return s.Status.String()
}
