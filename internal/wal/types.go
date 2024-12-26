package wal

// TIRTHRAJ IF YOURE STALKING THIS FUCK YOU GET A LIFE BITCH

type LogType int

const (
	LogTypeBeginTx LogType = iota // this will mark the start of a transaction
	LogTypeCommitTx              // this will mark the end of a transaction
	LogTypeAbortTx               // this will mark the abort of a transaction
	LogTypeInsert                // this will mark the insertion of a record
	LogTypeUpdate                // this will mark the update of a record
	LogTypeDelete                // this will mark the deletion of a record
)

// LSN (Log Sequence Number) is used to uniquely identify log records
type LSN uint64
// LSN tracks what the last operation was that modifies a page
//tracks which page has been modified
// PageLSN tracks the latest LSN that modified each page
type PageLSN struct {
	PageID uint64
	LSN    LSN
}

// LogRecord does the actual data being logged
type LogRecord struct {
	Before []byte // Data change chya adhi - used for undo
	After  []byte // Data change nantar - used for redo
}


//Example of logging an update to a record

// logEntry := &LogEntry{
//     Type: LogTypeUpdate,
//     Record: LogRecord{
//         Before: []byte("old value"),
//         After:  []byte("new value"),
//     }
// }


// pageLSN := PageLSN{
//     PageID: 1,
//     LSN: 1234,  
// }
// this is a simple example of a pageLSN
