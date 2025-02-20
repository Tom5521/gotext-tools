test:
  go test -v ./pkg/po
  go test -v ./pkg/po/parse/lexer
  go test -v ./pkg/po/parse/ast
  go test -v ./pkg/go/parse
bench:
  go test -v -bench=. ./pkg/go/parse
  go test -v -bench=. ./pkg/po/parse
  go test -v -bench=. ./pkg/po/parse/lexer
  go test -v -bench=. ./pkg/po/parse/ast
gen-uml:
  goplantuml -recursive ./pkg/go/parse > ./pkg/go/parse/structure.puml
  goplantuml -recursive ./pkg/po/parse > ./pkg/po/parse/structure.puml
  goplantuml ./pkg/po/compiler/ > ./pkg/po/compiler/structure.puml
  goplantuml ./pkg/po/ > ./pkg/po/structure.puml

  plantuml -theme spacelab ./pkg/po/compiler/structure.puml
  plantuml -theme spacelab ./pkg/po/structure.puml
  plantuml -theme spacelab ./pkg/po/parse/structure.puml
  plantuml -theme spacelab ./pkg/go/parse/structure.puml
