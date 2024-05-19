package csv

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const (
	Comma string = ","
)

type CsvHeader []string
type CsvRecord []string

// CSV
type Csv struct {
	// CSVのデリミタ
	Delimiter string
	// ユニークキー
	Unique string

	// ソート有無フラグ
	HasSort bool
	// ソートキー
	Sort string
	// ソート方向
	IsAccending bool

	// CSVヘッダー
	Header CsvHeader
	// CSVレコード
	Records []CsvRecord
}

type CsvHandle interface {
	// Read a csv file to create a CsvFile structure.
	CsvRead(p string) (*Csv, error)

	// Get index of unique key.
	UniqueKeyIndex() int

	// Get index of sort key.
	SortKeyIndex() int

	// Remove duplicates.
	Distint() error
}

// Initial Csv
func InitCsv() *Csv {
	c := new(Csv)
	c.Delimiter = Comma
	c.Header = make(CsvHeader, 0)
	c.Records = make([]CsvRecord, 0)

	return c
}

func (c Csv) UniqueKeyIndex() int {
	var result int = -1
	for idx, v := range c.Header {
		if v == c.Unique {
			result = idx
			break
		}
	}
	return result
}

func (c Csv) SortKeyIndex() int {
	var result int = -1
	if !c.HasSort {
		return -1
	}
	for idx, v := range c.Header {
		if v == c.Sort {
			result = idx
			break
		}
	}
	return result
}

func CsvRead(p string) (*Csv, error) {
	cf := InitCsv()

	f, err := os.Open(p)
	if err != nil {
		log.Fatal(err.Error())
		return cf, err
	}
	defer f.Close()

	isHeader := true
	s := bufio.NewScanner(f)
	for s.Scan() {
		if isHeader {
			isHeader = false
			cf.Header = append(cf.Header, strings.Split(s.Text(), Comma)...)
		} else {
			cf.Records = append(cf.Records, strings.Split(s.Text(), Comma))
		}
	}
	return cf, nil
}

type collectDist struct {
	value CsvRecord
}

func collect(m map[string]collectDist) []CsvRecord {
	var result []CsvRecord = []CsvRecord{}
	for _, v := range m {
		result = append(result, v.value)
	}
	return result
}

func (c *Csv) Distinct() error {

	m := map[string]collectDist{}
	pk := c.UniqueKeyIndex()
	sort := c.SortKeyIndex()

	records := c.Records
	for _, newRecord := range records {
		name := newRecord[pk]
		if _, ok := m[name]; !ok {
			m[name] = collectDist{newRecord}
		} else {
			// rather then origin
			if strings.Compare(newRecord[sort], m[name].value[sort]) == 1 {
				m[name] = collectDist{newRecord}
			}
		}
	}
	c.Records = collect(m)
	return nil
}
