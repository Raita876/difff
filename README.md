# Difff

This CLI compares files located in two directories and outputs the differences.

## Install

### go install

```bash
go install github.com/Raita876/difff/cmd/difff@latest
```

### binary install

Check [difff/releases](https://github.com/Raita876/difff/releases) for the latest version.

```bash
DOWNLOAD_OS=Linux
DOWNLOAD_ARCH=x86_64
DOWNLOAD_VERSION=0.5.2
curl -L \
  https://github.com/Raita876/difff/releases/download/${DOWNLOAD_VERSION}/difff_${DOWNLOAD_OS}_${DOWNLOAD_ARCH}.tar.gz \
  -o /tmp/difff_${DOWNLOAD_OS}_${DOWNLOAD_ARCH}.tar.gz
tar -C /tmp -xzf /tmp/difff_${DOWNLOAD_OS}_${DOWNLOAD_ARCH}.tar.gz
chmod 755 /tmp/difff
mv /tmp/difff /usr/local/bin/
```

## Example

```bash
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
$ difff ./e2e/data/source ./e2e/data/target
{
  "source": {
    "path": "e2e/data/source",
    "num": 6
  },
  "target": {
    "path": "e2e/data/target",
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
```
