package action

import (
	"embed"
	"fmt"
	"path"
	"strings"

	"github.com/pkg/errors"

	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/util"
)

var scripts embed.FS

func Assets() embed.FS {
	return scripts
}

type Script struct {
	BaseAction
	Name        string
	File        string
	Args        []string
	PrintOutput bool
}

func (s *Script) GetName() string {
	return "script"
}

func (s *Script) Execute(runtime connector.Runtime) error {
	scriptFileName := path.Join(constants.WorkDir, common.Scripts, s.File)
	if !util.IsExist(scriptFileName) {
		return errors.New(fmt.Sprintf("script file %s not exist", s.File))
	}
	var cmd = fmt.Sprintf("bash %s %s", scriptFileName, strings.Join(s.Args, " "))
	_, _, err := util.Exec(cmd, s.PrintOutput)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), fmt.Sprintf("exec script %s failed, args: %v", s.File, s.Args))
	}

	return nil
}
