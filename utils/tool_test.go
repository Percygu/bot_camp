package utils

import (
	"fmt"
	"os/user"
	"testing"
)

func TestTool(t *testing.T) {
	fmt.Println(IsTestEnv())
	u, _ := user.Current()
	fmt.Println(u)
}
