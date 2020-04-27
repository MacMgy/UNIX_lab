#!/bin/bash -e

NORMALTEXT="\033[0m"

function succeful_message() {
    GREEN='\033[0;32m'
    echo ""
    printf "${GREEN} *** %s ${NORMALTEXT}\n" $1
}

function error() {
  RED='\033[0;31m'
  echo ""
  printf "${RED} %s ${NORMALTEXT}" $1
  echo ""
  exit 1
}

function clean_up() {
  rm -rf -- $mktemp_name
  [ ! -z $1 ] && error $1
}

function handle_SIGNALS() {
  clean_up
  error "User kill process."
}

filename="$1"

[ -z "$filename" ] && error 'First argument must be a file name'
[ ! -e "$filename" ] && error 'There is no file, you are looking for'
[ ! -r "$filename" ] && error 'Faild to get permission, when try read file'

OUTPUT_REGEX="s/^[[:space:]]*\/\/[[:space:]]*Output[[:space:]]*\([^ ]*\)$/\1/p"

executable_file_name=$(sed -n -e "$OUTPUT_REGEX" "$filename" | grep -m 1 "")
[ -z "$executable_file_name" ] && error 'Output name is not found or empty'

echo "Create temporary folder..."
mktemp_name=$(mktemp -d -t temp) || error 'Failed to create temp folder'

trap handle_SIGNALS HUP INT QUIT PIPE TERM

echo "Copying src to temporary folder..."
cp "$filename" $mktemp_name || clean_up 'Failed to copy file.' 

echo "Build src file..."
current_path=$(pwd)
cd "$mktemp_name"
go build -o "$executable_file_name" "$filename" || clean_up 'Failed compiling src file.'

echo "Move executable file to current path..."
cp "$executable_file_name" "$current_path" || clean_up 'Failed to move executable file'

clean_up
succeful_message "Succeful."
