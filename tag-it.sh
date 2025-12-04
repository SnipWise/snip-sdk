#!/bin/bash
: <<'COMMENT'
Make a tag
COMMENT

find . -name '.DS_Store' -type f -delete

git add .
git commit -m "🏷️ save step: ${1}"

git tag -a ${1} -m "🏷️ save step"

