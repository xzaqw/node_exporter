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

type psLineDesc struct {
	pcpu float64
	pmem float64
	pid uint
	zoneid uint
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
	cpuGaugeList := [PROCESSES_NUM]*prometheus.GaugeVec{}
	memGaugeList := [PROCESSES_NUM]*prometheus.GaugeVec{}

	for i := range cpuGaugeList {
		cpuGaugeList[i] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: fmt.Sprintf("node_ps_cpu_process_%d", i),
			Help: fmt.Sprintf("Process index %d from processes list sorted by CPU utilization.", i),
		}, []string{
			"pid",
			"zoneid",
			"args"})
	}
	for i := range memGaugeList {
		memGaugeList[i] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: fmt.Sprintf("node_ps_cpu_process_%d", i),
			Help: fmt.Sprintf("Process index %d from processes list sorted by memory utilization.", i),
		}, []string{
			"pid",
			"zoneid",
			"args"})
	}

	return &PsCollector {
		psCpu: cpuGaugeList,
		psMem: memGaugeList,
		logger: logger,
	}, nil
}

func (c *PsCollector) Update(ch chan<- prometheus.Metric) error {
	c.getPsOut()

	for i := range c.psCpu {
		c.psCpu[i].Collect(ch)
		//c.psMem[i].Collect(ch)
	}

	return nil
}

func (c *PsCollector) Describe(ch chan<- *prometheus.Desc) {
	for i := range c.psCpu {
		c.psCpu[i].Describe(ch)
	}
	for i := range c.psMem {
		c.psMem[i].Describe(ch)
	}
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

		args = ""
		for i := range parsed_line[4:len(parsed_line)] {
			args += parsed_line[4+i] + " "
		}

		out = append(out, psLineDesc {
			pcpu: pcpu,
			pmem: pmem,
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

//[from Gatherer #2] collected metric node_ps_cpu_process_9 label:{name:\"args\"  value:\"/usr/sbin/sshd -R \"}  label:{name:\"pid\"  value:\"2019\"}  label:{name:\"zoneid\"  value:\"0\"}  gauge:{value:0} has help \"Process index 9 from processes list sorted by memory utilization.\" but should have \"Process index 9 from processes list sorted by CPU utilization.\""

func (c *PsCollector) getPsOut() error {
	out, eerr := exec.Command("ps", "-eo", "pcpu,pmem,pid,zoneid,args").Output()
	if eerr != nil {
		level.Error(c.logger).Log("error on executing ps: %v", eerr)
	} else {
		psOutputParsed, perr := parsePsOutput(string(out))
		if perr != nil {
			level.Error(c.logger).Log("error on parsing ps out: %v", perr)
		}

		psCpuOutput := psOutputParsed
		psMemOutput := psOutputParsed
		psLen := len(psOutputParsed)

		slices.SortFunc(psCpuOutput, cmpPsCpu)
		slices.SortFunc(psMemOutput, cmpPsMem)

		for i := 0; (i < len(c.psCpu)) && (i < psLen); i++ {
			line := psCpuOutput [psLen - 1 - i]
			c.psCpu[i].With(prometheus.Labels{
				"pid": 		fmt.Sprintf("%d", line.pid), 
				"zoneid": 	fmt.Sprintf("%d", line.zoneid),
				"args":		line.args,
			}).Set(line.pcpu)
		}

		for i := 0; (i < len(c.psMem)) && (i < psLen); i++ {
			line := psMemOutput [psLen - 1 - i]
			c.psMem[i].With(prometheus.Labels{
				"pid": 		fmt.Sprintf("%d", line.pid), 
				"zoneid": 	fmt.Sprintf("%d", line.zoneid),
				"args":		line.args,
			}).Set(line.pmem)
		}
	}
	return nil
}
