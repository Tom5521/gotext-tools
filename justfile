test:
  go test -v ./pkg/po
  go test -v ./pkg/go/parse
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
  antlr4 -Dlanguage=Go -package parse ./pkg/po/parse/Po.g4
