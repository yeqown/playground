package onoff

import "strconv"

type (
	evtTyp int8

	state int8
)

func (s state) String() string {
	return strconv.Itoa(int(s))
}

const (
	_ONLINE  state = 1
	_OFFLINE state = 0
)

const (
	_ENTER evtTyp = 1
	_LEAVE evtTyp = -1
)

type Evt struct {
	Typ evtTyp
	Seq uint8
}

type stateMachine struct {
	_state state
	seq    uint8
}

func newStateMachine() *stateMachine {
	return &stateMachine{_state: _OFFLINE, seq: 0}
}

// Handle .
func (m *stateMachine) Handle(evtSequences ...Evt) {
	for _, evt := range evtSequences {
		if evt.Seq < m.seq {
			// drop evt and no change seq
			continue
		}

		// seq is bigger or equal
		switch m._state {
		case _ONLINE:
			m.onlineTransition(evt)
		case _OFFLINE:
			m.offlineTransition(evt)
		}
	}
}

func (m *stateMachine) onlineTransition(evt Evt) (changed bool) {
	switch evt.Typ {
	case _ENTER:
		// E(seq >= k) drop evt but update seq
		m.seq = evt.Seq
		return false
	case _LEAVE:
		// L(seq >= k)
		m.seq = evt.Seq
		m._state = _OFFLINE
		return true
	}

	// invalid path, but do not change state as default
	return false
}

func (m *stateMachine) offlineTransition(evt Evt) (changed bool) {
	switch evt.Typ {
	case _LEAVE:
		// L(seq >= k) drop evt but update seq
		m.seq = evt.Seq
		return false
	case _ENTER:
		// E(seq >= k)

		// if machine.seq is equal to event, should drop this event too.
		if m.seq == evt.Seq {
			// drop, if receive and old enter event
			// special logic, NOTICE HERE !!!!
			return
		}

		// evt.Seq is bigger
		m.seq = evt.Seq
		m._state = _ONLINE
		return true
	}

	// invalid path, but do not change state as default
	return false
}

func (m *stateMachine) state() state {
	return m._state
}
