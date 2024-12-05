package storage

import (
	"errors"
	"os"
)

type Database struct {
	Path          string
	File          *os.File
	PageSize      uint16 // 4kb
	Cache         *Cache
	Tables        map[string]*Table
	RecordManager *RecordManager
}

// Page is the smallest unit of storage in the database. Data is stored in pages (fixed size blocks) rather than one continuous block.
// reading writing fixed size chunks of data is more efficient than reading writing variable size chunks of data.
// easier to manage and more space efficient.
type Page struct {
	ID      uint64 // page id for unique identifications
	Data    []byte // actual data stored in the page
	IsDirty bool   // bool to indicate if page needs to be written to disk
}

func NewDatabase(path string) (*Database, error) {
	db := &Database{
		Path:     path,
		PageSize: 4096,           // Standard page size 4kb
		Cache:    NewCache(1000), // LRU cache with 1000 pages
		Tables:   make(map[string]*Table),
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	db.File = file
	db.RecordManager = NewRecordManager(db)
	return db, nil
}

func (db *Database) allocatePage() (*Page, error) {
	pageID := db.getNextPageID()
	page := &Page{
		ID:      pageID,
		Data:    make([]byte, db.PageSize),
		IsDirty: false,
	}
	return page, nil
}

func (db *Database) readPage(pageID uint64) (*Page, error) {
	offset := int64(pageID) * int64(db.PageSize)
	data := make([]byte, db.PageSize)
	_, err := db.File.ReadAt(data, offset)
	if err != nil {
		return nil, err
	}
	return &Page{ID: pageID, Data: data}, nil
}

func (db *Database) writePage(page *Page) error {
	if !page.IsDirty {
		return nil
	}
	offset := int64(page.ID) * int64(db.PageSize)
	_, err := db.File.WriteAt(page.Data, offset)
	if err != nil {
		return err
	}
	page.IsDirty = false
	return nil
}

func (db *Database) getNextPageID() uint64 {
	// Simplified example: calculate next page ID based on file size
	fileInfo, err := db.File.Stat()
	if err != nil {
		return 0
	}
	return uint64(fileInfo.Size() / int64(db.PageSize))
}

func (db *Database) CreateTable(name string, columns []Column) error {
	if _, exists := db.Tables[name]; exists {
		return errors.New("table already exists")
	}

	table := NewTable(name, columns)
	db.Tables[name] = table
	return nil
}

func (db *Database) GetPage(pageID uint64) (*Page, error) {
	// First, try to get the page from cache
	if page, found := db.Cache.Get(pageID); found {
		return page, nil
	}

	// If not in cache, read from disk
	page, err := db.readPage(pageID)
	if err != nil {
		return nil, err
	}

	// Add to cache for future use
	db.Cache.Put(page)

	return page, nil
}
