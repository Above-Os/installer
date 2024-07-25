/*
 Copyright 2021 The KubeSphere Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package connector

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/core/cache"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/util"
)

type BaseHost struct {
	Name            string `yaml:"name,omitempty" json:"name,omitempty"`
	Address         string `yaml:"address,omitempty" json:"address,omitempty"`
	InternalAddress string `yaml:"internalAddress,omitempty" json:"internalAddress,omitempty"`
	Port            int    `yaml:"port,omitempty" json:"port,omitempty"`
	User            string `yaml:"user,omitempty" json:"user,omitempty"`
	Password        string `yaml:"password,omitempty" json:"password,omitempty"`
	PrivateKey      string `yaml:"privateKey,omitempty" json:"privateKey,omitempty"`
	PrivateKeyPath  string `yaml:"privateKeyPath,omitempty" json:"privateKeyPath,omitempty"`
	Arch            string `yaml:"arch,omitempty" json:"arch,omitempty"`
	Timeout         int64  `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Minikube        bool   `yaml:"minikube,omitempty" json:"minikube,omitempty"`
	MiniKubeProfile string `json:"minikubeProfileName,omitempty" json:"minikubeProfileName,omitempty"`

	Roles     []string        `json:"-"`
	RoleTable map[string]bool `json:"-"`
	Cache     *cache.Cache    `json:"-"`
}

func NewHost() *BaseHost {
	return &BaseHost{
		Roles:     make([]string, 0, 0),
		RoleTable: make(map[string]bool),
		Cache:     cache.NewCache(),
	}
}

func (b *BaseHost) GetName() string {
	return b.Name
}

func (b *BaseHost) SetName(name string) {
	b.Name = name
}

func (b *BaseHost) GetAddress() string {
	return b.Address
}

func (b *BaseHost) SetAddress(str string) {
	b.Address = str
}

func (b *BaseHost) GetInternalAddress() string {
	return b.InternalAddress
}

func (b *BaseHost) SetInternalAddress(str string) {
	b.InternalAddress = str
}

func (b *BaseHost) GetPort() int {
	return b.Port
}

func (b *BaseHost) SetPort(port int) {
	b.Port = port
}

func (b *BaseHost) GetUser() string {
	return b.User
}

func (b *BaseHost) SetUser(u string) {
	b.User = u
}

func (b *BaseHost) GetPassword() string {
	return b.Password
}

func (b *BaseHost) SetPassword(password string) {
	b.Password = password
}

func (b *BaseHost) GetPrivateKey() string {
	return b.PrivateKey
}

func (b *BaseHost) SetPrivateKey(privateKey string) {
	b.PrivateKey = privateKey
}

func (b *BaseHost) GetPrivateKeyPath() string {
	return b.PrivateKeyPath
}

func (b *BaseHost) SetPrivateKeyPath(path string) {
	b.PrivateKeyPath = path
}

func (b *BaseHost) GetArch() string {
	return b.Arch
}

func (b *BaseHost) SetArch(arch string) {
	b.Arch = arch
}

func (b *BaseHost) SetMinikube(minikube bool) {
	b.Minikube = minikube
}
func (b *BaseHost) GetMinikube() bool {
	return b.Minikube
}
func (b *BaseHost) SetMinikubeProfile(profile string) {
	b.MiniKubeProfile = profile
}
func (b *BaseHost) GetMinikubeProfile() string {
	return b.MiniKubeProfile
}

func (b *BaseHost) GetTimeout() int64 {
	return b.Timeout
}

func (b *BaseHost) SetTimeout(timeout int64) {
	b.Timeout = timeout
}

func (b *BaseHost) GetRoles() []string {
	return b.Roles
}

func (b *BaseHost) SetRoles(roles []string) {
	b.Roles = roles
}

func (b *BaseHost) SetRole(role string) {
	b.RoleTable[role] = true
	b.Roles = append(b.Roles, role)
}

func (b *BaseHost) IsRole(role string) bool {
	if res, ok := b.RoleTable[role]; ok {
		return res
	}
	return false
}

func (b *BaseHost) GetCache() *cache.Cache {
	return b.Cache
}

func (b *BaseHost) SetCache(c *cache.Cache) {
	b.Cache = c
}

func (b *BaseHost) Echo() {
}

func (b *BaseHost) Exec(cmd string, printOutput bool, printLine bool) (stdout string, code int, err error) {
	stdout, code, err = util.Exec(cmd, printOutput, printLine)
	if err != nil {
		logger.Errorf("[exec] %s CMD: %s, ERROR: %s", b.GetName(), cmd, err)
	}

	if printOutput {
		logger.Debugf("[exec] %s CMD: %s, OUTPUT: \n%s", b.GetName(), cmd, stdout)
	}

	logger.Infof("[exec] %s CMD: %s, OUTPUT: %s", b.GetName(), cmd, stdout)

	return stdout, code, err
}

func (b *BaseHost) Cmd(cmd string, printOutput bool, printLine bool) (string, error) {
	stdout, _, err := b.Exec(cmd, printOutput, printLine)
	if err != nil {
		return stdout, err
	}
	return stdout, nil
}

func (b *BaseHost) CmdExt(cmd string, printOutput bool, printLine bool) (string, error) {
	stdout, _, err := b.Exec(cmd, printOutput, printLine)
	if printOutput {
		logger.Debugf("[exec] %s CMD: %s, OUTPUT: \n%s", b.GetName(), cmd, stdout)
	}

	logger.Infof("[exec] %s CMD: %s, OUTPUT: %s", b.GetName(), cmd, stdout)

	return stdout, err
}

func (b *BaseHost) Scp(local, remote string) error {
	var cmd = fmt.Sprintf("minikube -p %s cp %s %s", b.GetMinikubeProfile(), local, remote)
	_, _, err := b.Exec(cmd, false, false)
	return err
}
