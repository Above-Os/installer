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

package precheck

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	kubekeyapiv1alpha2 "bytetrade.io/web3os/installer/apis/kubekey/v1alpha2"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/action"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/files"
	"bytetrade.io/web3os/installer/pkg/utils"
	"bytetrade.io/web3os/installer/pkg/version/kubernetes"
	"bytetrade.io/web3os/installer/pkg/version/kubesphere"
	"github.com/pkg/errors"
	versionutil "k8s.io/apimachinery/pkg/util/version"
)

// ~ PatchDeps
// install socat and contrack on other systemc
type PatchDeps struct {
	action.BaseAction
}

func (t *PatchDeps) GetName() string {
	return "PatchDeps"
}

func (t *PatchDeps) Execute(runtime connector.Runtime) error {
	// 如果是特殊的系统，需要通过源代码来安装 socat 和 contrack
	switch constants.OsPlatform {
	case common.Ubuntu, common.Debian, common.Raspbian, common.CentOs, common.Fedora, common.RHEl:
		return nil
	}

	socat := files.NewKubeBinary("socat", constants.OsArch, kubekeyapiv1alpha2.DefaultSocatVersion, runtime.GetDependDir())
	contrack := files.NewKubeBinary("contrack", constants.OsArch, kubekeyapiv1alpha2.DefaultContrackVersion, runtime.GetDependDir())

	binaries := []*files.KubeBinary{socat, contrack}
	binariesMap := make(map[string]*files.KubeBinary)

	for _, binary := range binaries {
		if err := binary.CreateBaseDir(); err != nil {
			return errors.Wrapf(errors.WithStack(err), "create file %s base dir failed", binary.FileName)
		}

		logger.Infof("%s downloading %s %s %s ...", common.LocalHost, constants.OsArch, binary.ID, binary.Version)

		binariesMap[binary.ID] = binary
		if util.IsExist(binary.Path()) {
			p := binary.Path()
			if err := binary.SHA256Check(); err != nil {
				_ = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s", p)).Run()
			} else {
				continue
			}
		}

		if err := binary.Download(); err != nil {
			return fmt.Errorf("Failed to download %s binary: %s error: %w ", binary.ID, binary.Url, err)
		}
	}

	t.PipelineCache.Set(common.KubeBinaries+"-"+constants.OsArch, binariesMap)
	return nil
}

// ~ GetSysInfoTask
type GetSysInfoTask struct {
	action.BaseAction
}

func (t *GetSysInfoTask) GetName() string {
	return "GetSysInfo"
}

func (t *GetSysInfoTask) Execute(runtime connector.Runtime) error {
	host, err := util.GetHost()
	if err != nil {
		return err
	}
	constants.HostName = host[0]
	constants.HostId = host[1]
	constants.OsType = host[2]
	constants.OsPlatform = host[3]
	constants.OsVersion = host[4]
	constants.OsArch = host[5]

	cpuModel, cpuLogicalCount, cpuPhysicalCount, err := util.GetCpu()
	if err != nil {
		return err
	}
	constants.CpuModel = cpuModel
	constants.CpuLogicalCount = cpuLogicalCount
	constants.CpuPhysicalCount = cpuPhysicalCount

	diskTotal, diskFree, err := util.GetDisk()
	if err != nil {
		return err
	}
	constants.DiskTotal = diskTotal
	constants.DiskFree = diskFree

	memTotal, memFree, err := util.GetMem()
	if err != nil {
		return err
	}
	constants.MemTotal = memTotal
	constants.MemFree = memFree

	logger.Debugf("MACHINE, hostname: %s, cpu: %d, mem: %d, disk: %d",
		constants.HostName, constants.CpuPhysicalCount, utils.FormatBytes(int64(constants.MemTotal)), utils.FormatBytes(int64(constants.DiskTotal)))
	logger.Debugf("SYSTEM, os: %s, platform: %s, arch: %s, version: %s",
		constants.OsType, constants.OsPlatform, constants.OsArch, constants.OsVersion)

	logger.Infof("host info, hostname: %s, hostid: %s, os: %s, platform: %s, version: %s, arch: %s",
		constants.HostName, constants.HostId, constants.OsType, constants.OsPlatform, constants.OsVersion, constants.OsArch)
	logger.Infof("cpu info, model: %s, logical count: %d, physical count: %d",
		constants.CpuModel, constants.CpuLogicalCount, constants.CpuPhysicalCount)
	logger.Infof("disk info, total: %d, free: %d", utils.FormatBytes(int64(constants.DiskTotal)), utils.FormatBytes(int64(constants.DiskFree)))
	logger.Infof("mem info, total: %d, free: %d", utils.FormatBytes(int64(constants.MemTotal)), utils.FormatBytes(int64(constants.MemFree)))

	return nil
}

