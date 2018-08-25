package tabular

import (
	"fmt"
	"strings"
)

// Tabular is a representation fo a Database Table, with it's Name and columns listed in
// a serial order.
type Tabular struct {
	Name   string
	Fields []string
}

// New constructs a Tabular given a table-name and column-fields.
func New(name string, fields ...string) (tabular Tabular) {
	tabular.Name = name
	tabular.Fields = fields
	return
}

// Insertion generates an "INSERT INTO [Name] (...) VALUES (?...)" query.
func (tabular Tabular) Insertion(queryFmt string) (query string) {
	return fmt.Sprintf(queryFmt, fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES (%s)",
		tabular.Name,
		fmt.Sprintf("`%s`", strings.Join(tabular.Fields, "`, `")),
		fmt.Sprintf("?%s", strings.Repeat(", ?", len(tabular.Fields)-1)),
	))
}

// Selection generates a list of all field-names including the joined ones in serial
// order. It is a fragment of SQL, and is hence unexecutable by itself. It needs to be
// prefixed with a "SELECT " and suffixed with a " FROM...", that includes and defines
// the specific JOIN relationships. The fragment is incomplete by design, as it allows
// the user to define virtual fields, etc. in the query if required.
func (tabular Tabular) Selection(queryFmt string, joined ...Tabular) (query string) {
	q := func(t Tabular) string {
		return fmt.Sprintf("`%s`.`%s`", t.Name, strings.Join(
			t.Fields, fmt.Sprintf("`, `%s`.`", t.Name),
		))
	}

	query = q(tabular)

	for _, j := range joined {
		query = fmt.Sprintf("%s, %s", query, q(j))
	}

	query = fmt.Sprintf(queryFmt, query)
	return
}
