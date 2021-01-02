package tests

import (
	"github.com/mlinhard/ga4gh-search-go/sql"
	"testing"
)

func Test_SimpleQuery(t *testing.T) {
	query := parse(t, "SELECT * FROM PERSONS").
		Query().
		AssertOnlyDef("Node", "QueryBody").
		QueryBody().
		QuerySpecification()
	query.Select().
		AssertNumItems(1).
		SelectItem(0).
		AllColumns().
		AssertOnlyDef("Node")
	query.From().
		Table().
		AssertName("PERSONS")
}

func parse(t *testing.T, statementSql string) *StatementTester {
	statement := sql.Parse(statementSql)
	return NewStatementTester(t, statement)
}
