#!/bin/bash -e

function error() {
  BIGWHITETEXT="\033[1;37m"
  BGRED='\033[41m'
  NORMAL="\033[0m"
  echo ""
  printf "${BIGWHITETEXT}${BGRED} %s ${NORMAL}" $1
  echo ""
  exit 1
}

function succeful_message() {
    BLACK="\033[30m"
    GREEN='\033[0;32m'
    NORMAL="\033[0m"
    BGGREEN="\033[42m"
    printf "${GREEN} *** %s ${NORMAL}\n" $1
}

function handle_SIGNALS() {
  rm -rf -- $mktemp_name
  error "User kill process."
}

file_name="$1"

[ -z "$file_name" ] && error 'First argument must be a file name'
[ ! -e "$file_name" ] && error 'File does not exist'
[ ! -r "$file_name" ] && error 'File can not be read'

OUTPUT_REGEX="s/^[[:space:]]*\/\/[[:space:]]*Output[[:space:]]*\([^ ]*\)$/\1/p"

executable_file_name=$(sed -n -e "$OUTPUT_REGEX" "$file_name" | grep -m 1 "")
[ -z "$executable_file_name" ] && error 'Output name is not found'

echo "Create temporary folder..."
mktemp_name=$(mktemp -d -t temp) || error 'Failed to create temp folder'

trap handle_SIGNALS HUP INT QUIT PIPE TERM

echo "Copying src to temporary folder..."
cp "$file_name" $mktemp_name || { rm -rf -- "$mktemp_name"; error 'Failed to copy file.'; }

echo "Build src file..."
current_path=$(pwd)
cd "$mktemp_name"
go build -o "$executable_file_name" "$file_name" || { rm -rf -- "$mktemp_name"; error "Failed compiling src file."; }

echo "Move executable file to current path..."
cp "$executable_file_name" "$current_path" || { rm -rf -- "$mktemp_name"; error "Failed to move executable file"; }

rm -rf -- "$mktemp_name"
succeful_message "Succeful."
