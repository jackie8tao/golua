package parser

import (
	"fmt"
	"strconv"

	"github.com/jackie8tao/golua/pkg/ast"
	"github.com/jackie8tao/golua/pkg/lexer"
)

type Parser struct {
	curr  ast.Token
	next  ast.Token
	lexer *lexer.Lexer
}

func NewParser(lexer *lexer.Lexer) *Parser {
	return &Parser{lexer: lexer}
}

func (p *Parser) advance() {
	p.curr = p.next
	tk, err := p.lexer.Scan()
	if err != nil {
		panic(err)
	}
	p.next = tk
}

func (p *Parser) Parse() (ast.Node, error) {
	p.advance()
	n, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	return n, nil
}

// block -> chunk
func (p *Parser) parseBlock() (*ast.Block, error) {
	chunk, err := p.parseChunk()
	if err != nil {
		return nil, err
	}
	return &ast.Block{Chunk: chunk}, nil
}

// chunk -> { stmt [';'] }
func (p *Parser) parseChunk() (*ast.Chunk, error) {
	chunk := &ast.Chunk{
		Stmts: make([]ast.Stmt, 0),
	}
	for {
		switch p.next.Type {
		case ast.TokenEOF, ast.TokenEnd, ast.TokenElse, ast.TokenElseif,
			ast.TokenUntil:
			goto exit
		default:
		}

		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		chunk.Stmts = append(chunk.Stmts, stmt)
		switch p.next.Type {
		case ast.TokenSemicolon:
			p.advance() // eat ';'
		default:
		}
	}
exit:
	return chunk, nil
}

// stmt -> varlist '=' explist
func (p *Parser) parseStmt() (ast.Stmt, error) {
	switch p.next.Type {
	case ast.TokenLocal:
		p.advance() // eat 'local'
		if p.next.Type == ast.TokenFunction {
			return p.parseLocalFuncDefStmt()
		}
		return p.parseLocalAssignStmt()
	case ast.TokenBreak:
		p.advance() // eat 'break'
		return &ast.BreakStmt{}, nil
	case ast.TokenFunction:
		return p.parseFuncDefStmt()
	case ast.TokenDo:
		return p.parseDoStmt()
	case ast.TokenWhile:
		return p.parseWhileStmt()
	case ast.TokenRepeat:
		return p.parseRepeatStmt()
	case ast.TokenIf:
		return p.parseIfStmt()
	case ast.TokenReturn:
		return p.parseReturnStmt()
	case ast.TokenFor:
		return p.parseForStmt()
	default:
		prefixExpr, err := p.parsePrefixExpr()
		if err != nil {
			return nil, err
		}
		if p.next.Type == ast.TokenAssign || p.next.Type == ast.TokenComma {
			return p.parseAssignStmt(prefixExpr)
		}
		switch prefixExpr.(type) {
		case *ast.FuncCallExpr, *ast.MethodCallExpr:
			return &ast.FuncCallStmt{
				Expr: prefixExpr,
			}, nil
		default:
			return nil, fmt.Errorf("invalid statement")
		}
	}
}

func (p *Parser) parseExpr() (ast.Expr, error) {
	left, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}

	for {
		switch p.next.Type {
		case ast.TokenAdd, ast.TokenSub, ast.TokenMul, ast.TokenDiv,
			ast.TokenPow, ast.TokenDotDot, ast.TokenLt, ast.TokenLeq,
			ast.TokenGt, ast.TokenGeq, ast.TokenEq, ast.TokenNeq,
			ast.TokenAnd, ast.TokenOr:
			p.advance() // eat operator
			op := p.curr.Type
			right, err := p.parsePrimaryExpr()
			if err != nil {
				return nil, err
			}
			left = &ast.BinOpExpr{
				Op:  op,
				LHS: left,
				RHS: right,
			}
			continue
		default:
		}
		break
	}

	return left, nil
}

