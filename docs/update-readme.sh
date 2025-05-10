#!/bin/sh
set -ex

TMP_PATH=/tmp/pdfoptimize-build

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd "${DIR}/.." || exit 1

rm -f "${TMP_PATH}"
go build -o "${TMP_PATH}"

io --allow-exec --template "docs/README.md.tpl" --input "{\"binary_path\":\"${TMP_PATH}\"}" --output "README.md"

rm -f "${TMP_PATH}"