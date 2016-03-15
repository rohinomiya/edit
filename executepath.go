// http://qiita.com/zetamatta/items/e5fb297099455fe558b6
package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32              = syscall.NewLazyDLL("kernel32")
	procGetModuleFileName = kernel32.NewProc("GetModuleFileNameW")
)

func ExecuteFilePath() string {
	var wpath [syscall.MAX_PATH]uint16
	r1, _, err := procGetModuleFileName.Call(0, uintptr(unsafe.Pointer(&wpath[0])), uintptr(len(wpath)))
	if r1 == 0 {
		fmt.Fprintln(os.Stderr, err)
	}
	return syscall.UTF16ToString(wpath[:])
}