func (p *Parser) parsePrimaryExpr() (ast.Expr, error) {
	switch p.next.Type {
	case ast.TokenNil:
		p.advance() // eat nil
		return &ast.NilExpr{}, nil
	case ast.TokenFalse:
		p.advance() // eat false
		return &ast.BoolExpr{
			Value: false,
		}, nil
	case ast.TokenTrue:
		p.advance() // eat true
		return &ast.BoolExpr{
			Value: true,
		}, nil
	case ast.TokenNumber:
		p.advance() // eat number
		val, err := strconv.ParseFloat(p.curr.Str, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", p.curr.Str)
		}
		return &ast.NumExpr{
			Value: val,
		}, nil
	case ast.TokenString:
		p.advance() // eat string
		return &ast.StrExpr{
			Value: p.curr.Str,
		}, nil
	case ast.TokenFunction:
		p.advance() // eat 'function'
		funcBody, err := p.parseFuncBody()
		if err != nil {
			return nil, err
		}
		return &ast.FuncExpr{
			Params:     funcBody.Params,
			HasVariant: funcBody.HasVariant,
			Block:      funcBody.Block,
		}, nil
	case ast.TokenSub, ast.TokenNot:
		p.advance() // eat '-' or 'not'
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryOpExpr{
			Op:  p.curr.Type,
			RHS: expr,
		}, nil
	case ast.TokenLbrace:
		expr, err := p.parseTableExpr()
		if err != nil {
			return nil, err
		}
		return expr, nil
	case ast.TokenIdentifier, ast.TokenLparen:
		return p.parsePrefixExpr()
	default:
		return nil, fmt.Errorf("invalid expression")
	}
}

// primary -> Name | '(' exp ')'
// suffix -> '[' exp ']' | '.' Name | args | ':' Name args
// prefixexp -> primary { suffix }
func (p *Parser) parsePrefixExpr() (ast.Expr, error) {
	var left ast.Expr
	switch p.next.Type {
	case ast.TokenLparen:
		p.advance() // eat '('
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.next.Type != ast.TokenRparen {
			return nil, fmt.Errorf("expected ')' after expression")
		}
		p.advance() // eat ')'
		left = expr
	case ast.TokenIdentifier:
		p.advance() // eat identifier
		left = &ast.IdentExpr{
			Name: p.curr.Str,
		}
	default:
		return nil, fmt.Errorf("expected '(' or identifier")
	}

	for {
		switch p.next.Type {
		case ast.TokenLbracket: // '['
			p.advance() // eat '['
			//if p.next.Type != ast.TokenIdentifier {
			//	return nil, fmt.Errorf("expected identifier after left bracket")
			//}
			index, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if p.next.Type != ast.TokenRbracket {
				return nil, fmt.Errorf("expected rbracket after left bracket")
			}
			left = &ast.IndexExpr{
				Table: left,
				Index: index,
			}
			p.advance() // eat ']'
		case ast.TokenDot: // '.'
			p.advance() // eat '.'
			if p.next.Type != ast.TokenIdentifier {
				return nil, fmt.Errorf("expected identifier after dot")
			}
			p.advance() // eat identifier
			left = &ast.IndexExpr{
				Table: left,
				Index: &ast.IdentExpr{
					Name: p.curr.Str,
				},
			}
		case ast.TokenColon: // ':'
			p.advance() // eat ':'
			if p.next.Type != ast.TokenIdentifier {
				return nil, fmt.Errorf("expected identifier after colon")
			}
			p.advance() // eat identifier
			method := p.curr.Str

			args, err := p.parseFuncArgs()
			if err != nil {
				return nil, err
			}
			left = &ast.MethodCallExpr{
				Receiver: left,
				Method:   method,
				Args:     args,
			}
		case ast.TokenLparen, ast.TokenLbrace, ast.TokenString:
			args, err := p.parseFuncArgs()
			if err != nil {
				return nil, err
			}
			left = &ast.FuncCallExpr{
				Func: left,
				Args: args,
			}
		default: // empty suffix, return
			return left, nil
		}
	}
}

