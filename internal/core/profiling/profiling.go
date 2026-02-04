package profiling

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

func StartProfiling(profile *string) func() {
	fmt.Println("Starting profiling")
	f, err := os.Create(*profile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)

	return func() {
		fmt.Println("Stopping profiling")
		pprof.StopCPUProfile()
		f.Close()
	}
}