// ~ GetLocalIpTask
type GetLocalIpTask struct {
	action.BaseAction
}

func (t *GetLocalIpTask) GetName() string {
	return "GetLocalIpTask"
}

func (t *GetLocalIpTask) Execute(runtime connector.Runtime) error {
	pingCmd := fmt.Sprintf("ping -c 1 %s", constants.HostName)
	pingCmdRes, _, err := util.Exec(pingCmd, false)
	if err != nil {
		return err
	}

	pingIps, err := utils.ExtractIP(pingCmdRes)
	if err != nil {
		return err
	}

	logger.Debugf("GetLocalIpHook, local ip: %s", pingIps)
	constants.LocalIp = pingIps

	return nil
}

// ~ TerminusGreetingsTask
type TerminusGreetingsTask struct {
	action.BaseAction
}

func (t *TerminusGreetingsTask) GetName() string {
	return "TerminusGreetingsTask"
}

func (h *TerminusGreetingsTask) Execute(runtime connector.Runtime) error {
	stdout, _, err := util.Exec("echo 'Greetings, Terminus!!!!!' ", false)
	if err != nil {
		return err
	}
	logger.Infof("TerminusGreetingsTask %s", stdout)
	return nil
}

// ~ GreetingsTask
type GreetingsTask struct {
	action.BaseAction
}

func (t *GreetingsTask) GetName() string {
	return "GreetingsTask"
}

func (h *GreetingsTask) Execute(runtime connector.Runtime) error {
	hello, err := runtime.GetRunner().SudoCmd("echo 'Greetings, KubeKey!!!!! hahahaha!!!!'", false)
	if err != nil {
		return err
	}
	logger.Infof("%s %s", runtime.RemoteHost().GetName(), hello)
	return nil
}

// ~ NodePreCheck
type NodePreCheck struct {
	common.KubeAction
}

func (t *NodePreCheck) GetName() string {
	return "NodePreCheck"
}

func (n *NodePreCheck) Execute(runtime connector.Runtime) error {
	var results = make(map[string]string)
	results["name"] = runtime.RemoteHost().GetName()
	for _, software := range baseSoftware {
		var (
			cmd string
		)

		switch software {
		case docker:
			cmd = "docker version --format '{{.Server.Version}}'"
		case containerd:
			cmd = "containerd --version | cut -d ' ' -f 3"
		default:
			cmd = fmt.Sprintf("which %s", software)
		}

		switch software {
		case sudo:
			// sudo skip sudo prefix
		default:
			cmd = connector.SudoPrefix(cmd)
		}

		res, err := runtime.GetRunner().Cmd(cmd, false)
		switch software {
		case showmount:
			software = nfs
		case rbd:
			software = ceph
		case glusterfs:
			software = glusterfs
		}
		if err != nil || strings.Contains(res, "not found") {
			results[software] = ""
		} else {
			// software in path
			if strings.Contains(res, "bin/") {
				results[software] = "y"
			} else {
				// get software version, e.g. docker, containerd, etc.
				results[software] = res
			}
		}
	}

	output, err := runtime.GetRunner().Cmd("date +\"%Z %H:%M:%S\"", false)
	if err != nil {
		results["time"] = ""
	} else {
		results["time"] = strings.TrimSpace(output)
	}

	host := runtime.RemoteHost()
	if res, ok := host.GetCache().Get(common.NodePreCheck); ok {
		m := res.(map[string]string)
		m = results
		host.GetCache().Set(common.NodePreCheck, m)
	} else {
		host.GetCache().Set(common.NodePreCheck, results)
	}
	return nil
}

// ~ GetKubeConfig
type GetKubeConfig struct {
	common.KubeAction
}

func (t *GetKubeConfig) GetName() string {
	return "GetKubeConfig"
}

