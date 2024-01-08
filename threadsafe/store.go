package threadsafe

// Register registers a new mock in the current mocks stack.
func (g *Gock) Register(mock Mock) {
	if g.Exists(mock) {
		return
	}

	// Make ops thread safe
	g.storeMutex.Lock()
	defer g.storeMutex.Unlock()

	// Expose mock in request/response for delegation
	mock.Request().Mock = mock
	mock.Response().Mock = mock

	// Registers the mock in the global store
	g.mocks = append(g.mocks, mock)
}

// GetAll returns the current stack of registered mocks.
func (g *Gock) GetAll() []Mock {
	g.storeMutex.RLock()
	defer g.storeMutex.RUnlock()
	return g.mocks
}

// Exists checks if the given Mock is already registered.
func (g *Gock) Exists(m Mock) bool {
	g.storeMutex.RLock()
	defer g.storeMutex.RUnlock()
	for _, mock := range g.mocks {
		if mock == m {
			return true
		}
	}
	return false
}

// Remove removes a registered mock by reference.
func (g *Gock) Remove(m Mock) {
	for i, mock := range g.mocks {
		if mock == m {
			g.storeMutex.Lock()
			g.mocks = append(g.mocks[:i], g.mocks[i+1:]...)
			g.storeMutex.Unlock()
		}
	}
}

// Flush flushes the current stack of registered mocks.
func (g *Gock) Flush() {
	g.storeMutex.Lock()
	defer g.storeMutex.Unlock()
	g.mocks = []Mock{}
}

// Pending returns an slice of pending mocks.
func (g *Gock) Pending() []Mock {
	g.Clean()
	g.storeMutex.RLock()
	defer g.storeMutex.RUnlock()
	return g.mocks
}

// IsDone returns true if all the registered mocks has been triggered successfully.
func (g *Gock) IsDone() bool {
	return !g.IsPending()
}

// IsPending returns true if there are pending mocks.
func (g *Gock) IsPending() bool {
	return len(g.Pending()) > 0
}

// Clean cleans the mocks store removing disabled or obsolete mocks.
func (g *Gock) Clean() {
	g.storeMutex.Lock()
	defer g.storeMutex.Unlock()

	buf := []Mock{}
	for _, mock := range g.mocks {
		if mock.Done() {
			continue
		}
		buf = append(buf, mock)
	}

	g.mocks = buf
}
