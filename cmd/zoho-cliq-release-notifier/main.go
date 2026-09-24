package main

import (
	"fmt"
	"os"

	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/action"
	"github.com/iamsadjad/zoho-cliq-release-notifier/internal/gha"
)

func main() {
	runner := &action.Runner{}
	if _, err := runner.Run(); err != nil {
		gha.Error("%s", err.Error())
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
