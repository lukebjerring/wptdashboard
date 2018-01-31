#!/usr/bin/env bash

# Helper script for getting the standardized staging version name for the repo,
# for use when deploying.

BRANCH_NAME="$(git rev-parse --abbrev-ref HEAD)"
USER="$(git remote -v get-url origin | sed -E 's#(https?:\/\/|git@)github.com(\/|:)##' | sed 's#/.*$##')"
VERSION="${USER}-${BRANCH_NAME}"

print_version() {
  echo "${VERSION}"
}