// args -> '(' explist ')' | tableconstructor | Literal
// tableconstructor -> '{' fieldlist '}'
func (p *Parser) parseFuncArgs() ([]ast.Expr, error) {
	switch p.next.Type {
	case ast.TokenLparen: // '('
		p.advance() // eat '('

		// empty args
		if p.next.Type == ast.TokenRparen {
			p.advance() // eat ')'
			return []ast.Expr{}, nil
		}

		exprs, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		if p.next.Type != ast.TokenRparen {
			return nil, fmt.Errorf("expected ')' after function arguments")
		}
		p.advance() // eat ')'
		return exprs, nil
	case ast.TokenLbrace: // '{'
		expr, err := p.parseTableExpr()
		if err != nil {
			return nil, err
		}
		return []ast.Expr{expr}, nil
	case ast.TokenString:
		p.advance() // eat string
		return []ast.Expr{
			&ast.StrExpr{Value: p.curr.Str},
		}, nil
	default:
		return []ast.Expr{}, nil
	}
}

// tableconstructor -> '{' fieldlist '}'
// fieldlist -> field {fieldsep field} [fieldsep]
// fieldsep -> ',' | ';'
// field -> '[' exp ']' '=' exp | name '=' exp
func (p *Parser) parseTableExpr() (ast.Expr, error) {
	if p.next.Type != ast.TokenLbrace {
		return nil, fmt.Errorf("expected '{' for table constructor")
	}
	p.advance() // eat '{'

	fields := make([]ast.TableField, 0)
	hasFieldSep := false
	for {
		switch p.next.Type {
		case ast.TokenLbracket: // '['
			p.advance() // eat '['
			key, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if p.next.Type != ast.TokenRbracket {
				return nil, fmt.Errorf("expected ']' after table key")
			}
			p.advance() // eat ']'
			if p.next.Type != ast.TokenAssign {
				return nil, fmt.Errorf("expected '=' after table key")
			}
			p.advance() // eat '='
			val, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			fields = append(fields, ast.TableField{
				Key:   key,
				Value: val,
			})
			hasFieldSep = false
		case ast.TokenIdentifier:
			p.advance() // eat identifier
			key := &ast.IdentExpr{
				Name: p.curr.Str,
			}
			if p.next.Type != ast.TokenAssign {
				return nil, fmt.Errorf("expected '=' after table key")
			}
			p.advance() // eat '='
			val, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			fields = append(fields, ast.TableField{
				Key:   key,
				Value: val,
			})
			hasFieldSep = false
		case ast.TokenComma, ast.TokenSemicolon: // ',' | ';'
			if hasFieldSep {
				return nil, fmt.Errorf("unexpected field separator")
			}
			p.advance() // eat ',' or ';'
			hasFieldSep = true
		case ast.TokenRbrace: // next token is '}', we should stop parsing fieldlist
			goto exit
		default:
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			fields = append(fields, ast.TableField{
				Key:   &ast.NilExpr{},
				Value: expr,
			})
			hasFieldSep = false
		}
	}
exit:
	if p.next.Type != ast.TokenRbrace {
		return nil, fmt.Errorf("expected '}' after table constructor")
	}
	p.advance() // eat '}'
	return &ast.TableExpr{
		Fields: fields,
	}, nil
}

// stmt -> 'do' block 'end'
func (p *Parser) parseDoStmt() (ast.Stmt, error) {
	if p.next.Type != ast.TokenDo {
		return nil, fmt.Errorf("expected 'do' keyword")
	}

	p.advance() // eat 'do'
	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	if p.next.Type != ast.TokenEnd {
		return nil, fmt.Errorf("expected 'end' after do block")
	}
	p.advance() // eat 'end'

	return &ast.DoStmt{
		Block: block,
	}, nil
}

// stmt -> 'while' exp 'do' block 'end'
func (p *Parser) parseWhileStmt() (ast.Stmt, error) {
	if p.next.Type != ast.TokenWhile {
		return nil, fmt.Errorf("expected 'while' keyword")
	}
	p.advance() // eat 'while'
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.next.Type != ast.TokenDo {
		return nil, fmt.Errorf("expected 'do' keyword after while condition")
	}
	p.advance() // eat 'do'
	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	if p.next.Type != ast.TokenEnd {
		return nil, fmt.Errorf("expected 'end' after do block")
	}
	p.advance() // eat 'end'
	return &ast.WhileStmt{
		Cond:  expr,
		Block: block,
	}, nil
}

