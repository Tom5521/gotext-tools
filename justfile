test:
  go test -v ./pkg/po
  go test -v ./pkg/go/parse
  go test -v ./internal/util
  go test -v ./pkg/po/compiler
  go test -v ./pkg/po/parse
bench:
  go test -v -bench=. ./pkg/go/parse
gen-uml:
  goplantuml -recursive ./pkg/go/parse > ./pkg/go/parse/structure.puml
  goplantuml ./pkg/po/compiler/ > ./pkg/po/compiler/structure.puml
  goplantuml ./pkg/po/ > ./pkg/po/structure.puml

  plantuml -theme spacelab ./pkg/po/compiler/structure.puml
  plantuml -theme spacelab ./pkg/po/structure.puml
  plantuml -theme spacelab ./pkg/go/parse/structure.puml
gen-parser:
  antlr4 -Dlanguage=Go -package parse ./internal/antlr-po/Po.g4
clean:
  rm *.po* *.log
  rm ./internal/antlr-po/*.interp ./internal/antlr-po/*.tokens
gen-diff:
  git diff --staged > diff.log
