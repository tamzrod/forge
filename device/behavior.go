// Package device provides the device model.
// A device is a deterministic memory image with executable behavior and protocol interfaces.
package device

// Behavior is the interface for device behaviors.
// Behaviors read and write device memory.
type Behavior interface {
	// ID returns a unique identifier for this behavior.
	ID() string

	// Attach connects the behavior to a device.
	// Called when the behavior is added to a device.
	Attach(d *Device)

	// Detach disconnects the behavior from a device.
	// Called when the behavior is removed from a device.
	Detach()

	// Tick executes one simulation step.
	// The behavior should read from and write to device memory.
	Tick()
}

// BehaviorFunc is an adapter that allows a function to satisfy Behavior.
type BehaviorFunc struct {
	idFunc     func() string
	attachFunc func(*Device)
	detachFunc func()
	tickFunc   func()
}

// ID calls the underlying function.
func (f *BehaviorFunc) ID() string {
	if f.idFunc == nil {
		return ""
	}
	return f.idFunc()
}

// Attach calls the underlying function.
func (f *BehaviorFunc) Attach(d *Device) {
	if f.attachFunc != nil {
		f.attachFunc(d)
	}
}

// Detach calls the underlying function.
func (f *BehaviorFunc) Detach() {
	if f.detachFunc != nil {
		f.detachFunc()
	}
}

// Tick calls the underlying function.
func (f *BehaviorFunc) Tick() {
	if f.tickFunc != nil {
		f.tickFunc()
	}
}

// NewBehaviorFunc creates a Behavior from functions.
func NewBehaviorFunc(
	id string,
	tick func(),
) *BehaviorFunc {
	return &BehaviorFunc{
		idFunc:   func() string { return id },
		tickFunc: tick,
	}
}
