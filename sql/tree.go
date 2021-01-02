package sql

import "strings"

type NodeLocation struct {
	Line   int
	Column int
}

type Node struct {
	Location *NodeLocation
}

type Query struct {
	Node
	Limit     interface{} // Limit or FetchRows
	Offset    *Offset
	OrderBy   *OrderBy
	QueryBody QueryBody
	With      *With
}

type Limit struct {
	RowCount Expression
}

type FetchFirst struct {
	RowCount Expression
	WithTies *bool
}

type Offset struct {
	RowCount Expression
}

type OrderBy struct {
	SortItems []*SortItem
}

type QueryBody interface{} // Table, Values, QuerySpecification

type Relation interface{} // QueryBody

type QuerySpecification struct {
	Node
	From    Relation
	GroupBy *GroupBy
	Having  Expression
	Limit   interface{}
	Offset  *Offset
	OrderBy *OrderBy
	Select  *Select
	Where   Expression
}

type GroupBy struct {
	Node
	Distinct         *bool
	GroupingElements []*GroupingElement
}

type Select struct {
	Node
	Distinct    bool
	SelectItems []SelectItem
}

type SelectItem interface{} // SingleColumn or AllColumns

type SingleColumn struct {
	Node
	Alias      *Identifier
	Expression Expression
}

type AllColumns struct {
	Node
	Aliases []*Identifier
	Target  Expression
}

type ShowTables struct {
	Node
	Escape      *string
	LikePattern *string
	Schema      *QualifiedName
}

type QualifiedName struct {
	OriginalParts []*Identifier
}

func (qn *QualifiedName) String() string {
	s := make([]string, 0, len(qn.OriginalParts))
	for _, n := range qn.OriginalParts {
		s = append(s, n.Name)
	}
	return strings.Join(s, ".")
}

type Identifier struct {
	Node
	Delimited bool
	Name      string
}

type ShowSchemas struct {
	Node
	Catalog     *Identifier
	Escape      *string
	LikePattern *string
}

type LongLiteral struct {
	Node
	Value *int64
}

type Parameter struct {
	Node
	Position *int
}

type Expression interface{}

type SortItem struct {
	NullOrdering *NullOrdering
	Ordering     *Ordering
	SortKey      *Expression
}

type Ordering int

const (
	ASCENDING Ordering = iota
	DESCENDING
)

type NullOrdering int

const (
	FIRST NullOrdering = iota
	LAST
	UNDEFINED
)

type GroupingElement interface { // Cube, Rollup, SimpleGroupBy, GroupingSets
}

type GroupingSets struct {
	Sets [][]*Expression
}

type SimpleGroupBy struct {
	Columns []*Expression
}

type Rollup struct {
	Columns []*Expression
}

type Cube struct {
	Columns []*Expression
}

type With struct {
	Queries   []*WithQuery
	Recursive *bool
}

type WithQuery struct {
	ColumnNames []*Identifier
	Name        *Identifier
	Query       *Query
}

type Statement interface{} // Query, ShowTables, ShowSchemas

type Table struct {
	Name *QualifiedName
}

type Values struct {
	Rows []*Expression
}

/////////////// PROCESSED /////////////////////////
//
//type PathElement struct {
//	Catalog *Identifier
//	Schema  *Identifier
//}
//
//type SimpleCaseExpression struct {
//	DefaultValue *Expression
//	Operand      *Expression
//	WhenClauses  []*WhenClause
//}
//
//type TableSubquery struct {
//	QuerySpecification *QuerySpecification
//}
//
//type JoinOn struct {
//	Expression *Expression
//}
//
//type TypeParameter struct {
//	Value UNKNOWN
//}

//type Join struct {
//	Criteria *JoinCriteria
//	Left     *Relation
//	Right    *Relation
//	Type     *Type
//}
//type BindExpression struct {
//	Function *Expression
//	Values   []*Expression
//}
//type ShowRoles struct {
//	Catalog *Identifier
//	Current *bool
//}
//type PrincipalSpecification struct {
//	Name *Identifier
//	Type *Type
//}
//type CallArgument struct {
//	Name  *string
//	Value *Expression
//}
//type DecimalLiteral struct {
//	Value *string
//}
//
//type RefreshMaterializedView struct {
//	Name *QualifiedName
//}
//type Property struct {
//	Name  *Identifier
//	Value *Expression
//}
//
//type JoinUsing struct {
//	Columns []*Identifier
//	Nodes   UNKNOWN
//}
//type DropMaterializedView struct {
//	Exists *bool
//	Name   *QualifiedName
//}

