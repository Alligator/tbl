# tbl

tbl is a simple way to print text tables.

## Installing

```
go get -u github.com/alligator/tbl
```

## Example

```go
package main

import (
	"fmt"

	"github.com/alligator/tbl"
)

type row struct {
	id   int
	name string
}

func main() {
	rows := []row{
		{1, "gate"},
		{2, "boop"},
	}

	t := tbl.NewTable()

	for _, row := range rows {
		t.NewRow()

		t.NewCol("Id")
		t.Printf("%d", row.id)

		t.NewCol("Name")
		t.Print(row.name)
	}

	fmt.Print(t.String())
}
```
