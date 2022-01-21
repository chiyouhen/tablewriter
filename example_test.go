package tablewriter_test

import (
	"github.com/chiyouhen/tablewriter"
	"os"
)

func ExampleWriter() {
	type User struct {
		Name  string
		Age   int
		Speed float64
	}

	users := []User{
		{
			"John", 15, 14.3,
		},
		{
			"Mike", 14, 11.5,
		},
	}
	wr := tablewriter.NewWriter(os.Stdout, []string{"Name", "Age", "Speed"})
	for _, u := range users {
		wr.Write(u)
	}
	wr.Flush()
	// output
	// Name Age Speed
	// John  15 14.300000
	// Mike  14 11.500000
}