func (g *GetKubeConfig) Execute(runtime connector.Runtime) error {
	if exist, err := runtime.GetRunner().FileExist("$HOME/.kube/config"); err != nil {
		return err
	} else {
		if exist {
			return nil
		} else {
			if exist, err := runtime.GetRunner().FileExist("/etc/kubernetes/admin.conf"); err != nil {
				return err
			} else {
				if exist {
					if _, err := runtime.GetRunner().Cmd("mkdir -p $HOME/.kube", false); err != nil {
						return err
					}
					if _, err := runtime.GetRunner().SudoCmd("cp /etc/kubernetes/admin.conf $HOME/.kube/config", false); err != nil {
						return err
					}
					userId, err := runtime.GetRunner().Cmd("echo $(id -u)", false)
					if err != nil {
						return errors.Wrap(errors.WithStack(err), "get user id failed")
					}

					userGroupId, err := runtime.GetRunner().Cmd("echo $(id -g)", false)
					if err != nil {
						return errors.Wrap(errors.WithStack(err), "get user group id failed")
					}

					chownKubeConfig := fmt.Sprintf("chown -R %s:%s $HOME/.kube", userId, userGroupId)
					if _, err := runtime.GetRunner().SudoCmd(chownKubeConfig, false); err != nil {
						return errors.Wrap(errors.WithStack(err), "chown user kube config failed")
					}
				}
			}
		}
	}
	return errors.New("kube config not found")
}

// ~ GetAllNodesK8sVersion
type GetAllNodesK8sVersion struct {
	common.KubeAction
}

func (t *GetAllNodesK8sVersion) GetName() string {
	return "GetAllNodesK8sVersion"
}

func (g *GetAllNodesK8sVersion) Execute(runtime connector.Runtime) error {
	var nodeK8sVersion string
	kubeletVersionInfo, err := runtime.GetRunner().SudoCmd("/usr/local/bin/kubelet --version", false)
	if err != nil {
		return errors.Wrap(err, "get current kubelet version failed")
	}
	nodeK8sVersion = strings.Split(kubeletVersionInfo, " ")[1]

	host := runtime.RemoteHost()
	if host.IsRole(common.Master) {
		apiserverVersion, err := runtime.GetRunner().SudoCmd(
			"cat /etc/kubernetes/manifests/kube-apiserver.yaml | grep 'image:' | rev | cut -d ':' -f1 | rev",
			false)
		if err != nil {
			return errors.Wrap(err, "get current kube-apiserver version failed")
		}
		nodeK8sVersion = apiserverVersion
	}
	host.GetCache().Set(common.NodeK8sVersion, nodeK8sVersion)
	return nil
}

// ~ CalculateMinK8sVersion
type CalculateMinK8sVersion struct {
	common.KubeAction
}

func (g *CalculateMinK8sVersion) GetName() string {
	return "CalculateMinK8sVersion"
}

func (g *CalculateMinK8sVersion) Execute(runtime connector.Runtime) error {
	versionList := make([]*versionutil.Version, 0, len(runtime.GetHostsByRole(common.K8s)))
	for _, host := range runtime.GetHostsByRole(common.K8s) {
		version, ok := host.GetCache().GetMustString(common.NodeK8sVersion)
		if !ok {
			return errors.Errorf("get node %s Kubernetes version failed by host cache", host.GetName())
		}
		if versionObj, err := versionutil.ParseSemantic(version); err != nil {
			return errors.Wrap(err, "parse node version failed")
		} else {
			versionList = append(versionList, versionObj)
		}
	}

	minVersion := versionList[0]
	for _, version := range versionList {
		if !minVersion.LessThan(version) {
			minVersion = version
		}
	}
	g.PipelineCache.Set(common.K8sVersion, fmt.Sprintf("v%s", minVersion))
	return nil
}

// ~ CheckDesiredK8sVersion
type CheckDesiredK8sVersion struct {
	common.KubeAction
}

func (k *CheckDesiredK8sVersion) GetName() string {
	return "CheckDesiredK8sVersion"
}

func (k *CheckDesiredK8sVersion) Execute(_ connector.Runtime) error {
	if ok := kubernetes.VersionSupport(k.KubeConf.Cluster.Kubernetes.Version); !ok {
		return errors.New(fmt.Sprintf("does not support upgrade to Kubernetes %s",
			k.KubeConf.Cluster.Kubernetes.Version))
	}
	k.PipelineCache.Set(common.DesiredK8sVersion, k.KubeConf.Cluster.Kubernetes.Version)
	return nil
}

