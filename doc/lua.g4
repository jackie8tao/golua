grammar lua;

chunk : block
    ;

block: stat* retstat?
    ;

stat: ':'
    | varlist '=' explist
    | functioncall
    | label
    | 'break'
    | 'goto' NAME
    | 'do' block 'end'
    | 'while' exp 'do' block 'end'
    | 'repeat' block 'until' exp
    | 'if' exp 'then' block ('elseif' exp 'then' block)* ('else' block)? 'end'
    | 'for' NAME '=' exp ',' exp (',' exp)? 'do' block 'end'
    | 'for' namelist 'in' explist 'do' block 'end'
    | 'function' funcname funcbody
    | 'local' 'function' NAME funcbody
    | 'local' attnamelist ('=' explist)
    ;

attnamelist: NAME attrib (',' NAME attrib)*
    ;

attrib: ('<' NAME '>')?
    ;

retstat: 'return' explist? ';'?
    ;

label: '::' NAME '::'
    ;

funcname: NAME ('.' NAME)* (':' NAME)?
    ;

varlist: var (',' var)
    ;

var: NAME
    | prefixexp ('[' exp ']' | '.' NAME)
    ;

namelist: NAME (','  NAME)*
    ;

explist: exp (',' exp)*
    ;

exp: 'nil'
    | 'false'
    | 'true'
    | number
    | string
    | '...'
    | functiondef
    | prefixexp
    | tableconstructor
    | exp binop exp
    | unop exp
    ;

prefixexp: primaryexp ('[' exp ']' | '.' NAME | (':' NAME)? args)*
    ;

functioncall: prefixexp (':' NAME)? args
    ;

primaryexp: NAME
    | '(' exp ')'
    ;

args: '(' explist? ')'
    | tableconstructor
    | string
    ;

functiondef: 'function' funcbody
    ;

funcbody: '(' parlist? ')' block 'end'
    ;

parlist: namelist (',' '...')?
    | '...'
    ;

tableconstructor: '{' fieldlist? '}'
    ;

fieldlist: field (fieldsep field)* fieldsep?
    ;

field: '[' exp ']' '=' exp
    | NAME '=' exp
    | exp
    ;

fieldsep: ','
    | ';'
    ;

binop: '+'
    | '-'
    | '*'
    | '/'
    | '//'
    | '^'
    | '%'
    | '&'
    | '~'
    | '|'
    | '>>'
    | '<<'
    | '..'
    | '<'
    | '<='
    | '>'
    | '>='
    | '=='
    | '~='
    | 'and'
    | 'or'
    ;

unop: '-'
    | 'not'
    | '#'
    | '~'
    ;

number: INT
    | HEX
    | FLOAT
    | HEX_FLOAT
    ;

string: NORMALSTRING
    | CHARSTRING
    | LONGSTRING
    ;

// LEXER
NAME
    : [a-zA-Z_][a-zA-Z_0-9]*
    ;

NORMALSTRING
    : '"' ( EscapeSequence | ~('\\'|'"') )* '"'
    ;

CHARSTRING
    : '\'' ( EscapeSequence | ~('\''|'\\') )* '\''
    ;

LONGSTRING
    : '[' NESTED_STR ']'
    ;

fragment
NESTED_STR
    : '=' NESTED_STR '='
    | '[' .*? ']'
    ;

INT
    : Digit+
    ;

HEX
    : '0' [xX] HexDigit+
    ;

FLOAT
    : Digit+ '.' Digit* ExponentPart?
    | '.' Digit+ ExponentPart?
    | Digit+ ExponentPart
    ;

HEX_FLOAT
    : '0' [xX] HexDigit+ '.' HexDigit* HexExponentPart?
    | '0' [xX] '.' HexDigit+ HexExponentPart?
    | '0' [xX] HexDigit+ HexExponentPart
    ;

fragment
ExponentPart
    : [eE] [+-]? Digit+
    ;

fragment
HexExponentPart
    : [pP] [+-]? Digit+
    ;

fragment
EscapeSequence
    : '\\' [abfnrtvz"'\\]
    | '\\' '\r'? '\n'
    | DecimalEscape
    | HexEscape
    | UtfEscape
    ;
fragment
DecimalEscape
    : '\\' Digit
    | '\\' Digit Digit
    | '\\' [0-2] Digit Digit
    ;
fragment
HexEscape
    : '\\' 'x' HexDigit HexDigit
    ;
fragment
UtfEscape
    : '\\' 'u{' HexDigit+ '}'
    ;
fragment
Digit
    : [0-9]
    ;
fragment
HexDigit
    : [0-9a-fA-F]
    ;
COMMENT
    : '--[' NESTED_STR ']' -> channel(HIDDEN)
    ;
LINE_COMMENT
    : '--'
    (                                               // --
    | '[' '='*                                      // --[==
    | '[' '='* ~('='|'['|'\r'|'\n') ~('\r'|'\n')*   // --[==AA
    | ~('['|'\r'|'\n') ~('\r'|'\n')*                // --AAA
    ) ('\r\n'|'\r'|'\n'|EOF)
    -> channel(HIDDEN)
    ;
WS
    : [ \t\u000C\r\n]+ -> skip
    ;
SHEBANG
    : '#' '!' ~('\n'|'\r')* -> channel(HIDDEN)
    ;

