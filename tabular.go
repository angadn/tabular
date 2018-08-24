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

func (tabular Tabular) Insertion() (sql string) {
	return fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES (%s)",
		tabular.Name,
		fmt.Sprintf("`%s`", strings.Join(tabular.Fields, "`, `")),
		fmt.Sprintf("?%s", strings.Repeat(", ?", len(tabular.Fields)-1)),
	)
}

func (tabular Tabular) Selection(joined ...Tabular) (sql string) {
	sql = fmt.Sprintf("`%s`.`%s`", tabular.Name, strings.Join(
		tabular.Fields, fmt.Sprintf("`, `%s`.`", tabular.Name),
	))

	if len(joined) > 0 {
		sql = fmt.Sprintf("%s, %s", sql, joined[0].Selection(joined[1:]...))
	}

	return
}