// stmt -> 'repeat' block 'until' exp
func (p *Parser) parseRepeatStmt() (ast.Stmt, error) {
	if p.next.Type != ast.TokenRepeat {
		return nil, fmt.Errorf("expected 'repeat' keyword")
	}
	p.advance() // eat 'repeat'
	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	if p.next.Type != ast.TokenUntil {
		return nil, fmt.Errorf("expected 'until' keyword after repeat block")
	}
	p.advance() // eat 'until'
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &ast.RepeatStmt{
		Cond:  expr,
		Block: block,
	}, nil
}

// stmt -> 'if' exp 'then' block { elseif exp 'then' block } ['else' block]
func (p *Parser) parseIfStmt() (ast.Stmt, error) {
	if p.next.Type != ast.TokenIf {
		return nil, fmt.Errorf("expected 'if' keyword")
	}
	p.advance() // eat 'if'
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.next.Type != ast.TokenThen {
		return nil, fmt.Errorf("expected 'then' after if condition")
	}
	p.advance() // eat 'then'
	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	ifStmt := &ast.IfStmt{
		Cond:   expr,
		Then:   block,
		ElseIf: make([]*ast.ElseIfSeg, 0),
	}

	// optional 'elseif'
	for {
		if p.next.Type != ast.TokenElseif {
			break
		}
		p.advance() // eat 'elseif'
		elseIfExpr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.next.Type != ast.TokenThen {
			return nil, fmt.Errorf("expected 'then' after elseif condition")
		}
		p.advance() // eat 'then'
		elseIfBlock, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		ifStmt.ElseIf = append(ifStmt.ElseIf, &ast.ElseIfSeg{
			Cond: elseIfExpr,
			Then: elseIfBlock,
		})
	}

	// optional 'else'
	if p.next.Type == ast.TokenElse {
		p.advance() // eat 'else'
		elseBlock, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		ifStmt.Else = elseBlock
	}
	if p.next.Type != ast.TokenEnd {
		return nil, fmt.Errorf("expected 'end' after if statement")
	}
	p.advance() // eat 'end'
	return ifStmt, nil
}

// stmt -> 'return' [explist]
func (p *Parser) parseReturnStmt() (ast.Stmt, error) {
	if p.next.Type != ast.TokenReturn {
		return nil, fmt.Errorf("expected 'return' keyword")
	}
	p.advance() // eat 'return'
	switch p.next.Type {
	case ast.TokenEOF, ast.TokenSemicolon, ast.TokenEnd, ast.TokenElse,
		ast.TokenElseif, ast.TokenUntil: // empty return explist
		return &ast.ReturnStmt{
			Exprs: make([]ast.Expr, 0),
		}, nil
	default:
	}

	exprs, err := p.parseExprList()
	if err != nil {
		return nil, err
	}
	return &ast.ReturnStmt{
		Exprs: exprs,
	}, nil
}

// explist -> { exp ',' } exp
func (p *Parser) parseExprList() ([]ast.Expr, error) {
	exprs := make([]ast.Expr, 0)
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
		if p.next.Type != ast.TokenComma {
			break
		}
		p.advance() // eat ','
	}
	return exprs, nil
}

