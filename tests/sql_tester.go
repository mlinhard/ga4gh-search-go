package tests

import (
	"github.com/mlinhard/ga4gh-search-go/sql"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type StatementTester struct {
	t         *testing.T
	Statement sql.Statement
	typeQuery reflect.Type
}

type QueryTester struct {
	t     *testing.T
	Query *sql.Query
}

type QueryBodyTester struct {
	t         *testing.T
	QueryBody sql.QueryBody
	typeQuery reflect.Type
}

type QuerySpecificationTester struct {
	t                  *testing.T
	QuerySpecification *sql.QuerySpecification
}

type SelectTester struct {
	t      *testing.T
	Select *sql.Select
}

type RelationTester struct {
	t        *testing.T
	Relation sql.Relation
}

type SelectItemTester struct {
	t          *testing.T
	SelectItem sql.SelectItem
}

type SingleColumnTester struct {
	t            *testing.T
	SingleColumn *sql.SingleColumn
}

type AllColumnsTester struct {
	t          *testing.T
	AllColumns *sql.AllColumns
}

type ExpressionTester struct {
	t          *testing.T
	Expression sql.Expression
}

type TableTester struct {
	t     *testing.T
	Table *sql.Table
}

func (t *SingleColumnTester) AssertHasExpression() *SingleColumnTester {
	assert.NotNil(t.t, t.SingleColumn.Expression)
	return t
}

func (t *SingleColumnTester) Expression() *SingleColumnTester {
	t.AssertHasExpression()
	return t
}

func NewStatementTester(t *testing.T, statement sql.Statement) *StatementTester {
	if statement == nil {
		t.Error("Statement cannot be nil")
	}
	return &StatementTester{
		t:         t,
		Statement: statement,
		typeQuery: reflect.TypeOf(&sql.Query{}),
	}
}

func (s *StatementTester) Query() *QueryTester {
	s.AssertQuery()
	return &QueryTester{s.t, s.Statement.(*sql.Query)}
}

func (s *StatementTester) AssertQuery() *StatementTester {
	t := reflect.TypeOf(s.Statement)
	assert.Equal(s.t, s.typeQuery, t, "%v != %v", s.typeQuery.Name(), t.Name())
	return s
}

func (t *QueryTester) AssertHasWith() *QueryTester {
	assert.NotNil(t.t, t.Query.With)
	return t
}

func (t *QueryTester) AssertHasQueryBody() *QueryTester {
	assert.NotNil(t.t, t.Query.QueryBody, "QueryBody")
	return t
}

func (t *QueryTester) AssertOnlyDef(definedProperties ...string) *QueryTester {
	AssertOnlyDef(t.t, reflect.ValueOf(t.Query), definedProperties...)
	return t
}

func (t *QueryTester) QueryBody() *QueryBodyTester {
	t.AssertHasQueryBody()
	return &QueryBodyTester{
		t:         t.t,
		QueryBody: t.Query.QueryBody,
		typeQuery: reflect.TypeOf(&sql.QuerySpecification{}),
	}
}

func (t *QueryBodyTester) AssertQuerySpecification() *QueryBodyTester {
	qt := reflect.TypeOf(t.QueryBody)
	assert.Equal(t.t, t.typeQuery, qt, "%v != %v", t.typeQuery.Name(), qt.Name())
	return t
}

func (t *QueryBodyTester) QuerySpecification() *QuerySpecificationTester {
	t.AssertQuerySpecification()
	return &QuerySpecificationTester{
		t:                  t.t,
		QuerySpecification: t.QueryBody.(*sql.QuerySpecification),
	}
}

func (t *QuerySpecificationTester) AssertHasSelect() *QuerySpecificationTester {
	assert.NotNil(t.t, t.QuerySpecification.Select, "select")
	return t
}

func (t *QuerySpecificationTester) AssertHasFrom() *QuerySpecificationTester {
	assert.NotNil(t.t, t.QuerySpecification.From, "from")
	return t
}

func (t *QuerySpecificationTester) Select() *SelectTester {
	t.AssertHasSelect()
	return &SelectTester{
		t:      t.t,
		Select: t.QuerySpecification.Select,
	}
}

func (t *QuerySpecificationTester) From() *RelationTester {
	t.AssertHasFrom()
	return &RelationTester{
		t:        t.t,
		Relation: t.QuerySpecification.From,
	}
}

func (t *SelectTester) AssertDistinct() *SelectTester {
	assert.NotNil(t.t, t.Select.Distinct)
	assert.True(t.t, t.Select.Distinct)
	return t
}

func (t *SelectTester) AssertNumItems(expectedNumItems int) *SelectTester {
	assert.NotNil(t.t, t.Select.SelectItems)
	assert.Equal(t.t, expectedNumItems, len(t.Select.SelectItems))
	return t
}

func (t *SelectTester) SelectItem(index int) *SelectItemTester {
	assert.NotNil(t.t, t.Select.SelectItems)
	assert.Greater(t.t, len(t.Select.SelectItems), index)
	return &SelectItemTester{
		t:          t.t,
		SelectItem: t.Select.SelectItems[index],
	}
}

func (t *SelectItemTester) AssertSingleColumn() *SelectItemTester {
	typeExpect := reflect.TypeOf(&sql.SingleColumn{})
	typeActual := reflect.TypeOf(t.SelectItem)
	assert.Equal(t.t, typeExpect, typeActual, "%v != %v", typeExpect.Name(), typeActual.Name())
	return t
}

func (t *SelectItemTester) AssertAllColumns() *SelectItemTester {
	typeExpect := reflect.TypeOf(&sql.AllColumns{})
	typeActual := reflect.TypeOf(t.SelectItem)
	assert.Equal(t.t, typeExpect, typeActual, "%v != %v", typeExpect.Name(), typeActual.Name())
	return t
}

func (t *SelectItemTester) SingleColumn() *SingleColumnTester {
	t.AssertSingleColumn()
	return &SingleColumnTester{
		t:            t.t,
		SingleColumn: t.SelectItem.(*sql.SingleColumn),
	}
}

func (t *SelectItemTester) AllColumns() *AllColumnsTester {
	t.AssertAllColumns()
	return &AllColumnsTester{
		t:          t.t,
		AllColumns: t.SelectItem.(*sql.AllColumns),
	}
}

func (t *AllColumnsTester) AssertOnlyDef(definedProperties ...string) *AllColumnsTester {
	AssertOnlyDef(t.t, reflect.ValueOf(t.AllColumns), definedProperties...)
	return t
}

func (t *RelationTester) AssertTable() *RelationTester {
	typeExpect := reflect.TypeOf(&sql.Table{})
	typeActual := reflect.TypeOf(t.Relation)
	assert.Equal(t.t, typeExpect, typeActual, "%v != %v", typeExpect.Name(), typeActual.Name())
	return t
}

func (t *RelationTester) Table() *TableTester {
	t.AssertTable()
	return &TableTester{
		t:     t.t,
		Table: t.Relation.(*sql.Table),
	}
}

func (t *TableTester) AssertName(expectedName string) *TableTester {
	assert.Equal(t.t, expectedName, t.Table.Name.String())
	return t
}
