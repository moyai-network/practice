package command

import (
	"os"
	"runtime/pprof"

	"github.com/df-mc/dragonfly/server/cmd"
)

type Pprof struct {
	ProfileType profileType          `cmd:"profileType"`
	Path        cmd.Optional[string] `cmd:"path"`
}

var cpuProfile bool
var cpuFile *os.File

func (p Pprof) Run(src cmd.Source, o *cmd.Output) {
	switch p.ProfileType {
	case "cpu":
		if cpuProfile {
			o.Print("Stopping cpu profile")
			pprof.StopCPUProfile()
			_ = cpuFile.Close()
			cpuFile = nil
			cpuProfile = false
		} else {
			path, ok := p.Path.Load()
			if !ok {
				o.Error("No path given")
				return
			}
			var err error
			cpuFile, err = os.Create(path)
			if err != nil {
				o.Error(err)
				return
			}
			cpuProfile = true
			_ = pprof.StartCPUProfile(cpuFile)
			o.Print("Starting cpu profile at ", path)
		}
	case "memory":
		path, ok := p.Path.Load()
		if !ok {
			o.Error("No path given")
			return
		}
		f, err := os.Create(path)
		if err != nil {
			o.Error(err)
			return
		}
		defer f.Close()
		_ = pprof.WriteHeapProfile(f)
		o.Print("Wrote memory profile to ", path)
	case "goroutine":
		path, ok := p.Path.Load()
		if !ok {
			o.Error("No path given")
			return
		}
		f, err := os.Create(path)
		if err != nil {
			o.Error(err)
			return
		}
		defer f.Close()
		// TODO: figure out how to do a goroutine profile from code
		o.Print("Wrote goroutine profile to ", path)
	}
}

// Allow ...
func (Pprof) Allow(s cmd.Source) bool {
	return allow(s, true)
}

// response ...
type profileType string

// Type ...
func (p profileType) Type() string {
	return "profileType"
}

// Options ...
func (p profileType) Options(cmd.Source) []string {
	return []string{"cpu", "memory", "goroutine"}
}
