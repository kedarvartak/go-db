package transaction

type LockManager struct{}
type TransactionStatus int
type Operation struct{}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{
		transactions: make(map[uint64]*Transaction),
		lockManager:  &LockManager{},
	}
}
