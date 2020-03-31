package cmd

import (
	"fmt"
	"os"
)

// outErr prints errors (if not nil) to stderr
func outErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}