// stmt -> 'for' Name '=' exp ',' exp [',' exp] 'do' block 'end'
// stmt -> 'for' Name {',' Name} 'in' explist 'do' block 'end'
func (p *Parser) parseForStmt() (ast.Stmt, error) {
	if p.next.Type != ast.TokenFor {
		return nil, fmt.Errorf("expected 'for' keyword")
	}
	p.advance() // eat 'for'
	if p.next.Type != ast.TokenIdentifier {
		return nil, fmt.Errorf("expected variable name after 'for'")
	}
	p.advance() // eat variable name
	varName := p.curr.Str

	if p.next.Type == ast.TokenAssign {
		p.advance() // eat '='
		varExpr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.next.Type != ast.TokenComma {
			return nil, fmt.Errorf("expected ',' after for variable")
		}
		p.advance() // eat ','
		condExpr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt := &ast.NumericForStmt{
			VarName:  varName,
			VarExpr:  varExpr,
			CondExpr: condExpr,
		}

		// optional step expr
		if p.next.Type == ast.TokenComma {
			p.advance() // eat ','
			stepExpr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			stmt.StepExpr = stepExpr
		}
		if p.next.Type != ast.TokenDo {
			return nil, fmt.Errorf("expected 'do' keyword after for loop")
		}
		p.advance() // eat 'do'
		block, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		if p.next.Type != ast.TokenEnd {
			return nil, fmt.Errorf("expected 'end' after for loop")
		}
		p.advance() // eat 'end'
		stmt.Block = block
		return stmt, nil
	}

	stmt := &ast.GenericForStmt{
		Names: []string{varName},
		Exprs: make([]ast.Expr, 0),
	}
	for {
		if p.next.Type != ast.TokenComma {
			break
		}
		p.advance() // eat ','
		if p.next.Type != ast.TokenIdentifier {
			return nil, fmt.Errorf("expected variable name after ','")
		}
		p.advance() // eat variable name
		stmt.Names = append(stmt.Names, p.curr.Str)
	}
	if p.next.Type != ast.TokenIn {
		return nil, fmt.Errorf("expected 'in' keyword after for loop")
	}
	p.advance() // eat 'in'
	exprs, err := p.parseExprList()
	if err != nil {
		return nil, err
	}
	stmt.Exprs = exprs

	if p.next.Type != ast.TokenDo {
		return nil, fmt.Errorf("expected 'do' keyword after for loop")
	}
	p.advance() // eat 'do'
	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	if p.next.Type != ast.TokenEnd {
		return nil, fmt.Errorf("expected 'end' after for loop")
	}
	p.advance() // eat 'end'
	stmt.Block = block
	return stmt, nil
}

// local keyword has eaten
// stmt -> 'local' 'function' name funcbody
func (p *Parser) parseLocalFuncDefStmt() (ast.Stmt, error) {
	if p.next.Type != ast.TokenFunction {
		return nil, fmt.Errorf("expected 'function' keyword")
	}
	p.advance() // eat 'function'
	if p.next.Type != ast.TokenIdentifier {
		return nil, fmt.Errorf("expected function name after 'function'")
	}
	p.advance() // eat function name
	stmt := &ast.FuncDefStmt{
		IsLocal: true,
		Names:   []string{p.curr.Str},
	}
	if p.next.Type != ast.TokenLparen {
		return nil, fmt.Errorf("expected '(' after function name")
	}
	body, err := p.parseFuncBody()
	if err != nil {
		return nil, fmt.Errorf("failed to parse function body: %w", err)
	}
	//if p.next.Type != ast.TokenEnd {
	//	return nil, fmt.Errorf("expected 'end' after local function definition")
	//}
	//p.advance() // eat 'end'
	stmt.Params = body.Params
	stmt.Block = body.Block
	return stmt, nil
}

// stmt -> 'function' funcname funcbody
// funcname -> name { '.' name } [ ':' name ]
func (p *Parser) parseFuncDefStmt() (ast.Stmt, error) {
	if p.next.Type != ast.TokenFunction {
		return nil, fmt.Errorf("expected 'function' keyword")
	}
	p.advance() // eat 'function'
	if p.next.Type != ast.TokenIdentifier {
		return nil, fmt.Errorf("expected function name after 'function'")
	}
	p.advance() // eat function name
	stmt := &ast.FuncDefStmt{
		Names:      []string{p.curr.Str},
		SuffixName: "",
		Params:     nil,
		HasVariant: false,
		Block:      nil,
	}
	for {
		if p.next.Type != ast.TokenDot {
			break
		}
		p.advance() // eat '.'
		if p.next.Type != ast.TokenIdentifier {
			return nil, fmt.Errorf("expected function name after '.'")
		}
		p.advance() // eat function name
		stmt.Names = append(stmt.Names, p.curr.Str)
	}
	if p.next.Type == ast.TokenColon {
		p.advance() // eat ':'
		if p.next.Type != ast.TokenIdentifier {
			return nil, fmt.Errorf("expected suffix name after ':'")
		}
		p.advance() // eat suffix name
		stmt.SuffixName = p.curr.Str
	}
	body, err := p.parseFuncBody()
	if err != nil {
		return nil, err
	}
	//if p.next.Type != ast.TokenEnd {
	//	return nil, fmt.Errorf("expected 'end' after function definition")
	//}
	//p.advance() // eat 'end'
	stmt.Params = body.Params
	stmt.Block = body.Block
	return stmt, nil
}

