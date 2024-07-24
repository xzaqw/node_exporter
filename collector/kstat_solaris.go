package collector

import (
	"fmt"
//	"strconv"
	"github.com/go-kit/log"
	"github.com/illumos/go-kstat"
	"github.com/prometheus/client_golang/prometheus"
	//"C"
)

// #include <unistd.h>
import "C"

type kstatValue struct {
	kstat kstat.KStat
	desc typedDesc
}

type kstatCollector struct {
	values []kstatValue
	logger log.Logger
}

func init() {
	registerCollector("kstat", defaultEnabled, NewKstatCollector)
}

func NewKstatCollector(logger log.Logger) (Collector, error) {
	var (	c kstatCollector
	)
	c = c
	for _,module := range statistics {	
		fmt.Print(module, "\n")

	}
	c.logger = logger
	return &c, nil
}

func (c *kstatCollector) Update(ch chan<- prometheus.Metric) error {
/*
	var (	kstatValue *kstat.Named
		err error
	)
	ncpus := C.sysconf(C._SC_NPROCESSORS_ONLN)

	kstatValue = kstatValue
	ncpus = ncpus

	tok, err := kstat.Open()
	if err != nil {
		return err
	}

	defer tok.Close()
*/

/*
	for cpu := 0; cpu < int(ncpus); cpu++ {
		ksCPU, err := tok.Lookup("cpu", cpu, "sys")
		if err != nil { goto exit }

		for k, v := range map[string]string{
			"idle":   "cpu_nsec_idle",
			"kernel": "cpu_nsec_kernel",
			"user":   "cpu_nsec_user",
			"intr":   "cpu_nsec_intr",
			"dtrace": "cpu_nsec_dtrace",
		} {
			kstatValue, err = ksCPU.GetNamed(v)
			if (err != nil) { goto exit }
			ch <- c.cpu_seconds.mustNewConstMetric(
				float64(kstatValue.UintVal)/1e9, strconv.Itoa(cpu), k)
		}
	}
*/

/*
	for _,v := range kstatInit {
		fmt.Print(v, "\n")
	}

	goto exit
exit:
	if err != nil {
		return err
	}
*/
	return nil
}
