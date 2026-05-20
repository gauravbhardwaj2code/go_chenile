#!/usr/bin/env sh
set -eu

if [ "$#" -ne 1 ]; then
  echo "usage: $0 /absolute/path/to/go_chenile/chenile-framework" >&2
  exit 2
fi

framework_root=$1
for module in base bdd-utils chenile config core http owiz packager; do
  if [ ! -f "$framework_root/$module/go.mod" ]; then
    echo "missing $framework_root/$module/go.mod" >&2
    exit 1
  fi
  go mod edit -replace "github.com/gauravbhardwaj2code/go_chenile/chenile-framework/$module=$framework_root/$module"
done

go mod tidy
