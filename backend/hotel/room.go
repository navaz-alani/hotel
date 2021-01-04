package hotel

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Possible room states.
const (
	StateOccupied    RoomState = "OCCUPIED"
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

// `RoomNumber` is the ID/room number of a room.
type RoomNumber uint

// `RoomAttribute` is a property that a room can have.
type RoomAttribute string

// `RoomState` indicates the current state of the `Room`.
type RoomState string

// `Room` is a room in a hotel. It has an `ID` (the room number), a price, a
// current state and a set of attributes.
type Room struct {
	mu    *sync.RWMutex
	id    RoomNumber
	price uint
	state RoomState
	attrs map[RoomAttribute]struct{}
}

// `NewRoom` returns a pointer to a `Room` with the given `id` (room number).
func NewRoom(id RoomNumber) *Room {
	return &Room{
		mu:    &sync.RWMutex{},
		id:    id,
		attrs: make(map[RoomAttribute]struct{}),
	}
}

func NewRoomFromRecord(record []string, validAttributes []RoomAttribute) (*Room, error) {
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
	var state RoomState
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
	roomAttrs := make(map[RoomAttribute]struct{})
	for _, attr := range strings.Split(record[EntryAttributes], ",") {
		roomAttrs[RoomAttribute(attr)] = struct{}{}
	}
	room := &Room{
		mu:    &sync.RWMutex{},
		id:    RoomNumber(uint(id64)),
		price: uint(price64),
		state: state,
		attrs: roomAttrs,
	}
	return room, nil
}

// `ID` returns the ID (room number) of the room.
func (r *Room) ID() RoomNumber {
	// id is immutable - do not need to lock mutex for read (no writers exist)
	return r.id
}

// `AddAttribute` adds the given `RoomAttribute`, `attr`, to the room.
func (r *Room) AddAttribute(attr RoomAttribute) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.attrs[attr] = struct{}{}
}

// `Satisfies` returns whether the room satisfies the given attributes `attrs`.
func (r *Room) Satisfies(attrs []RoomAttribute) bool {
	r.mu.RLock()
	defer r.mu.RLocker()
	for _, attr := range attrs {
		if _, ok := r.attrs[attr]; !ok {
			return false
		}
	}
	return true
}
