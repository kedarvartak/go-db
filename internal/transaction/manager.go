package transaction

import (
	"time"
)

type TransactionManager struct {
	transactions map[uint64]*Transaction
	lockManager  *LockManager
}

type Transaction struct {
	id        uint64
	status    TransactionStatus
	timestamp time.Time
	operations []Operation
} 