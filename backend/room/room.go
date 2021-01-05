package room

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Possible room states.
const (
	StateOccupied    State = "OCCUPIED"
	StateUnavailable           = "UNAVAILABLE"
	StateFree                  = "FREE"
)

// Parts of a record.
const (
	EntryID = iota
	EntryPrice
	EntryState
	EntryAttributes
)

// `Number` is the ID/room number of a room.
type Number uint

// `Attribute` is a property that a room can have.
type Attribute string

// `State` indicates the current state of the `Room`.
type State string

// `Room` is a room in a hotel. It has an `ID` (the room number), a price, a
// current state and a set of attributes.
type Room struct {
	mu    *sync.RWMutex
	id    Number
	price uint
	state State
	attrs map[Attribute]struct{}
}

// `NewRoom` returns a pointer to a `Room` with the given `id` (room number).
func NewRoom(id Number) *Room {
	return &Room{
		mu:    &sync.RWMutex{},
		id:    id,
		attrs: make(map[Attribute]struct{}),
	}
}

func NewRoomFromRecord(record []string, validAttributes []Attribute) (*Room, error) {
	const recordLen = 4
	if len(record) != recordLen {
		return nil, fmt.Errorf("invalid record: expected %d entries", recordLen)
	}
	id64, err := strconv.ParseUint(record[EntryID], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid record (id: '%s'): %s", record[0], err.Error())
	}
	price64, err := strconv.ParseUint(record[EntryPrice], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid record (price: '%s'): %s", record[1], err.Error())
	}
	var state State
	stateStr := record[EntryState]
	switch stateStr {
	case "OCCUPIED":
		state = StateOccupied
	case "UNAVAILABLE":
		state = StateUnavailable
	case "FREE":
		state = StateFree
	default:
		return nil, fmt.Errorf("invalid record (state: '%s;): unrecognized state", stateStr)
	}
	roomAttrs := make(map[Attribute]struct{})
	for _, attr := range strings.Split(record[EntryAttributes], ",") {
		roomAttrs[Attribute(attr)] = struct{}{}
	}
	room := &Room{
		mu:    &sync.RWMutex{},
		id:    Number(uint(id64)),
		price: uint(price64),
		state: state,
		attrs: roomAttrs,
	}
	return room, nil
}

// `ID` returns the ID (room number) of the room.
func (r *Room) ID() Number {
	// id is immutable - do not need to lock mutex for read (no writers exist)
	return r.id
}

// `AddAttribute` adds the given `RoomAttribute`, `attr`, to the room.
func (r *Room) AddAttribute(attr Attribute) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.attrs[attr] = struct{}{}
}

// `Satisfies` returns whether the room satisfies the given attributes `attrs`.
func (r *Room) Satisfies(attrs []Attribute) bool {
	r.mu.RLock()
	defer r.mu.RLocker()
	for _, attr := range attrs {
		if _, ok := r.attrs[attr]; !ok {
			return false
		}
	}
	return true
}
