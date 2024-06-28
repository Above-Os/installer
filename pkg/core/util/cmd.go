package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"bytetrade.io/web3os/installer/pkg/core/logger"
	"github.com/pkg/errors"
)

func Exec(name string, printOutput bool, printLine bool) (stdout string, code int, err error) {
	exitCode := 0

	cmd := exec.Command("/bin/sh", "-c", name)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return "", exitCode, err
	}

	// logger.Infof("exec cmd: %s", cmd.String())
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		exitCode = -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		return "", exitCode, err
	}

	var outputBuffer bytes.Buffer
	r := bufio.NewReader(out)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("read error:", err)
			}
			break
		}

		outputBuffer.WriteString(line)
	}

	err = cmd.Wait()
	if err != nil {
		exitCode = -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	res := outputBuffer.String()
	res = strings.TrimSpace(res)

	if printOutput {
		logger.Debugf("[exec] OUT: %s, CMD: %s", res, cmd.String())
	}

	return res, exitCode, errors.Wrapf(err, "Failed to exec command: %s \n%s", cmd, res)
}

// 用于全量包安装测试
func ExecWithChannel(name string, printOutput bool, printLine bool, output chan string) (stdout string, code int, err error) {
	defer close(output)
	exitCode := 0

	cmd := exec.Command("/bin/sh", "-c", name)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return "", exitCode, err
	}

	// logger.Infof("exec cmd: %s", cmd.String())
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		exitCode = -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		return "", exitCode, err
	}

	var outputBuffer bytes.Buffer
	r := bufio.NewReader(out)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("read error:", err)
			}
			break
		}

		if strings.Contains(line, "[INFO]") {
			output <- line
		}
		if printLine {
			fmt.Println(line)
		}

		outputBuffer.WriteString(line)
	}

	err = cmd.Wait()
	if err != nil {
		exitCode = -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	res := outputBuffer.String()
	res = strings.TrimSpace(res)

	if printOutput {
		logger.Debugf("[exec] CMD: %s, OUTPUT: \n%s", cmd.String(), res)
	}

	return res, exitCode, errors.Wrapf(err, "Failed to exec command: %s \n%s", cmd, res)
}
