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
	PROCESSES_NUM = 10
)

type PsCollector struct {
	psCpu [PROCESSES_NUM]*prometheus.GaugeVec
	psMem [PROCESSES_NUM]*prometheus.GaugeVec
	logger	log.Logger
}

func init() {
	registerCollector("ps", defaultEnabled, NewPsCollector)
}

func NewPsCollector(logger log.Logger) (Collector, error) {
	cpuGaugeList := [PROCESSES_NUM]*prometheus.GaugeVec{}
	memGaugeList := [PROCESSES_NUM]*prometheus.GaugeVec{}

	for i := range cpuGaugeList {
		cpuGaugeList[i] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: fmt.Sprintf("node_ps_cpu_process_%d", i),
			Help: fmt.Sprintf("Process index %d from processes list sorted by CPU utilization.", i),
		}, []string{"pid","zoneid","args"})
	}
	for i := range memGaugeList {
		memGaugeList[i] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: fmt.Sprintf("node_ps_cpu_process_%d", i),
			Help: fmt.Sprintf("Process index %d from processes list sorted by memory utilization.", i),
		}, []string{"pid","zoneid","args"})
	}

	return &PsCollector {
		psCpu: cpuGaugeList,
		psMem: memGaugeList,
		logger: logger,
	}, nil
}

func (c *PsCollector) Update(ch chan<- prometheus.Metric) error {
	c.getPsCpu()

	for i := range c.psCpu {
		c.psCpu[i].Collect(ch)
	}

	return nil
}

func (c *PsCollector) Describe(ch chan<- *prometheus.Desc) {
	for i := range c.psCpu {
		c.psCpu[i].Describe(ch)
	}
}

type psCpuLineDesc struct {
	pcpu float64
	pid uint
	zoneid uint
	args string
}

func parsePsCpuOutput(psCpuOut string) ([]psCpuLineDesc , error) {
	var err error
	var out []psCpuLineDesc 
	var args string

	psOutLines := strings.Split(psCpuOut, "\n")

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

		pid, err := strconv.ParseUint(parsed_line[1], 10, 32)
		if err != nil { goto exit }

		zoneid, err := strconv.ParseUint(parsed_line[2], 10, 32)
		if err != nil { goto exit }

		args = ""
		for i := range parsed_line[3:len(parsed_line)] {
			args += parsed_line[3+i] + " "
		}

		out = append(out, psCpuLineDesc {
			pcpu: pcpu,
			pid: uint(pid),
			zoneid: uint(zoneid),
			args: args,
		})
	}
exit:
	if err != nil {
		return nil, err
	}

	return out, nil
}

func cmpPsCpu(a, b psCpuLineDesc) int {
	return cmp.Compare(a.pcpu, b.pcpu)
}

func (c *PsCollector) getPsCpu() error {
	out, eerr := exec.Command("ps", "-eo", "pcpu,pid,zoneid,args").Output()
	if eerr != nil {
		level.Error(c.logger).Log("error on executing ps: %v", eerr)
	} else {
		psCpuParsed, perr := parsePsCpuOutput(string(out))

		if perr != nil {
			level.Error(c.logger).Log("error on parsing ps out: %v", perr)
		}

		slices.SortFunc(psCpuParsed, cmpPsCpu)

		outLen := len(psCpuParsed)
		for i := 0; (i < len(c.psCpu)) && (i < outLen); i++ {
			line := psCpuParsed[outLen - 1 - i]
			c.psCpu[i].With(prometheus.Labels{
				"pid": 		fmt.Sprintf("%d", line.pid), 
				"zoneid": 	fmt.Sprintf("%d", line.zoneid),
				"args":		line.args,
			}).Set(line.pcpu)
		}
	}
	return nil
}
