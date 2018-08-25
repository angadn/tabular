package tabular

import "database/sql"

// Scanner contains pointers to fields that we can Scan SQL values into.
type Scanner struct {
	Fields []interface{}
}

// NewScanner is a constructor for Scanner.
func NewScanner(fields ...interface{}) (scanner Scanner) {
	scanner.Fields = fields
	return
}

// Scan rows into the fields for the Scanner, chaining with any other scanner in the case
// of joins.
func (scanner Scanner) Scan(rows *sql.Rows, joined ...Scanner) (err error) {
	fields := scanner.Fields
	for i := 0; i < len(joined); i++ {
		fields = append(fields, joined[i].Fields...)
	}

	values := make([]sql.RawBytes, len(fields))
	scanArgs := make([]interface{}, len(fields))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	if err = rows.Scan(scanArgs...); err != nil {
		return
	}

	for i, v := range values {
		if v == nil {
			fields[i] = &Scapegoat{}
		}
	}

	err = rows.Scan(fields...)
	return
}

// Scapegoat is to receive values in Scan that we're not interested in.
type Scapegoat struct{}

// Scan does nothing upon receiving data.
func (goat Scapegoat) Scan(src interface{}) (err error) {
	return
}
