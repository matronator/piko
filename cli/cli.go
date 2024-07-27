package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var version string

func printVersion() {
	if version == "" {
		version = "development"
	}

	println("PIKOlang version " + version)
}

var versionFlag = flag.Bool("version", false, "Print the current version and exit")
var helpFlag = flag.Bool("help", false, "Print help message and exit")

func ParseFlags() {
	flag.BoolVar(versionFlag, "v", false, "Print the current version and exit")
	flag.BoolVar(helpFlag, "h", false, "Print help message and exit")

	flag.Parse()

	if *versionFlag {
		printVersion()
		os.Exit(0)
	}

	if *helpFlag || len(flag.Args()) >= 0 && flag.Arg(0) == "help" {
		println("Usage:\n\tpiko [file]")
		println("\n[file] The file to execute")
		println("\nFlags: (optional)")
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func GetFileContentsFromArgs() []byte {
	filename := os.Args[1]

	if filename == "" {
		println("No filename provided")
		os.Exit(1)
	}

	if !strings.HasSuffix(filename, ".piko") {
		println("Unsupported file. Files must have .piko extension")
		os.Exit(1)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("File does not exist:", filename)
		os.Exit(1)
	}

	file, err := os.ReadFile(filename)

	if err != nil {
		println("Error reading file", filename)
		os.Exit(1)
	}

	return file
}
