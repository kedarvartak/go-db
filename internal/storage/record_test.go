package storage

import (
	"fmt"
	"os"
	"testing"
)

func TestRecordOperations(t *testing.T) {
	// Setup
	dbPath := "test_db.db"
	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer func() {
		db.File.Close()
		os.Remove(dbPath)
	}()

	// Create a test table
	table := NewTable("test_table", []Column{
		{Name: "id", DataType: TypeInteger},
		{Name: "name", DataType: TypeVarchar, Length: 50},
		{Name: "age", DataType: TypeInteger},
	})

	t.Run("Insert and Retrieve Record", func(t *testing.T) {
		// Create a test record
		record := &Record{
			Values: []interface{}{
				int(1),     // Explicitly use int type
				"John Doe", // string
				int(30),    // Explicitly use int type
			},
		}

		// Insert record
		rid, err := db.RecordManager.InsertRecord(table, record)
		if err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}
		if rid == nil {
			t.Fatal("Expected record ID, got nil")
		}

		// Retrieve record
		retrieved, err := db.RecordManager.GetRecord(table, rid)
		if err != nil {
			t.Fatalf("Failed to retrieve record: %v", err)
		}

		// Verify retrieved record
		if len(retrieved.Values) != len(record.Values) {
			t.Errorf("Expected %d values, got %d", len(record.Values), len(retrieved.Values))
		}

		// Check each value
		for i, val := range record.Values {
			if retrieved.Values[i] != val {
				t.Errorf("Value mismatch at index %d: expected %v, got %v",
					i, val, retrieved.Values[i])
			}
		}
	})

	t.Run("Multiple Records", func(t *testing.T) {
		records := []*Record{
			{Values: []interface{}{2, "Jane Doe", 25}},
			{Values: []interface{}{3, "Bob Smith", 35}},
			{Values: []interface{}{4, "Alice Johnson", 28}},
		}

		var rids []*RecordID
		// Insert multiple records
		for _, record := range records {
			rid, err := db.RecordManager.InsertRecord(table, record)
			if err != nil {
				t.Fatalf("Failed to insert record: %v", err)
			}
			rids = append(rids, rid)
		}

		// Retrieve and verify each record
		for i, rid := range rids {
			retrieved, err := db.RecordManager.GetRecord(table, rid)
			if err != nil {
				t.Fatalf("Failed to retrieve record %d: %v", i, err)
			}

			// Verify record values
			for j, val := range records[i].Values {
				if retrieved.Values[j] != val {
					t.Errorf("Record %d: value mismatch at index %d: expected %v, got %v",
						i, j, val, retrieved.Values[j])
				}
			}
		}
	})

	t.Run("Page Management", func(t *testing.T) {
		// Create records until we need a new page
		var lastPageID uint64
		for i := 0; i < 100; i++ { // Arbitrary number that should force page creation
			record := &Record{
				Values: []interface{}{
					i + 1000,                       // id
					fmt.Sprintf("Test User %d", i), // Fixed: proper string formatting
					20 + i%30,                      // age
				},
			}

			rid, err := db.RecordManager.InsertRecord(table, record)
			if err != nil {
				t.Fatalf("Failed to insert record %d: %v", i, err)
			}

			if rid.PageID != lastPageID {
				t.Logf("New page created: %d", rid.PageID)
				lastPageID = rid.PageID
			}
		}

		if len(table.PageIDs) <= 1 {
			t.Error("Expected multiple pages to be created")
		}
	})

	t.Run("Error Cases", func(t *testing.T) {
		// Test invalid record retrieval
		_, err := db.RecordManager.GetRecord(table, &RecordID{
			PageID:  999999, // Invalid page ID
			SlotNum: 1,
		})
		if err == nil {
			t.Error("Expected error when retrieving invalid record")
		}

		// Test inserting invalid record
		_, err = db.RecordManager.InsertRecord(table, &Record{
			Values: []interface{}{}, // Empty record
		})
		if err == nil {
			t.Error("Expected error when inserting empty record")
		}
	})
}
