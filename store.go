package gock

// mocks is internally used to store registered mocks.
var mocks = []Mock{}

// Register registers a new mock in the current mocks stack.
func Register(mock Mock) {
	if Exists(mock) {
		return
	}

	// Make ops thread safe
	mutex.Lock()
	defer mutex.Unlock()

	// Expose mock in request/response for delegation
	mock.Request().Mock = mock
	mock.Response().Mock = mock

	// Registers the mock in the global store
	mocks = append(mocks, mock)
}

// GetAll returns the current stack of registed mocks.
func GetAll() []Mock {
	mutex.Lock()
	defer mutex.Unlock()
	return mocks
}

// Exists checks if the given Mock is already registered.
func Exists(m Mock) bool {
	mutex.Lock()
	defer mutex.Unlock()

	for _, mock := range mocks {
		if mock == m {
			return true
		}
	}
	return false
}

// Remove removes a registered mock by reference.
func Remove(m Mock) {
	for i, mock := range mocks {
		if mock == m {
			mutex.Lock()
			mocks = append(mocks[:i], mocks[i+1:]...)
			mutex.Unlock()
		}
	}
}

// Flush flushes the current stack of registered mocks.
func Flush() {
	mutex.Lock()
	defer mutex.Unlock()
	mocks = []Mock{}
}

// Pending returns an slice of pending mocks.
func Pending() []Mock {
	Clean()
	return mocks
}

// IsDone returns true if all the registered mocks has been triggered successfully.
func IsDone() bool {
	return !IsPending()
}

// IsPending returns true if there are pending mocks.
func IsPending() bool {
	return len(Pending()) > 0
}

// Clean cleans the mocks store removing disabled or obsolete mocks.
func Clean() {
	mutex.Lock()
	defer mutex.Unlock()

	buf := []Mock{}
	for _, mock := range mocks {
		if mock.Done() {
			continue
		}
		buf = append(buf, mock)
	}

	mocks = buf
}
