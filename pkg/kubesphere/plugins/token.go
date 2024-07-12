package plugins

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/utils"
	"github.com/pkg/errors"
)

// ~ GenerateKubeSphereToken
type GenerateKubeSphereToken struct {
	common.KubeAction
}

func (t *GenerateKubeSphereToken) Execute(runtime connector.Runtime) error {
	var random, err = utils.GeneratePassword(32)
	if err != nil {
		logger.Errorf("failed to generate password: %v", err)
		return err
	}

	token, err := util.EncryptToken(random)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "create kubesphere token failed")
	}

	fmt.Println("---1---", random)
	fmt.Println("---2---", token)

	return nil
}

// +++++

// ~ CreateKubeSphereSecretModule
type CreateKubeSphereSecretModule struct {
	common.KubeModule
}

func (m *CreateKubeSphereSecretModule) Init() {
	m.Name = "CreateKubeSphereSecret"

	generateKubeSphereToken := &task.RemoteTask{
		Name:  "GenerateKubeSphereToken",
		Hosts: m.Runtime.GetHostsByRole(common.Master),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
			new(NotEqualDesiredVersion),
		},
		Action:   new(GenerateKubeSphereToken),
		Parallel: false,
		Retry:    0,
	}

	m.Tasks = []task.Interface{generateKubeSphereToken}
}
