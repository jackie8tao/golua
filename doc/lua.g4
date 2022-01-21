grammar lua;

chunk : block;

// 语法块
block: stat* retstat?;

// 语句
stat: ':'
    | varlist '=' explist
    | functioncall
    | label
    | ;

//