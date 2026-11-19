package main

import (
	"context"
	"fmt"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/server"
	"os"
)

func main() {
	ctx := context.Background()
	if err := server.Run(ctx, os.Args, os.Getenv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
