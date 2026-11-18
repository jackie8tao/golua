# GoLua 改进路线图与技术方案

本文档基于对 `golua` 仓库当前实现（`cmd/`、`pkg/ast/`、`pkg/lexer/`、`pkg/parser/`、`pkg/vm/`、`pkg/libs/`、`tools/`、`tests/`）的代码审查与实测结果整理，将后续优化与演进拆分为四个可执行阶段。

---

## 背景与目标定位

当前仓库包含 Lua 5.0 手册（`docs/lua-5.0.pdf`）与测试集（`tests/lua50/`），词法与语法层面主要围绕 Lua 5.0 设计。在进入重构与功能补全前，建议先明确以下两点基线：

1. **目标 Lua 版本**：优先以 **Lua 5.0** 为第一阶段完整兼容目标（与现有 `tests/lua50/` 回归集对齐），待核心链路稳定后再决定是否向 Lua 5.1+（如 `#` 取长度、`%` 取模、十六进制数字、`[=[...]=]` 长括号）演进。
2. **验收标准**：以 `go test ./...` 单元测试 + `tests/lua50/*.lua` 端到端比对标准 Lua 输出作为最终验证基线。

---

## 阶段一：前端正确性修复与位置信息补全（P0）

目标：消除词法与语法分析器中的解析错误、死循环和 `panic`，为后续编译器提供准确可靠的 AST。

### 1.1 词法分析器修复（`pkg/lexer/`）

- **修复注释解析逻辑**（`pkg/lexer/lexer.go` `skipComments`）
  - 处理完 `--[[ ... ]]` 多行注释后应直接 `return nil`，避免继续落入单行注释分支而吞掉 `]]` 同一行之后的代码。
  - 单行注释循环中增加 `ch == runeEOF` 退出条件，防止文件末尾无换行符时的单行注释触发死循环。
  - 处理 `--[` 后紧跟非 `[` 字符的情况（如 `--[abc`），避免误消费换行符。
  - 根据目标版本完善嵌套 `[[...]]` 或长括号 `[=[...]=]` 支持。
- **统一运算符与 Token 映射**（`pkg/lexer/lexer.go`、`pkg/lexer/utils.go`）
  - `Scan()` 中扫描了 `%` 和 `#`，但 `operators` 映射表中未包含它们，导致返回零值 `TokenEOF`。需按目标版本决定移除或补全 `TokenMod`、`TokenLen` 及其映射。
- **修复字面量与位置记录**
  - `scanEscape()` 中的十进制转义 `\ddd` 应校验范围 `0..255` 并使用 `buf.WriteByte(byte(val))`，避免 `> 127` 时被 UTF-8 编码为多字节。
  - `scanDecimal()` 补充对 `.5` 等小数形式的支持，并禁止 `1.2e3.4` 这类指数后再次出现小数点。
  - 在 `Scan()` 开始匹配 Token 前记录起始行列号，避免当前在 `finally` 处记录结束位置导致报错偏移。

### 1.2 语法分析器修复（`pkg/parser/parser.go`）

- **引入运算符优先级与结合性**
  - 使用 **Pratt Parsing** 或**优先级爬升算法（Precedence Climbing）** 重写 `parseExpr()`，正确处理 Lua 运算符优先级表，以及 `^`、`..` 的右结合性。
  - 修复 `parsePrimaryExpr()` 中一元运算符（`-`、`not`）直接递归调用 `parseExpr()` 导致的结合错误（例如 `-1 + 2` 应为 `(-1) + 2` 而非 `-(1 + 2)`），确保一元运算符优先级低于 `^`、高于 `*` 和 `/`。
- **修复表构造器与语句解析缺陷**
  - 修复 `parseTableExpr()`：`TokenIdentifier` 分支目前强制要求后续为 `=`，导致 `{ x, y }`、`{ f() }` 等数组形式报错；需向前预读判断是否为 `name = expr`，否则回退按普通表达式解析，同时规范分隔符校验。
  - 修复 `parseAssignStmt()` 中 `validateExpr` 错误检查 `firstExpr.(type)` 的问题，改为校验传入参数 `expr.(type)`。
  - 修复 `parseLocalFuncDefStmt()` 和 `parseFuncDefStmt()` 遗漏赋值 `stmt.HasVariant = body.HasVariant` 的问题；修复 `parseFuncBody()` 中 `switch` 内 `break` 未能跳出外层 `for` 循环的问题（`...` 后不应继续接收形参）。
  - 修复 `parseFuncArgs()` 的 `default` 分支静默返回空参数列表的问题（如 `a:b` 缺失参数列表应报错）。
