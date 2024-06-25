package action

import (
	"embed"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/util"
)

var scripts embed.FS

func Assets() embed.FS {
	return scripts
}

type Script struct {
	BaseAction
	Name string
	Args []string
}

func (s *Script) Execute(runtime connector.Runtime) error {
	logger.Debugf("[action] Script: %s", s.Name)
	var cmd = fmt.Sprintf("bash %s %s", s.Name, strings.Join(s.Args, " "))
	_, _, err := util.Exec(cmd, false)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), fmt.Sprintf("exec script %s failed, args: %v", s.Name, s.Args))
	}

	return nil
}
