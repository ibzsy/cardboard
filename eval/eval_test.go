package eval

import (
	"testing"

	"github.com/ibzsy/cardboard/lexer"
	"github.com/ibzsy/cardboard/object"
	"github.com/ibzsy/cardboard/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10}}
	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
func TestEvalComplexIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"+5;", 5},
		{"-5;", -5},
		{"-50 + 100 + -50", 0},
		{"20 - 5;", 15},
		{"-15 - 100;", -115},
		{"5 + (5 - 10);", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalUnboxStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"unbox 10 + 2; 100;", 12},
		{"5 + (5 - 10); unbox (90 + 10); -15 - 100;", 100},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalPutStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"put a = 5; a;", 5},
		{"put a = 5 + 5; a;", 10},
		{"put a = 5; put b = a; b;", 5},
		{"put a = 5; put b = a; put c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionExprEval(t *testing.T) {
	test := "box(a,b){return a + b;}"
	evaluated := testEval(test, t)
	box, ok := evaluated.(*object.Box)

	if !ok {
		t.Fatalf("Test failed. Expected return type of <Box>. Got <%s>", evaluated.Type())
	}

	if len(box.ParameterList) != 2 {
		t.Fatalf("Test failed. Expected parameter length of 2. Got length of <%d>", len(box.ParameterList))
	}

	if box.ParameterList[0].String() != "a" || box.ParameterList[1].String() != "b" {
		t.Fatalf("Test failed. Expected parameter list of (a, b). Got <(%s, %s)>",
			box.ParameterList[0].String(),
			box.ParameterList[1].String())
	}

	if box.Body.String() != "{return(a+b)}" {
		t.Fatalf("Test failed. Expected function body of '{return(a+b)}'. Got <%s>", box.Body.String())
	}

}

func TestFullFunctionEval(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"put identity = box(x) { x; }; identity(5);", 5},
		{"put identity = box(x) { unbox x; }; identity(5);", 5},
		{"put double = box(x) { x + x; }; double(5);", 10},
		{"put add = box(x, y) { x + y; }; add(5, 5);", 10},
		{"put add = box(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"box(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input, t), tt.expected)
	}
}
func TestClosures(t *testing.T) {
	input := `
   put newAdder = box(x) {
     box(y) { x + y };
};
   put addTwo = newAdder(2);
   addTwo(2);`

	testIntegerObject(t, testEval(input, t), 4)
}

func TestFunctionsAsArguments(t *testing.T) {
	input := `
	put add = box(a, b) { a + b };
	put sub = box(a, b) { a - b };
	put applyFunc = box(a, b, func) { func(a, b) };
	applyFunc(2, 2, add);
	`
	testIntegerObject(t, testEval(input, t), 4)
}

func testEval(input string, t *testing.T) object.Object {
	l := lexer.CreateLexer(input)
	p := parser.CreateParser(l)
	program := p.ParseCardBoard()
	checkParserErrors(t, p)
	env := object.CreateEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("object is not Integer. got=%T", obj)
		return false
	}
	if result.Value != expected {
		t.Fatalf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errs := p.GetErrors()
	if len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
		t.FailNow()
	}
}
