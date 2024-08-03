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

type DiffResponse struct {
	Source Dir  `json:"source"`
	Target Dir  `json:"target"`
	Diff   Diff `json:"diff"`
}

type Dir struct {
	Path string `json:"path"`
	Num  uint64 `json:"num"`
}

type Diff struct {
	Source DiffInfo `json:"source"`
	Target DiffInfo `json:"target"`
}

type DiffInfo struct {
	Num     uint64  `json:"num"`
	Results Results `json:"results"`
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

func countFiles(dir string) (uint64, error) {
	var count uint64 = 0

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			count++
		}

		return nil
	})

	return count, err
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

	count1, err := countFiles(source)
	if err != nil {
		return err
	}

	count2, err := countFiles(target)
	if err != nil {
		return err
	}

	di := DiffResponse{
		Source: Dir{
			Path: source,
			Num:  count1,
		},
		Target: Dir{
			Path: target,
			Num:  count2,
		},
		Diff: Diff{
			Source: DiffInfo{
				Num:     uint64(len(diff1)),
				Results: diff1,
			},
			Target: DiffInfo{
				Num:     uint64(len(diff2)),
				Results: diff2,
			},
		},
	}

	b, err := json.MarshalIndent(di, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
