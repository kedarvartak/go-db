package storage

// Think of tables like excel spreadsheet with different columns
// Column represents a table column definition
type Column struct {
	Name     string
	DataType DataType
	Length   int // For VARCHAR etc.
	NotNull  bool
}

// Table represents a database table structure
type Table struct {
	Name       string
	Columns    []Column
	PrimaryKey int      // Index of primary key column
	PageIDs    []uint64 // Pages containing table data
}

// DataType represents supported data types
type DataType int

const (
	TypeInteger DataType = iota
	TypeVarchar
	TypeBoolean
	TypeTimestamp
)

func NewTable(name string, columns []Column) *Table {
	return &Table{
		Name:    name,
		Columns: columns,
		PageIDs: make([]uint64, 0),
	}
}

func (t *Table) AddPage(pageID uint64) {
	t.PageIDs = append(t.PageIDs, pageID)
}

// Serialize table metadata for storage
func (t *Table) Serialize() []byte {
	// Implementation for serializing table metadata
	// This would include table name, column definitions, etc.
	return nil // Placeholder
}

// Deserialize table metadata from storage
func DeserializeTable(data []byte) (*Table, error) {
	// Implementation for deserializing table metadata
	return nil, nil // Placeholder
}
