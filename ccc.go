package main

import (
	"fmt"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

func main() {
	var bar = progressbar.NewOptions64(100,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionOnCompletion(func() {
			fmt.Println("")
		}),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan][1/3][downloading] %s  ", "hello")),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	defer bar.Close()

	for i := 0; i < 10; i++ {
		bar.Set64(int64(10 * (i + 1)))
		time.Sleep(300 * time.Millisecond)
	}

	// bar.Describe("[cyan][1/3][downloading] done!")

	time.Sleep(2 * time.Second)
	bar.Reset()
	bar.ChangeMax(-1)
	// bar.Set(-1)
	bar.Describe("[cyan][2/3][installing] haha")
	for i := 1; i <= 5; i++ {
		bar.Set(i)
		bar.Add(i)
		time.Sleep(1 * time.Second)
	}

	bar.Finish()
	for {
	}

}
