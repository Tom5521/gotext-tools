#!/usr/bin/env bash

for path in $(./scripts/app-paths.sh); do
  basename "$path"
done
