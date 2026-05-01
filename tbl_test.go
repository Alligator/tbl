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
	tbl := NewTable()
	expected := `| id | name                 |
| -- | -------------------- |
| 1  | gate (owner)         |
| 2  | tooty (contributor)  |
| 3  | sponge (contributor) |
`
	test(t, tbl, expected)
}

func TestTblMinimal(t *testing.T) {
	tbl := NewTable()
	tbl.Style = StyleMinimal
	expected := `ID  NAME
1   gate (owner)
2   tooty (contributor)
3   sponge (contributor)
`
	test(t, tbl, expected)
}

func test(t *testing.T, table *Table, expected string) {
	rows := []row{
		{"1", "gate", "owner"},
		{"2", "tooty", "contributor"},
		{"3", "sponge", "contributor"},
	}

	for _, r := range rows {
		table.NewRow()
		table.NewCol("id")
		table.Print(r.id)

		table.NewCol("name")
		table.Printf("%s (%s)", r.name, r.role)
	}

	result := table.String()
	if result != expected {
		t.Errorf("result (len %d)\n%s\ndoes not match expected (len %d)\n%s\n", len(result), result, len(expected), expected)
	}
}
