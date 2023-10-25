// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"errors"
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/microsoft/yardl/tooling/internal/validation"
)

var (
	expressionLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "String", Pattern: `("([^\\"]|\\.)*")|('(\\'|[^'])*')`},
		{Name: "Float", Pattern: `(((\d+\.\d*)|(\d*\.\d+))(e[-+]?[0-9]+)?)|(\d+(e[-+]?[0-9]+))`},
		{Name: "Int", Pattern: `((0[xX][0-9A-Fa-f]+)|\d+)`},
		{Name: "As", Pattern: `as`},
		{Name: "Ident", Pattern: `[a-zA-Z_]\w*`},
		{Name: "OpenParen", Pattern: `\(`},
		{Name: "CloseParen", Pattern: `\)`},
		{Name: "Comma", Pattern: `,`},
		{Name: "OpenBracket", Pattern: `\[`},
		{Name: "CloseBracket", Pattern: `\]`},
		{Name: "Dot", Pattern: `\.`},
		{Name: "Pow", Pattern: `\*\*`},
		{Name: "Plus", Pattern: `\+`},
		{Name: "Minus", Pattern: `-`},
		{Name: "Star", Pattern: `\*`},
		{Name: "Slash", Pattern: `/`},
		{Name: "Colon", Pattern: `:`},
		{Name: "UnterminatedString", Pattern: `["']`},
		{Name: "OtherPunct", Pattern: `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
		{Name: "whitespace", Pattern: `[ \t]+`},
	})

	TokenTypeString             = expressionLexer.Symbols()["String"]
	TokenTypeFloat              = expressionLexer.Symbols()["Float"]
	TokenTypeInt                = expressionLexer.Symbols()["Int"]
	TokenTypeAs                 = expressionLexer.Symbols()["As"]
	TokenTypeIdent              = expressionLexer.Symbols()["Ident"]
	TokenTypeOpenParen          = expressionLexer.Symbols()["OpenParen"]
	TokenTypeCloseParen         = expressionLexer.Symbols()["CloseParen"]
	TokenTypeComma              = expressionLexer.Symbols()["Comma"]
	TokenTypeOpenBracket        = expressionLexer.Symbols()["OpenBracket"]
	TokenTypeCloseBracket       = expressionLexer.Symbols()["CloseBracket"]
	TokenTypeDot                = expressionLexer.Symbols()["Dot"]
	TokenTypePow                = expressionLexer.Symbols()["Pow"]
	TokenTypePlus               = expressionLexer.Symbols()["Plus"]
	TokenTypeMinus              = expressionLexer.Symbols()["Minus"]
	TokenTypeStar               = expressionLexer.Symbols()["Star"]
	TokenTypeSlash              = expressionLexer.Symbols()["Slash"]
	TokenTypeColon              = expressionLexer.Symbols()["Colon"]
	TokenTypeUnterminatedString = expressionLexer.Symbols()["UnterminatedString"]

	operatorInfo = map[lexer.TokenType]struct {
		IsRightAssociative bool
		Precedence         int
		IsBinary           bool
	}{
		TokenTypeOpenBracket: {Precedence: 5, IsBinary: false}, // for subscript
		TokenTypeOpenParen:   {Precedence: 5, IsBinary: false}, // for function call
		TokenTypeDot:         {Precedence: 5, IsBinary: true},
		TokenTypeAs:          {Precedence: 4, IsBinary: true},
		TokenTypePow:         {Precedence: 3, IsBinary: true, IsRightAssociative: true},
		TokenTypeStar:        {Precedence: 2, IsBinary: true},
		TokenTypeSlash:       {Precedence: 2, IsBinary: true},
		TokenTypePlus:        {Precedence: 1, IsBinary: true},
		TokenTypeMinus:       {Precedence: 1, IsBinary: true},
	}
)

var expressionParser = participle.MustBuild[Expression](
	participle.Lexer(expressionLexer),
	participle.Unquote("String"),
	participle.ParseTypeWith[Expression](parseExpr),
)

func ParseExpression(input string, lineOffset, columnOffet int) (Expression, error) {
	exp, err := expressionParser.ParseString("", input)
	if err != nil {
		if err, ok := err.(participle.Error); ok {
			line := lineOffset + err.Position().Line - 1
			column := columnOffet + err.Position().Column - 1
			return nil, validation.ValidationError{
				Message: errors.New(err.Message()),
				Line:    &line,
				Column:  &column,
			}
		}
		panic(fmt.Errorf("unexpected error type %T: %v", err, err))
	}

	if exp == nil {
		return nil, err
	}

	// adjust positions based on the host node
	Visit(*exp, func(self Visitor, node Node) {
		meta := node.GetNodeMeta()
		meta.Line += lineOffset - 1
		meta.Column += columnOffet - 1
		self.VisitChildren(node)
	})

	return *exp, err
}

func parseExpr(lex *lexer.PeekingLexer) (Expression, error) {
	return parseExprWithPrecedence(lex, 0)
}

func parseExprWithPrecedence(lex *lexer.PeekingLexer, minPrec int) (Expression, error) {
	lhs, err := parseAtom(lex)
	if err != nil {
		return nil, err
	}
	for {
		tok := lex.Peek()
		if tok.EOF() {
			break
		}

		opInfo, isExpectedOp := operatorInfo[tok.Type]
		if !isExpectedOp || opInfo.Precedence < minPrec {
			break
		}
		switch {
		case opInfo.IsBinary:
			nextMinPrec := opInfo.Precedence
			if !opInfo.IsRightAssociative {
				nextMinPrec++
			}
			lex.Next()
			var err error
			if err != nil {
				return nil, err
			}

			rhs, err := parseExprWithPrecedence(lex, nextMinPrec)
			if err != nil {
				return nil, err
			}

			lhs, err = combineOperands(lhs, tok, rhs)
			if err != nil {
				return nil, err
			}
		case tok.Type == TokenTypeOpenParen:
			lhs, err = parseCall(lex, lhs)
			if err != nil {
				return nil, err
			}
		case tok.Type == TokenTypeOpenBracket:
			if minPrec >= 5 {
				break
			}
			lhs, err = parseSubscript(lex, lhs)
			if err != nil {
				return nil, err
			}
		}
	}

	return lhs, nil
}

func combineOperands(lhs Expression, tok *lexer.Token, rhs Expression) (Expression, error) {
	switch tok.Type {
	case TokenTypeAs:
		ma, ok := rhs.(*MemberAccessExpression)
		if !ok || ma.Target != nil {
			return nil, &participle.ParseError{Pos: getLexerPosFromNodeMeta(*rhs.GetNodeMeta()), Msg: "The right-hand side of an 'as' operator must be a type reference"}
		}

		// for now, we know that this will be a simple type ([] and * tokens would have been rejected)
		return &TypeConversionExpression{
			NodeMeta: nodeMetaFromPosition(tok.Pos),
			Type: &SimpleType{
				NodeMeta: nodeMetaFromPosition(tok.Pos),
				Name:     ma.Member,
			},
			Expression: lhs,
		}, nil
	case TokenTypeDot:
		ma, ok := rhs.(*MemberAccessExpression)
		if !ok || ma.Target != nil {
			return nil, &participle.ParseError{Pos: getLexerPosFromNodeMeta(*rhs.GetNodeMeta()), Msg: "The right-hand side of a '.' operator must be an identifier"}
		}

		return &MemberAccessExpression{
			NodeMeta: nodeMetaFromPosition(tok.Pos),
			Target:   lhs,
			Member:   ma.Member,
		}, nil
	}

	binaryExpression := BinaryExpression{
		NodeMeta: nodeMetaFromPosition(tok.Pos),
		Left:     lhs,
		Right:    rhs,
	}

	switch tok.Type {
	case TokenTypePlus:
		binaryExpression.Operator = BinaryOpAdd
	case TokenTypeMinus:
		binaryExpression.Operator = BinaryOpSub
	case TokenTypeStar:
		binaryExpression.Operator = BinaryOpMul
	case TokenTypeSlash:
		binaryExpression.Operator = BinaryOpDiv
	case TokenTypePow:
		binaryExpression.Operator = BinaryOpPow
	default:
		panic(fmt.Sprintf("unexpected token type %v", tok.Type))
	}

	return &binaryExpression, nil
}

func parseAtom(lex *lexer.PeekingLexer) (Expression, error) {
	tok := lex.Next()
	if tok.EOF() {
		return nil, &participle.UnexpectedTokenError{Unexpected: *tok}
	}

	switch tok.Type {
	case TokenTypePlus:
		// discard unary plus
		return parseAtom(lex)
	case TokenTypeMinus:
		// unary minus
		expr, err := parseAtom(lex)
		if err != nil {
			return nil, err
		}

		switch expr := expr.(type) {
		case *IntegerLiteralExpression:
			expr.Value.Neg(&expr.Value)
			return expr, nil
		case *FloatingPointLiteralExpression:
			expr.Value = "-" + expr.Value
			return expr, nil
		default:
			return &UnaryExpression{
				NodeMeta:   nodeMetaFromPosition(tok.Pos),
				Operator:   UnaryOpNegate,
				Expression: expr,
			}, nil
		}
	case TokenTypeInt:
		i := IntegerLiteralExpression{
			NodeMeta: nodeMetaFromPosition(tok.Pos),
		}
		if err := i.Value.UnmarshalText([]byte(tok.Value)); err != nil {
			return nil, err
		}
		return &i, nil
	case TokenTypeFloat:
		return &FloatingPointLiteralExpression{
			NodeMeta: nodeMetaFromPosition(tok.Pos),
			Value:    tok.Value,
		}, nil
	case TokenTypeString:
		return &StringLiteralExpression{
			NodeMeta: nodeMetaFromPosition(tok.Pos),
			Value:    tok.Value,
		}, nil
	case TokenTypeUnterminatedString:
		return nil, &lexer.Error{Pos: tok.Pos, Msg: "unterminated string"}
	case TokenTypeOpenParen:
		// subexpression
		expr, err := parseExprWithPrecedence(lex, 0)
		if err != nil {
			return nil, err
		}
		tok = lex.Next()
		if tok.EOF() || tok.Type != TokenTypeCloseParen {
			return nil, &participle.UnexpectedTokenError{Unexpected: *tok, Expect: "closing parenthesis"}
		}
		return expr, nil
	case TokenTypeIdent:
		identifier := tok.Value
		return &MemberAccessExpression{
			NodeMeta: nodeMetaFromPosition(tok.Pos),
			Member:   identifier,
		}, nil
	}

	return nil, &participle.UnexpectedTokenError{Unexpected: *tok}
}

func parseCall(lex *lexer.PeekingLexer, target Expression) (*FunctionCallExpression, error) {
	if lex.Next().Type != TokenTypeOpenParen {
		panic("expected open paren")
	}

	ma, ok := target.(*MemberAccessExpression)
	if !ok || ma.Target != nil {
		return nil, &participle.ParseError{Pos: getLexerPosFromNodeMeta(*target.GetNodeMeta()), Msg: "The target of a function call must be an identifier, e.g. `size(...)`"}
	}

	args := make([]Expression, 0)
	for {
		if lex.Peek().Type == TokenTypeCloseParen {
			lex.Next()
			break
		}
		if len(args) > 0 {
			tok := lex.Next()
			if tok.EOF() || tok.Type != TokenTypeComma {
				return nil, &participle.UnexpectedTokenError{Unexpected: *tok, Expect: "a comma"}
			}
		}
		switch lex.Peek().Type {
		case TokenTypeCloseParen, TokenTypeComma:
			return nil, &participle.UnexpectedTokenError{Unexpected: *lex.Next(), Expect: "an expression"}
		}
		arg, err := parseExpr(lex)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return &FunctionCallExpression{
		NodeMeta:     *target.GetNodeMeta(),
		FunctionName: ma.Member,
		Arguments:    args,
	}, nil
}

func parseSubscript(lex *lexer.PeekingLexer, target Expression) (*SubscriptExpression, error) {
	if lex.Next().Type != TokenTypeOpenBracket {
		panic("expected open bracket")
	}

	args := make([]*SubscriptArgument, 0)
	for {
		if lex.Peek().Type == TokenTypeCloseBracket {
			lex.Next()
			break
		}
		if len(args) > 0 {
			tok := lex.Next()
			if tok.EOF() || tok.Type != TokenTypeComma {
				return nil, &participle.UnexpectedTokenError{Unexpected: *tok, Expect: "a comma"}
			}
		}
		switch lex.Peek().Type {
		case TokenTypeCloseBracket, TokenTypeComma:
			return nil, &participle.UnexpectedTokenError{Unexpected: *lex.Next(), Expect: "a subscript argument"}
		}
		arg, err := parseSubscriptArg(lex)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return &SubscriptExpression{
		NodeMeta:  *target.GetNodeMeta(),
		Target:    target,
		Arguments: args,
	}, nil
}

func parseSubscriptArg(lex *lexer.PeekingLexer) (*SubscriptArgument, error) {
	expr, err := parseExpr(lex)
	if err != nil {
		return nil, err
	}
	if lex.Peek().Type == TokenTypeColon {
		lex.Next()
		label := expr
		expr, err = parseExpr(lex)
		if err != nil {
			return nil, err
		}

		ma, ok := label.(*MemberAccessExpression)
		if !ok {
			return nil, &participle.ParseError{Pos: getLexerPosFromNodeMeta(*label.GetNodeMeta()), Msg: "expected label to be an identifier"}
		}
		return &SubscriptArgument{
			NodeMeta: *label.GetNodeMeta(),
			Label:    ma.Member,
			Value:    expr,
		}, nil
	}

	return &SubscriptArgument{
		NodeMeta: *expr.GetNodeMeta(),
		Value:    expr,
	}, nil
}

func nodeMetaFromPosition(pos lexer.Position) NodeMeta {
	return NodeMeta{}
}

func getLexerPosFromNodeMeta(meta NodeMeta) lexer.Position {
	return lexer.Position{Filename: meta.File, Line: meta.Line, Column: meta.Column}
}
