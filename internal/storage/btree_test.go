package storage

import (
	"testing"
)

func TestBTree(t *testing.T) {
	t.Run("Node Splitting", func(t *testing.T) {
		btree := NewBTree(3)
		if btree == nil {
			t.Fatal("Failed to create B-tree")
		}

		// Insert keys in ascending order to force splits
		for i := 1; i <= 10; i++ {
			t.Logf("Inserting key %d", i)
			btree.Insert(i, RecordID{uint64(i), 1})
			t.Logf("Tree after insertion:\n%s", btree.String())

			// Verify all previously inserted keys are still findable
			for j := 1; j <= i; j++ {
				_, found := btree.Search(j)
				if !found {
					t.Errorf("Key %d not found after inserting key %d\nTree state:\n%s",
						j, i, btree.String())
				} else {
					t.Logf("Successfully found key %d", j)
				}
			}
		}
	})

	t.Run("Basic Insert and Search", func(t *testing.T) {
		btree := NewBTree(3) // Degree 3
		if btree == nil {
			t.Fatal("Failed to create B-tree")
		}

		// Insert some test records
		records := []struct {
			key int
			rid RecordID
		}{
			{5, RecordID{1, 1}},
			{3, RecordID{1, 2}},
			{7, RecordID{1, 3}},
			{1, RecordID{1, 4}},
			{9, RecordID{1, 5}},
		}

		// Insert records
		for _, r := range records {
			btree.Insert(r.key, r.rid)
			// Verify immediate insertion
			rid, found := btree.Search(r.key)
			if !found {
				t.Errorf("Key %d not found immediately after insertion", r.key)
			}
			if rid != nil && *rid != r.rid {
				t.Errorf("For key %d, got RecordID %v, want %v", r.key, rid, r.rid)
			}
		}

		// Search for non-existent key
		_, found := btree.Search(100)
		if found {
			t.Error("Found non-existent key 100")
		}
	})
}
