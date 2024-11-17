package core

import "os"

func corePanic(_ error) {
	os.Exit(1)
}
