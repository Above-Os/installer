package cluster

import (
	kubekeyapiv1alpha2 "bytetrade.io/web3os/installer/apis/kubekey/v1alpha2"

	"bytetrade.io/web3os/installer/pkg/addons"
	"bytetrade.io/web3os/installer/pkg/artifact"
	"bytetrade.io/web3os/installer/pkg/binaries"
	"bytetrade.io/web3os/installer/pkg/bootstrap/confirm"
	"bytetrade.io/web3os/installer/pkg/bootstrap/os"
	"bytetrade.io/web3os/installer/pkg/bootstrap/precheck"
	"bytetrade.io/web3os/installer/pkg/certs"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/container"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/etcd"
	"bytetrade.io/web3os/installer/pkg/filesystem"
	"bytetrade.io/web3os/installer/pkg/images"
	"bytetrade.io/web3os/installer/pkg/kubernetes"
	"bytetrade.io/web3os/installer/pkg/kubesphere"
	"bytetrade.io/web3os/installer/pkg/loadbalancer"
	"bytetrade.io/web3os/installer/pkg/plugins"
	"bytetrade.io/web3os/installer/pkg/plugins/dns"
	"bytetrade.io/web3os/installer/pkg/plugins/network"
	"bytetrade.io/web3os/installer/pkg/plugins/storage"
)

func NewCreateClusterPhase(runtime *common.KubeRuntime) []module.Module {
	noArtifact := runtime.Arg.Artifact == ""
	skipPushImages := runtime.Arg.SKipPushImages || noArtifact || (!noArtifact && runtime.Cluster.Registry.PrivateRegistry == "")
	skipLocalStorage := true
	if runtime.Arg.DeployLocalStorage != nil {
		skipLocalStorage = !*runtime.Arg.DeployLocalStorage
	} else if runtime.Cluster.KubeSphere.Enabled {
		skipLocalStorage = false
	}

	m := []module.Module{
		&precheck.GreetingsModule{},
		&precheck.NodePreCheckModule{},
		&confirm.InstallConfirmModule{Skip: runtime.Arg.SkipConfirmCheck},
		&artifact.UnArchiveModule{Skip: noArtifact},
		&os.RepositoryModule{Skip: noArtifact || !runtime.Arg.InstallPackages},
		&os.ConfigureOSModule{},
		&binaries.NodeBinariesModule{},
		&kubernetes.StatusModule{},
		&container.InstallContainerModule{},
		&images.CopyImagesToRegistryModule{Skip: skipPushImages},
		&images.PullModule{Skip: runtime.Arg.SkipPullImages},
		&etcd.PreCheckModule{Skip: runtime.Cluster.Etcd.Type != kubekeyapiv1alpha2.KubeKey},
		&etcd.CertsModule{},
		&etcd.InstallETCDBinaryModule{Skip: runtime.Cluster.Etcd.Type != kubekeyapiv1alpha2.KubeKey},
		&etcd.ConfigureModule{Skip: runtime.Cluster.Etcd.Type != kubekeyapiv1alpha2.KubeKey},
		&etcd.BackupModule{Skip: runtime.Cluster.Etcd.Type != kubekeyapiv1alpha2.KubeKey},
		&kubernetes.InstallKubeBinariesModule{},
		&loadbalancer.KubevipModule{Skip: !runtime.Cluster.ControlPlaneEndpoint.IsInternalLBEnabledVip()},
		&kubernetes.InitKubernetesModule{},
		&dns.ClusterDNSModule{},
		&kubernetes.StatusModule{},
		&kubernetes.JoinNodesModule{},
		&loadbalancer.KubevipModule{Skip: !runtime.Cluster.ControlPlaneEndpoint.IsInternalLBEnabledVip()},
		&loadbalancer.HaproxyModule{Skip: !runtime.Cluster.ControlPlaneEndpoint.IsInternalLBEnabled()},
		&network.DeployNetworkPluginModule{},
		&kubernetes.ConfigureKubernetesModule{},
		&filesystem.ChownModule{},
		&certs.AutoRenewCertsModule{Skip: !runtime.Cluster.Kubernetes.EnableAutoRenewCerts()},
		&kubernetes.SecurityEnhancementModule{Skip: !runtime.Arg.SecurityEnhancement},
		&kubernetes.SaveKubeConfigModule{},
		&plugins.DeployPluginsModule{},
		&addons.AddonsModule{},
		&storage.DeployLocalVolumeModule{Skip: skipLocalStorage},
		&kubesphere.DeployModule{Skip: !runtime.Cluster.KubeSphere.Enabled},
		&kubesphere.CheckResultModule{Skip: !runtime.Cluster.KubeSphere.Enabled},
	}

	return m
}
