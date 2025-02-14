test:
  go test -v ./pkg/po/types
  go test -v ./pkg/po/parse/generator
  go test -v ./pkg/po/parse/lexer
  go test -v ./pkg/po/parse/ast
  go test -v ./pkg/go/parse/
bench:
  go test -v -bench=. ./pkg/go/parse
  go test -v -bench=. ./pkg/po/parse
  go test -v -bench=. ./pkg/po/parse/lexer
  go test -v -bench=. ./pkg/po/parse/ast

