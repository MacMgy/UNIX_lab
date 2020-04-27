#!/bin/bash -e

NORMALTEXT="\033[0m"

function handle_SIGNALS() {
  rm -rf -- $mktemp_name
  error "User kill process."
}

function succeful_message() {
    GREEN='\033[0;32m'
    echo ""
    printf "${GREEN} *** %s ${NORMALTEXT}\n" $1
    echo ""
}

function error() {
  RED='\033[0;31m'
  echo ""
  printf "${RED} %s ${NORMALTEXT}" $1
  echo ""
  exit 1
}

filename="$1"

[ -z "$filename" ] && error 'First argument must be a file name'
[ ! -e "$filename" ] && error 'File does not exist'
[ ! -r "$filename" ] && error 'File can not be read'

OUTPUT_REGEX="s/^[[:space:]]*\/\/[[:space:]]*Output[[:space:]]*\([^ ]*\)$/\1/p"

executable_file_name=$(sed -n -e "$OUTPUT_REGEX" "$filename" | grep -m 1 "")
[ -z "$executable_file_name" ] && error 'Output name is not found'

echo "Create temporary folder..."
mktemp_name=$(mktemp -d -t temp) || error 'Failed to create temp folder'

trap handle_SIGNALS HUP INT QUIT PIPE TERM

echo "Copying src to temporary folder..."
cp "$filename" $mktemp_name || { rm -rf -- "$mktemp_name"; error 'Failed to copy file.'; }

echo "Build src file..."
current_path=$(pwd)
cd "$mktemp_name"
go build -o "$executable_file_name" "$filename" || { rm -rf -- "$mktemp_name"; error "Failed compiling src file."; }

echo "Move executable file to current path..."
cp "$executable_file_name" "$current_path" || { rm -rf -- "$mktemp_name"; error "Failed to move executable file"; }

rm -rf -- "$mktemp_name"
succeful_message "Succeful."
