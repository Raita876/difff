package difff

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
)

func Test_getResults(t *testing.T) {
	tempDir := t.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "temp")
	if err != nil {
		t.Fatal("failed create temp file.")
	}
	fileName, err := filepath.Rel(tempDir, tempFile.Name())
	if err != nil {
		t.Fatal("failed get file name.")
	}

	type args struct {
		root            string
		excludePatterns []string
	}
	tests := []struct {
		name    string
		args    args
		want    ResultsInfo
		wantErr bool
	}{
		{
			name: "getResults",
			args: args{
				root:            tempDir,
				excludePatterns: []string{},
			},
			want: ResultsInfo{
				root: tempDir,
				results: Results{
					Result{
						Path: fileName,
						Hash: "d41d8cd98f00b204e9800998ecf8427e",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "getResults error",
			args: args{
				root:            "non existent directory",
				excludePatterns: []string{},
			},
			want:    ResultsInfo{},
			wantErr: true,
		},
		{
			name: "getResults with excludePatterns",
			args: args{
				root:            tempDir,
				excludePatterns: []string{"^.*temp.*$"},
			},
			want: ResultsInfo{
				root:    tempDir,
				results: Results{},
			},
			wantErr: false,
		},
		{
			name: "getResults error with excludePatterns",
			args: args{
				root:            tempDir,
				excludePatterns: []string{"hoge(fuga"}, // This is not a valid regular expression.
			},
			want:    ResultsInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getResults(tt.args.root, tt.args.excludePatterns)
			if (err != nil) != tt.wantErr {
				t.Errorf("getResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getResults() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getHash(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "getHash",
			args: args{
				r: strings.NewReader("Hello World!"),
			},
			want:    "ed076287532e86365e841e92bfc50d8c",
			wantErr: false,
		},
		{
			name: "getHash error",
			args: args{
				r: iotest.ErrReader(fmt.Errorf("io error")),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getHash(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tempDir1 := t.TempDir()
	tempDir2 := t.TempDir()

	type args struct {
		source          string
		target          string
		ft              FormatType
		excludePatterns []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Run",
			args: args{
				source:          tempDir1,
				target:          tempDir2,
				ft:              JSON,
				excludePatterns: []string{},
			},
			wantErr: false,
		},
		{
			name: "Run error (source is non existent directory)",
			args: args{
				source:          "non existent directory",
				target:          tempDir2,
				ft:              JSON,
				excludePatterns: []string{},
			},
			wantErr: true,
		},
		{
			name: "Run error (target is non existent directory)",
			args: args{
				source:          tempDir1,
				target:          "non existent directory",
				ft:              JSON,
				excludePatterns: []string{},
			},
			wantErr: true,
		},
		{
			name: "Run error (UNKNOWN format)",
			args: args{
				source:          tempDir1,
				target:          tempDir2,
				ft:              "UNKNOWN",
				excludePatterns: []string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(tt.args.source, tt.args.target, tt.args.ft, tt.args.excludePatterns); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_countFiles(t *testing.T) {
	tempDir := t.TempDir()
	for range 3 {
		_, err := os.CreateTemp(tempDir, "temp")
		if err != nil {
			t.Fatal("failed create temp file.")
		}
	}

	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "countFiles",
			args: args{
				dir: tempDir,
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "countFiles",
			args: args{
				dir: "non existent directory",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := countFiles(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("countFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("countFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_run(t *testing.T) {
	tempDir1 := t.TempDir()
	tempFile1, err := os.CreateTemp(tempDir1, "temp")
	if err != nil {
		t.Fatal("failed create temp file.")
	}
	fileName1, err := filepath.Rel(tempDir1, tempFile1.Name())
	if err != nil {
		t.Fatal("failed get file name.")
	}

	tempDir2 := t.TempDir()
	tempFile2, err := os.CreateTemp(tempDir2, "temp")
	if err != nil {
		t.Fatal("failed create temp file.")
	}
	fileName2, err := filepath.Rel(tempDir2, tempFile2.Name())
	if err != nil {
		t.Fatal("failed get file name.")
	}

	type args struct {
		source          string
		target          string
		ft              FormatType
		excludePatterns []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "run",
			args: args{
				source:          tempDir1,
				target:          tempDir2,
				ft:              JSON,
				excludePatterns: []string{},
			},
			want: fmt.Sprintf(`{
  "source": {
    "path": "%s",
    "num": 1
  },
  "target": {
    "path": "%s",
    "num": 1
  },
  "exclude": [],
  "diff": {
    "source": {
      "num": 1,
      "results": [
        {
          "path": "%s",
          "hash": "d41d8cd98f00b204e9800998ecf8427e"
        }
      ]
    },
    "target": {
      "num": 1,
      "results": [
        {
          "path": "%s",
          "hash": "d41d8cd98f00b204e9800998ecf8427e"
        }
      ]
    }
  }
}`, tempDir1, tempDir2, fileName1, fileName2,
			),
			wantErr: false,
		},
		{
			name: "run",
			args: args{
				source:          tempDir1,
				target:          tempDir2,
				ft:              YAML,
				excludePatterns: []string{},
			},
			want: fmt.Sprintf(`source:
  path: %s
  num: 1
target:
  path: %s
  num: 1
exclude: []
diff:
  source:
    num: 1
    results:
    - path: %s
      hash: d41d8cd98f00b204e9800998ecf8427e
  target:
    num: 1
    results:
    - path: %s
      hash: d41d8cd98f00b204e9800998ecf8427e
`, tempDir1, tempDir2, fileName1, fileName2,
			),
			wantErr: false,
		},
		{
			name: "run",
			args: args{
				source:          tempDir1,
				target:          tempDir2,
				ft:              XML,
				excludePatterns: []string{},
			},
			want: fmt.Sprintf(`<DiffResponse>
  <source>
    <path>%s</path>
    <num>1</num>
  </source>
  <target>
    <path>%s</path>
    <num>1</num>
  </target>
  <diff>
    <source>
      <num>1</num>
      <results>
        <path>%s</path>
        <hash>d41d8cd98f00b204e9800998ecf8427e</hash>
      </results>
    </source>
    <target>
      <num>1</num>
      <results>
        <path>%s</path>
        <hash>d41d8cd98f00b204e9800998ecf8427e</hash>
      </results>
    </target>
  </diff>
</DiffResponse>`, tempDir1, tempDir2, fileName1, fileName2,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := run(tt.args.source, tt.args.target, tt.args.ft, tt.args.excludePatterns)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_marshal(t *testing.T) {
	type args struct {
		dr *DiffResponse
		ft FormatType
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "marshal",
			args: args{
				dr: &DiffResponse{
					Source: Dir{
						Path: "source_path",
						Num:  1,
					},
					Target: Dir{
						Path: "target_path",
						Num:  1,
					},
					Exclude: []string{},
					Diff: Diff{
						Source: DiffInfo{
							Num:     0,
							Results: []Result{},
						},
						Target: DiffInfo{
							Num:     0,
							Results: []Result{},
						},
					},
				},
				ft: JSON,
			},
			want: []byte(`{
  "source": {
    "path": "source_path",
    "num": 1
  },
  "target": {
    "path": "target_path",
    "num": 1
  },
  "exclude": [],
  "diff": {
    "source": {
      "num": 0,
      "results": []
    },
    "target": {
      "num": 0,
      "results": []
    }
  }
}`),
			wantErr: false,
		},
		{
			name: "marshal",
			args: args{
				dr: &DiffResponse{},
				ft: "UNKNONW",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := marshal(tt.args.dr, tt.args.ft)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("marshal() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
