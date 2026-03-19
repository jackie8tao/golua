package vm

import (
	"fmt"

	"github.com/jackie8tao/golua/pkg/ast"
)

type LocalVar struct {
	Name  string
	Depth int
}

type Compiler struct {
	binOpTable map[ast.TokenType]uint8
	locals     []LocalVar
	depth      int
}

func NewCompiler() *Compiler {
	return &Compiler{
		locals: make([]LocalVar, 0),
		depth:  0,
		binOpTable: map[ast.TokenType]uint8{
			ast.TokenAdd: OpAdd,
			ast.TokenSub: OpSub,
			ast.TokenMul: OpMul,
			ast.TokenDiv: OpDiv,
			ast.TokenPow: OpPow,
		},
	}
}

func (c *Compiler) Compile(root ast.Node) (*FuncProto, error) {
	fp := NewFuncProto()
	err := c.visit(fp, root)
	if err != nil {
		return nil, err
	}
	return fp, nil
}

func (c *Compiler) visit(fp *FuncProto, node ast.Node) error {
	switch n := node.(type) {
	case *ast.Chunk:
		for i := 0; i < len(n.Stmts); i++ {
			err := c.visit(fp, n.Stmts[i])
			if err != nil {
				return err
			}
		}
	case *ast.LocalAssignStmt:
		for i := 0; i < len(n.Exprs); i++ {
			if err := c.visit(fp, n.Exprs[i]); err != nil {
				return err
			}
		}
		for i := len(n.Names) - 1; i >= 0; i-- {
			idx := c.setLocal(n.Names[i])
			if len(c.locals) >= fp.maxLocals {
				fp.maxLocals = len(c.locals)
			}
			c.writeCode(fp, OpSetLocal, convToUint8(idx)...)
		}
	case *ast.AssignStmt:
		for i := 0; i < len(n.RHS); i++ {
			if err := c.visit(fp, n.RHS[i]); err != nil {
				return err
			}
		}
		for i := len(n.LHS) - 1; i >= 0; i-- {
			switch l := n.LHS[i].(type) {
			case *ast.IdentExpr:
				if idx, found := c.getLocal(l.Name); found {
					c.writeCode(fp, OpSetLocal, convToUint8(idx)...)
				} else {
					idx := convToUint8(c.setConstant(fp, &LString{Value: l.Name}))
					c.writeCode(fp, OpSetGlobal, idx...)
				}
			default:
				return fmt.Errorf("invalid assignment target")
			}
		}
	case *ast.BinOpExpr:
		if err := c.visit(fp, n.LHS); err != nil {
			return err
		}
		if err := c.visit(fp, n.RHS); err != nil {
			return err
		}
		if code, ok := c.binOpTable[n.Op]; !ok {
			return fmt.Errorf("invalid binary operator")
		} else {
			c.writeCode(fp, code)
		}
	case *ast.NumExpr:
		idx := convToUint8(c.setConstant(fp, &LNumber{Value: n.Value}))
		c.writeCode(fp, OpConstant, idx[0], idx[1])
	case *ast.IdentExpr:
		if idx, found := c.getLocal(n.Name); found {
			c.writeCode(fp, OpGetLocal, convToUint8(idx)...)
		} else {
			idx := c.setConstant(fp, &LString{Value: n.Name})
			c.writeCode(fp, OpGetGlobal, convToUint8(idx)...)
		}
	default:
		return fmt.Errorf("invalid ast node")
	}
	return nil
}

func (c *Compiler) setLocal(name string) uint16 {
	c.locals = append(c.locals, LocalVar{
		Name:  name,
		Depth: c.depth,
	})
	return uint16(len(c.locals) - 1)
}

func (c *Compiler) getLocal(name string) (uint16, bool) {
	for i := len(c.locals) - 1; i >= 0; i-- {
		if c.locals[i].Depth != c.depth {
			continue
		}
		if c.locals[i].Name == name {
			return uint16(i), true
		}
	}
	return 0, false
}

// used for compiler to write code to function instruction list
func (c *Compiler) writeCode(fn *FuncProto, code uint8, operands ...uint8) {
	fn.codes = append(fn.codes, code)
	for _, op := range operands {
		fn.codes = append(fn.codes, op)
	}
}

// used for compiler to write constant
func (c *Compiler) setConstant(fn *FuncProto, val LValue) uint16 {
	for i := 0; i < len(fn.constants); i++ {
		if fn.constants[i].Type() != val.Type() {
			continue
		}
		switch l := fn.constants[i].(type) {
		case *LNumber:
			if l.Value == val.(*LNumber).Value {
				return uint16(i)
			}
		case *LString:
			if l.Value == val.(*LString).Value {
				return uint16(i)
			}
		default:
		}
	}
	fn.constants = append(fn.constants, val)
	return uint16(len(fn.constants) - 1)
}

func (c *Compiler) Disassemble(fp *FuncProto) error {
	for {
		if fp.pc >= len(fp.codes) {
			break
		}
		switch fp.codes[fp.pc] {
		case OpAdd, OpSub, OpMul, OpDiv, OpPow:
			fmt.Printf("%s\n", opCodeNames[fp.codes[fp.pc]])
			fp.pc++
		case OpConstant, OpGetGlobal, OpSetGlobal, OpGetLocal, OpSetLocal:
			name := opCodeNames[fp.codes[fp.pc]]
			args := make([]uint8, 0)
			for i := 0; i < 2; i++ {
				fp.pc++
				args = append(args, fp.codes[fp.pc])
			}
			idx := convToUint16([2]uint8{args[0], args[1]})
			fmt.Printf("%s %d\n", name, idx)
			fp.pc++
		default:
			return fmt.Errorf("invalid opcode")
		}
	}

	return nil
}
