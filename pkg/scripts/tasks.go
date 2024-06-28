package scripts

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/action"
	"bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/util"
)

// ~ Greeting
type Greeting struct {
	action.BaseAction
}

func (t *Greeting) GetName() string {
	return "Greeting Script"
}

func (t *Greeting) Execute(runtime connector.Runtime) error {
	p := fmt.Sprintf("%s/%s/%s", constants.WorkDir, common.Scripts, common.GreetingShell)
	if ok := util.IsExist(p); ok {
		outstd, _, err := util.Exec(p, false, false)
		if err != nil {
			return err
		}
		logger.Debugf("script CMD: %s, OUTPUT: \n%s", p, outstd)
	}
	return nil
}

// ~ Copy
type Copy struct {
	action.BaseAction
}

func (t *Copy) GetName() string {
	return "Copy Scripts"
}

func (t *Copy) Execute(runtime connector.Runtime) error {
	p := fmt.Sprintf("%s/%s", constants.WorkDir, common.Scripts)
	if ok := util.IsExist(p); !ok {
		if err := util.CreateDir(p); err != nil {
			return err
		}
	}

	all := Assets()
	if err := fs.WalkDir(all, "files", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel("files", path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(p, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, os.ModePerm)
		}

		data, err := all.ReadFile(path)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(destPath, data, fs.ModePerm)
	}); err != nil {
		logger.Errorf("copy scripts files failed: %v", err)
		return err
	}

	return nil
}
