package tablewriter

import (
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	type User struct {
		Name string
		ID   int
		Sex  string
	}
	users := []User{
		{
			"AA", 123456, "M",
		},
		{
			"BB", 1423, "F",
		},
	}
	wr := NewWriter(os.Stdout, []string{"Name", "ID", "Sex"})

	for _, u := range users {
		wr.Write(u)
	}
	wr.Flush()
}
