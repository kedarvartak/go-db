package query

import (
	"testing"
)

func TestParseSQL(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantType QueryType
		wantErr  bool
	}{
		{
			name:     "Simple SELECT",
			sql:      "SELECT name,age FROM users",
			wantType: QuerySelect,
			wantErr:  false,
		},
		{
			name:     "SELECT with *",
			sql:      "SELECT * FROM products",
			wantType: QuerySelect,
			wantErr:  false,
		},
		{
			name:     "Simple INSERT",
			sql:      "INSERT INTO users VALUES ('John',30)",
			wantType: QueryInsert,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := ParseSQL(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if query.Type != tt.wantType {
					t.Errorf("ParseSQL() got type = %v, want %v", query.Type, tt.wantType)
				}
				t.Logf("Successfully parsed query: %+v", query)
			}
		})
	}
}
