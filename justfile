gopath := `go env GOPATH`
goos := `go env GOOS`
goarch := `go env GOARCH`

default app:
  just build {{app}} {{goos}} {{goarch}}
test:
  go clean -testcache
  go test -v $(dirname $(find . -name "*_test.go"))
benchmark:
  go test -bench=. $(dirname $(find . -name "*_benchmark_test.go"))
bench path:
  go test -bench=. {{path}}
puml:
  goplantuml ./pkg/go/parse > ./pkg/go/parse/structure.puml
  goplantuml ./pkg/po/compile/ > ./pkg/po/compile/structure.puml
  goplantuml ./pkg/po/ > ./pkg/po/structure.puml
  goplantuml ./pkg/po/parse/ > ./pkg/po/parse/structure.puml

  plantuml -theme spacelab ./pkg/po/compile/structure.puml
  plantuml -theme spacelab ./pkg/po/structure.puml
  plantuml -theme spacelab ./pkg/go/parse/structure.puml
  plantuml -theme spacelab ./pkg/po/parse/structure.puml
clean:
  # Cleaning...
  @rm -rf \
  $(find . \( -name "*.po" -o -name "*.mo" -o -name "*.pot" -o -name "*.log" \)) \
  builds
  # Cleaned!
diff:
  git diff --staged > diff.log
go-install app:
  go install -v -ldflags '-s -w' ./cli/{{app}}
[windows]
go-uninstall app:
  del {{gopath}}/bin/{{app}}.exe
[unix]
go-uninstall app:
  rm {{gopath}}/bin/{{app}} -f
[unix]
local-install app:
  just build {{app}} {{goos}} {{goarch}}
  mv ./builds/{{app}}-{{goos}}-{{goarch}} ~/.local/bin/{{app}}
[unix]
local-uninstall app:
  rm ~/.local/bin/{{app}}
[unix]
root-install app:
  just build {{app}} {{goos}} {{goarch}}
  sudo mv ./builds/{{app}}-{{goos}}-{{goarch}} /usr/local/bin/{{app}}
[unix]
root-uninstall app:
  sudo rm /usr/local/bin/{{app}}
build app os arch:
  # building {{app}} for {{os}}-{{arch}}...
  @GOOS={{os}} GOARCH={{arch}} \
  go build -o \
  builds/{{app}}-{{os}}-{{arch}} \
  -ldflags '-s -w' \
  ./cli/{{app}}
[private]
@win-build app arch:
  just build {{app}} windows {{arch}}
  mv ./builds/{{app}}-windows-{{arch}} ./builds/{{app}}-windows-{{arch}}.exe
[private]
@build-all-unix app os:
  just build {{app}} {{os}} 386
  just build {{app}} {{os}} amd64
  just build {{app}} {{os}} arm
  just build {{app}} {{os}} arm64
@build-all-app app:
  # ---- building {{app}} -----
  just build-all-unix {{app}} linux
  just build-all-unix {{app}} openbsd
  just build-all-unix {{app}} netbsd

  just win-build {{app}} 386
  just win-build {{app}} amd64
  just win-build {{app}} arm64

  just build {{app}} darwin amd64
  just build {{app}} darwin arm64
@build-all: clean
  just build-all-app msgomerge
  just build-all-app xgotext
[confirm]
release: clean
  just build-all
  gh release create {{`git describe --tags --abbrev=0`}} \
  --generate-notes --fail-on-no-commits builds/*
