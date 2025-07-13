#!/usr/bin/env bash

for app in ./cli/*; do
  if ! ([[ -d "$app" ]] && find "$app" -maxdepth 1 -name "*.go" | grep -q .); then
    continue
  fi
  echo "$app"
done
