package main

import (
	"fmt"

	"godb/internal/query"
	"godb/internal/storage"
	"godb/internal/transaction"
	"godb/internal/wal"
)

func main() {
	// Test Storage Engine
	db, err := storage.NewDatabase("testdb.db")
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer db.File.Close() // Close the database file when done
	fmt.Println("Database initialized successfully.")

	// Test B-tree Initialization
	btree := storage.NewBTree(3)
	if btree == nil {
		fmt.Println("Error initializing B-tree.")
		return
	}
	fmt.Println("B-tree initialized successfully.")

	// Test Write-Ahead Logging
	walLog, err := wal.NewWAL("testlog.wal")
	if err != nil {
		fmt.Println("Error creating WAL file:", err)
		return
	}
	defer walLog.File.Close() // Close the WAL file when done

	// Test WAL writing
	err = walLog.Write([]byte("test data"))
	if err != nil {
		fmt.Println("Error writing to WAL:", err)
		return
	}
	fmt.Println("WAL initialized and test write successful.")

	// Test Transaction Manager
	tm := transaction.NewTransactionManager()
	if tm == nil {
		fmt.Println("Error initializing Transaction Manager.")
		return
	}
	fmt.Println("Transaction Manager initialized successfully.")

	// Test Query Processor
	qp := query.QueryProcessor{}
	// Test query execution
	_, err = qp.Execute("SELECT * FROM test")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	fmt.Println("Query Processor initialized and test query executed successfully.")
}
