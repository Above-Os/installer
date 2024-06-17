package util

import (
	"bufio"
	"os/exec"

	"bytetrade.io/web3os/installer/pkg/core/logger"
)

func Exec(name string, args []string) {
	cmd := exec.Command(name, args...)

	// 获取标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Errorf("Error getting StdoutPipe:", err)
		return
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		logger.Errorf("Error starting command:", err)
		return
	}

	// 创建一个扫描器读取标准输出
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		// 每次读取一行并打印
		logger.Info(scanner.Text())
	}

	// 检查扫描错误
	if err := scanner.Err(); err != nil {
		logger.Errorf("Error reading standard output:", err)
	}

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		logger.Errorf("Error waiting for command to finish:", err)
	}
}
