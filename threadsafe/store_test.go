package threadsafe

import (
	"testing"

	"github.com/nbio/st"
)

func TestStoreRegister(t *testing.T) {
	g := NewGock()
	defer after(g)
	st.Expect(t, len(g.mocks), 0)
	mock := g.New("foo").Mock
	g.Register(mock)
	st.Expect(t, len(g.mocks), 1)
	st.Expect(t, mock.Request().Mock, mock)
	st.Expect(t, mock.Response().Mock, mock)
}

func TestStoreGetAll(t *testing.T) {
	g := NewGock()
	defer after(g)
	st.Expect(t, len(g.mocks), 0)
	mock := g.New("foo").Mock
	store := g.GetAll()
	st.Expect(t, len(g.mocks), 1)
	st.Expect(t, len(store), 1)
	st.Expect(t, store[0], mock)
}

func TestStoreExists(t *testing.T) {
	g := NewGock()
	defer after(g)
	st.Expect(t, len(g.mocks), 0)
	mock := g.New("foo").Mock
	st.Expect(t, len(g.mocks), 1)
	st.Expect(t, g.Exists(mock), true)
}

func TestStorePending(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("foo")
	st.Expect(t, g.mocks, g.Pending())
}

func TestStoreIsPending(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("foo")
	st.Expect(t, g.IsPending(), true)
	g.Flush()
	st.Expect(t, g.IsPending(), false)
}

func TestStoreIsDone(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("foo")
	st.Expect(t, g.IsDone(), false)
	g.Flush()
	st.Expect(t, g.IsDone(), true)
}

func TestStoreRemove(t *testing.T) {
	g := NewGock()
	defer after(g)
	st.Expect(t, len(g.mocks), 0)
	mock := g.New("foo").Mock
	st.Expect(t, len(g.mocks), 1)
	st.Expect(t, g.Exists(mock), true)

	g.Remove(mock)
	st.Expect(t, g.Exists(mock), false)

	g.Remove(mock)
	st.Expect(t, g.Exists(mock), false)
}

func TestStoreFlush(t *testing.T) {
	g := NewGock()
	defer after(g)
	st.Expect(t, len(g.mocks), 0)

	mock1 := g.New("foo").Mock
	mock2 := g.New("foo").Mock
	st.Expect(t, len(g.mocks), 2)
	st.Expect(t, g.Exists(mock1), true)
	st.Expect(t, g.Exists(mock2), true)

	g.Flush()
	st.Expect(t, len(g.mocks), 0)
	st.Expect(t, g.Exists(mock1), false)
	st.Expect(t, g.Exists(mock2), false)
}
