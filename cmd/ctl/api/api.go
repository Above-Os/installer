package api

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"bytetrade.io/web3os/installer/cmd/ctl/options"
	"bytetrade.io/web3os/installer/pkg/apiserver"
	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/phase/startup"
	"bytetrade.io/web3os/installer/pkg/utils"
	"github.com/spf13/cobra"
)

type ApiServerOptions struct {
	ApiOptions *options.ApiOptions
}

func NewApiServerOptions() *ApiServerOptions {
	return &ApiServerOptions{
		ApiOptions: options.NewApiOptions(),
	}
}

func NewCmdApi() *cobra.Command {
	o := NewApiServerOptions()
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Terminus Api Server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(constants.Logo)

			if err := InitLog(o.ApiOptions); err != nil {
				fmt.Println("init logger failed", err)
				os.Exit(1)
			}

			// fmt.Println("---1---")
			// logger.Debugf("this is hahah %s, max: %d", "123321", 900)
			// logger.Error("[error]", "error world!")
			// logger.Info("[info]", "hello world!!!")
			// fmt.Println("---end---")

			if err := GetCurrentUser(); err != nil {
				logger.Errorf(err.Error())
				os.Exit(1)
			}

			logger.Debugf("current user: %s", constants.CurrentUser)

			if constants.CurrentUser != "root" {
				logger.Error("Current user is not root!! exit ...")
				os.Exit(0)
			}

			logger.Info("Terminus Installer starting ...")

			if err := GetMachineInfo(); err != nil {
				logger.Errorf("failed to get machine info: %+v", err)
				os.Exit(1)
			}

			if err := Run(o.ApiOptions); err != nil {
				logger.Errorf("failed to run installer api server: %+v", err)
				os.Exit(1)
			}
		},
	}

	o.ApiOptions.AddFlags(cmd)

	return cmd
}

func InitLog(option *options.ApiOptions) error {
	workDir, err := utils.WorkDir()
	if err != nil {
		fmt.Println("fetch working path error", err)
		os.Exit(1)
	}

	constants.WorkDir = workDir
	constants.ApiServerListenAddress = option.Port
	constants.Proxy = option.Proxy

	logDir := path.Join(workDir, "logs")
	logger.InitLog(logDir, option.LogLevel)
	return nil
}

func GetCurrentUser() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	constants.CurrentUser = u.Username
	return nil
}

func GetMachineInfo() error {
	if err := startup.GetMachineInfo(); err != nil {
		return err
	}
	return nil
}

func Run(option *options.ApiOptions) error {

	logger.Infow("[Installer] API Server startup flags",
		"enabled", option.Enabled,
		"port", option.Port,
		"log-level", option.LogLevel,
	)

	s, err := apiserver.New()
	if err != nil {
		return err
	}

	if err = s.PrepareRun(); err != nil {
		return err
	}

	return s.Run()
}
