package tabular

import "database/sql"

type Scanner struct {
	Fields []interface{}
}

func NewScanner(fields ...interface{}) (scanner Scanner) {
	scanner.Fields = fields
	return
}

func (scanner Scanner) Scan(row *sql.Row, joined ...Scanner) (err error) {
	fields := scanner.Fields
	for i := 0; i < len(joined); i++ {
		fields = append(fields, joined[i].Fields...)
	}

	values := make([]sql.RawBytes, len(scanner.Fields))
	scanArgs := make([]interface{}, len(scanner.Fields))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	if err = row.Scan(scanArgs); err != nil {
		return
	}

	for i, v := range values {
		if v == nil {
			fields[i] = &scapegoat{}
		}
	}

	err = row.Scan(fields)
	return
}

type scapegoat struct{}

func (goat scapegoat) Scan(src interface{}) (err error) {
	return
}