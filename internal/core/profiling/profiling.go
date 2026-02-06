package profiling

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

func StartProfiling(profile *string) func() {
	fmt.Println(" ** STARTING PROFILING ** ")
	f, err := os.Create(*profile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)

	return func() {
		fmt.Println(" ** STOPPING PROFILING ** ")
		pprof.StopCPUProfile()
		f.Close()
	}
}
