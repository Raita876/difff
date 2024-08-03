#!/bin/bash
set -euo pipefail

THIS_SCRIPT_DIR=$(
  cd "$(dirname "$0")"
  pwd
)
readonly THIS_SCRIPT_DIR

readonly BIN_PATH="${THIS_SCRIPT_DIR}/../${CLI_BIN}"
readonly SOURCE_DIR_PATH="${THIS_SCRIPT_DIR}/data/source"
readonly TARGET_DIR_PATH="${THIS_SCRIPT_DIR}/data/target"

WANT_JSON=$(
  cat <<EOS
{
  "source": {
    "path": "${SOURCE_DIR_PATH}",
    "num": 6
  },
  "target": {
    "path": "${TARGET_DIR_PATH}",
    "num": 6
  },
  "diff": {
    "source": {
      "num": 2,
      "results": [
        {
          "path": "e/f/g.txt",
          "hash": "211c102123b4a41bd5227dcc84952349"
        },
        {
          "path": "e/i.txt",
          "hash": "1b08ef3ea73ce6fd8b2ef57f54073b5a"
        }
      ]
    },
    "target": {
      "num": 2,
      "results": [
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
  }
}
EOS
)
readonly WANT_JSON

WANT_YAML=$(
  cat <<EOS
source:
  path: ${SOURCE_DIR_PATH}
  num: 6
target:
  path: ${TARGET_DIR_PATH}
  num: 6
diff:
  source:
    num: 2
    results:
    - path: e/f/g.txt
      hash: 211c102123b4a41bd5227dcc84952349
    - path: e/i.txt
      hash: 1b08ef3ea73ce6fd8b2ef57f54073b5a
  target:
    num: 2
    results:
    - path: e/f/g.txt
      hash: e4727cb9315a4fddec71e1a85cad6c09
    - path: e/f/j.txt
      hash: f6c79025f3b5bedac7cd769f0847e36a
EOS
)
readonly WANT_YAML

WANT_XML=$(
  cat <<EOS
<DiffResponse>
  <source>
    <path>${SOURCE_DIR_PATH}</path>
    <num>6</num>
  </source>
  <target>
    <path>${TARGET_DIR_PATH}</path>
    <num>6</num>
  </target>
  <diff>
    <source>
      <num>2</num>
      <results>
        <path>e/f/g.txt</path>
        <hash>211c102123b4a41bd5227dcc84952349</hash>
      </results>
      <results>
        <path>e/i.txt</path>
        <hash>1b08ef3ea73ce6fd8b2ef57f54073b5a</hash>
      </results>
    </source>
    <target>
      <num>2</num>
      <results>
        <path>e/f/g.txt</path>
        <hash>e4727cb9315a4fddec71e1a85cad6c09</hash>
      </results>
      <results>
        <path>e/f/j.txt</path>
        <hash>f6c79025f3b5bedac7cd769f0847e36a</hash>
      </results>
    </target>
  </diff>
</DiffResponse>
EOS
)
readonly WANT_XML

function main() {
  local result_json
  result_json=$(
    ${BIN_PATH} --format JSON \
      "${SOURCE_DIR_PATH}" \
      "${TARGET_DIR_PATH}"
  )

  diff -u \
    <(echo "${result_json}") \
    <(echo "${WANT_JSON}")

  local result_yaml
  result_yaml=$(
    ${BIN_PATH} --format YAML \
      "${SOURCE_DIR_PATH}" \
      "${TARGET_DIR_PATH}"
  )

  diff -u \
    <(echo "${result_yaml}") \
    <(echo "${WANT_YAML}")

  local result_xml
  result_xml=$(
    ${BIN_PATH} --format XML \
      "${SOURCE_DIR_PATH}" \
      "${TARGET_DIR_PATH}"
  )

  diff -u \
    <(echo "${result_xml}") \
    <(echo "${WANT_XML}")

}

main "$@"
