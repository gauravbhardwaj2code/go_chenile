package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("chenile-packager", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifest := flags.String("manifest", "chenile-packager.yaml", "packager manifest path")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if err := checkManifest(*manifest); err != nil {
		fmt.Fprintf(stderr, "manifest %q not found: %v\n", *manifest, err)
		return 1
	}
	fmt.Fprintf(stdout, "manifest %q exists; use packager.NewWebApp(...) from a mainweb app to combine services\n", *manifest)
	return 0
}

func checkManifest(path string) error {
	_, err := os.Stat(path)
	return err
}
