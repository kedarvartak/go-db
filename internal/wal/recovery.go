package wal

import (
    "encoding/binary"
    "io"
)

// TIRTHRAJ IF YOURE STALKING THIS FUCK YOU GET A LIFE BITCH

type Recovery struct {
    wal *WAL
    activeTxs map[uint64]struct{} // set of active transactions during crash
    redoLog []LogEntry            // Log entries that need to be replayed
}

func (w *WAL) StartRecovery() *Recovery {
    return &Recovery{
        wal: w,
        activeTxs: make(map[uint64]struct{}),
    }
}

//The recovery happens in three phases (known as ARIES recovery protocol)
//1. Analysis phase: scan log to identify active transactions
//2. Redo phase: replay all changes
//3. Undo phase: rollback incomplete transactions


func (r *Recovery) Recover() error {
    // Reset file position
    if _, err := r.wal.file.Seek(0, 0); err != nil {
        return err
    }

    // Analysis phase: scan log to identify active transactions
	// analysis phase madhe entire log scan hoto which identifies which transactions were active during crash
	// this is done by scanning the log and identifying the log entries that are of type LogTypeBeginTx and LogTypeCommitTx
    if err := r.analysisPhase(); err != nil {
        return err
    }

    // Redo phase: replay all changes
	// Should implement:
    // 1. Replay all changes in the log
    // 2. Bring database to state it was in before crash
    // 3. Apply changes even for transactions that didn't commit
    if err := r.redoPhase(); err != nil {
        return err
    }

    // Undo phase: rollback incomplete transactions
	 // Should implement:
    // 1. Rollback all active transactions found in analysis phase
    // 2. Use 'Before' images to restore original data
    // 3. Write compensation log records
    return r.undoPhase()
}

func (r *Recovery) analysisPhase() error {
    for {
        entry, err := r.readLogEntry()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }

        switch entry.Type {
        case LogTypeBeginTx:
            r.activeTxs[entry.TxID] = struct{}{}
        case LogTypeCommitTx, LogTypeAbortTx:
            delete(r.activeTxs, entry.TxID)
        }
    }
    return nil
}

func (r *Recovery) readLogEntry() (*LogEntry, error) {
    // Read and deserialize log entry
    // Implementation details here
    return nil, nil
} 