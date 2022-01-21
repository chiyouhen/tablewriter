package tablewriter

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/text/width"
)

type Cell struct {
	Value        interface{}
	StringLength int
}

func NewCell(v reflect.Value) (Cell, error) {
	for i := 0; i < 100; i++ {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				break
			}
			v = v.Elem()
		} else {
			break
		}
	}
	c := Cell{
		Value:        "None",
		StringLength: 4,
	}
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	var s string
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		i := v.Float()
		if float64(int64(i)) == i {
			c.Value = int64(i)
			s = fmt.Sprintf("%d", c.Value.(int64))
		} else {
			c.Value = i
			s = fmt.Sprintf("%0.6f", c.Value.(float64))
		}
	case reflect.String:
		c.Value = v.String()
		s = c.Value.(string)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		c.Value = v.Int()
		s = fmt.Sprintf("%d", c.Value.(int64))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		c.Value = v.Uint()
		s = fmt.Sprintf("%d", c.Value.(uint64))
	}
	buf := []byte(s)
	ra := []rune(s)
	n := len(buf)
	if n != len(ra) {
		n = 0
		for _, r := range ra {
			p := width.LookupRune(r)
			k := p.Kind()
			if k == width.EastAsianWide {
				n += 2
			} else {
				n += 1
			}
		}
	}
	c.StringLength = n
	return c, nil
}

type Column struct {
	Name       string
	Length     int
	IntAsFloat bool
}

type Writer struct {
	columns []Column
	rows    [][]Cell
	ww      io.Writer
}

func NewWriter(ww io.Writer, columns []string) *Writer {
	wr := &Writer{
		ww:      ww,
		columns: make([]Column, len(columns)),
		rows:    make([][]Cell, 0),
	}
	for i, s := range columns {
		c := Column{
			Name:   s,
			Length: len(s),
		}
		wr.columns[i] = c
	}
	return wr
}

func (wr *Writer) Write(d interface{}) error {
	row := make([]Cell, len(wr.columns))
	v := reflect.ValueOf(d)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Array && v.Kind() != reflect.Struct && v.Kind() != reflect.Map {
		return fmt.Errorf("unsupported data type %s", v.Kind())
	}
	for i, c := range wr.columns {
		key := reflect.ValueOf(c.Name)
		var value reflect.Value
		switch v.Kind() {
		case reflect.Array:
			value = v.Index(i)
		case reflect.Map:
			value = v.MapIndex(key)
		case reflect.Struct:
			value = v.FieldByName(c.Name)
		}
		cell, err := NewCell(value)
		if err != nil {
			return err
		}
		row[i] = cell
		if cell.StringLength > c.Length {
			c.Length = cell.StringLength
			wr.columns[i] = c
		}
		if _, ok := cell.Value.(float64); ok {
			c.IntAsFloat = true
		}
	}
	wr.rows = append(wr.rows, row)

	return nil
}

func (wr *Writer) flushRow(row []Cell) error {
	items := make([]string, len(row))
	var s string
	for i, c := range wr.columns {
		cell := row[i]

		switch v := cell.Value.(type) {
		case int64, uint64:
			if c.IntAsFloat {
				s = fmt.Sprintf("% "+strconv.Itoa(c.Length)+".6f", v)
			} else {
				s = fmt.Sprintf("% "+strconv.Itoa(c.Length)+"d", v)
			}
		case float64:
			s = fmt.Sprintf("% "+strconv.Itoa(c.Length)+".6f", v)
		case string:
			printLength := c.Length
			ra := []rune(v)
			for _, r := range ra {
				p := width.LookupRune(r)
				k := p.Kind()
				if k == width.EastAsianWide {
					printLength -= 1
				}
			}
			s = fmt.Sprintf("%-"+strconv.Itoa(printLength)+"s", v)
		}
		items[i] = s

	}
	fmt.Println(strings.Join(items, " "))
	return nil
}

func (wr *Writer) Flush() error {
	headerRows := make([]Cell, len(wr.columns))
	for i, c := range wr.columns {
		headerRows[i] = Cell{
			Value: c.Name,
		}
	}
	wr.flushRow(headerRows)
	for _, row := range wr.rows {
		wr.flushRow(row)
	}
	return nil
}