//type ShowCatalogs struct {
//	Escape      *string
//	LikePattern *string
//}
//type ExistsPredicate struct {
//	Subquery *Expression
//}
//type NaturalJoin struct {
//	Nodes UNKNOWN
//}
//type Literal struct {
//}
//type DescribeInput struct {
//	Name *Identifier
//}
//type WindowFrame struct {
//	End   *FrameBound
//	Start *FrameBound
//	Type  *Type
//}
//type TransactionAccessMode struct {
//	ReadOnly *bool
//}
//type RenameTable struct {
//	Exists *bool
//	Source *QualifiedName
//	Target *QualifiedName
//}
//type SetSchemaAuthorization struct {
//	Principal UNKNOWN
//	Source    UNKNOWN
//}
//type ExplainFormat struct {
//	Type *Type
//}
//type CurrentPath struct {
//}
//
//type ExpressionTreeRewriter struct {
//}
//type Call struct {
//	Arguments []*CallArgument
//	Name      *QualifiedName
//}
//type Extract struct {
//	Expression *Expression
//	Field      *Field
//}
//type Intersect struct {
//	Distinct  *bool
//	Relations []*Relation
//}
//type NotExpression struct {
//	Value *Expression
//}
//type NumericParameter struct {
//	Value *string
//}
//type DropView struct {
//	Exists *bool
//	Name   *QualifiedName
//}
//type RenameColumn struct {
//	ColumnExists *bool
//	Source       *Identifier
//	Table        *QualifiedName
//	TableExists  *bool
//	Target       *Identifier
//}
//type ResetSession struct {
//	Name *QualifiedName
//}
//type ArrayConstructor struct {
//	Values []*Expression
//}
//type ShowStats struct {
//	Relation *Relation
//}
//type FrameBound struct {
//	Type  *Type
//	Value *Expression
//}
//type DereferenceExpression struct {
//	Base  *Expression
//	Field *Identifier
//}
//type Prepare struct {
//	Name      *Identifier
//	Statement *Statement
//}
//type BinaryLiteral struct {
//	Value *Slice
//}
//type GrantRoles struct {
//	AdminOption *bool
//	Grantees    []*PrincipalSpecification
//	Grantor     *GrantorSpecification
//	Roles       []*Identifier
//}
//type ArithmeticUnaryExpression struct {
//	Sign  *Sign
//	Value *Expression
//}
//type IntervalDayTimeDataType struct {
//	From *Field
//	To   *Field
//}
//type Explain struct {
//	Analyze   *bool
//	Options   []*ExplainOption
//	Statement *Statement
//	Verbose   *bool
//}
//type IntervalLiteral struct {
//	EndField    *IntervalField
//	Sign        *Sign
//	StartField  *IntervalField
//	Value       *string
//	YearToMonth *bool
//}
//type ComparisonExpression struct {
//	Left     *Expression
//	Operator *Operator
//	Right    *Expression
//}
//type LambdaExpression struct {
//	Arguments []*LambdaArgumentDeclaration
//	Body      *Expression
//}
//type ShowCreate struct {
//	Name *QualifiedName
//	Type *Type
//}
//type GenericLiteral struct {
//	Type  *string
//	Value *string
//}
//type ExplainType struct {
//	Type *Type
//}
//type QuantifiedComparisonExpression struct {
//	Operator   *Operator
//	Quantifier *Quantifier
//	Subquery   *Expression
//	Value      *Expression
//}
//type DescribeOutput struct {
//	Name *Identifier
//}
//type CharLiteral struct {
//	Slice *Slice
//	Value *string
//}
//type DateTimeDataType struct {
//	Precision    *DataTypeParameter
//	Type         *Type
//	WithTimeZone *bool
//}
//
//type CreateSchema struct {
//	NotExists  *bool
//	Principal  *PrincipalSpecification
//	Properties []*Property
//	SchemaName *QualifiedName
//}
//type ShowRoleGrants struct {
//	Catalog *Identifier
//}
//type LikeClause struct {
//	PropertiesOption *PropertiesOption
//	TableName        *QualifiedName
//}
//type DoubleLiteral struct {
//	Value *float64
//}
//type SetViewAuthorization struct {
//	Principal UNKNOWN
//	Source    UNKNOWN
//}
//type GrantorSpecification struct {
//	Principal *PrincipalSpecification
//	Type      *Type
//}
//type TryExpression struct {
//	InnerExpression *Expression
//}

