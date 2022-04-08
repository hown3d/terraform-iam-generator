package terraform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newCommand(t *testing.T) {
	type args struct {
		mode string
		opts Options
	}
	tests := []struct {
		name    string
		args    args
		want    *exec.Cmd
		wantErr bool
		t       *testing.T
	}{
		{
			name: "destroy",
			args: args{
				mode: "destroy",
				opts: Options{
					Directory: "test",
				},
			},
			want: &exec.Cmd{
				Path:   "terraform",
				Args:   []string{"terraform", "destroy", "-auto-approve"},
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			},
			wantErr: false,
		},
		{
			name: "apply",
			args: args{
				mode: "apply",
				opts: Options{
					Directory: "test",
				},
			},
			want: &exec.Cmd{
				Path:   "terraform",
				Args:   []string{"terraform", "apply", "-auto-approve"},
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			},
			wantErr: false,
		},
		{
			name: "apply with vars",
			args: args{
				mode: "apply",
				opts: Options{
					Directory: "test",
					Vars: []struct {
						Key   string
						Value string
					}{
						{
							Key:   "hello",
							Value: "world",
						},
						{
							Key:   "foo",
							Value: "bar",
						},
					},
				},
			},
			want: &exec.Cmd{
				Path:   "terraform",
				Args:   []string{"terraform", "apply", "-auto-approve", "-var", "'hello=world'", "-var", "'foo=bar'"},
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			},
			wantErr: false,
		},
		{
			name: "apply with var-files",
			args: args{
				mode: "apply",
				opts: Options{
					Directory: "test",
					VarsFiles: []string{"./helloWorld.tfvars", "./fooBar.tfvars"},
				},
			},
			want: &exec.Cmd{
				Path: "terraform",
				Args: []string{
					"terraform",
					"apply",
					"-auto-approve",
					fmt.Sprintf("-var-file=%s", "./helloWorld.tfvars"),
					fmt.Sprintf("-var-file=%s", "./fooBar.tfvars"),
				},
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t = t
			abs, err := filepath.Abs(tt.args.opts.Directory)
			if err != nil {
				t.Fatal(err)
			}
			tt.want.Dir = abs
			tt.want.Path = lookPath(t, tt.want.Path)
			got, err := newCommand(tt.args.mode, tt.args.opts)
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func lookPath(t *testing.T, binary string) string {
	path, err := exec.LookPath("terraform")
	if err != nil {
		t.Fatal(err)
	}
	return path
}
