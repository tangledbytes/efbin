package main

import (
	"os"

	"github.com/utkarsh-pro/efbin/pkg/ef"
	"github.com/utkarsh-pro/efbin/pkg/util"
)

func main() {
	if !util.IsWrappedBinPresent(util.GetBinaryName()) {
		util.Exit(1, "failed to execute binary: binary not found - ", util.GetBinaryName())
	}

	util.PreventFuckUp()

	ef.Run(os.Args[1:])
}
