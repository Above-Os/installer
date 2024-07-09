package scripts

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
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

func (t *Copy) Execute(runtime connector.Runtime) error {
	p := path.Join(runtime.GetRootDir(), common.Scripts)
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

// ~ CopyUninstallScriptTask
type CopyUninstallScriptTask struct {
	action.BaseAction
}

func (t *CopyUninstallScriptTask) Execute(runtime connector.Runtime) error {
	dest := path.Join(runtime.GetPackageDir(), common.InstallDir)

	if ok := util.IsExist(dest); !ok {
		return fmt.Errorf("directory %s not exist", dest)
	}

	all := Assets()
	fileContent, err := all.ReadFile(path.Join("files", common.UninstallOsScript))
	if err != nil {
		return fmt.Errorf("read file %s failed: %v", common.UninstallOsScript, err)
	}

	dstFile := path.Join(dest, common.UninstallOsScript)
	err = ioutil.WriteFile(dstFile, fileContent, common.FileMode0755)
	if err != nil {
		log.Fatalf("failed to write file %s to target directory: %v", dstFile, err)
	}

	return nil
}
