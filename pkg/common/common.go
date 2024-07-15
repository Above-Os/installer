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

package common

const (
	DefaultK8sVersion        = "v1.22.10"
	DefaultK3sVersion        = "v1.22.16-k3s"
	DefaultKubeSphereVersion = "v3.3.0"
)

const (
	K3s        = "k3s"
	K8e        = "k8e"
	Kubernetes = "kubernetes"

	LocalHost = "localhost"

	AllInOne    = "allInOne"
	File        = "file"
	Operator    = "operator"
	CommandLine = "commandLine"

	Master        = "master"
	Worker        = "worker"
	ETCD          = "etcd"
	K8s           = "k8s"
	Registry      = "registry"
	KubeKey       = "kubekey"
	Harbor        = "harbor"
	DockerCompose = "compose"

	KubeBinaries = "KubeBinaries"

	TmpDir                       = "/tmp/kubekey"
	BinDir                       = "/usr/local/bin"
	KubeConfigDir                = "/etc/kubernetes"
	KubeAddonsDir                = "/etc/kubernetes/addons"
	KubeCertDir                  = "/etc/kubernetes/pki"
	KubeManifestDir              = "/etc/kubernetes/manifests"
	KubeScriptDir                = "/usr/local/bin/kube-scripts"
	KubeletFlexvolumesPluginsDir = "/usr/libexec/kubernetes/kubelet-plugins/volume/exec"
	PreloadK3sImageDir           = "/var/lib/images"

	InstallerScriptsDir = "scripts"

	ETCDCertDir     = "/etc/ssl/etcd/ssl"
	RegistryCertDir = "/etc/ssl/registry/ssl"

	HaproxyDir = "/etc/kubekey/haproxy"

	IPv4Regexp = "[\\d]+\\.[\\d]+\\.[\\d]+\\.[\\d]+"
	IPv6Regexp = "[a-f0-9]{1,4}(:[a-f0-9]{1,4}){7}|[a-f0-9]{1,4}(:[a-f0-9]{1,4}){0,7}::[a-f0-9]{0,4}(:[a-f0-9]{1,4}){0,7}"

	Calico  = "calico"
	Flannel = "flannel"
	Cilium  = "cilium"
	Kubeovn = "kubeovn"

	Docker     = "docker"
	Crictl     = "crictl"
	Containerd = "containerd"
	Crio       = "crio"
	Isula      = "isula"
	Runc       = "runc"

	// global cache key
	// PreCheckModule
	NodePreCheck           = "nodePreCheck"
	K8sVersion             = "k8sVersion"        // current k8s version
	MaxK8sVersion          = "maxK8sVersion"     // max k8s version of nodes
	KubeSphereVersion      = "kubeSphereVersion" // current KubeSphere version
	ClusterNodeStatus      = "clusterNodeStatus"
	ClusterNodeCRIRuntimes = "ClusterNodeCRIRuntimes"
	DesiredK8sVersion      = "desiredK8sVersion"
	PlanK8sVersion         = "planK8sVersion"
	NodeK8sVersion         = "NodeK8sVersion"

	// ETCDModule
	ETCDCluster = "etcdCluster"
	ETCDName    = "etcdName"
	ETCDExist   = "etcdExist"

	// KubernetesModule
	ClusterStatus = "clusterStatus"
	ClusterExist  = "clusterExist"

	// CertsModule
	Certificate   = "certificate"
	CaCertificate = "caCertificate"

	// Artifact pipeline
	Artifact = "artifact"

	SkipMasterNodePullImages = "skipMasterNodePullImages"
)

const (
	Linux   = "linux"
	Darwin  = "darwin"
	Windows = "windows"

	Intel64 = "x86_64"
	Amd64   = "amd64"
	Arm     = "arm" // todo 这里要注意下，数据可能是 arm7  arm  armv7l ...
	Arm7    = "arm7"
	Armv7l  = "Armv7l"
	Armhf   = "armhf"
	Arm64   = "arm64"
	PPC64el = "ppc64el"
	PPC64le = "ppc64le"
	S390x   = "s390x"
	Riscv64 = "riscv64"

	Ubuntu   = "ubuntu"
	Debian   = "debian"
	CentOs   = "centos"
	Fedora   = "fedora"
	RHEl     = "rhel"
	Raspbian = "raspbian"
)

const (
	AliYun = "aliyun"
	AWS    = "aws"
)

const (
	RaspbianCmdlineFile  = "/boot/cmdline.txt"
	RaspbianFirmwareFile = "/boot/firmware/cmdline.txt"
)

const (
	CommandIptables  = "iptables"
	CommandGPG       = "gpg"
	CommandSocat     = "socat"
	CommandConntrack = "conntrack"
	CommandNtpdate   = "ntpdate"
	CommandHwclock   = "hwclock"
	CommandKubectl   = "kubectl"
)

const (
	CacheKubeletVersion = "version_kubelet"

	CacheKubectlKey = "cmd_kubectl"

	CacheSTSAccessKey = "sts_access_key"
	CacheSTSSecretKey = "sts_secret_key"
	CacheSTSToken     = "sts_token"
	CacheSTSClusterId = "sts_cluster_id"
	CacheProxy        = "proxy"

	CacheStorageClass  = "storage_class"
	CacheEnableHA      = "enable_ha"
	CacheMasterNum     = "master_num"
	CacheRedisPassword = "redis_password"
	CacheSecretsNum    = "secrets_num"
	CacheJwtSecret     = "jwt_secret"
	CacheCrdsNUm       = "users_iam_num"
)
