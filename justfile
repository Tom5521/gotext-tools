gopath := `go env GOPATH`
goos := `go env GOOS`
goarch := `go env GOARCH`

@default :
  just build msgomerge
  just build xgotext
run app args:
  GOOS={{goos}} GOARCH={{goarch}} go run -v ./cli/{{app}} {{args}}
test:
  go clean -testcache
  go test ./...
@benchmark:
  just bench ./...
bench path:
  go test -bench=. {{path}}
@puml:
  #!/usr/bin/env bash
  set -euxo pipefail
 
  paths=( 
    ./pkg/go/parse
    ./pkg/po/parse
    ./pkg/po/compile
    ./pkg/po
  )

  for path in "${paths[@]}"; do
    goplantuml "$path" > "$path/structure.puml"
    plantuml -theme spacelab "$path/structure.puml"
  done
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
  rm {{gopath}}/bin/{{app}}.exe
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
build app:
  # building {{app}} for {{goos}}-{{goarch}}...
  @GOOS={{goos}} GOARCH={{goarch}} \
  go build -o \
  builds/{{app}}-{{goos}}-{{goarch}}\
  $([[ "{{goos}}" == "windows" ]] && echo ".exe") \
  -ldflags '-s -w' \
  ./cli/{{app}}
[private]
build-all-app app:
  #!/usr/bin/env bash
  set -euo pipefail

  archs=(
    "386"
    "amd64"
    "arm"
    "arm64"
  )
  oses=(
    "windows"
    "linux"
    "netbsd"
    "openbsd"
    "plan9"
  )

  for os in "${oses[@]}"; do
    for arch in "${archs[@]}"; do
      if echo $(go tool dist list) | grep -q "$os/$arch"; then  
        just goos="$os" goarch="$arch" build {{app}}
      fi
    done
  done
build-all: clean
  #!/usr/bin/env bash
  set -euo pipefail
  for app in ./cli/*; do
    just build-all-app "$(basename "$app")"
  done
[confirm]
release: clean build-all
  gh release create {{`git describe --tags --abbrev=0`}} \
  --generate-notes --fail-on-no-commits builds/*