- **数值精度、错误处理与位置信息**
  - 将 `strconv.ParseFloat(p.curr.Str, 32)` 改为 `64` 位解析，保留 `float64` 精度。
  - 移除 `advance()` 中直接 `panic(err)` 的做法，将词法错误作为 `error` 正常向上返回。
  - 为 `pkg/ast` 中的 `Stmt` / `Expr` 节点补充 `ast.Position` 字段，为编译错误和运行时堆栈提供定位能力。

---

## 阶段二：重构虚拟机值系统与执行模型（P0）

目标：修复当前 VM 架构中的根本性运行时缺陷，跑通基础函数调用、嵌套调用、递归与表达式求值。

### 2.1 值系统与表语义重构（`pkg/vm/value.go`）

- **值语义化基础类型**
  - 将 `LNil`、`LBool`、`LNumber`、`LString` 改为不可变的值类型（如 `type LNumber float64`、`type LString string`），消除小对象频繁堆分配，并解决不同指针实例在 Go `map` 中无法作为同一 Key 命中的问题。
- **完善值方法与辅助函数**
  - 修复 `LNumber.String()` 固定使用 `%.2f` 的问题，改为符合 Lua 习惯的格式化方式（如 `strconv.FormatFloat(float64(l), 'g', 14, 64)`）。
  - 为 `LValue` 补充统一真值判断（Lua 中仅 `nil` 和 `false` 为假）与数值/字符串隐式转换辅助函数。
  - 为 `LTable` 封装 `RawGet`、`RawSet`、`Next`、`Len` 等方法，校验 `nil`/`NaN` 作为表键时的合法性。

### 2.2 修复调用帧与核心指令（`pkg/vm/function.go`、`pkg/vm/vm.go`）

- **迁移指令指针 `pc` 至 `CallFrame`**
  - 将 `Execute()` 及所有指令处理函数中的 `v.curFrame.fn.pc` 替换为 `v.curFrame.pc`，保证 `FuncProto` 为只读不可变字节码原型，修复递归调用、重复调用以及 `Disassemble()` 对 `FuncProto.pc` 的状态污染。
- **修复调用帧与栈平衡 Bug**
  - 修复 `opCall` 中 `newBp := v.sp - uint(argCount)` 遗漏加 `v.bp` 的错误（应基于当前绝对栈帧基准计算）。
  - 修复 `opReturn`：
    - 避免无返回值或根帧返回时无条件 `StPop()` 及 `v.curFrame.bp - 1` 触发 `uint` 下溢。
    - 当 `v.curFrame.parent == nil` 时正常终止 `Execute()` 循环，避免下一轮循环解引用空 `curFrame` 引发 panic。
- **修正算术操作数顺序与全局环境**
  - 修复 `opArithmetic`：栈顶为 `RHS`、次栈顶为 `LHS`，应按 `lhs <op> rhs` 计算（修正 `-`、`/`、`^` 方向相反的问题）。
  - 将全局变量表 `_G`（`Env *LTable`）从 `FuncProto.globals` 提升至 `Vm` / `LState` 实例级别，使顶层脚本、子函数与标准库共享同一全局环境。
  - 减少 `fnOpArg()` 每次取指令操作数时的切片堆分配，直接读取小端 `uint16`。

---

## 阶段三：补全编译器与字节码指令集（P1）

目标：完整打通 AST 到字节码的语义覆盖，支持控制流、表、闭包与元表等核心 Lua 特性。

### 3.1 字节码架构决策

当前为 8-bit 变长指令的栈式 VM，继续演进有两种路线可选：

| 维度 | 方案 A：完善现有栈式 VM | 方案 B：升级为 Lua 5.0 寄存器式 VM |
| :--- | :--- | :--- |
| **改动成本** | 低，沿用现有指令分发框架逐步扩展 | 中高，需重新设计 32 位定长指令编码与寄存器分配 |
| **实现复杂度** | 需小心维护求值栈与局部变量的栈平衡及多返回值对齐 | 局部与临时变量统一映射至寄存器窗口，尾调用与多返回值处理更清晰 |
| **性能与规范对齐** | 指令数较多、栈操作频繁 | 与 `docs/lua-5.0.pdf` 官方设计一致，执行效率更高 |

> **建议**：若优先追求快速跑通并降低学习门槛，可先采用 **方案 A**；若目标是完整还原 Lua 5.0 内部机制（尾调用、Upvalue 复用、高还原度反汇编），建议在阶段二稳定后转向 **方案 B**。

### 3.2 指令集与编译器语义补全（`pkg/vm/opcode.go`、`pkg/vm/compiler.go`）

