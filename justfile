gopath := env("GOPATH", `go env GOPATH`)
goos := env("GOOS", `go env GOOS`)
goarch := env("GOARCH", `go env GOARCH`)
gocmd := env("GOCMD", "go")
verbose := env("VERBOSE", "0")
ext := if goos == "windows" { ".exe" } else { "" }

default:
    #!/usr/bin/env bash
    set -euo pipefail

    export VERBOSE={{ verbose }}
    export GOCMD={{ gocmd }}

    for app in ./cli/*; do
        if ! ([[ -d "$app" ]] && find "$app" -maxdepth 1 -name "*.go"| grep -q .); then
            continue
        fi
        echo -n {{ BOLD }}$(basename $app)-{{ NORMAL }}
        just build $(basename $app)
    done

run app args:
    GOOS={{ goos }} GOARCH={{ goarch }} {{ gocmd }} run -C ./cli \
    $([[ "{{ verbose }}" == "1" ]] && echo "-v") \
    ./{{ app }} {{ args }}

test:
    {{ gocmd }} clean -testcache
    {{ gocmd }} test \
    $([ "{{ verbose }}" -eq "1" ] && echo "-v") \
    ./...

test-versions:
    just test
    just gocmd=go1.18 test

@benchmark:
    just verbose={{ verbose }} bench ./...

bench path:
    #!/usr/bin/env bash
    {{ gocmd }} test \
    $([[ "{{ verbose }}" == "1" ]] && echo "-v") \
    -bench=. {{ path }}

puml:
    #!/usr/bin/env bash

    paths=( 
        ./pkg/go/parse
        ./pkg/po/parse
        ./pkg/po/compile
        ./pkg/po
    )

    for path in "${paths[@]}"; do
        echo -n {{ BOLD }}"Generating PUML's of $path..."{{ NORMAL }}
        structure="$path/structure.puml"
        goplantuml "$path" > "$structure"
        plantuml -theme spacelab "$structure"

        if [[ $? == 0 ]]; then
            echo {{ GREEN }}"OK"{{ NORMAL }}
        else
            echo {{ RED }}"ERROR: $?"{{ NORMAL }}
        fi
    done

clean:
    #!/usr/bin/env bash
    echo -n {{ BOLD }}"Cleaning..."{{ NORMAL }}
    rm -rf \
    $(find . \( -name "*.po" -o -name "*.mo" -o -name "*.pot" -o -name "*.log" \)) \
    builds
    cd cli && goreleaser --clean
    echo {{ GREEN }}"OK"{{ NORMAL }}

diff:
    git diff --staged > diff.log

go-install app:
    {{ gocmd }} install -C ./cli -v -ldflags '-s -w' ./{{ app }}

go-uninstall app:
    rm {{ gopath }}/bin/{{ app }}{{ ext }} -f

[unix]
local-install app:
    just build {{ app }} {{ goos }} {{ goarch }}
    mv ./builds/{{ app }}-{{ goos }}-{{ goarch }} ~/.local/bin/{{ app }}

[unix]
local-uninstall app:
    rm ~/.local/bin/{{ app }}

[unix]
root-install app:
    just build {{ app }} {{ goos }} {{ goarch }}
    sudo mv ./builds/{{ app }}-{{ goos }}-{{ goarch }} /usr/local/bin/{{ app }}

[unix]
root-uninstall app:
    sudo rm /usr/local/bin/{{ app }}
cli-docs:
    #!/usr/bin/env bash
    cd cli
    rm -rf docs
    go run ./gotext-tools doc-tree -D docs
release:
    cd cli && goreleaser release --clean
