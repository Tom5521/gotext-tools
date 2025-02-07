

[private]
run-test path:
  go test -v {{path}}

test:
  @just run-test ./pkg/po/types
  @just run-test ./pkg/po/parse/generator
  @just run-test ./pkg/po/parse/lexer
  @just run-test ./pkg/po/parse/ast
  @just run-test ./pkg/go/parse/
