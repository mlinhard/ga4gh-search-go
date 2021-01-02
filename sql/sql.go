package sql

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/mlinhard/ga4gh-search-go/sql/parser"
	"strconv"
	"strings"
	"unicode"
)

type Visitor struct {
	parser.BaseSqlBaseVisitor
	parameterPosition int
	errors            []error
}

func NewVisitor() *Visitor {
	visitor := new(Visitor)
	visitor.BaseSqlBaseVisitor = parser.BaseSqlBaseVisitor{}
	visitor.parameterPosition = 0
	visitor.errors = make([]error, 0, 10)
	return visitor
}

func Parse(statement string) Statement {
	input := NewCaseChangingStream(antlr.NewInputStream(statement), true)
	lexer := parser.NewSqlBaseLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewSqlBaseParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	listener := NewVisitor()
	checkVisitor(listener)
	return listener.Visit(p.Statement())
}

func checkVisitor(visitor parser.SqlBaseVisitor) {

}

// CaseChangingStream wraps an existing CharStream, but upper cases, or
// lower cases the input before it is tokenized.
type CaseChangingStream struct {
	antlr.CharStream

	upper bool
}

// NewCaseChangingStream returns a new CaseChangingStream that forces
// all tokens read from the underlying stream to be either upper case
// or lower case based on the upper argument.
func NewCaseChangingStream(in antlr.CharStream, upper bool) *CaseChangingStream {
	return &CaseChangingStream{in, upper}
}

// LA gets the value of the symbol at offset from the current position
// from the underlying CharStream and converts it to either upper case
// or lower case.
func (is *CaseChangingStream) LA(offset int) int {
	in := is.CharStream.LA(offset)
	if in < 0 {
		// Such as antlr.TokenEOF which is -1
		return in
	}
	if is.upper {
		return int(unicode.ToUpper(rune(in)))
	}
	return int(unicode.ToLower(rune(in)))
}

