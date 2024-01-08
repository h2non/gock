package gock

import (
	"testing"

	"github.com/nbio/st"
)

func TestStoreRegister(t *testing.T) {
	defer after()
	st.Expect(t, len(GetAll()), 0)
	mock := New("foo").Mock
	Register(mock)
	st.Expect(t, len(GetAll()), 1)
	st.Expect(t, mock.Request().Mock, mock)
	st.Expect(t, mock.Response().Mock, mock)
}

func TestStoreGetAll(t *testing.T) {
	defer after()
	st.Expect(t, len(GetAll()), 0)
	mock := New("foo").Mock
	store := GetAll()
	st.Expect(t, len(GetAll()), 1)
	st.Expect(t, len(store), 1)
	st.Expect(t, store[0], mock)
}

func TestStoreExists(t *testing.T) {
	defer after()
	st.Expect(t, len(GetAll()), 0)
	mock := New("foo").Mock
	st.Expect(t, len(GetAll()), 1)
	st.Expect(t, Exists(mock), true)
}

func TestStorePending(t *testing.T) {
	defer after()
	New("foo")
	st.Expect(t, GetAll(), Pending())
}

func TestStoreIsPending(t *testing.T) {
	defer after()
	New("foo")
	st.Expect(t, IsPending(), true)
	Flush()
	st.Expect(t, IsPending(), false)
}

func TestStoreIsDone(t *testing.T) {
	defer after()
	New("foo")
	st.Expect(t, IsDone(), false)
	Flush()
	st.Expect(t, IsDone(), true)
}

func TestStoreRemove(t *testing.T) {
	defer after()
	st.Expect(t, len(GetAll()), 0)
	mock := New("foo").Mock
	st.Expect(t, len(GetAll()), 1)
	st.Expect(t, Exists(mock), true)

	Remove(mock)
	st.Expect(t, Exists(mock), false)

	Remove(mock)
	st.Expect(t, Exists(mock), false)
}

func TestStoreFlush(t *testing.T) {
	defer after()
	st.Expect(t, len(GetAll()), 0)

	mock1 := New("foo").Mock
	mock2 := New("foo").Mock
	st.Expect(t, len(GetAll()), 2)
	st.Expect(t, Exists(mock1), true)
	st.Expect(t, Exists(mock2), true)

	Flush()
	st.Expect(t, len(GetAll()), 0)
	st.Expect(t, Exists(mock1), false)
	st.Expect(t, Exists(mock2), false)
}
