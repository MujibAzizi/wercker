#!/bin/bash

print_usage() {
  echo "Usage: ./release.sh [options] [version]" >&2
}


main() {
  # Define options
  local version
  local skip_branch="false"
  local skip_dirty="false"
  local skip_changelog="false"
  local branch="master"

  # Do some flag parsing
  while getopts "b:" opt; do
    case $opt in
      b)
        branch="$OPTARG"
        ;;
      \?)
        fail
        ;;
    esac
  done

  shift $((OPTIND-1))
  version="$1"

  if [ -z "$version" ]; then
    fail "version is required"
  fi

  printf "Performing sanity checks before releasing version v%s\n\n" "$version"

  check_branch "$skip_branch" "$branch"
  check_dirty "$skip_dirty"
  check_changelog "$skip_changelog" "$version"
}

check_branch() {
  printf "Checking if current branch matches \"%s\"... " "$2"

  if [ "$1" == "true" ]; then
    echo "skipped"
    return
  fi

  currentBranch=$(git rev-parse --abbrev-ref HEAD)

  if [ "$currentBranch" != "$2" ]; then
    echo "fail"

    printf "\n\tRelease is only allowed from branch \"%s\", current branch \"%s\"\n\n" "$2" "$currentBranch"
    exit 1
  fi

  echo "ok"
}

check_dirty() {
  printf "Checking if current git commit is dirty... "

  if [ "$1" == "true" ]; then
    echo "skipped"
    return
  fi

  if ! git diff --quiet; then
    echo "fail"

    printf "\n\tWorking directory is dirty.\n\tDid you forget to add all files?\n\n"
    exit 1
  fi

  if ! git diff --cached --quiet; then
    echo "fail"

    printf "\n\tIndex is dirty.\n\tDid you forget commit your files?\n\n"
    exit 1
  fi


  echo "ok"
}

check_changelog() {
  printf "Checking changelog for version %s... " "$2"

  if [ "$1" == "true" ]; then
    echo "skipped"
    return
  fi

  if head -n 1 changelog.md | grep -q -v "v$2"; then
    echo "failed"

    printf "\n\tVersion: \"v%s\" was not found as the latest version in the changelog\n\tDid you forget to update the changelog?\n\n" "$2"
    exit 1
  fi

  echo "ok"
}

fail() {
  echo "$@"
  print_usage
  exit 1
}

main "$@";
