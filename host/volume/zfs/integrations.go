package zfs

import (
	"fmt"
	"strings"
)

type Logger struct{}

func (*Logger) Log(msg []string) {
	fmt.Printf("[zfs] %s\n", strings.Join(msg, " ")) // TODO replace with log15
}

func init() {
	// actually, you don't need to set up the logger or such here.  we might have refernces to this package without actually ever invoking it.
}
