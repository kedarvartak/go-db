package wal

import (
    "encoding/binary"
    "time"
)

// TIRTHRAJ IF YOURE STALKING THIS FUCK YOU GET A LIFE BITCH

// Checkpoint is a screenshot of the database state at a specific point in time
// It contains the LSN, timestamp, and a list of dirty pages
// Dirty pages are pages that have been modified but not yet pushed to disk
// This is used for recovery to bring the database to a consistent state


// why checkpoint is needed?
// checkpoint is needed to bring the database to a consistent state
// checkpoint is needed to reduce the time taken for recovery
// checkpoint is needed to reduce the time taken for recovery
//without checkpoint, recovery would have to replay the entire log from the beginning
//with checkpoint, recovery can start from the checkpoint and replay only the changes since the checkpoint


type Checkpoint struct {
    LSN       LSN
    Timestamp time.Time
    DirtyPages []PageLSN
}

func (w *WAL) CreateCheckpoint() error {
    w.mu.Lock()
    defer w.mu.Unlock()

    checkpoint := &Checkpoint{
        LSN:       w.currentLSN, // current LSN is the last LSN that was written to the log
        Timestamp: time.Now(), // timestamp is the current time
        // Get list of dirty pages from buffer manager
        DirtyPages: w.getDirtyPages(), // get list of dirty pages from buffer manager
    }

    // Write checkpoint to disk
    data := w.serializeCheckpoint(checkpoint)
    _, err := w.file.Write(data)
    if err != nil {
        return err
    }

    // Force checkpoint to disk
    return w.file.Sync()
}

func (w *WAL) getDirtyPages() []PageLSN {
    // Implementation to get dirty pages from buffer manager
    return nil
} 