func (v *Visitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

func (v *Visitor) VisitStatementDefault(ctx *parser.StatementDefaultContext) interface{} {
	return ctx.Query().Accept(v).(*Query)
}

func (v *Visitor) getNode(ctx *antlr.BaseParserRuleContext) Node {
	return Node{Location: v.getLocation(ctx)}
}

func (v *Visitor) getNodeFromToken(ctx antlr.TerminalNode) Node {
	return Node{Location: v.getLocationFromTerminal(ctx)}
}

func (v *Visitor) VisitQuery(ctx *parser.QueryContext) interface{} {
	query := ctx.QueryNoWith().Accept(v).(*Query)
	return &Query{
		Node:      v.getNode(ctx.BaseParserRuleContext),
		Limit:     query.Limit,
		Offset:    query.Offset,
		OrderBy:   query.OrderBy,
		QueryBody: query.QueryBody,
		With:      v.getWith(ctx.With()),
	}
}

func (v *Visitor) getWith(withCtx parser.IWithContext) *With {
	if withCtx == nil {
		return nil
	}
	w := withCtx.Accept(v)
	if w == nil {
		return nil
	}
	return w.(*With)
}

func (v *Visitor) VisitShowTables(ctx *parser.ShowTablesContext) interface{} {
	escape := ctx.ESCAPE().GetText()
	like := ctx.LIKE().GetText()
	return &ShowTables{
		Node:        v.getNode(ctx.BaseParserRuleContext),
		Escape:      &escape,
		LikePattern: &like,
		Schema:      v.VisitQualifiedName(ctx.QualifiedName().(*parser.QualifiedNameContext)).(*QualifiedName),
	}
}

func (v *Visitor) VisitQualifiedName(ctx *parser.QualifiedNameContext) interface{} {
	return &QualifiedName{
		OriginalParts: v.getIdentifiers(ctx.AllIdentifier()),
	}
}

func (v *Visitor) getIdentifiers(identifiers []parser.IIdentifierContext) []*Identifier {
	r := make([]*Identifier, 0, len(identifiers))
	for _, iid := range identifiers {
		id := iid.(*parser.IdentifierContext)
		r = append(r, v.getIdentifier(id))
	}
	return r
}

func (v *Visitor) getIdentifier(ctx parser.IIdentifierContext) *Identifier {
	if ctx == nil {
		return nil
	}
	return ctx.Accept(v).(*Identifier)
}

func (v *Visitor) getLocation(ctx *antlr.BaseParserRuleContext) *NodeLocation {
	if ctx == nil {
		return nil
	}
	start := ctx.GetStart()
	if start == nil {
		return nil
	}
	return &NodeLocation{
		Line:   start.GetLine(),
		Column: start.GetColumn(),
	}
}

func (v *Visitor) getLocationFromTerminal(ctx antlr.TerminalNode) *NodeLocation {
	if ctx == nil {
		return nil
	}
	symbol := ctx.GetSymbol()
	if symbol == nil {
		return nil
	}
	return &NodeLocation{
		Line:   symbol.GetLine(),
		Column: symbol.GetColumn(),
	}
}

func (v *Visitor) VisitShowSchemas(ctx *parser.ShowSchemasContext) interface{} {
	escape := ctx.ESCAPE().GetText()
	like := ctx.LIKE().GetText()
	return &ShowSchemas{
		Node:        v.getNode(ctx.BaseParserRuleContext),
		Catalog:     v.getIdentifier(ctx.Identifier()),
		Escape:      &escape,
		LikePattern: &like,
	}
}

func (v *Visitor) VisitShowCatalogs(ctx *parser.ShowCatalogsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitShowColumns(ctx *parser.ShowColumnsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitShowStats(ctx *parser.ShowStatsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitShowStatsForQuery(ctx *parser.ShowStatsForQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitShowFunctions(ctx *parser.ShowFunctionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitShowSession(ctx *parser.ShowSessionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSetSession(ctx *parser.SetSessionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitResetSession(ctx *parser.ResetSessionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitWith(ctx *parser.WithContext) interface{} {
	return &With{
		Queries:   v.getQueries(ctx.AllNamedQuery()),
		Recursive: nil,
	}
}

func (v *Visitor) getQueries(queries []parser.INamedQueryContext) []*WithQuery {
	r := make([]*WithQuery, len(queries))
	for _, iid := range queries {
		id := iid.(*parser.NamedQueryContext)
		r = append(r, v.VisitNamedQuery(id).(*WithQuery))
	}
	return r
}

func (v *Visitor) getLimit(ctx *parser.QueryNoWithContext) interface{} {
	limitCtx := ctx.GetLimit()
	if limitCtx != nil {
		return limitCtx.Accept(v)
	}
	return nil
}

func (v *Visitor) VisitQueryNoWith(ctx *parser.QueryNoWithContext) interface{} {
	return &Query{
		Node:      v.getNode(ctx.BaseParserRuleContext),
		Limit:     v.getLimit(ctx),
		Offset:    nil,
		OrderBy:   nil,
		QueryBody: v.getQueryBody(ctx.QueryTerm()),
		With:      nil,
	}
}

func (v *Visitor) VisitLimitRowCount(ctx *parser.LimitRowCountContext) interface{} {
	return &Limit{
		RowCount: ctx.RowCount().Accept(v).(*Expression),
	}
}

func (v *Visitor) VisitRowCount(ctx *parser.RowCountContext) interface{} {
	if ctx.INTEGER_VALUE() != nil {
		value, err := strconv.ParseInt(ctx.INTEGER_VALUE().GetText(), 10, 64)
		if err != nil {
			v.processError(err)
			return nil
		}
		return &LongLiteral{Node: v.getNodeFromToken(ctx.INTEGER_VALUE()), Value: &value}
	} else if ctx.PARAMETER() != nil {
		pos := v.parameterPosition
		v.parameterPosition++
		return &Parameter{Node: v.getNodeFromToken(ctx.INTEGER_VALUE()), Position: &pos}
	} else {
		return nil
	}
}

func (v *Visitor) processError(error error) {
	v.errors = append(v.errors, error)
}

func (v *Visitor) VisitQueryTermDefault(ctx *parser.QueryTermDefaultContext) interface{} {
	return ctx.QueryPrimary().Accept(v)
}

func (v *Visitor) VisitSetOperation(ctx *parser.SetOperationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitQueryPrimaryDefault(ctx *parser.QueryPrimaryDefaultContext) interface{} {
	return ctx.QuerySpecification().Accept(v)
}

func (v *Visitor) VisitTable(ctx *parser.TableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitInlineTable(ctx *parser.InlineTableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSubquery(ctx *parser.SubqueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSortItem(ctx *parser.SortItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitQuerySpecification(ctx *parser.QuerySpecificationContext) interface{} {
	return &QuerySpecification{
		Node:    v.getNode(ctx.BaseParserRuleContext),
		From:    v.getFrom(ctx.AllRelation()),
		GroupBy: nil,
		Having:  nil,
		Limit:   nil,
		Offset:  nil,
		OrderBy: nil,
		Select:  v.getSelect(ctx),
		Where:   nil,
	}
}

func (v *Visitor) VisitGroupBy(ctx *parser.GroupByContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSingleGroupingSet(ctx *parser.SingleGroupingSetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitRollup(ctx *parser.RollupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitCube(ctx *parser.CubeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitMultipleGroupingSets(ctx *parser.MultipleGroupingSetsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitGroupingSet(ctx *parser.GroupingSetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitNamedQuery(ctx *parser.NamedQueryContext) interface{} {
	return &WithQuery{
		ColumnNames: ctx.ColumnAliases().Accept(v).([]*Identifier),
		Name:        v.getIdentifier(ctx.Identifier()),
		Query:       ctx.Query().Accept(v).(*Query),
	}
}

func (v *Visitor) VisitSetQuantifier(ctx *parser.SetQuantifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSelectSingle(ctx *parser.SelectSingleContext) interface{} {
	return &SingleColumn{
		Node:       v.getNode(ctx.BaseParserRuleContext),
		Alias:      v.getIdentifier(ctx.Identifier()),
		Expression: v.getExpression(ctx.Expression()),
	}
}

func (v *Visitor) VisitSelectAll(ctx *parser.SelectAllContext) interface{} {
	return &AllColumns{
		Node:    v.getNode(ctx.BaseParserRuleContext),
		Aliases: nil,
		Target:  nil,
	}
}

func (v *Visitor) VisitRelationDefault(ctx *parser.RelationDefaultContext) interface{} {
	return ctx.Accept(v)
}

func (v *Visitor) VisitJoinRelation(ctx *parser.JoinRelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitJoinType(ctx *parser.JoinTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitJoinCriteria(ctx *parser.JoinCriteriaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSampledRelation(ctx *parser.SampledRelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSampleType(ctx *parser.SampleTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitAliasedRelation(ctx *parser.AliasedRelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitColumnAliases(ctx *parser.ColumnAliasesContext) interface{} {
	return v.getIdentifiers(ctx.AllIdentifier())
}

func (v *Visitor) VisitTableName(ctx *parser.TableNameContext) interface{} {
	return &Table{
		Name: ctx.QualifiedName().Accept(v).(*QualifiedName),
	}
}

func (v *Visitor) VisitSubqueryRelation(ctx *parser.SubqueryRelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitUnnest(ctx *parser.UnnestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitLateral(ctx *parser.LateralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitParenthesizedRelation(ctx *parser.ParenthesizedRelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitExpression(ctx *parser.ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitLogicalNot(ctx *parser.LogicalNotContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitPredicated(ctx *parser.PredicatedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitLogicalBinary(ctx *parser.LogicalBinaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitComparison(ctx *parser.ComparisonContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitQuantifiedComparison(ctx *parser.QuantifiedComparisonContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitBetween(ctx *parser.BetweenContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitInList(ctx *parser.InListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitInSubquery(ctx *parser.InSubqueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitLike(ctx *parser.LikeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitNullPredicate(ctx *parser.NullPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDistinctFrom(ctx *parser.DistinctFromContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitValueExpressionDefault(ctx *parser.ValueExpressionDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitConcatenation(ctx *parser.ConcatenationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitArithmeticBinary(ctx *parser.ArithmeticBinaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitArithmeticUnary(ctx *parser.ArithmeticUnaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitAtTimeZone(ctx *parser.AtTimeZoneContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDereference(ctx *parser.DereferenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitTypeConstructor(ctx *parser.TypeConstructorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSpecialDateTimeFunction(ctx *parser.SpecialDateTimeFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSubstring(ctx *parser.SubstringContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitCast(ctx *parser.CastContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitLambda(ctx *parser.LambdaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitParenthesizedExpression(ctx *parser.ParenthesizedExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitParameter(ctx *parser.ParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitNormalize(ctx *parser.NormalizeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitIntervalLiteral(ctx *parser.IntervalLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitNumericLiteral(ctx *parser.NumericLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitBooleanLiteral(ctx *parser.BooleanLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSimpleCase(ctx *parser.SimpleCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitColumnReference(ctx *parser.ColumnReferenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitNullLiteral(ctx *parser.NullLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitRowConstructor(ctx *parser.RowConstructorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSubscript(ctx *parser.SubscriptContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitCurrentPath(ctx *parser.CurrentPathContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSubqueryExpression(ctx *parser.SubqueryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitBinaryLiteral(ctx *parser.BinaryLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitCurrentUser(ctx *parser.CurrentUserContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitExtract(ctx *parser.ExtractContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitStringLiteral(ctx *parser.StringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitArrayConstructor(ctx *parser.ArrayConstructorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitFunctionCall(ctx *parser.FunctionCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitExists(ctx *parser.ExistsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitPosition(ctx *parser.PositionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitSearchedCase(ctx *parser.SearchedCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitGroupingOperation(ctx *parser.GroupingOperationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitBasicStringLiteral(ctx *parser.BasicStringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitUnicodeStringLiteral(ctx *parser.UnicodeStringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitTimeZoneInterval(ctx *parser.TimeZoneIntervalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitTimeZoneString(ctx *parser.TimeZoneStringContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitComparisonOperator(ctx *parser.ComparisonOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitComparisonQuantifier(ctx *parser.ComparisonQuantifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitBooleanValue(ctx *parser.BooleanValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitInterval(ctx *parser.IntervalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitIntervalField(ctx *parser.IntervalFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitNormalForm(ctx *parser.NormalFormContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitRowType(ctx *parser.RowTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitIntervalType(ctx *parser.IntervalTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitArrayType(ctx *parser.ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDoublePrecisionType(ctx *parser.DoublePrecisionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitLegacyArrayType(ctx *parser.LegacyArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitGenericType(ctx *parser.GenericTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDateTimeType(ctx *parser.DateTimeTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitLegacyMapType(ctx *parser.LegacyMapTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitRowField(ctx *parser.RowFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitTypeParameter(ctx *parser.TypeParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitWhenClause(ctx *parser.WhenClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitFilter(ctx *parser.FilterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitUnquotedIdentifier(ctx *parser.UnquotedIdentifierContext) interface{} {
	return &Identifier{
		Node:      Node{Location: v.getLocation(ctx.BaseParserRuleContext)},
		Delimited: false,
		Name:      ctx.IDENTIFIER().GetText(),
	}
}

func (v *Visitor) VisitQuotedIdentifier(ctx *parser.QuotedIdentifierContext) interface{} {
	return &Identifier{
		Node:      Node{Location: v.getLocation(ctx.BaseParserRuleContext)},
		Delimited: true,
		Name:      ctx.QUOTED_IDENTIFIER().GetText(),
	}
}

func (v *Visitor) VisitBackQuotedIdentifier(ctx *parser.BackQuotedIdentifierContext) interface{} {
	return &Identifier{
		Node:      Node{Location: v.getLocation(ctx.BaseParserRuleContext)},
		Delimited: true,
		Name:      ctx.BACKQUOTED_IDENTIFIER().GetText(),
	}
}

func (v *Visitor) VisitDigitIdentifier(ctx *parser.DigitIdentifierContext) interface{} {
	return &Identifier{
		Node:      Node{Location: v.getLocation(ctx.BaseParserRuleContext)},
		Delimited: false,
		Name:      ctx.DIGIT_IDENTIFIER().GetText(),
	}
}

func (v *Visitor) VisitDecimalLiteral(ctx *parser.DecimalLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitDoubleLiteral(ctx *parser.DoubleLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitIntegerLiteral(ctx *parser.IntegerLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) VisitNonReserved(ctx *parser.NonReservedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Visitor) getQueryBody(term parser.IQueryTermContext) QueryBody {
	if term == nil {
		return nil
	}
	return term.Accept(v)
}

func (v *Visitor) getSelect(ctx *parser.QuerySpecificationContext) *Select {
	return &Select{
		Node:        v.getNode(ctx.BaseParserRuleContext),
		Distinct:    v.getDistinct(ctx.SetQuantifier()),
		SelectItems: v.getSelectItems(ctx.AllSelectItem()),
	}
}

func (v *Visitor) getDistinct(setQuantifier parser.ISetQuantifierContext) bool {
	if setQuantifier == nil {
		return false
	} else {
		return strings.ToUpper(setQuantifier.GetText()) == "DISTINCT"
	}
}

func (v *Visitor) getSelectItems(items []parser.ISelectItemContext) []SelectItem {
	r := make([]SelectItem, 0, len(items))
	for _, i := range items {
		r = append(r, v.getSelectItem(i))
	}
	return r
}

func (v *Visitor) getSelectItem(item parser.ISelectItemContext) SelectItem {
	return item.Accept(v)
}

func (v *Visitor) getExpression(expression parser.IExpressionContext) Expression {
	return expression.Accept(v)
}

func (v *Visitor) getFrom(relation []parser.IRelationContext) Relation {
	//TODO parse multiple relations, if more than one, create IMPLICIT join (as in presto)
	return nil
}


