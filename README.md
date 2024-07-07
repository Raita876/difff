# Difff

This CLI compares files located in two directories and outputs the differences.

## Install

### go install

```
$ go install github.com/Raita876/difff/cmd/difff@latest
```

## Usage

```
$ difff <source_path> <target_path>
```

## Example

```
$ find ./e2e/data/source -type f
./e2e/data/source/e/h.txt
./e2e/data/source/e/i.txt
./e2e/data/source/e/f/g.txt
./e2e/data/source/a.txt
./e2e/data/source/b.txt
./e2e/data/source/c/d.txt
$ find ./e2e/data/target -type f
./e2e/data/target/e/h.txt
./e2e/data/target/e/f/g.txt
./e2e/data/target/e/f/j.txt
./e2e/data/target/a.txt
./e2e/data/target/b.txt
./e2e/data/target/c/d.txt
$ ./bin/linux/amd64/difff ./e2e/data/source ./e2e/data/target
{
  "source": "./e2e/data/source",
  "target": "./e2e/data/target",
  "diff": {
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
}
```
