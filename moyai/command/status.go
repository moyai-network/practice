package command

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/hako/durafmt"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/carrot/role"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/shirou/gopsutil/cpu"
	"golang.org/x/text/language"
)

// Status is a command that displays information about the server such as uptime and memory usage.
type Status struct {
	startTime time.Time
}

// statusData represents the machine status data displayed in the command.
type statusData struct {
	Memory struct {
		Free  string
		Used  string
		Total string
	}
	CPU struct {
		Model string
	}
}

// NewStatus ...
func NewStatus(startTime time.Time) Status {
	return Status{startTime: startTime}
}

// Run ...
func (st Status) Run(s cmd.Source, o *cmd.Output) {
	l := language.English
	sb := &strings.Builder{}
	sb.WriteString(lang.Translatef(l, "command.status.header"))
	add := func(name string, value any) {
		sb.WriteString(text.Colourf("<grey>%s: </grey><green>%v</green>\n", name, value))
	}
	data := getStatusData()
	data.formatStorage()

	add(lang.Translatef(l, "command.status.entry.uptime"), durafmt.Parse(time.Since(st.startTime).Round(time.Second)).String())
	add(lang.Translatef(l, "command.status.entry.cpu"), data.CPU.Model)
	cpuUsage, _ := cpu.Percent(0, false)
	add(lang.Translatef(l, "command.status.entry.cpu-usage"), text.Colourf("<gold>%.2f%%</gold>", cpuUsage[0]))
	add(lang.Translatef(l, "command.status.entry.memory"), text.Colourf("<red>%v</red><grey>/</grey><gold>%v</gold>", data.Memory.Used, data.Memory.Total))

	o.Print(sb.String())
}

// getStatusData ...
func getStatusData() (s statusData) {
	if runtime.GOOS == "linux" {
		if d, err := os.ReadFile("/proc/meminfo"); err == nil {
			parseFields(d, map[string]*string{"MemAvailable": &s.Memory.Free, "MemTotal": &s.Memory.Total}, true)
		}
		if d, err := os.ReadFile("/proc/cpuinfo"); err == nil {
			parseFields(d, map[string]*string{"model name": &s.CPU.Model}, false)
		}
		return s
	}
	panic("unsupported os")
}

// parseFields ...
func parseFields(data []byte, fields map[string]*string, memory bool) {
	for _, line := range strings.Split(string(data), "\n") {
		index := strings.IndexRune(line, ':')
		if index == -1 {
			continue
		}
		if val, ok := fields[strings.TrimSpace(line[:index])]; ok {
			if memory {
				*val = strings.TrimSpace(strings.TrimRight(line[index+1:], "kB"))
			} else {
				if len(line) > index+2 {
					*val = line[index+2:]
				} else {
					*val = "Unavailable"
				}
			}
		}
	}
}

// formatStorage ...
func (d *statusData) formatStorage() {
	var memory []int
	f := func(s string) string {
		if n, err := strconv.Atoi(s); err == nil {
			memory = append(memory, n/1000)
			return strconv.Itoa(n/1000) + " MB"
		}
		return "Unavailable"
	}
	d.Memory.Free, d.Memory.Total = f(d.Memory.Free), f(d.Memory.Total)
	if len(memory) >= 2 {
		d.Memory.Used = strconv.Itoa(memory[1]-memory[0]) + " MB"
	} else {
		d.Memory.Used = "Unavailable"
	}
}

// Allow ...
func (Status) Allow(s cmd.Source) bool {
	return allow(s, true, role.Admin{})
}
