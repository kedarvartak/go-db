package storage

import "errors"

// manages how records are inserted, retrieved, and deleted from the database

// Think of records like rows in the spreadsheet
// Record represents a single row in a table
type Record struct {
	Values []interface{}
}

// RecordID identifies a record's location
type RecordID struct {
	PageID  uint64
	SlotNum uint16
}

// RecordManager handles record operations
type RecordManager struct {
	db *Database
}

func NewRecordManager(db *Database) *RecordManager {
	return &RecordManager{
		db: db,
	}
}

func (rm *RecordManager) InsertRecord(table *Table, record *Record) (*RecordID, error) {
	// Serialize the record
	recordData, err := SerializeRecord(record)
	if err != nil {
		return nil, err
	}

	// Find a page with enough space
	pageID, err := rm.findPageWithSpace(table, len(recordData))
	if err != nil {
		return nil, err
	}

	// Get the page from cache or disk
	page, err := rm.db.GetPage(pageID)
	if err != nil {
		return nil, err
	}

	// Insert record into page
	slotNum, err := rm.insertIntoPage(page, recordData)
	if err != nil {
		return nil, err
	}

	return &RecordID{
		PageID:  pageID,
		SlotNum: slotNum,
	}, nil
}

func (rm *RecordManager) GetRecord(table *Table, rid *RecordID) (*Record, error) {
	// Get the page containing the record
	page, err := rm.db.GetPage(rid.PageID)
	if err != nil {
		return nil, err
	}

	// Extract record from page
	record, err := rm.extractRecord(page, rid.SlotNum, table)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// Helper methods
func (rm *RecordManager) findPageWithSpace(table *Table, recordSize int) (uint64, error) {
	// Check existing pages
	for _, pageID := range table.PageIDs {
		page, err := rm.db.GetPage(pageID)
		if err != nil {
			continue
		}

		layout := DeserializePageLayout(page.Data)
		if layout.getFreeSpace() >= uint32(recordSize+SlotEntrySize) {
			return pageID, nil
		}
	}

	// No existing page has enough space, create new page
	newPage := &Page{
		ID:      uint64(len(table.PageIDs)),
		Data:    make([]byte, rm.db.PageSize),
		IsDirty: true,
	}

	// Initialize new page layout
	layout := NewPageLayout(rm.db.PageSize)
	newPage.Data = layout.Serialize()

	// Add page to table
	table.AddPage(newPage.ID)

	// Add page to cache
	rm.db.Cache.Put(newPage)

	return newPage.ID, nil
}

func (rm *RecordManager) insertIntoPage(page *Page, recordData []byte) (uint16, error) {
	// Get or create page layout
	layout := DeserializePageLayout(page.Data)

	// Find a slot for the record
	slotNum, err := layout.findFreeSlot(uint16(len(recordData)))
	if err != nil {
		return 0, err
	}

	// Write record data to page
	slot := layout.slots[slotNum-1]
	copy(page.Data[slot.Offset:], recordData)

	// Update page layout
	page.Data = layout.Serialize()
	page.IsDirty = true

	return slotNum, nil
}

func (rm *RecordManager) extractRecord(page *Page, slotNum uint16, table *Table) (*Record, error) {
	// Get page layout
	layout := DeserializePageLayout(page.Data)

	// Validate slot number
	if int(slotNum) > len(layout.slots) {
		return nil, errors.New("invalid slot number")
	}

	// Get slot entry
	slot := layout.slots[slotNum-1]

	// Check if record is deleted
	if slot.Flags&1 != 0 {
		return nil, errors.New("record deleted")
	}

	// Extract record data
	recordData := page.Data[slot.Offset : slot.Offset+uint32(slot.Length)]

	// Deserialize record
	return DeserializeRecord(recordData)
}
