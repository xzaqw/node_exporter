//go:build !nops

package collector

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"os/exec"
	"strings"
	"strconv"
	"cmp"
	"slices"
)

const (
	// The maximum threshold that we'll never reach regardless of config 
	MAX_PROCESSES_NUM = 200
	//TODO: move to config
	CONFIG_PROCESSES_NUM = 15
	CONFIG_LOW_CPU_THRESHOLD = float64(5)
	CONFIG_LOW_MEM_THRESHOLD = float64(5)
)

type PsCollector struct {
	psCpu *prometheus.GaugeVec
	psMem *prometheus.GaugeVec
	logger	log.Logger
}

type psLineDesc struct {
	pcpu float64
	pmem float64
	pid uint
	zoneid uint
	comm string
	args string
}

func cmpPsCpu(a, b psLineDesc) int {
	return cmp.Compare(a.pcpu, b.pcpu)
}

func cmpPsMem(a, b psLineDesc) int {
	return cmp.Compare(a.pmem, b.pmem)
}

func init() {
	registerCollector("ps", defaultEnabled, NewPsCollector)
}

func NewPsCollector(logger log.Logger) (Collector, error) {
	return &PsCollector {
		psCpu: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "node_ps_top_cpu_percents",
			Help: "Process of top CPU consumption processes.",
		}, []string{"index", "pid", "zoneid", "comm", "args"}),
		psMem: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "node_ps_top_mem_percents",
			Help: "Process of top memory consumption processes.",
		}, []string{"index", "pid", "zoneid", "comm", "args"}),
		logger: logger,
	}, nil
}

func (c *PsCollector) Update(ch chan<- prometheus.Metric) error {
	c.getPsOut()

	c.psCpu.Collect(ch)
	c.psMem.Collect(ch)

	return nil
}

func (c *PsCollector) Describe(ch chan<- *prometheus.Desc) {
	c.psCpu.Describe(ch)
	c.psMem.Describe(ch)
}

func parsePsOutput(psOut string) ([]psLineDesc , error) {
	var err error
	var out []psLineDesc 
	var args string

	psOutLines := strings.Split(psOut, "\n")

	for _,line := range psOutLines {
		//Filter out the header of ps output
		if strings.HasPrefix(line, "%") {
			continue
		}
		parsed_line := strings.Fields(line)
		if len(parsed_line) == 0 {
			continue
		}

		pcpu, err := strconv.ParseFloat(parsed_line[0], 64)
		if err != nil { goto exit }

		pmem, err := strconv.ParseFloat(parsed_line[0], 64)
		if err != nil { goto exit }

		pid, err := strconv.ParseUint(parsed_line[2], 10, 32)
		if err != nil { goto exit }

		zoneid, err := strconv.ParseUint(parsed_line[3], 10, 32)
		if err != nil { goto exit }

		comm := parsed_line[4]

		args = ""
		for i := range parsed_line[5:len(parsed_line)] {
			args += parsed_line[5+i] + " "
		}

		out = append(out, psLineDesc {
			pcpu: pcpu,
			pmem: pmem,
			pid: uint(pid),
			zoneid: uint(zoneid),
			comm: comm,
			args: args,
		})
	}
exit:
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *PsCollector) getPsOut() error {
	out, eerr := exec.Command("ps", "-eo", "pcpu,pmem,pid,zoneid,comm,args").Output()
	if eerr != nil {
		level.Error(c.logger).Log("error on executing ps: %v", eerr)
	} else {
		psOutputParsed, perr := parsePsOutput(string(out))
		if perr != nil {
			level.Error(c.logger).Log("error on parsing ps out: %v", perr)
		}

		psCpuOutput := psOutputParsed
		psMemOutput := psOutputParsed
		psLen := min(len(psOutputParsed), MAX_PROCESSES_NUM, CONFIG_PROCESSES_NUM)

		slices.SortFunc(psCpuOutput, cmpPsCpu)
		slices.SortFunc(psMemOutput, cmpPsMem)

		//We must only leave the `psLen` right items
		psCpuOutput = psCpuOutput[len(psOutputParsed) - psLen : len(psOutputParsed)]
		psMemOutput = psMemOutput[len(psOutputParsed) - psLen : len(psOutputParsed)]

		c.psCpu.Reset()
		c.psMem.Reset()

		index := 0
		for _,l:= range psCpuOutput {
			if l.pcpu < CONFIG_LOW_CPU_THRESHOLD { continue }
			c.psCpu.With(prometheus.Labels{
				"index":	fmt.Sprintf("%d", index),
				"pid": 		fmt.Sprintf("%d", l.pid), 
				"zoneid": 	fmt.Sprintf("%d", l.zoneid),
				"comm":		l.comm,
				"args":		l.args,
			}).Set(l.pcpu)
			index++
		}

		index = 0
		for _,l:= range psMemOutput {
			if l.pcpu < CONFIG_LOW_MEM_THRESHOLD { continue }
			c.psMem.With(prometheus.Labels{
				"index":	fmt.Sprintf("%d", index),
				"pid": 		fmt.Sprintf("%d", l.pid), 
				"zoneid": 	fmt.Sprintf("%d", l.zoneid),
				"comm":		l.comm,
				"args":		l.args,
			}).Set(l.pmem)
			index++
		}
	}
	return nil
}
