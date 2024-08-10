#!/bin/sh

# Ensure a commit hash is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <commit-hash>"
  exit 1
fi

commit_hash=$1

# Get the commit date in ISO 8601 format
commit_date=$(git show -s --format=%ci $commit_hash)
if [ -z "$commit_date" ]; then
  echo "Invalid commit hash."
  exit 1
fi

# Convert the commit date to UTC and format it as YYYYMMDDHHMMSS
datetime=$(date -u -j -f "%Y-%m-%d %H:%M:%S %z" "$commit_date" +"%Y%m%d%H%M%S")
if [ -z "$datetime" ]; then
  echo "Failed to convert commit date to UTC."
  exit 1
fi

# Construct the pseudo-version
pseudo_version="${datetime}-${commit_hash:0:12}"
echo $pseudo_version