/*
 Copyright 2020 The KubeSphere Authors.

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

package files

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/log"
	"bytetrade.io/web3os/installer/pkg/utils"
	"github.com/cavaliergopher/grab/v3"
	"github.com/pkg/errors"
)

type InstallerPackage struct {
	Type     string
	ID       string // todo 这个字段貌似没啥用？
	FileName string
	Os       string
	Arch     string
	Version  string
	Url      string
	BaseDir  string
	getCmd   func(path, url string) string
}

// + 需要整合第三方的下载库，这样才能监控下载进度
// todo
// todo 估计需要约定一个下载目录，比如 terminus-installer
// todo 也可能需要结合 terminus.sh 的目录，在那里面来操作
// todo 需要判断内网还是外网？
func NewInstallerPackage(name, os, arch, version, prePath string) *InstallerPackage {
	component := new(InstallerPackage)
	component.ID = name
	component.Os = os
	component.Arch = arch
	component.Version = version
	// component.getCmd = getCmd

	// install-wizard-linux-amd64-v1.0.0-full.tar.gz
	// install-wizard-linux-arm64-v1.0.0-full.tar.gz
	// install-wizard-linux-amd64-v1.0.0-mini.tar.gz
	// install-wizard-linux-arm64-v1.0.0-mini.tar.gz
	// install-wizard-{os}-{arch}-v{version}-{full|mini}.tar.gz
	fileName := fmt.Sprintf("%s-%s-%s-v%s", common.DefaultPackageName, os, arch, version)

	// todo 暂定两种包；后面肯定还要区分平台，版本
	switch name {
	case common.MiniPackage:
		component.Type = common.MiniPackage
		component.FileName = fmt.Sprintf("%s-%s.tar.gz", fileName, common.DefaultMiniPackageName)
		component.Url = fmt.Sprintf("") // todo 需要考虑内、外网的场景
	case common.FullPackage:
		component.Type = common.MiniPackage
		component.FileName = fmt.Sprintf("%s-%s.tar.gz", fileName, common.DefaultFullPackageName)
		component.Url = fmt.Sprintf("") // todo 需要考虑内、外网的场景
	default:
		log.Fatalf("unsupported installer packages %s", name)
	}

	if component.BaseDir == "" {
		component.BaseDir = filepath.Join(prePath, component.Type, component.Version, component.Arch)
	}

	component.Url = "http://192.168.50.32/install-wizard-full.tar.gz"

	return component
}

func (p *InstallerPackage) CreateBaseDir() error {
	if err := util.CreateDir(p.BaseDir); err != nil {
		return err
	}
	return nil
}

func (p *InstallerPackage) Path() string {
	return filepath.Join(p.BaseDir, p.FileName)
}

func (p *InstallerPackage) Download() error {
	return nil
}

// - split -

const (
	kubeadm    = "kubeadm"
	kubelet    = "kubelet"
	kubectl    = "kubectl"
	kubecni    = "kubecni"
	etcd       = "etcd"
	helm       = "helm"
	amd64      = "amd64"
	arm64      = "arm64"
	k3s        = "k3s"
	k8e        = "k8e"
	docker     = "docker"
	crictl     = "crictl"
	registry   = "registry"
	harbor     = "harbor"
	compose    = "compose"
	containerd = "containerd"
	runc       = "runc"

	// todo 安装包会进行拆分，可能不会再有 full 包了
	// todo 所以我可以假设 f1.tar.gz f2.tar.gz f3.tar.gz ...
	file1 = "file1"
	file2 = "file2"
	file3 = "file3"
)

// KubeBinary Type field const
const (
	CNI        = "cni"
	CRICTL     = "crictl"
	DOCKER     = "docker"
	ETCD       = "etcd"
	HELM       = "helm"
	KUBE       = "kube"
	REGISTRY   = "registry"
	CONTAINERD = "containerd"
	RUNC       = "runc"
	// todo installer package
	INSTALLER = "installer"
)

type KubeBinary struct {
	Type     string
	ID       string
	FileName string
	Arch     string
	Version  string
	Url      string
	BaseDir  string
	Zone     string
	getCmd   func(path, url string) string
}

func NewKubeBinary(name, arch, version, prePath string, getCmd func(path, url string) string) *KubeBinary {
	component := new(KubeBinary)
	component.ID = name
	component.Arch = arch
	component.Version = version
	component.Zone = os.Getenv("KKZONE")
	component.getCmd = getCmd

	switch name {
	case etcd:
		component.Type = ETCD
		component.FileName = fmt.Sprintf("etcd-%s-linux-%s.tar.gz", version, arch)
		component.Url = fmt.Sprintf("https://github.com/coreos/etcd/releases/download/%s/etcd-%s-linux-%s.tar.gz", version, version, arch)
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf(
				"https://kubernetes-release.pek3b.qingstor.com/etcd/release/download/%s/etcd-%s-linux-%s.tar.gz",
				component.Version, component.Version, component.Arch)
		}
	case kubeadm:
		component.Type = KUBE
		component.FileName = kubeadm
		component.Url = fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/%s/kubeadm", version, arch)
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/release/%s/bin/linux/%s/kubeadm", version, arch)
		}
	case kubelet:
		component.Type = KUBE
		component.FileName = kubelet
		component.Url = fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/%s/kubelet", version, arch)
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/release/%s/bin/linux/%s/kubelet", version, arch)
		}
	case kubectl:
		component.Type = KUBE
		component.FileName = kubectl
		component.Url = fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/%s/kubectl", version, arch)
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/release/%s/bin/linux/%s/kubectl", version, arch)
		}
	case kubecni:
		component.Type = CNI
		component.FileName = fmt.Sprintf("cni-plugins-linux-%s-%s.tgz", arch, version)
		component.Url = fmt.Sprintf("https://github.com/containernetworking/plugins/releases/download/%s/cni-plugins-linux-%s-%s.tgz", version, arch, version)
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf("https://containernetworking.pek3b.qingstor.com/plugins/releases/download/%s/cni-plugins-linux-%s-%s.tgz", version, arch, version)
		}
	case helm:
		component.Type = HELM
		component.FileName = helm
		component.Url = fmt.Sprintf("https://get.helm.sh/helm-%s-linux-%s.tar.gz", version, arch)
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf("https://kubernetes-helm.pek3b.qingstor.com/linux-%s/%s/helm", arch, version)
		}
	case docker:
		component.Type = DOCKER
		component.FileName = fmt.Sprintf("docker-%s.tgz", version)
		component.Url = fmt.Sprintf("https://download.docker.com/linux/static/stable/%s/docker-%s.tgz", utils.ArchAlias(arch), version)
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf("https://mirrors.aliyun.com/docker-ce/linux/static/stable/%s/docker-%s.tgz", utils.ArchAlias(arch), version)
		}
	case crictl:
		component.Type = CRICTL
		component.FileName = fmt.Sprintf("crictl-%s-linux-%s.tar.gz", version, arch)
		component.Url = fmt.Sprintf("https://github.com/kubernetes-sigs/cri-tools/releases/download/%s/crictl-%s-linux-%s.tar.gz", version, version, arch)
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/cri-tools/releases/download/%s/crictl-%s-linux-%s.tar.gz", version, version, arch)
		}
	case k3s:
		component.Type = KUBE
		component.FileName = k3s
		component.Url = fmt.Sprintf("https://github.com/k3s-io/k3s/releases/download/%s+k3s1/k3s", version)
		if arch == arm64 {
			component.Url = fmt.Sprintf("https://github.com/k3s-io/k3s/releases/download/%s+k3s1/k3s-%s", version, arch)
		}
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/k3s/releases/download/%s+k3s1/linux/%s/k3s", version, arch)
		}
	case k8e:
		component.Type = KUBE
		component.FileName = k8e
		component.Url = fmt.Sprintf("https://github.com/xiaods/k8e/releases/download/%s+k8e2/k8e", version)
		if arch == arm64 {
			component.Url = fmt.Sprintf("https://github.com/xiaods/k8e/releases/download/%s+k8e2/k8e-%s", version, arch)
		}
	case registry:
		component.Type = REGISTRY
		component.FileName = fmt.Sprintf("registry-%s-linux-%s.tar.gz", version, arch)
		component.Url = fmt.Sprintf("https://github.com/kubesphere/kubekey/releases/download/v2.0.0-alpha.1/registry-%s-linux-%s.tar.gz", version, arch)
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/registry/%s/registry-%s-linux-%s.tar.gz", version, version, arch)
		}
		component.BaseDir = filepath.Join(prePath, component.Type, component.ID, component.Version, component.Arch)
	case harbor:
		component.Type = REGISTRY
		component.FileName = fmt.Sprintf("harbor-offline-installer-%s.tgz", version)
		component.Url = fmt.Sprintf("https://github.com/goharbor/harbor/releases/download/%s/harbor-offline-installer-%s.tgz", version, version)
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/harbor/releases/download/%s/harbor-offline-installer-%s.tgz", version, version)
		}
		component.BaseDir = filepath.Join(prePath, component.Type, component.ID, component.Version, component.Arch)
	case compose:
		component.Type = REGISTRY
		component.FileName = "docker-compose-linux-x86_64"
		component.Url = fmt.Sprintf("https://github.com/docker/compose/releases/download/%s/docker-compose-linux-x86_64", version)
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/docker/compose/releases/download/%s/docker-compose-linux-x86_64", version)
		}
		component.BaseDir = filepath.Join(prePath, component.Type, component.ID, component.Version, component.Arch)
	case containerd:
		component.Type = CONTAINERD
		component.FileName = fmt.Sprintf("containerd-%s-linux-%s.tar.gz", version, arch)
		component.Url = fmt.Sprintf("https://github.com/containerd/containerd/releases/download/v%s/containerd-%s-linux-%s.tar.gz", version, version, arch)
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/containerd/containerd/releases/download/v%s/containerd-%s-linux-%s.tar.gz", version, version, arch)
		}
	case runc:
		component.Type = RUNC
		component.FileName = fmt.Sprintf("runc.%s", arch)
		component.Url = fmt.Sprintf("https://github.com/opencontainers/runc/releases/download/%s/runc.%s", version, arch)
		if component.Zone == "cn" {
			component.Url = fmt.Sprintf("https://kubernetes-release.pek3b.qingstor.com/opencontainers/runc/releases/download/%s/runc.%s", version, arch)
		}
	case file1: // todo 这里是新增的文件类型，后面应该都是跟安装包有关
		// todo 要区分内外网了，公网肯定是走 CDN
		component.Type = INSTALLER
		component.FileName = fmt.Sprintf("file1_%s_v%s.tar.gz", arch, version)
		component.Url = "http://192.168.50.32/containerd/containerd/releases/download/v1.6.4/containerd-1.6.4-linux-amd64.tar.gz"
	case file2:
		component.Type = INSTALLER
		component.FileName = fmt.Sprintf("file2_%s_v%s.tar.gz", arch, version)
		component.Url = "http://192.168.50.32/coreos/etcd/releases/download/v3.4.13/etcd-v3.4.13-linux-amd64.tar.gz"
	case file3:
		component.Type = INSTALLER
		component.FileName = fmt.Sprintf("file3_%s_v%s.tar.gz", arch, version)
		component.Url = "http://192.168.50.32/kubernetes-release/release/v1.22.10/bin/linux/amd64/kubelet"
	default:
		log.Fatalf("unsupported kube binaries %s", name)
	}

	if component.BaseDir == "" {
		component.BaseDir = filepath.Join(prePath, component.Type, component.Version, component.Arch)
	}

	return component
}

func (b *KubeBinary) CreateBaseDir() error {
	if err := utils.CreateDir(b.BaseDir); err != nil {
		return err
	}
	return nil
}

func (b *KubeBinary) Path() string {
	return filepath.Join(b.BaseDir, b.FileName)
}

func (b *KubeBinary) GetCmd() string {
	cmd := b.getCmd(b.Path(), b.Url)

	if b.ID == helm && b.Zone != "cn" {
		get := b.getCmd(filepath.Join(b.BaseDir, fmt.Sprintf("helm-%s-linux-%s.tar.gz", b.Version, b.Arch)), b.Url)
		cmd = fmt.Sprintf("%s && cd %s && tar -zxf helm-%s-linux-%s.tar.gz && mv linux-%s/helm . && rm -rf *linux-%s*",
			get, b.BaseDir, b.Version, b.Arch, b.Arch, b.Arch)
	}
	return cmd
}

func (b *KubeBinary) GetSha256() string {
	s := FileSha256[b.ID][b.Arch][b.Version]
	return s
}

func (b *KubeBinary) GetFileSize() (int64, error) {
	resp, err := http.Head(b.Url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad status: %s", resp.Status)
	}

	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (b *KubeBinary) DownloadEx() error {
	for i := 5; i > 0; i-- {
		totalSize, err := b.GetFileSize()
		if err != nil {
			return err
		}

		client := grab.NewClient()
		req, _ := grab.NewRequest(fmt.Sprintf("%s/%s", b.BaseDir, b.FileName), b.Url)
		req.RateLimiter = NewLimiter(1024 * 4096) // ! debug
		req.HTTPRequest = req.HTTPRequest.WithContext(context.Background())
		ctx, cancel := context.WithTimeout(req.HTTPRequest.Context(), 5*time.Minute)
		defer cancel()

		req.HTTPRequest = req.HTTPRequest.WithContext(ctx)
		resp := client.Do(req)

		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()

	Loop:
		for {
			select {
			case <-t.C:
				downloaded := resp.BytesComplete()
				result := float64(downloaded) / float64(totalSize)
				fmt.Printf("  transferred %d / %d bytes (%.2f%%)\n",
					resp.BytesComplete(),
					totalSize,
					math.Round(result*10000)/100)
			case <-resp.Done:
				break Loop
			}
		}

		if err := resp.Err(); err != nil {
			// fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
			logger.Errorf("Download failed: %v", err)
			if i == 1 {
				log.Error("All download attempts failed")
				return err
			}
			time.Sleep(2 * time.Second)
			continue
		}

		if err := b.SHA256Check(); err != nil {
			log.Errorf("SHA256 check failed: %v", err)
			if i == 1 {
				return err
			}
			path := b.Path()
			_ = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s", path)).Run()
			time.Sleep(2 * time.Second)
			continue
		}

		log.Info("Download succeeded")
		break
	}

	return nil

	// 	totalSize, err := b.GetFileSize()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	client := grab.NewClient()
	// 	req, _ := grab.NewRequest(fmt.Sprintf("%s/%s", b.BaseDir, b.FileName), b.Url)

	// 	req.RateLimiter = NewLimiter(1024 * 2048) // ! debug

	// 	resp := client.Do(req)

	// 	t := time.NewTicker(500 * time.Millisecond)
	// 	defer t.Stop()

	// Loop:
	// 	for {
	// 		select {
	// 		case <-t.C:
	// 			downloaded := resp.BytesComplete()
	// 			result := float64(downloaded) / float64(totalSize)
	// 			fmt.Printf("  transferred %d / %d bytes (%.2f%%)\n",
	// 				resp.BytesComplete(),
	// 				totalSize,
	// 				math.Round(result*10000)/100)
	// 		case <-resp.Done:
	// 			break Loop
	// 		}
	// 	}

	// 	if err := resp.Err(); err != nil {
	// 		// fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
	// 		logger.Errorf("Download failed: %v", err)
	// 		return err
	// 	}

	// return nil
}

func (b *KubeBinary) Download() error {
	for i := 5; i > 0; i-- {
		cmd := exec.Command("/bin/sh", "-c", b.GetCmd())
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		cmd.Stderr = cmd.Stdout

		if err = cmd.Start(); err != nil {
			return err
		}
		for {
			tmp := make([]byte, 1024)
			_, err := stdout.Read(tmp)
			fmt.Print(string(tmp)) // Get the output from the pipeline in real time and print it to the terminal
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				logger.Error(err)
				break
			}
		}
		if err = cmd.Wait(); err != nil {
			if os.Getenv("KKZONE") != "cn" {
				logger.Warn("Having a problem with accessing https://storage.googleapis.com? You can try again after setting environment 'export KKZONE=cn'")
			}
			return err
		}

		if err := b.SHA256Check(); err != nil {
			if i == 1 {
				return err
			}
			path := b.Path()
			_ = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s", path)).Run()
			continue
		}
		break
	}
	return nil
}

// SHA256Check is used to hash checks on downloaded binary. (sha256)
func (b *KubeBinary) SHA256Check() error {
	output, err := sha256sum(b.Path())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to check SHA256 of %s", b.Path()))
	}

	if strings.TrimSpace(b.GetSha256()) == "" {
		return errors.New(fmt.Sprintf("No SHA256 found for %s. %s is not supported.", b.ID, b.Version))
	}
	if output != b.GetSha256() {
		return errors.New(fmt.Sprintf("SHA256 no match. %s not equal %s", b.GetSha256(), output))
	}
	return nil
}

func sha256sum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

var (
	// FileSha256 is a hash table the storage the checksum of the binary files.
	FileSha256 = map[string]map[string]map[string]string{
		kubeadm: {
			amd64: {
				"v1.19.0": "88ce7dc5302d8847f6e679aab9e4fa642a819e8a33d70731fb7bc8e110d8659f",
				"v1.19.8": "9c6646cdf03efc3194afc178647205195da4a43f58d0b70954953f566fa15c76",
			},
			arm64: {
				"v1.19.0": "db1c432646e6e6484989b6f7191f3610996ac593409f12574290bfc008ea11f5",
				"v1.19.8": "dfb838ffb88d79e4d881326f611ae5e5999accb54cdd666c75664da264b5d58e",
			},
		},
		kubelet: {
			amd64: {
				"v1.19.0": "3f03e5c160a8b658d30b34824a1c00abadbac96e62c4d01bf5c9271a2debc3ab",
				"v1.19.8": "f5cad5260c29584dd370ec13e525c945866957b1aaa719f1b871c31dc30bcb3f",
			},
			arm64: {
				"v1.19.0": "d8fa5a9739ecc387dfcc55afa91ac6f4b0ccd01f1423c423dbd312d787bbb6bf",
				"v1.19.8": "a00146c16266d54f961c40fc67f92c21967596c2d730fa3dc95868d4efb44559",
			},
		},
		kubectl: {
			amd64: {
				"v1.19.0": "79bb0d2f05487ff533999a639c075043c70a0a1ba25c1629eb1eef6ebe3ba70f",
				"v1.19.8": "a0737d3a15ca177816b6fb1fd59bdd5a3751bfdc66de4e08dffddba84e38bf3f",
			},
			arm64: {
				"v1.19.0": "d4adf1b6b97252025cb2f7febf55daa3f42dc305822e3da133f77fd33071ec2f",
				"v1.19.8": "8f037ab2aa798bbc66ebd1d52653f607f223b07813bcf98d9c1d0c0e136910ec",
			},
		},
		etcd: {
			amd64: {
				"v3.4.13": "2ac029e47bab752dacdb7b30032f230f49e2f457cbc32e8f555c2210bb5ff107",
			},
			arm64: {
				"v3.4.13": "1934ebb9f9f6501f706111b78e5e321a7ff8d7792d3d96a76e2d01874e42a300",
			},
		},
		helm: {
			amd64: {
				"v3.2.1": "98c57f2b86493dd36ebaab98990e6d5117510f5efbf21c3344c3bdc91a4f947c",
				"v3.6.3": "6e5498e0fa82ba7b60423b1632dba8681d629e5a4818251478cb53f0b71b3c82",
			},
			arm64: {
				"v3.2.1": "20bb9d66e74f618cd104ca07e4525a8f2f760dd6d5611f7d59b6ac574624d672",
				"v3.6.3": "fce1f94dd973379147bb63d8b6190983ad63f3a1b774aad22e54d2a27049414f",
			},
		},
		kubecni: {
			amd64: {
				"v0.8.2": "21283754ffb953329388b5a3c52cef7d656d535292bda2d86fcdda604b482f85",
				"v0.8.6": "994fbfcdbb2eedcfa87e48d8edb9bb365f4e2747a7e47658482556c12fd9b2f5",
			},
			arm64: {
				"v0.8.6": "43fbf750c5eccb10accffeeb092693c32b236fb25d919cf058c91a677822c999",
				"v0.9.1": "ef17764ffd6cdcb16d76401bac1db6acc050c9b088f1be5efa0e094ea3b01df0",
			},
		},
		k3s: {
			amd64: {
				"v1.20.2": "ce3055783cf115ee68fc00bb8d25421d068579ece2fafa4ee1d09f3415aaeabf",
				"v1.20.4": "1c7b68b0b7d54f21a9c1727545a7db181668115f161a3986bc137261dd817e98",
			},
			arm64: {
				"v1.21.4": "b7f8c026c5346b3e894d731f1dc2490cd7281687549f34c28a849f58c62e3e48",
				"v1.21.6": "1f06a2da0e1e8596220a5504291ce69237979ebf520e2458c2d72573945a9c1d",
			},
		},
		containerd: {
			amd64: {
				"1.6.2": "3d94f887de5f284b0d6ee61fa17ba413a7d60b4bb27d756a402b713a53685c6a",
				"1.6.4": "f23c8ac914d748f85df94d3e82d11ca89ca9fe19a220ce61b99a05b070044de0",
			},
			arm64: {
				"1.6.2": "a4b24b3c38a67852daa80f03ec2bc94e31a0f4393477cd7dc1c1a7c2d3eb2a95",
				"1.6.4": "0205bd1907154388dc85b1afeeb550cbb44c470ef4a290cb1daf91501c85cae6",
			},
		},
		runc: {
			amd64: {
				"v1.1.1": "5798c85d2c8b6942247ab8d6830ef362924cd72a8e236e77430c3ab1be15f080",
			},
			arm64: {
				"v1.1.1": "20c436a736547309371c7ac2a335f5fe5a42b450120e497d09c8dc3902c28444",
			},
		},
		crictl: {
			amd64: {
				"v1.22.0": "45e0556c42616af60ebe93bf4691056338b3ea0001c0201a6a8ff8b1dbc0652a",
				"v1.23.0": "b754f83c80acdc75f93aba191ff269da6be45d0fc2d3f4079704e7d1424f1ca8",
			},
			arm64: {
				"v1.22.0": "a713c37fade0d96a989bc15ebe906e08ef5c8fe5e107c2161b0665e9963b770e",
				"v1.23.0": "91094253e77094435027998a99b9b6a67b0baad3327975365f7715a1a3bd9595",
			},
		},
	}
)