//type DropColumn struct {
//	Column       *Identifier
//	ColumnExists *bool
//	Table        *QualifiedName
//	TableExists  *bool
//}
//type Row struct {
//	Items []*Expression
//}
//type BooleanLiteral struct {
//	Value *bool
//}
//type Insert struct {
//	Columns []*Identifier
//	QuerySpecification   *QuerySpecification
//	Target  *QualifiedName
//}
//type Commit struct {
//}
//type ExpressionRewriter struct {
//}
//type SetSession struct {
//	Name  *QualifiedName
//	Value *Expression
//}
//type Delete struct {
//	Table *Table
//	Where *Expression
//}
//type SubqueryExpression struct {
//	QuerySpecification *QuerySpecification
//}
//type Window struct {
//	Frame       *WindowFrame
//	OrderBy     *OrderBy
//	PartitionBy []*Expression
//}
//type Isolation struct {
//	Level *Level
//}

//type IfExpression struct {
//	Condition  *Expression
//	FalseValue *Expression
//	TrueValue  *Expression
//}
//type BetweenPredicate struct {
//	Max   *Expression
//	Min   *Expression
//	Value *Expression
//}
//
//type RenameView struct {
//	Source *QualifiedName
//	Target *QualifiedName
//}
//type ArithmeticBinaryExpression struct {
//	Left     *Expression
//	Operator *Operator
//	Right    *Expression
//}
//type Union struct {
//	Distinct  *bool
//	Relations []*Relation
//}
//type DropSchema struct {
//	Cascade    *bool
//	Exists     *bool
//	SchemaName *QualifiedName
//}
//type CurrentTime struct {
//	Function  *Function
//	Precision *int
//}
//
//type SampledRelation struct {
//	Relation         *Relation
//	SamplePercentage *Expression
//	Type             *Type
//}
//type Field struct {
//	Name *Identifier
//	Type *DataType
//}
//
//type ExplainOption struct {
//}
//type CreateView struct {
//	Comment  *string
//	Name     *QualifiedName
//	QuerySpecification    *QuerySpecification
//	Replace  *bool
//	Security *Security
//}
//type SetAuthorizationStatement struct {
//	Principal *PrincipalSpecification
//	Source    *QualifiedName
//}
//type CreateTableAsSelect struct {
//	ColumnAliases []*Identifier
//	Comment       *string
//	Name          *QualifiedName
//	NotExists     *bool
//	Properties    []*Property
//	QuerySpecification         *QuerySpecification
//	WithData      *bool
//}
//type CoalesceExpression struct {
//	Operands []*Expression
//}
//type InListExpression struct {
//	Values []*Expression
//}
//type CreateMaterializedView struct {
//	Comment    *string
//	Name       *QualifiedName
//	NotExists  *bool
//	Properties []*Property
//	QuerySpecification      *QuerySpecification
//	Replace    *bool
//}
//type RowDataType struct {
//	Fields []*Field
//}

//type RewritingVisitor struct {
//}
//type SubscriptExpression struct {
//	Base  *Expression
//	Index *Expression
//}
//type WhenClause struct {
//	Operand *Expression
//	Result  *Expression
//}
//type Format struct {
//	Arguments []*Expression
//}
//type ShowFunctions struct {
//	Escape      *string
//	LikePattern *string
//}
//type JoinCriteria struct {
//	Nodes UNKNOWN
//}
//type Execute struct {
//	Name       *Identifier
//	Parameters []*Expression
//}
//type DataType struct {
//}
//type LogicalBinaryExpression struct {
//	Left     *Expression
//	Operator *Operator
//	Right    *Expression
//}
//type
//struct {
//}
//type SetTableAuthorization struct {
//	Principal UNKNOWN
//	Source    UNKNOWN
//}
//
//type FieldReference struct {
//	FieldIndex *int
//}

//type ShowSession struct {
//	Escape      *string
//	LikePattern *string
//}
//type SearchedCaseExpression struct {
//	DefaultValue *Expression
//	WhenClauses  []*WhenClause
//}
//type DataTypeParameter struct {
//}
//type PathSpecification struct {
//	Path []*PathElement
//}

//type IsNotNullPredicate struct {
//	Value *Expression
//}
//type ColumnDefinition struct {
//	Comment    *string
//	Name       *Identifier
//	Nullable   *bool
//	Properties []*Property
//	Type       *DataType
//}

