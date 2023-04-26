package gopprof

import (
	//"go.uber.org/zap"
	"os"
	"runtime"
	"runtime/pprof"
	"td_report/pkg/logger"
)

func StartCpuProf() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		logger.Logger.Error("create cpu profile file error: ", err)
		return
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		logger.Logger.Error("can not start cpu profile,  error: ", err)
		f.Close()
	}
}

func StopCpuProf() {
	pprof.StopCPUProfile()
}

// ProfGc --------Mem
func ProfGc() {
	runtime.GC() // get up-to-date statistics
}

func SaveMemProf() {
	f, err := os.Create("mem.prof")
	if err != nil {
		logger.Logger.Error("create mem profile file error: ", err)
		return
	}

	if err := pprof.WriteHeapProfile(f); err != nil {
		logger.Logger.Error("could not write memory profile: ", err)
	}

	f.Close()
}

// SaveBlockProfile goroutine block
func SaveBlockProfile() {
	f, err := os.Create("block.prof")
	if err != nil {
		logger.Logger.Error("create mem profile file error: ", err)
		return
	}

	if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
		logger.Logger.Error("could not write block profile: ", err)
	}
	f.Close()
}
