package common

import (
	"fmt"
	"strconv"

	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
	"github.com/pkg/errors"
)

// ~ Stop
type Stop struct {
	prepare.BasePrepare
}

func (p *Stop) PreCheck(runtime connector.Runtime) (bool, error) {
	return true, nil
	// return false, fmt.Errorf("STOP !!!!!!")
}

// ~ GetCommandKubectl
type GetCommandKubectl struct {
	prepare.BasePrepare
}

func (p *GetCommandKubectl) PreCheck(runtime connector.Runtime) (bool, error) {
	cmd, err := runtime.GetRunner().Host.GetCommand(CommandKubectl)
	if err != nil {
		return true, nil
	}
	if cmd != "" {
		p.PipelineCache.Set(CacheKubectlKey, cmd)
	}
	return true, nil
}

type GetKubeVersion struct {
	prepare.BasePrepare
}

func (p *GetKubeVersion) PreCheck(runtime connector.Runtime) (bool, error) {
	return true, nil
}

// ~ GetMasterNum
type GetMasterNum struct {
	prepare.BasePrepare
}

func (p *GetMasterNum) PreCheck(runtime connector.Runtime) (bool, error) {
	var cmd = fmt.Sprintf("/usr/local/bin/kubectl get node | awk '{if(NR>1){print $3}}' | grep master | wc -l")
	var stdout, err = runtime.GetRunner().SudoCmd(cmd, false, false)
	if err != nil {
		return false, errors.Wrap(errors.WithStack(err), "get master num failed")
	}

	masterNum, _ := strconv.ParseInt(stdout, 10, 64)

	p.PipelineCache.Set(CacheMasterNum, masterNum)

	return true, nil
}

// ~ GetNodeNum
type GetNodeNum struct {
	prepare.BasePrepare
}

func (p *GetNodeNum) PreCheck(runtime connector.Runtime) (bool, error) {
	var cmd = fmt.Sprintf("/usr/local/bin/kubectl get node | wc -l")
	var stdout, err = runtime.GetRunner().SudoCmd(cmd, false, false)
	if err != nil {
		return false, errors.Wrap(errors.WithStack(err), "get node num failed")
	}

	nodeNum, _ := strconv.ParseInt(stdout, 10, 64)

	p.PipelineCache.Set(CacheNodeNum, nodeNum)

	return true, nil
}
