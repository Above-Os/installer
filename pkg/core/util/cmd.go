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
				logger.Errorf("[exec] read error: %s", err)
			}
			break
		}

		if printLine {
			fmt.Println(strings.TrimSuffix(line, "\n"))
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
		fmt.Printf("[exec] CMD: %s, OUTPUT: \n%s\n", cmd.String(), res)
	}

	logger.Infof("[exec] CMD: %s, OUTPUT: %s", cmd.String(), res)
	return res, exitCode, errors.Wrapf(err, "Failed to exec command: %s \n%s", cmd, res)
}

// ! only test for install_cmd.sh
func ExecWithChannel(name string, printOutput bool, printLine bool, output chan []interface{}) (stdout string, code int, err error) {
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
	var step int64 = 3

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("read error:", err)
			}
			break
		}
		l := strings.TrimSuffix(line, "\n")
		l = strings.TrimSpace(l)

		if strings.Contains(l, "[INFO]") { // ! only for debug
			output <- []interface{}{l, step}
			step++
		}
		if printLine {
			fmt.Println(l)
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
		fmt.Printf("[exec] CMD: %s, OUTPUT: \n%s", cmd.String(), res)
	}

	logger.Infof("[exec] CMD: %s, OUTPUT: \n%s", cmd.String(), res)
	return res, exitCode, errors.Wrapf(err, "Failed to exec command: %s \n%s", cmd, res)
}
