package terraform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Options struct {
	Env       []string
	Directory string
	VarsFiles []string
	Vars      []struct {
		Key   string
		Value string
	}
	AdditionalArgs []string
}

func Apply(opts Options) error {
	addAutoApprove(&opts)
	runInit(opts.Directory)
	cmd, err := newCommand("apply", opts)
	if err != nil {
		return fmt.Errorf("creating terraform apply command: %w", err)
	}
	return cmd.Run()
}

func Destroy(opts Options) error {
	addAutoApprove(&opts)
	cmd, err := newCommand("destroy", opts)
	if err != nil {
		return fmt.Errorf("creating terraform destroy command: %w", err)
	}
	return cmd.Run()
}

func runInit(dir string) error {
	cmd, err := newCommand("init", Options{Directory: dir})
	if err != nil {
		return fmt.Errorf("creating terraform init command: %w", err)
	}
	return cmd.Run()
}

func newCommand(mode string, opts Options) (*exec.Cmd, error) {
	args := append([]string{mode}, opts.AdditionalArgs...)
	cmd := exec.Command("terraform", args...)
	if opts.Directory != "" {
		dir, err := filepath.Abs(opts.Directory)
		if err != nil {
			return nil, fmt.Errorf("getting fullpath for %s: %w", opts.Directory, err)
		}
		cmd.Dir = dir
	}
	for _, file := range opts.VarsFiles {
		cmd.Args = append(cmd.Args, fmt.Sprintf("-var-file=%s", file))
	}
	for _, v := range opts.Vars {
		cmd.Args = append(cmd.Args, "-var", fmt.Sprintf("'%s=%s'", v.Key, v.Value))
	}
	cmd.Env = opts.Env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}

func addAutoApprove(opts *Options) {
	opts.AdditionalArgs = append([]string{"-auto-approve"}, opts.AdditionalArgs...)
}
