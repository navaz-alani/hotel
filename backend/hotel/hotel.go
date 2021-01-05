package hotel

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/navaz-alani/hotel/room"
)

type Hotel struct {
	mu        *sync.RWMutex
	numRooms  uint
	rooms     map[room.Number]*room.Room
	roomAttrs []room.Attribute
}

// `NewHotelFromData` creates a new `Hotel` from the attributes data contained
// in `attrData` and the room data contained in `roomData`. Any fatal errors
// encountered are returned by default, however with `strict` set to true, any
// errors encountered while parsing will be returned.
//
// Check the 'record_formats' directory for the formats of these two data files.
func NewHotelFromData(attrData, roomData string, strict bool) (*Hotel, error) {
	hotel := &Hotel{
		mu:    &sync.RWMutex{},
		rooms: make(map[room.Number]*room.Room),
	}
	if err := hotel.loadAttributes(attrData); err != nil {
		return nil, err
	} else if err = hotel.loadRooms(roomData, strict); err != nil {
		return nil, err
	}
	hotel.numRooms = uint(len(hotel.rooms))
	return hotel, nil
}

// `loadRooms` loads `Room`s from the data in the file with name `roomData`. Any
// errors occurred while opening the `roomData` file or reading from it will be
// returned. Errors encountered while parsing scanned data into a `Room` will be
// ignored, unless the `strict` flag is true.
//
// The parsed rooms are loaded into the `Hotel`, `h`, directly. If an error is
// occurred, the state of `h` is unchanged.
//
// Full format specs in record_formats/room_list_format
func (h *Hotel) loadRooms(roomData string, strict bool) error {
	f, err := os.Open(roomData)
	if err != nil {
		return fmt.Errorf("rooms load err: %s", err.Error())
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	initialRecord := true
	rooms := make(map[room.Number]*room.Room)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("load err [fatal]: %s", err.Error())
		}
		if initialRecord { // header
			initialRecord = false
			continue
		}
		room, err := room.NewRoomFromRecord(record, h.roomAttrs)
		if err != nil && strict {
			return fmt.Errorf("load err: room parse err: %s", err.Error())
		}
		// this means that if there are multiple rooms in the room data file which
		// have the same room number, the last such record is the one that will
		// appear - room numbers must be unique.
		rooms[room.ID()] = room
	}

	// modifying hotel contents
	h.mu.Lock()
	defer h.mu.Unlock()
	// enter the parsed data into the hotel
	for k, v := range rooms {
		h.rooms[k] = v
	}

	return nil
}

// `loadAttributes` loads the attribues contained in the file with the name
// `attrData` and returns any errors encountered. It takes only the first word
// (consecutive non-whitespace string) on each line as the attribute - this
// means that there can be comments on each line after the attribute in addition
// to entire line comments i.e. lines which begin with "# ").
//
// The attributes are loaded into the `Hotel`, `h`. If an error occurs, the
// state of `h` is unchanged.
//
// Full format specs in record_formats/attr_list_format
func (h *Hotel) loadAttributes(attrData string) error {
	attrFile, err := os.Open(attrData)
	if err != nil {
		return fmt.Errorf("attributes load err: %s", err.Error())
	}

	// modifying hotel contents
	h.mu.Lock()
	defer h.mu.Unlock()

	scanner := bufio.NewScanner(attrFile)
	for scanner.Scan() {
		attr := strings.Split(scanner.Text(), " \t")[0]
		if attr == "" || attr == "#" {
			continue
		}
		h.roomAttrs = append(
			h.roomAttrs,
			room.Attribute(attr),
		)
	}
	return nil
}