- **补全字节码指令与分发**
  - 字面量与表：`OpLoadNil`、`OpLoadTrue`、`OpLoadFalse`、`OpNewTable`、`OpGetTable`、`OpSetTable`、`OpSetList`（数组部分批量初始化）。
  - 运算与比较：一元运算（`OpUnm`、`OpNot`）、字符串拼接（`OpConcat`）、关系比较（`OpEq`、`OpLt`、`OpLeq`）。
  - 控制流与逻辑：`OpJmp`、`OpJmpIfFalse`（支持 `and`/`or` 短路求值）、`OpForPrep`/`OpForLoop`、`OpTForLoop`。
  - 函数与闭包：`OpClosure`、`OpGetUpval`、`OpSetUpval`、`OpClose`、`OpTailCall`、`OpVararg`，并支持调用指令编码期望返回值个数（`nresults`）。
- **完善作用域、闭包与语句编译**
  - 实现块级词法作用域管理：启用 `enterBlock()` / `leaveBlock()`，允许内层作用域查找外层局部变量并支持同名遮蔽；块退出时回收局部槽位并闭合逃逸的 Upvalue。
  - 修复 `FuncDefStmt` 编译：不再编译期直接写入 `globals`，而是编译为 `OpClosure` 后通过指令存入全局或局部/表字段（支持 `function a.b()`、`function a:b()`、`local function f()`）。
  - 实现 `ReturnStmt`、`IfStmt`、`WhileStmt`、`RepeatStmt`、`NumericForStmt`、`GenericForStmt`、`BreakStmt`、`DoStmt`、`FuncExpr`、`IndexExpr`、`MethodCallExpr`、`TableExpr` 等全部剩余 AST 节点。
  - 处理多目标赋值与多返回值对齐（不足补 `nil`，多余丢弃；语句级函数调用清空返回值）。
  - 接入元表（Metatable）分发钩子（`__index`、`__newindex`、`__call`、算术/比较元方法等）。

---

## 阶段四：标准库、CLI 与自动化测试体系（P2）

目标：完善宿主交互与标准库能力，打通终端使用体验，建立自动化测试防线。

### 4.1 规范 Go/Lua 互调协议与标准库（`pkg/libs/`）

- **统一 `LGFunction` 栈协议**
  - 修正函数调用时函数对象与参数的栈布局，提供按 1-based 或 0-based 正序读取参数（`CheckNumber`、`CheckString` 等）及按序压入返回值的 API，替换目前 `basicPrint` 逆序 `StPop()` 的写法。
  - 修正 `print` 行为：使用 `\t` 分隔同一行参数、末尾换行，并通过 `tostring` 逻辑输出值。
- **分批补齐标准库模块**
  1. `basic.go`：`print`、`type`、`tostring`、`tonumber`、`pairs`、`ipairs`、`next`、`assert`、`error`、`pcall`、`xpcall`、`setmetatable`、`getmetatable`、`select`、`unpack`。
  2. `table.go` / `string.go` / `math.go`：先覆盖 `tests/lua50/` 依赖的核心函数（如 `table.insert`/`remove`/`sort`/`concat`、`string.format`/`find`/`sub`/`gsub`、`math.*`）。
  3. `os.go` / `io.go` / `load.go` / `coroutine.go` / `debug.go`：逐步补充文件 I/O、动态编译加载（`loadstring`/`dofile`）、协程（基于独立 `LState` / 栈帧）与调试钩子。

### 4.2 完善 CLI 与运行时错误报告（`cmd/main.go`）

- 打通 `lexer -> parser -> compiler -> vm` 主流程，支持：
  - `golua <script.lua> [args...]`：运行指定 Lua 脚本并注入 `arg` 表。
  - `golua -e "<code>"`：执行单行 Lua 代码片段。
  - 无参数启动交互式 **REPL**。
- 结合 AST 中的 `Position` 与 `CallFrame` 链，输出包含文件名、行号与调用栈（Stack Traceback）的错误信息。

### 4.3 测试体系与工程规范

- **Go 单元测试**：为 `pkg/lexer`、`pkg/parser`、`pkg/vm`、`pkg/libs` 补充单元测试，覆盖边界语法、优先级、作用域、闭包、异常路径与指令行为。
- **端到端集成测试**：编写自动化测试驱动 `tests/lua50/*.lua`，与预期输出或参考 Lua 5.0 行为进行回归比对。
- **仓库清理与文档**：添加 `.gitignore` 忽略 `.DS_Store` 及编译产物；完善 `README.md` 与 `README.zh-CN.md`，记录项目架构、构建方式、已支持特性与测试方法。
