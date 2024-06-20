package runuser

import (
	"errors"
	"fmt"
	"os/user"

	"bytetrade.io/web3os/installer/pkg/core/ending"
	"bytetrade.io/web3os/installer/pkg/core/module"
)

type RunUserCheckHook struct {
	Module module.Module
	Result *ending.ModuleResult
}

func (h *RunUserCheckHook) Init(module module.Module, result *ending.ModuleResult) {
	h.Module = module
	h.Result = result
}

func (h *RunUserCheckHook) Try() error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	if u.Username != "root" {
		return errors.New(fmt.Sprintf("Current user is %s. Please use root!", u.Username))
	}
	return nil
}

func (h *RunUserCheckHook) Catch(err error) error {
	return err
}

func (h *RunUserCheckHook) Finally() {
}
