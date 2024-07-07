#!/bin/bash
set -euo pipefail

THIS_SCRIPT_DIR=$(
    cd "$(dirname "$0")"
    pwd
)
readonly THIS_SCRIPT_DIR

readonly BIN_PATH="${THIS_SCRIPT_DIR}/../bin/${GOOS}/${GOARCH}/difff"
readonly SOURCE_DIR_PATH="${THIS_SCRIPT_DIR}/data/source"
readonly TARGET_DIR_PATH="${THIS_SCRIPT_DIR}/data/target"

WANT=$(
    cat <<EOS
{
  "source": [
    {
      "path": "e/f/g.txt",
      "hash": "211c102123b4a41bd5227dcc84952349"
    },
    {
      "path": "e/i.txt",
      "hash": "1b08ef3ea73ce6fd8b2ef57f54073b5a"
    }
  ],
  "target": [
    {
      "path": "e/f/g.txt",
      "hash": "e4727cb9315a4fddec71e1a85cad6c09"
    },
    {
      "path": "e/f/j.txt",
      "hash": "f6c79025f3b5bedac7cd769f0847e36a"
    }
  ]
}
EOS
)
readonly WANT

function main() {
    local result
    result=$(
        ${BIN_PATH} \
            "${SOURCE_DIR_PATH}" \
            "${TARGET_DIR_PATH}"
    )

    local result_diff
    result_diff=$(echo "${result}" | jq '.diff')

    diff -u \
        <(echo "${result_diff}") \
        <(echo "${WANT}")

}

main "$@"
