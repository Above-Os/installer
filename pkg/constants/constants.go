package constants

var Logo = `
  TERMINUS Installer
`

var (
	HostName         string
	HostId           string
	OsType           string
	OsPlatform       string
	OsVersion        string
	OsArch           string
	CpuModel         string
	CpuLogicalCount  int
	CpuPhysicalCount int
	MemTotal         uint64
	MemFree          uint64
	DiskTotal        uint64
	DiskFree         uint64

	LocalIp  []string
	PublicIp []string
)

var (
	WorkDir                string
	ApiServerListenAddress string
	Proxy                  string
	CurrentUser            string
)
