package storage

import (
	"encoding/binary"
	"errors"
)

// manages how data is stored in a page. Contains header (metadata) and slots (data)
const (
	// Page layout constants
	PageHeaderSize = 16 // Size of page header in bytes
	SlotEntrySize  = 8  // Size of each slot entry
	MinRecordSize  = 4  // Minimum size of a record

	// Header offsets
	OffsetSlotCount  = 0  // Number of slots
	OffsetFreeSpace  = 4  // Free space pointer
	OffsetLastSlotID = 8  // Last used slot ID
	OffsetFlags      = 12 // Flags/Reserved
)

// SlotEntry represents an entry in the slot directory
type SlotEntry struct {
	Offset uint32 // Offset from start of page
	Length uint16 // Length of record
	Flags  uint16 // Record flags (e.g., deleted, overflow)
}

// PageHeader represents the header section of a page
type PageHeader struct {
	SlotCount  uint32 // Number of slots in use
	FreeSpace  uint32 // Pointer to start of free space
	LastSlotID uint32 // ID of last used slot
	Flags      uint32 // Page flags
}

// PageLayout manages the internal layout of a page
type PageLayout struct {
	header PageHeader
	slots  []SlotEntry
	data   []byte
}

// NewPageLayout initializes a new page layout
func NewPageLayout(pageSize uint16) *PageLayout {
	return &PageLayout{
		header: PageHeader{
			SlotCount:  0,
			FreeSpace:  uint32(pageSize - PageHeaderSize),
			LastSlotID: 0,
			Flags:      0,
		},
		slots: make([]SlotEntry, 0),
		data:  make([]byte, pageSize),
	}
}

// getFreeSpace calculates available free space
func (pl *PageLayout) getFreeSpace() uint32 {
	usedBySlots := uint32(len(pl.slots) * SlotEntrySize)
	return pl.header.FreeSpace - usedBySlots
}

// findFreeSlot finds space for a new record
func (pl *PageLayout) findFreeSlot(recordSize uint16) (uint16, error) {
	if pl.getFreeSpace() < uint32(recordSize+SlotEntrySize) {
		return 0, errors.New("insufficient space in page")
	}

	// Find a suitable location in the page
	offset := pl.header.FreeSpace - uint32(recordSize)

	// Create new slot entry
	slotID := uint16(pl.header.LastSlotID + 1)
	pl.slots = append(pl.slots, SlotEntry{
		Offset: offset,
		Length: recordSize,
		Flags:  0,
	})

	pl.header.LastSlotID++
	pl.header.SlotCount++
	pl.header.FreeSpace = offset

	return slotID, nil
}

// Serialize converts the page layout to bytes
func (pl *PageLayout) Serialize() []byte {
	// Write header
	binary.LittleEndian.PutUint32(pl.data[OffsetSlotCount:], pl.header.SlotCount)
	binary.LittleEndian.PutUint32(pl.data[OffsetFreeSpace:], pl.header.FreeSpace)
	binary.LittleEndian.PutUint32(pl.data[OffsetLastSlotID:], pl.header.LastSlotID)
	binary.LittleEndian.PutUint32(pl.data[OffsetFlags:], pl.header.Flags)

	// Write slot directory
	slotOffset := PageHeaderSize
	for _, slot := range pl.slots {
		binary.LittleEndian.PutUint32(pl.data[slotOffset:], slot.Offset)
		binary.LittleEndian.PutUint16(pl.data[slotOffset+4:], slot.Length)
		binary.LittleEndian.PutUint16(pl.data[slotOffset+6:], slot.Flags)
		slotOffset += SlotEntrySize
	}

	return pl.data
}

// Deserialize reads the page layout from bytes
func DeserializePageLayout(data []byte) *PageLayout {
	pl := &PageLayout{
		data: data,
	}

	// Read header
	pl.header.SlotCount = binary.LittleEndian.Uint32(data[OffsetSlotCount:])
	pl.header.FreeSpace = binary.LittleEndian.Uint32(data[OffsetFreeSpace:])
	pl.header.LastSlotID = binary.LittleEndian.Uint32(data[OffsetLastSlotID:])
	pl.header.Flags = binary.LittleEndian.Uint32(data[OffsetFlags:])

	// Read slot directory
	pl.slots = make([]SlotEntry, pl.header.SlotCount)
	slotOffset := PageHeaderSize
	for i := uint32(0); i < pl.header.SlotCount; i++ {
		pl.slots[i].Offset = binary.LittleEndian.Uint32(data[slotOffset:])
		pl.slots[i].Length = binary.LittleEndian.Uint16(data[slotOffset+4:])
		pl.slots[i].Flags = binary.LittleEndian.Uint16(data[slotOffset+6:])
		slotOffset += SlotEntrySize
	}

	return pl
}
