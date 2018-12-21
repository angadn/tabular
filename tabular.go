package tabular

import (
	"fmt"
	"strings"
)

// Tabular is a representation fo a Database Table, with it's Name and columns listed in
// a serial order.
type Tabular struct {
	Name       string
	Fields     []string
	selectNull bool
}

// New constructs a Tabular given a table-name and column-fields.
func New(name string, fields ...string) (tabular Tabular) {
	tabular.Name = name
	tabular.Fields = fields
	return
}

// WithAlias returns a Tabular with it's Name set to an alias. This is useful when
// JOINing to the same table twice in a single query.
func (tabular Tabular) WithAlias(alias string) (aliased Tabular) {
	aliased = tabular
	aliased.Name = alias
	return
}

// WithNullSelection returns a Tabular that when selecting will return NULL once for
// each of it's fields. Useful when you want to perform UNIONs after JOINs.
func (tabular Tabular) WithNullSelection() (nulled Tabular) {
	nulled = tabular
	nulled.selectNull = true
	return
}

// Insertion generates an "INSERT INTO [Name] (...) VALUES (?...)" query. You can also
// configure specific keys to be evaluated as expressions. A common example
// is Insertion("%s", "created_at", "NOW()", "updated_at", "NOW()"). The returned query
// can be suffixed with "ON DUPLICATE KEY UPDATE..." for UPSERT-like behaviour.
func (tabular Tabular) Insertion(queryFmt string, keyval ...string) (query string) {
	query = tabular.BatchInsertion(queryFmt, 1, keyval...)
	return
}

// BatchInsertion generates an "INSERT INTO [Name] (...) VALUES (?...)" query just like
// Insertion, but for more than one Row.
func (tabular Tabular) BatchInsertion(
	queryFmt string, rows int, keyval ...string,
) (query string) {
	values := strings.Split(strings.Repeat("?", len(tabular.Fields)), "")
	for i := 0; i < len(keyval); i += 2 {
		for j, field := range tabular.Fields {
			if keyval[i] == field {
				values[j] = keyval[i+1]
				break
			}
		}
	}

	return fmt.Sprintf(queryFmt, fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES %s",
		tabular.Name,
		fmt.Sprintf("`%s`", strings.Join(tabular.Fields, "`, `")),
		strings.TrimRight(strings.Repeat(
			fmt.Sprintf("(%s), ", strings.Join(values, ", ")), rows,
		), ", ")),
	)
}

// Selection generates a list of all field-names including the joined ones in serial
// order. It is a fragment of SQL, and is hence unexecutable by itself. It needs to be
// prefixed with a "SELECT " and suffixed with a " FROM...", that includes and defines
// the specific JOIN relationships. The fragment is incomplete by design, as it allows
// the user to define virtual fields, etc. in the query if required.
func (tabular Tabular) Selection(queryFmt string, joined ...Tabular) (query string) {
	return tabular.selection(queryFmt, false, joined...)
}

// PrefixedSelection generates field-names same as Selection except here the field-names
// will be prefixed with their table name.
func (tabular Tabular) PrefixedSelection(queryFmt string, joined ...Tabular) (query string) {
	return tabular.selection(queryFmt, true, joined...)
}

func (tabular Tabular) selection(
	queryFmt string, tableAsPrefix bool, joined ...Tabular,
) (query string) {
	fieldFmt := "`%s`.`%s`"
	prefixFmt := "`%s_%s`"

	q := func(t Tabular) (ret string) {
		for _, field := range t.Fields {
			formattedField := fmt.Sprintf(fieldFmt, t.Name, field)
			if tableAsPrefix {
				prefixed := fmt.Sprintf(prefixFmt, t.Name, field)
				if t.selectNull {
					formattedField = fmt.Sprintf("NULL AS %s", prefixed)
				} else {
					formattedField = fmt.Sprintf("%s AS %s", formattedField, prefixed)
				}
			} else if t.selectNull {
				formattedField = fmt.Sprintf("NULL AS `%s`", field)
			}

			ret = fmt.Sprintf("%s, %s", ret, formattedField)
		}

		ret = strings.TrimLeft(ret, ", ")
		return
	}

	query = q(tabular)

	for _, j := range joined {
		query = fmt.Sprintf("%s, %s", query, q(j))
	}

	query = fmt.Sprintf(queryFmt, query)
	return
}
