package difff

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
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
		root string
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
				root: tempDir,
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
				root: "non existent directory",
			},
			want:    ResultsInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getResults(tt.args.root)
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
		source string
		target string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Run",
			args: args{
				source: tempDir1,
				target: tempDir2,
			},
			wantErr: false,
		},
		{
			name: "Run error",
			args: args{
				source: "non existent directory",
				target: "non existent directory",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(tt.args.source, tt.args.target); (err != nil) != tt.wantErr {
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
