package runner

import (
	"reflect"
	"testing"
	"time"
)

func TestNewCommand(t *testing.T) {
	to := time.Second * 1
	cmd := NewCommand("/bin/ls", []string{"-l", "-a"}, "/tmp", map[string]string{"TESTING": "TRUE", "CI": "false"}, to)

	if cmd.path != "/bin/ls" {
		t.Errorf("Unexpected path, expected: %+v, got: %+v", "/bin/ls", cmd.path)
	}

	if !reflect.DeepEqual([]string{"-l", "-a"}, cmd.arg) {
		t.Errorf("Unexpected arg, expected: %+v, got: %+v", []string{"-l", "-a"}, cmd.arg)
	}

	if cmd.workdir != "/tmp" {
		t.Errorf("Unexpected workdir, expected: %+v, got: %+v", "/tmp", cmd.workdir)
	}

	if !reflect.DeepEqual(cmd.envSlice(), []string{"TESTING=TRUE", "CI=false"}) {
		t.Errorf("Unexpected env variables, expected: %+v, got: %+v", []string{"TESTING=TRUE", "CI=false"}, cmd.envSlice())
	}

	if cmd.timeout != to {
		t.Errorf("Unexpected Timeout value, expected: %+v, got: %v", to, cmd.timeout)
	}
}

func TestCommand_envSlice(t *testing.T) {
	type fields struct {
		env      map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"nil", fields{env: nil}, nil},
		{"empty", fields{env: map[string]string{}}, nil},
		{"two", fields{env:map[string]string{"TESTING": "true", "CI": "false"}}, []string{"TESTING=true", "CI=false"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{
				path:     "/bin/ls",
				env:      tt.fields.env,
			}

			if got := c.envSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Command.envSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
