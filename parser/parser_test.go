package parser

import (
	"fmt"
	"testing"

	"github.com/vshalt/arbok/ast"
	"github.com/vshalt/arbok/lexer"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let z = y;", "z", "y"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
		}
		statement := program.Statements[0]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
		value := statement.(*ast.LetStatement).Value
		if !testLiteralExpression(t, value, tt.expectedValue) {
			return
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"foobar";`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	statement := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := statement.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expression is not ast.StringLiteral, got=%T", statement.Expression)
	}
	if literal.Value != "foobar" {
		t.Fatalf("literal value not %q, got=%q", "foobar", literal.Value)
	}
}

func TestArrayLiteralExpression(t *testing.T) {
	input := `[1, 2*2, 3 + 3];`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	statement := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := statement.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expression is not ast.ArrayLiteral, got=%T", statement.Expression)
	}
	if len(literal.Elements) != 3 {
		t.Fatalf("length of array not %d, got=%d", 3, len(literal.Elements))
	}
	testIntegerLiteral(t, literal.Elements[0], 1)
	testInfixExpression(t, literal.Elements[1], 2, "*", 2)
	testInfixExpression(t, literal.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpression(t *testing.T) {
	input := `myArray[1+1];`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	statement := program.Statements[0].(*ast.ExpressionStatement)
	expr, ok := statement.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("expression is not ast.IndexExpression, got=%T", statement.Expression)
	}
	if !testIdentifier(t, expr.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, expr.Index, 1, "+", 1) {
		return
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 1000;
`
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
	}
	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement not *ast.ReturnStatement, got=%T", statement)
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral() does not return 'return', got %q", returnStatement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}
	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	identifier, ok := expressionStatement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expressionStatement.Expression is not identifier, got=%T", expressionStatement.Expression)
	}
	if identifier.Value != "foobar" {
		t.Errorf("identifier.Value is not %q, got=%q", "foobar", identifier.Value)
	}
	if identifier.TokenLiteral() != "foobar" {
		t.Errorf("identifier.TokenLiteral() is not %q, got=%q", "foobar", identifier.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}
	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	literal, ok := expressionStatement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expressionStatement is not ast.integerLiteral, got=%T", program.Statements[0])
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value is not %d, got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() is not %q, got=%q", "5", literal.TokenLiteral())
	}
}

func TestBooleanLiteralExpression(t *testing.T) {
	input := `true;`
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}
	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	boolean, ok := expressionStatement.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("expressionStatement is not ast.Boolean, got=%T", program.Statements[0])
	}
	if boolean.Value != true {
		t.Errorf("literal.Value is not %t , got=%t", true, boolean.Value)
	}
	if boolean.TokenLiteral() != "true" {
		t.Errorf("literal.TokenLiteral() is not %s, got=%s", "true", boolean.TokenLiteral())
	}
}

func TestParsePrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statement does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}
		expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not a ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		prefixExpression, ok := expressionStatement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("statement is not ast.PrefixExpression, got=%T", expressionStatement.Expression)
		}
		if prefixExpression.Operator != tt.operator {
			t.Fatalf("expression.Operator is not %q, got=%q", tt.operator, prefixExpression.Operator)
		}
		if !testLiteralExpression(t, prefixExpression.Right, tt.value) {
			return
		}
	}
}

func TestParseInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not have %d statements. got=%d\n", 1, len(program.Statements))
		}

		expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement, got=%T", program.Statements[0])
		}
		expr, ok := expressionStatement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expressionStatement is not *ast.InfixExpression, got=%T", expressionStatement.Expression)
		}
		if !testLiteralExpression(t, expr.Left, tt.leftValue) {
			return
		}
		if expr.Operator != tt.operator {
			t.Fatalf("expr.Operator is not %q, got=%q", tt.operator, expr.Operator)
		}
		if !testLiteralExpression(t, expr.Right, tt.rightValue) {
			return
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) {x};`
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not have %d statements. got=%d\n", 1, len(program.Statements))
	}

	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	expr, ok := expressionStatement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expressionStatement is not *ast.IfExpression, got=%T", expressionStatement.Expression)
	}
	if !testInfixExpression(t, expr.Condition, "x", "<", "y") {
		return
	}
	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("Consequence is not 1 statements. got=%d\n", len(expr.Consequence.Statements))
	}
	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statements[0] is not *ast.ExpressionStatement, got=%T", expr.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if expr.Alternative != nil {
		t.Errorf("expr.Alternative.Statements was not nil, got=%+v", expr.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) {x} else {y};`
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not have %d statements. got=%d\n", 1, len(program.Statements))
	}

	expressionStatement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	expr, ok := expressionStatement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expressionStatement is not *ast.IfExpression, got=%T", expressionStatement.Expression)
	}
	if !testInfixExpression(t, expr.Condition, "x", "<", "y") {
		return
	}
	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("Consequence is not 1 statements. got=%d\n", len(expr.Consequence.Statements))
	}
	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statements[0] is not *ast.ExpressionStatement, got=%T", expr.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if expr.Alternative == nil {
		t.Fatalf("expr.Alternative.Statements was nil")
	}
	if len(expr.Alternative.Statements) != 1 {
		t.Errorf("Alternative is not 1 statements. got=%d\n", len(expr.Alternative.Statements))
	}
	alternative, ok := expr.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statements[0] is not *ast.ExpressionStatement, got=%T", expr.Alternative.Statements[0])
	}
	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func(x, y) {x + y;}`
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("Program does not contain %d statements, got=%d", 1, len(program.Statements))
	}
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	function, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not *ast.FunctionLiteral, got=%T", statement.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong, want=2, got=%d", len(function.Parameters))
	}
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements is not %d, got=%d", 1, len(function.Body.Statements))
	}

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not *ast.ExpressionStatement, got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {

	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "func() {};", expectedParams: []string{}},
		{input: "func(x) {};", expectedParams: []string{"x"}},
		{input: "func(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length of parameters wrong, got=%d", len(tt.expectedParams))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statement, got=%d", 1, len(program.Statements))
	}
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	callExpression, ok := statement.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("statement.Expression is not *ast.CallExpression, got=%T", statement.Expression)
	}

	if !testIdentifier(t, callExpression.Function, "add") {
		return
	}

	if len(callExpression.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(callExpression.Arguments))
	}
	testLiteralExpression(t, callExpression.Arguments[0], 1)
	testInfixExpression(t, callExpression.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExpression.Arguments[2], 4, "+", 5)
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4;-5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7*8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a+b+c*d/f+g)", "add((((a + b) + ((c * d) / f)) + g))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := statement.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("expression is not ast.HashLiteral, got=%T", statement.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length, got=%d", len(hash.Pairs))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral, got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{};"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := statement.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("expression is not ast.HashLiteral, got=%T", statement.Expression)
	}
	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length, got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralWithExpression(t *testing.T) {
	input := `{"one": 0+1, "two": 10-8, "three": 15/5};`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := statement.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("expression is not ast.HashLiteral, got=%T", statement.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length, got=%d", len(hash.Pairs))
	}
	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}
		testFunc(value)
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, expectedIdentifier string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral is not 'let', got=%q", statement.TokenLiteral())
		return false
	}
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement is not *ast.LetStatement, got=%T", statement)
		return false
	}
	if letStatement.Name.Value != expectedIdentifier {
		t.Errorf("letStatement.Name.Value is not %q, got=%q", expectedIdentifier, letStatement.Name.Value)
		return false
	}
	if letStatement.Name.TokenLiteral() != expectedIdentifier {
		t.Errorf("letStatement.Name.TokenLiteral() is not %q, got=%q", expectedIdentifier, letStatement.Name)
		return false
	}
	return true
}
func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)

	}
	t.Errorf("type of exp is not handled. got=%T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp is not *ast.Boolean, got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t, got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t, got=%s", value, bo.TokenLiteral())
		return false
	}
	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %s. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	identifier, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp is not *ast.Identifier, got=%T", exp)
		return false
	}
	if identifier.Value != value {
		t.Errorf("identifier.Value is not %q, got=%q", value, identifier.Value)
		return false
	}
	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral() is not %q, got=%q", value, identifier.TokenLiteral())
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not *ast.IntegerLiteral, got=%T", il)
		return false
	}
	if integer.Value != value {
		t.Errorf("integer.Value is not %d, got=%d", value, integer.Value)
		return false
	}
	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() is not %d, got=%s", value, integer.TokenLiteral())
		return false
	}
	return true
}
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
