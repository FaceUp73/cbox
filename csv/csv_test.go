package csv

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestInitCsv(t *testing.T) {
	tests := []struct {
		name string
		want Csv
	}{
		{"Init Test 01", Csv{Delimiter: Comma}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InitCsv()
			if got.Delimiter != tt.want.Delimiter {
				t.Errorf("InitCsv().Delimiter = %v, want.Delimiter %v", got.Delimiter, tt.want.Delimiter)
			}
			if got.HasSort != tt.want.HasSort {
				t.Errorf("InitCsv().HasSort = %v, want.HasSort %v", got.HasSort, tt.want.HasSort)
			}
			if got.IsAccending != tt.want.IsAccending {
				t.Errorf("InitCsv().IsAccending = %v, want.IsAccending %v", got.IsAccending, tt.want.IsAccending)
			}
		})
	}
}

func TestCsv_UniqueKeyIndex(t *testing.T) {
	type fields struct {
		Delimiter   string
		Unique      string
		HasSort     bool
		Sort        string
		IsAccending bool
		Header      CsvHeader
		Records     []CsvRecord
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"Unique Test 01", fields{HasSort: true, Unique: "Col3", Header: []string{"Col1", "Col2", "Col3"}}, 2},
		{"Unique Test 02", fields{HasSort: false, Unique: "Col2", Header: []string{"Col1", "Col2", "Col3"}}, 1},
		{"Unique Test 11", fields{HasSort: true, Unique: "Col0", Header: []string{"Col1", "Col2", "Col3"}}, -1},
		{"Unique Test 12", fields{HasSort: false, Unique: "Col4", Header: []string{"Col1", "Col2", "Col3"}}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Csv{
				Delimiter:   tt.fields.Delimiter,
				Unique:      tt.fields.Unique,
				HasSort:     tt.fields.HasSort,
				Sort:        tt.fields.Sort,
				IsAccending: tt.fields.IsAccending,
				Header:      tt.fields.Header,
				Records:     tt.fields.Records,
			}
			if got := c.UniqueKeyIndex(); got != tt.want {
				t.Errorf("Csv.UniqueKeyIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCsv_SortKeyIndex(t *testing.T) {
	type fields struct {
		Delimiter   string
		Unique      string
		HasSort     bool
		Sort        string
		IsAccending bool
		Header      CsvHeader
		Records     []CsvRecord
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"Sort Test 01", fields{HasSort: true, Sort: "Col1", Header: []string{"Col1", "Col2", "Col3"}}, 0},
		{"Sort Test 02", fields{HasSort: true, Sort: "Col2", Header: []string{"Col1", "Col2", "Col3"}}, 1},
		{"Sort Test 11", fields{HasSort: true, Sort: "Col4", Header: []string{"Col1", "Col2", "Col3"}}, -1},
		{"Sort Test 12", fields{HasSort: false, Sort: "Col2", Header: []string{"Col1", "Col2", "Col3"}}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Csv{
				Delimiter:   tt.fields.Delimiter,
				Unique:      tt.fields.Unique,
				HasSort:     tt.fields.HasSort,
				Sort:        tt.fields.Sort,
				IsAccending: tt.fields.IsAccending,
				Header:      tt.fields.Header,
				Records:     tt.fields.Records,
			}
			if got := c.SortKeyIndex(); got != tt.want {
				t.Errorf("Csv.SortKeyIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collect(t *testing.T) {
	want := []CsvRecord{
		{"01", "colValue"},
		{"02", "colValue2"},
	}
	type args struct {
		m map[string]collectDist
	}
	tests := []struct {
		name string
		args args
		want []CsvRecord
	}{
		{"Collect Test 01", args{map[string]collectDist{"01": collectDist{CsvRecord{"01", "colValue"}}, "02": collectDist{CsvRecord{"02", "colValue2"}}}}, want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := collect(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("distint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCsv_CsvRead(t *testing.T) {
	type fields struct {
		Delimiter   string
		Unique      string
		HasSort     bool
		Sort        string
		IsAccending bool
		Header      CsvHeader
		Records     []CsvRecord
	}
	type args struct {
		p string
	}

	w := InitCsv()
	w.Header = append(w.Header, strings.Split("col1,col2,col3", Comma)...)
	w.Records = append(w.Records, strings.Split("\"001\",\"col2val1\",\"col3val\"", Comma))
	w.Records = append(w.Records, strings.Split("\"002\",\"col2val1\",\"col3val\"", Comma))
	w.Records = append(w.Records, strings.Split("\"003\",\"col2val1\",\"col3val\"", Comma))
	w.Records = append(w.Records, strings.Split("\"003\",\"col2val2\",\"col3val\"", Comma))
	w.Records = append(w.Records, strings.Split("\"003\",\"col2val3\",\"col3val\"", Comma))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Csv
		wantErr bool
	}{
		{"CsvRead 01", fields{Delimiter: ",", HasSort: false}, args{"./sampledata/testdata1.csv"}, *w, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CsvRead(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("CsvFile.CsvRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("CsvFile.CsvRead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCsv_Distinct(t *testing.T) {
	type args struct {
		c *Csv
	}
	c := InitCsv()
	c.Unique = "col1"
	c.HasSort = true
	c.Sort = "col2"
	c.Header = CsvHeader{"col1", "col2", "col3"}
	c.Records = []CsvRecord{
		CsvRecord{"01", "val01", "001"},
		CsvRecord{"02", "val01", "001"},
		CsvRecord{"02", "val03", "001"},
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Distinct Test 01", args{c}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.c
			if err := c.Distinct(); (err != nil) != tt.wantErr {
				t.Errorf("Csv.Distint() error = %v, wantErr %v", c, tt.wantErr)
			}
			fmt.Printf("Csv.Distint() info = %v\n", c)
		})
	}
}
