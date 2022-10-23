package envs

import (
	"log"
	"os"
)

var (
	Lg = log.New(os.Stderr, "INFO -", 18)
)