// ~ KsVersionCheck
type KsVersionCheck struct {
	common.KubeAction
}

func (k *KsVersionCheck) GetName() string {
	return "KsVersionCheck"
}

func (k *KsVersionCheck) Execute(runtime connector.Runtime) error {
	ksVersionStr, err := runtime.GetRunner().SudoCmd(
		"/usr/local/bin/kubectl get deploy -n  kubesphere-system ks-console -o jsonpath='{.metadata.labels.version}'",
		false)
	if err != nil {
		if k.KubeConf.Cluster.KubeSphere.Enabled {
			return errors.Wrap(err, "get kubeSphere version failed")
		} else {
			ksVersionStr = ""
		}
	}

	ccKsVersionStr, ccErr := runtime.GetRunner().SudoCmd(
		"/usr/local/bin/kubectl get ClusterConfiguration ks-installer -n  kubesphere-system  -o jsonpath='{.metadata.labels.version}'",
		false)
	if ccErr == nil && ksVersionStr == "v3.1.0" {
		ksVersionStr = ccKsVersionStr
	}
	k.PipelineCache.Set(common.KubeSphereVersion, ksVersionStr)
	return nil
}

// ~ DependencyCheck
type DependencyCheck struct {
	common.KubeAction
}

func (d *DependencyCheck) GetName() string {
	return "DependencyCheck"
}

func (d *DependencyCheck) Execute(_ connector.Runtime) error {
	currentKsVersion, ok := d.PipelineCache.GetMustString(common.KubeSphereVersion)
	if !ok {
		return errors.New("get current KubeSphere version failed by pipeline cache")
	}
	desiredVersion := d.KubeConf.Cluster.KubeSphere.Version

	if d.KubeConf.Cluster.KubeSphere.Enabled {
		var version string
		if latest, ok := kubesphere.LatestRelease(desiredVersion); ok {
			version = latest.Version
		} else if ks, ok := kubesphere.DevRelease(desiredVersion); ok {
			version = ks.Version
		} else {
			r := regexp.MustCompile("v(\\d+\\.)?(\\d+\\.)?(\\*|\\d+)")
			version = r.FindString(desiredVersion)
		}

		ksInstaller, ok := kubesphere.VersionMap[version]
		if !ok {
			return errors.New(fmt.Sprintf("Unsupported version: %s", desiredVersion))
		}

		if currentKsVersion != desiredVersion {
			if ok := ksInstaller.UpgradeSupport(currentKsVersion); !ok {
				return errors.New(fmt.Sprintf("Unsupported upgrade plan: %s to %s", currentKsVersion, desiredVersion))
			}
		}

		if ok := ksInstaller.K8sSupport(d.KubeConf.Cluster.Kubernetes.Version); !ok {
			return errors.New(fmt.Sprintf("KubeSphere %s does not support running on Kubernetes %s",
				version, d.KubeConf.Cluster.Kubernetes.Version))
		}
	} else {
		ksInstaller, ok := kubesphere.VersionMap[currentKsVersion]
		if !ok {
			return errors.New(fmt.Sprintf("Unsupported version: %s", currentKsVersion))
		}

		if ok := ksInstaller.K8sSupport(d.KubeConf.Cluster.Kubernetes.Version); !ok {
			return errors.New(fmt.Sprintf("KubeSphere %s does not support running on Kubernetes %s",
				currentKsVersion, d.KubeConf.Cluster.Kubernetes.Version))
		}
	}
	return nil
}

// ~ GetKubernetesNodesStatus
type GetKubernetesNodesStatus struct {
	common.KubeAction
}

func (g *GetKubernetesNodesStatus) GetName() string {
	return "GetKubernetesNodesStatus"
}

func (g *GetKubernetesNodesStatus) Execute(runtime connector.Runtime) error {
	nodeStatus, err := runtime.GetRunner().SudoCmd("/usr/local/bin/kubectl get node -o wide", false)
	if err != nil {
		return err
	}
	g.PipelineCache.Set(common.ClusterNodeStatus, nodeStatus)

	cri, err := runtime.GetRunner().SudoCmd("/usr/local/bin/kubectl get node -o jsonpath=\"{.items[*].status.nodeInfo.containerRuntimeVersion}\"", false)
	if err != nil {
		return err
	}
	g.PipelineCache.Set(common.ClusterNodeCRIRuntimes, cri)
	return nil
}