// this function only used for parseLocalFuncDefStmt and parseFuncDefStmt
// funcbody -> '(' [parlist] ')' block 'end'
// parlist -> name {',' name} [',' '...'] | '...'
func (p *Parser) parseFuncBody() (*ast.FuncDefStmt, error) {
	if p.next.Type != ast.TokenLparen {
		return nil, fmt.Errorf("expected '(' for function body")
	}
	p.advance() // eat '('
	stmt := &ast.FuncDefStmt{
		Params: make([]string, 0),
	}
	// optional parlist
	if p.next.Type == ast.TokenIdentifier {
		p.advance() // eat first parameter name
		stmt.Params = append(stmt.Params, p.curr.Str)
		for {
			if p.next.Type != ast.TokenComma {
				break
			}
			p.advance() // eat ','
			switch p.next.Type {
			case ast.TokenDots:
				p.advance() // eat ...
				stmt.HasVariant = true
				break
			case ast.TokenIdentifier:
				p.advance() // eat parameter name
				stmt.Params = append(stmt.Params, p.curr.Str)
			default:
				return nil, fmt.Errorf("expected parameter name or dots after ','")
			}
		}
	} else if p.next.Type == ast.TokenDots {
		p.advance() // eat ...
		stmt.HasVariant = true
	}
	if p.next.Type != ast.TokenRparen {
		return nil, fmt.Errorf("expected ')' for function body")
	}
	p.advance() // eat ')'
	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	stmt.Block = block
	if p.next.Type != ast.TokenEnd {
		return nil, fmt.Errorf("expected 'end' after function body")
	}
	p.advance() // eat 'end'

	return stmt, nil
}

// local keyword has eaten
// stmt -> 'local' namelist [init]
// namelist -> name {',' name}
// init -> '=' explist
func (p *Parser) parseLocalAssignStmt() (ast.Stmt, error) {
	names := make([]string, 0)
	if p.next.Type != ast.TokenIdentifier {
		return nil, fmt.Errorf("expected variable name after 'local'")
	}
	p.advance() // eat variable name
	names = append(names, p.curr.Str)

	// optional names
	for {
		if p.next.Type != ast.TokenComma {
			break
		}
		p.advance() // eat ','
		if p.next.Type != ast.TokenIdentifier {
			return nil, fmt.Errorf("expected variable name after ','")
		}
		p.advance() // eat variable name
		names = append(names, p.curr.Str)
	}

	stmt := &ast.LocalAssignStmt{
		Names: names,
	}

	// optional init
	if p.next.Type == ast.TokenAssign {
		p.advance() // eat '='
		exprs, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		stmt.Exprs = exprs
	}

	return stmt, nil
}

func (p *Parser) parseAssignStmt(firstExpr ast.Expr) (ast.Stmt, error) {
	validateExpr := func(expr ast.Expr) error {
		switch firstExpr.(type) {
		case *ast.IndexExpr, *ast.IdentExpr:
			return nil
		default:
			return fmt.Errorf("expected identifier or index for assignment")
		}
	}
	err := validateExpr(firstExpr)
	if err != nil {
		return nil, err
	}

	lhs := []ast.Expr{firstExpr}
	for {
		if p.next.Type != ast.TokenComma {
			break
		}
		p.advance() // eat ','
		nextVar, err := p.parsePrefixExpr()
		if err != nil {
			return nil, err
		}
		err = validateExpr(nextVar)
		if err != nil {
			return nil, err
		}
		lhs = append(lhs, nextVar)
	}
	if p.next.Type != ast.TokenAssign {
		return nil, fmt.Errorf("expected '=' after assignment")
	}
	p.advance() // eat '='
	rhs, err := p.parseExprList()
	if err != nil {
		return nil, err
	}
	return &ast.AssignStmt{
		LHS: lhs,
		RHS: rhs,
	}, nil
}
