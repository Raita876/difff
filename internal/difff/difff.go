package difff

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/samber/lo"
)

type DiffInfo struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Diff   Diff   `json:"diff"`
}

type Diff struct {
	Source Results `json:"source"`
	Target Results `json:"target"`
}

type Result struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

type Results []Result

type ResultsInfo struct {
	root    string
	results Results
}

func getResults(root string) (ResultsInfo, error) {
	rs := Results{}
	cd, err := os.Getwd()
	if err != nil {
		return ResultsInfo{}, err
	}
	// TODO: error handling
	//nolint:errcheck
	defer os.Chdir(cd)

	err = os.Chdir(root)
	if err != nil {
		return ResultsInfo{}, err
	}

	err = filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		hash, err := getHash(f)
		if err != nil {
			return err
		}

		r := Result{
			Path: path,
			Hash: hash,
		}

		rs = append(rs, r)

		return nil
	})
	if err != nil {
		return ResultsInfo{}, err
	}

	return ResultsInfo{
		root:    root,
		results: rs,
	}, nil
}

func getHash(r io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func Run(source, target string) error {
	ri1, err := getResults(source)
	if err != nil {
		return err
	}

	ri2, err := getResults(target)
	if err != nil {
		return err
	}

	diff1, diff2 := lo.Difference(ri1.results, ri2.results)

	di := DiffInfo{
		Source: source,
		Target: target,
		Diff: Diff{
			Source: diff1,
			Target: diff2,
		},
	}

	b, err := json.MarshalIndent(di, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
