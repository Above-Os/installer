package util

import (
	"bufio"
	"fmt"
	"os/exec"
)

func Exec(name string, printOutput bool) error {
	cmd := exec.Command("/bin/sh", "-c", name)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return err
	}

	var (
		output []byte
		line   = ""
		r      = bufio.NewReader(stdout)
	)

	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		output = append(output, b)
		if b == byte('\n') {
			fmt.Println(line)
			line = ""
			continue
		}
		line += string(b)
		// tmp := make([]byte, 1024)
		// _, err := stdout.Read(tmp)
		// if errors.Is(err, io.EOF) {
		// 	break
		// } else if err != nil {
		// 	fmt.Println("read error", err)
		// 	break
		// }
	}

	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}
