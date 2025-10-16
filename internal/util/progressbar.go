package util

import (
	"github.com/schollz/progressbar/v3"
)

func NewProgressBar(size int64) *progressbar.ProgressBar {
	return progressbar.NewOptions(
		int(size),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("Downloading..."),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]#[reset]",
			SaucerHead:    "[green]#[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}
