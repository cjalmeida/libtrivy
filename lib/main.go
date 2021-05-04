package main

import "C"
import (
	"github.com/cjalmeida/libtrivy/pkg/scan"
)

//export TrivyScan
func TrivyScan(src, dst *C.char) *C.char {

	// go func() {
	// 	sigs := make(chan os.Signal, 1)
	// 	signal.Notify(sigs, syscall.SIGHUP)
	// 	buf := make([]byte, 1<<20)
	// 	for {
	// 		<-sigs
	// 		stacklen := runtime.Stack(buf, true)
	// 		fmt.Printf("=== received SIGHUP ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stacklen])
	// 	}
	// }()

	sourceFile := C.GoString(src)
	destFile := C.GoString(dst)
	err := scan.Scan(sourceFile, destFile)
	if err != nil {
		return C.CString("ERROR " + err.Error())
	}
	return C.CString("OK")
}

func main() {}
