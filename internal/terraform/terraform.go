package terraform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var environ = os.Environ

type Variable struct {
	Key   string
	Value string
}

func (v Variable) String() string {
	if v.Key == "" {
		return ""
	}
	return fmt.Sprintf("%v=%v", v.Key, v.Value)
}

type VariableFile string

func (v VariableFile) String() string {
	return string(v)
}

type Options struct {
	Directory      string
	VarsFiles      []VariableFile
	Vars           []Variable
	AdditionalArgs []string
	AutoApprove    bool
}

func Apply(opts Options) error {
	checkAutoApprove(&opts)
	runInit(opts)
	cmd, err := newCommand("apply", opts)
	if err != nil {
		return fmt.Errorf("creating terraform apply command: %w", err)
	}
	return cmd.Run()
}

func Destroy(opts Options) error {
	checkAutoApprove(&opts)
	cmd, err := newCommand("destroy", opts)
	if err != nil {
		return fmt.Errorf("creating terraform destroy command: %w", err)
	}
	return cmd.Run()
}
func checkAutoApprove(opts *Options) {
	if opts.AutoApprove {
		addAutoApprove(opts)
	}
}

func runInit(opts Options) error {
	cmd, err := newCommand("init", opts)
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
		cmd.Args = append(cmd.Args, "-var", fmt.Sprintf("%s=%s", v.Key, v.Value))
	}

	// activate aws client side monitoring to get all api calls
	cmd.Env = append(cmd.Env, "AWS_CSM_ENABLED=true")
	cmd.Env = append(cmd.Env, "AWS_CSM_PORT=31000")
	cmd.Env = append(cmd.Env, "AWS_CSM_HOST=localhost")
	cmd.Env = append(cmd.Env, environ()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}

func addAutoApprove(opts *Options) {
	opts.AdditionalArgs = append([]string{"-auto-approve"}, opts.AdditionalArgs...)
}
