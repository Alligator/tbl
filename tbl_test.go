package tbl

import (
	"testing"
)

type row struct {
	id   string
	name string
	role string
}

func TestTbl(t *testing.T) {
	rows := []row{
		{"1", "gate", "owner"},
		{"2", "tooty", "contributor"},
		{"3", "sponge", "contributor"},
	}

	table := NewTable()

	for _, r := range rows {
		table.NewRow()
		table.NewCol("id")
		table.Print(r.id)

		table.NewCol("name")
		table.Printf("%s (%s)", r.name, r.role)
	}

	result := table.String()
	expected := `| id | name                 |
| -- | -------------------- |
| 1  | gate (owner)         |
| 2  | tooty (contributor)  |
| 3  | sponge (contributor) |
`

	if result != expected {
		t.Errorf("result (len %d)\n%s\ndoes not match expected (len %d)\n%s\n", len(result), result, len(expected), expected)
	}
}
