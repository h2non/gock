package gock

// Register registers a new mock in the current mocks stack.
func Register(mock Mock) {
	g.Register(mock)
}

// GetAll returns the current stack of registered mocks.
func GetAll() []Mock {
	return g.GetAll()
}

// Exists checks if the given Mock is already registered.
func Exists(m Mock) bool {
	return g.Exists(m)
}

// Remove removes a registered mock by reference.
func Remove(m Mock) {
	g.Remove(m)
}

// Flush flushes the current stack of registered mocks.
func Flush() {
	g.Flush()
}

// Pending returns an slice of pending mocks.
func Pending() []Mock {
	return g.Pending()
}

// IsDone returns true if all the registered mocks has been triggered successfully.
func IsDone() bool {
	return g.IsDone()
}

// IsPending returns true if there are pending mocks.
func IsPending() bool {
	return g.IsPending()
}

// Clean cleans the mocks store removing disabled or obsolete mocks.
func Clean() {
	g.Clean()
}
