// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !nocpu
// +build !nocpu

package collector

import (
	"strconv"

	"github.com/go-kit/log"
	"github.com/illumos/go-kstat"
	"github.com/prometheus/client_golang/prometheus"
)

// #include <unistd.h>
import "C"

type cpuCollector struct {
	cpu_seconds typedDesc
	cpu_ticks typedDesc
	cpu_load_intr typedDesc
	cpumigrate typedDesc
	logger log.Logger
}

func init() {
	registerCollector("cpu", defaultEnabled, NewCpuCollector)
}

func NewCpuCollector(logger log.Logger) (Collector, error) {
	return &cpuCollector{
		cpu_seconds: typedDesc{nodeCPUSecondsDesc, 
			prometheus.CounterValue},

		cpu_ticks: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "ticks_total"),
				"Ticks the CPUs spent in each mode.",
				[]string{"cpu", "mode"}, nil,
			), prometheus.CounterValue},

		cpu_load_intr: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "load_intr_percents"),
				"Interrupt load factor, percents.",
				[]string{"cpu"}, nil,
			), prometheus.GaugeValue},

		cpumigrate: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "cpumigrate_total"),
				"CPU migrations by threads.",
				[]string{"cpu"}, nil,
			), prometheus.CounterValue},

		logger: logger,
	}, nil
}

func (c *cpuCollector) Update(ch chan<- prometheus.Metric) error {

	var (	kstatValue *kstat.Named
		err error
	)

	ncpus := C.sysconf(C._SC_NPROCESSORS_ONLN)

	tok, err := kstat.Open()
	if err != nil {
		return err
	}

	defer tok.Close()

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

		for k, v := range map[string]string{
			"idle":   "cpu_ticks_idle",
			"kernel": "cpu_ticks_kernel",
			"user":   "cpu_ticks_user",
			"intr":   "cpu_ticks_wait",
		} {
			kstatValue, err = ksCPU.GetNamed(v)
			if err != nil { goto exit }
			ch <- c.cpu_ticks.mustNewConstMetric(
				float64(kstatValue.UintVal), strconv.Itoa(cpu), k)
		}

		kstatValue, err = ksCPU.GetNamed("cpu_load_intr")
		if err != nil { goto exit }
		ch <- c.cpu_load_intr.mustNewConstMetric(
			float64(kstatValue.UintVal), strconv.Itoa(cpu))

		kstatValue, err = ksCPU.GetNamed("cpumigrate")
		if err != nil { goto exit }
		ch <- c.cpumigrate.mustNewConstMetric(
			float64(kstatValue.UintVal), strconv.Itoa(cpu))
	}
exit:
	if err != nil {
		return err
	}
	return nil
}
