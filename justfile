gopath := env("GOPATH",`go env GOPATH`)
goos := env("GOOS",`go env GOOS`)
goarch := env("GOARCH",`go env GOARCH`)
verbose := env("VERBOSE","0")
ext := if goos == "windows" { ".exe" } else { "" }

default :
  #!/bin/env bash
  set -euo pipefail
  
  for app in ./cli/*; do
    just build $(basename $app)
  done
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
    structure="$path/structure.puml"
    goplantuml "$path" > "$structure"
    plantuml -theme spacelab "$structure"

    if [[ $? == 0 ]]; then
    echo {{GREEN}}"OK"{{NORMAL}}
    else
      echo {{RED}}"ERROR: $?"{{NORMAL}}
    fi
  done
clean:
  #!/bin/env bash
  echo -n {{BOLD}}"Cleaning..."{{NORMAL}}
  rm -rf \
  $(find . \( -name "*.po" -o -name "*.mo" -o -name "*.pot" -o -name "*.log" \)) \
  builds
  echo {{GREEN}}"OK"{{NORMAL}}
diff:
  git diff --staged > diff.log
go-install app:
  go install -v -ldflags '-s -w' ./cli/{{app}}
go-uninstall app:
  rm {{gopath}}/bin/{{app}}{{ext}} -f
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
  echo -n {{BOLD}}"{{goos}}-{{goarch}}..."{{NORMAL}}

  GOOS={{goos}} GOARCH={{goarch}} \
  go build \
  $([[ "{{verbose}}" == "1" ]] && echo "-v") \
  -o builds/{{app}}-{{goos}}-{{goarch}}{{ext}} \
  -ldflags '-s -w' \
  ./cli/{{app}}

  if [[ $? == 0 ]]; then
    echo {{GREEN}}"OK"{{NORMAL}}
  else
    echo {{RED}}"ERROR: $?"{{NORMAL}}
  fi
build-all-app app:
  #!/usr/bin/env bash
  set -euo pipefail

  export VERBOSE={{verbose}}

  archs=(
    386 # Ahem...
    amd64
    arm
    arm64

    # WHO TF USE THIS ARCHITECTURES?!?!?!?!
    ppc64
    ppc64le
    riscv64
  )
  oses=(
    linux
    netbsd
    freebsd
    windows
    darwin

    # And... the distros that nobody uses
    openbsd
    plan9
    solaris
    dragonfly
  )

  valid=$(go tool dist list)

  for os in "${oses[@]}"; do
    for arch in "${archs[@]}"; do
      if echo $valid | grep -qw "$os/$arch"; then  
        just goos="$os" goarch="$arch" build {{app}}
      fi
    done
  done
build-all: clean
  #!/usr/bin/env bash 
  set -euo pipefail
  export VERBOSE={{verbose}}

  for app in ./cli/*; do
    name=$(basename "$app")
    echo {{BOLD}}{{BG_WHITE}}{{BLACK}}"----- $name -----"{{NORMAL}}
    just build-all-app "$name"
  done
  echo {{BOLD}}{{BG_WHITE}}{{BLACK}}"----- FINISHED -----"{{NORMAL}}
[confirm]
release: clean build-all
  gh release create {{`git describe --tags --abbrev=0`}} \
  --generate-notes --fail-on-no-commits builds/*
