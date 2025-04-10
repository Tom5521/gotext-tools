test:
  go clean -testcache
  go test -v $(dirname $(find . -name "*_test.go"))
benchmark:
  go test -bench=. $(dirname $(find . -name "*_benchmark_test.go"))
bench path:
  go test -bench=. {{path}}
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
  rm -rf $(find . -name "*.po") \
  $(find . -name "*.mo") \
  $(find . -name "*.log") \
  builds
gen-diff:
  git diff --staged > diff.log
build app os arch:
  @GOOS={{os}} GOARCH={{arch}} \
  go build -o \
  builds/{{app}}-{{os}}-{{arch}}\
  $([[ "{{os}}" == "windows" ]] && echo ".exe") \
  -ldflags '-s -w' \
  ./cli/{{app}}
[private]
build-all-unix app os:
  just build {{app}} {{os}} 386
  just build {{app}} {{os}} amd64
  just build {{app}} {{os}} arm
  just build {{app}} {{os}} arm64
build-all-app app:
  @just build-all-unix {{app}} linux
  @just build-all-unix {{app}} openbsd
  @just build-all-unix {{app}} netbsd

  just build {{app}} windows 386
  just build {{app}} windows amd64
  just build {{app}} windows arm64

  just build {{app}} darwin amd64
  just build {{app}} darwin arm64
build-all:
  @just build-all-app msgomerge
  @just build-all-app xgotext