//
//type DropTable struct {
//	Exists    *bool
//	TableName *QualifiedName
//}

//type SetPath struct {
//	PathSpecification *PathSpecification
//}
//type StartTransaction struct {
//	TransactionModes []*TransactionMode
//}
//type RevokeRoles struct {
//	AdminOption *bool
//	Grantees    []*PrincipalSpecification
//	Grantor     *GrantorSpecification
//	Roles       []*Identifier
//}
//type InPredicate struct {
//	Value     *Expression
//	ValueList *Expression
//}

//type
//struct {
//}
//type Unnest struct {
//	Expressions    []*Expression
//	WithOrdinality *bool
//}
//type RenameSchema struct {
//	Source *QualifiedName
//	Target *Identifier
//}
//type DropRole struct {
//	Name *Identifier
//}

//type TransactionMode struct {
//}
//type CurrentUser struct {
//}
//type TableElement struct {
//}
//type SetOperation struct {
//	Distinct  *bool
//	Relations UNKNOWN
//}

//type SymbolReference struct {
//	Name *string
//}

//type GenericDataType struct {
//	Arguments []*DataTypeParameter
//	Name      *Identifier
//}
//
//type GroupingOperation struct {
//	GroupingColumns []*Expression
//}
//type AtTimeZone struct {
//	TimeZone *Expression
//	Value    *Expression
//}
//type NullLiteral struct {
//}
//type Cast struct {
//	Expression *Expression
//	Safe       *bool
//	Type       *DataType
//	TypeOnly   *bool
//}

//type Deallocate struct {
//	Name *Identifier
//}
//type DefaultTraversalVisitor struct {
//}
//type Comment struct {
//	Comment *string
//	Name    *QualifiedName
//	Type    *Type
//}
//
//type AliasedRelation struct {
//	Alias       *Identifier
//	ColumnNames []*Identifier
//	Relation    *Relation
//}
//type Lateral struct {
//	QuerySpecification *QuerySpecification
//}
//type DefaultExpressionTraversalVisitor struct {
//}
//type AllRows struct {
//}
//type AstVisitor struct {
//}
//type Grant struct {
//	Grantee         *PrincipalSpecification
//	Name            *QualifiedName
//	Privileges      []*string
//	Type            *GrantOnType
//	WithGrantOption *bool
//}
//type ShowGrants struct {
//	Table     *bool
//	TableName *QualifiedName
//}
//type
//struct {
//}
//type ShowColumns struct {
//	Escape      *string
//	LikePattern *string
//	Table       *QualifiedName
//}
//type IsNullPredicate struct {
//	Value *Expression
//}
//type FunctionCall struct {
//	Arguments     []*Expression
//	Distinct      *bool
//	Filter        *Expression
//	Name          *QualifiedName
//	NullTreatment *NullTreatment
//	OrderBy       *OrderBy
//	Window        *Window
//}
//type Use struct {
//	Catalog *Identifier
//	Schema  *Identifier
//}

//type LikePredicate struct {
//	Escape  *Expression
//	Pattern *Expression
//	Value   *Expression
//}
//type AddColumn struct {
//	Column          *ColumnDefinition
//	ColumnNotExists *bool
//	Name            *QualifiedName
//	TableExists     *bool
//}
//type Context struct {
//	DefaultRewrite *bool
//}
//type Revoke struct {
//	GrantOptionFor *bool
//	Grantee        *PrincipalSpecification
//	Name           *QualifiedName
//	Privileges     []*string
//	Type           *GrantOnType
//}

//type TimestampLiteral struct {
//	Value *string
//}
//type TimeLiteral struct {
//	Value *string
//}
//type CreateRole struct {
//	Grantor *GrantorSpecification
//	Name    *Identifier
//}

//type CreateTable struct {
//	Comment    *string
//	Elements   []*TableElement
//	Name       *QualifiedName
//	NotExists  *bool
//	Properties []*Property
//}
//type Analyze struct {
//	Properties []*Property
//	TableName  *QualifiedName
//}
//type NullIfExpression struct {
//	First  *Expression
//	Second *Expression
//}
//type StringLiteral struct {
//	Slice *Slice
//	Value *string
//}
//type Except struct {
//	Distinct  *bool
//	Left      *Relation
//	Relations UNKNOWN
//	Right     *Relation
//}
//type LambdaArgumentDeclaration struct {
//	Name *Identifier
//}
