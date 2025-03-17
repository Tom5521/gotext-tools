test:
  go test -v ./pkg/po
  go test -v ./pkg/go/parse
  go test -v ./internal/util
  go test -v ./pkg/po/compiler
  go test -v ./pkg/po/parse
benchmark:
  go test -bench=. ./pkg/go/parse
  go test -bench=. ./pkg/po/parse
  go test -bench=. ./pkg/po/compiler
  go test -bench=. ./pkg/po
gen-uml:
  goplantuml ./pkg/go/parse > ./pkg/go/parse/structure.puml
  goplantuml ./pkg/po/compiler/ > ./pkg/po/compiler/structure.puml
  goplantuml ./pkg/po/ > ./pkg/po/structure.puml
  goplantuml ./pkg/po/parse/ > ./pkg/po/parse/structure.puml

  plantuml -theme spacelab ./pkg/po/compiler/structure.puml
  plantuml -theme spacelab ./pkg/po/structure.puml
  plantuml -theme spacelab ./pkg/go/parse/structure.puml
  plantuml -theme spacelab ./pkg/po/parse/structure.puml
clean:
  -rm *.po* *.log
  -rm ./internal/antlr-po/*.interp ./internal/antlr-po/*.tokens
gen-diff:
  git diff --staged > diff.log
