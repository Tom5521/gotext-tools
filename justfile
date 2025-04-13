gopath := `go env GOPATH`
goos := `go env GOOS`
goarch := `go env GOARCH`
verbose := env("VERBOSE","0")

@default :
  just build msgomerge
  just build xgotext
run app args:
  GOOS={{goos}} GOARCH={{goarch}} go run -v ./cli/{{app}} {{args}}
test:
  #!/bin/env bash
  go clean -testcache
  go test \
  $([[ {{verbose}} == 1 ]] && echo "-v") \
  ./...
@benchmark:
  just verbose={{verbose}} bench ./...
bench path:
  go test \
  $([[ "{{verbose}}" == "1" ]] && echo -v) \
  -bench=. {{path}}
puml:
  #!/usr/bin/env bash
 
  paths=( 
    ./pkg/go/parse
    ./pkg/po/parse
    ./pkg/po/compile
    ./pkg/po
  )

  for path in "${paths[@]}"; do
    echo -n {{BOLD}}"Generating PUML's of $path..."{{NORMAL}}
    goplantuml "$path" > "$path/structure.puml"
    plantuml -theme spacelab "$path/structure.puml"

    if [[ $? == 0 ]]; then
    echo {{GREEN}}"DONE"{{NORMAL}}
    else
      echo {{BG_RED}}"ERROR: $?"{{NORMAL}}
    fi
  done
clean:
  #!/bin/env bash
  echo -n {{BOLD}}"Cleaning..."{{NORMAL}}
  rm -rf \
  $(find . \( -name "*.po" -o -name "*.mo" -o -name "*.pot" -o -name "*.log" \)) \
  builds
  echo {{GREEN}}"DONE"{{NORMAL}}
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
  #!/usr/bin/env bash

  echo -n {{BOLD}}"Building {{app}} for {{goos}}-{{goarch}}..."{{NORMAL}}

  GOOS={{goos}} GOARCH={{goarch}} \
  go build \
  $([[ "{{verbose}}" == "1" ]] && echo "-v") \
  -o builds/{{app}}-{{goos}}-{{goarch}}\
  $([[ "{{goos}}" == "windows" ]] && echo ".exe") \
  -ldflags '-s -w' \
  ./cli/{{app}}

  if [[ $? == 0 ]]; then
    echo {{GREEN}}"OK"{{NORMAL}}
  else
    echo {{BG_RED}}"ERROR: $?"{{NORMAL}}
  fi
[private]
build-all-app app:
  #!/usr/bin/env bash

  export VERBOSE={{verbose}}

  archs=(
    "386"
    "amd64"
    "arm"
    "arm64"
  )
  oses=(
    "linux"
    "netbsd"
    "openbsd"
    "freebsd"
    "plan9"
    "windows"
    "darwin"
    "solaris"
  )

  for os in "${oses[@]}"; do
    for arch in "${archs[@]}"; do
      if echo $(go tool dist list) | grep -qw "$os/$arch"; then  
        just goos="$os" goarch="$arch" build {{app}}
      fi
    done
  done
build-all: clean
  #!/usr/bin/env bash 
  export VERBOSE={{verbose}}

  for app in ./cli/*; do
    just build-all-app "$(basename "$app")"
  done
[confirm]
release: clean build-all
  gh release create {{`git describe --tags --abbrev=0`}} \
  --generate-notes --fail-on-no-commits builds/*
