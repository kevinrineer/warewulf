package list

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/warewulf/warewulf/internal/pkg/api/routes/wwapiv1"
	"github.com/warewulf/warewulf/internal/pkg/testenv"
	"github.com/warewulf/warewulf/internal/pkg/warewulfd"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func Test_List(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		stdout   string
		inDb     string
		mockFunc func()
	}{
		{
			name: "image list test",
			args: []string{"-l"},
			stdout: `
IMAGE NAME  NODES  KERNEL VERSION  CREATION TIME        MODIFICATION TIME    SIZE
----------  -----  --------------  -------------        -----------------    ----
test        1      kernel          01 Jan 70 00:00 UTC  01 Jan 70 00:00 UTC  0 B
`,
			inDb: `
nodeprofiles:
  default: {}
nodes:
  n01:
    profiles:
    - default
`,
			mockFunc: func() {
				imageList = func() (imageInfo []*wwapiv1.ImageInfo, err error) {
					imageInfo = append(imageInfo, &wwapiv1.ImageInfo{
						Name:          "test",
						NodeCount:     1,
						KernelVersion: "kernel",
						CreateDate:    uint64(time.Unix(0, 0).Unix()),
						ModDate:       uint64(time.Unix(0, 0).Unix()),
					})
					return
				}
			},
		},
	}

	warewulfd.SetNoDaemon()
	for _, tt := range tests {
		env := testenv.New(t)
		defer env.RemoveAll()
		env.WriteFile("etc/warewulf/nodes.conf", tt.inDb)

		t.Logf("Running test: %s\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			buf := new(bytes.Buffer)
			baseCmd := GetCommand()
			baseCmd.SetArgs(tt.args)
			baseCmd.SetOut(buf)
			baseCmd.SetErr(buf)
			wwlog.SetLogWriter(buf)
			err := baseCmd.Execute()
			assert.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(tt.stdout), strings.TrimSpace(buf.String()))
		})
	}
}
