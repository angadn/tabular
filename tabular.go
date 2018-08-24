package tabular

import (
	"fmt"
	"strings"
)

type Tabular struct {
	Name   string
	Fields []string
}

func New(name string, fields ...string) (tabular Tabular) {
	tabular.Name = name
	tabular.Fields = fields
	return
}

func (tabular Tabular) Insertion(queryFmt string) (query string) {
	return fmt.Sprintf(queryFmt, fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES (%s)",
		tabular.Name,
		fmt.Sprintf("`%s`", strings.Join(tabular.Fields, "`, `")),
		fmt.Sprintf("?%s", strings.Repeat(", ?", len(tabular.Fields)-1)),
	))
}

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
