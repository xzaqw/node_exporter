package collector

import (
//	"fmt"
	"strconv"
	"strings"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/illumos/go-kstat"
	"github.com/prometheus/client_golang/prometheus"
)

type kstatStat struct {
	ID string
	scaleFactor float64
	desc typedDesc
}

type kstatName struct {
	ID string
	stats []kstatStat
}

type kstatModule struct {
	ID string
	names []kstatName
}

type kstatCollector struct {
	modules []kstatModule
	logger log.Logger
}

func init() {
	registerCollector("kstat", defaultEnabled, NewKstatCollector)
}

func NewKstatCollector(logger log.Logger) (Collector, error) {
	var (	c kstatCollector
		cfg kstatConfig
		label string
	)

	/*
	cfg, err = cfg.init()
	if err != nil {
		return nil, err 
	}
	*/

	cfg = kstatConfigInstance
	for _, cfgModule := range cfg.KstatModules {
		module := kstatModule{}
		module.ID = cfgModule.ID
		for _,cfgName := range cfgModule.KstatNames {
			name := kstatName{}
			name.ID = cfgName.ID
			for _, cfgStat := range cfgName.KstatStats {
				if cfgStat.LabelString == "" {
					label = "instance"
				} else {
					label = cfgStat.LabelString
				}
				stat := kstatStat{}
				stat.ID = cfgStat.ID
				desc := prometheus.NewDesc(
					prometheus.BuildFQName(
						namespace, 
						"kstat_" + cfgModule.ID + "_" + cfgName.ID,
						cfgStat.ID + "_" + cfgStat.Suffix),
						cfgStat.Help, []string{label}, nil, )
				stat.desc = typedDesc{ desc, prometheus.CounterValue }
				stat.scaleFactor = cfgStat.ScaleFactor
				name.stats = append(name.stats, stat)
			}

			//Snaptime is separate kind because of 
			//different way to retrieve this metric
			stat := kstatStat{}
			stat.ID = "snaptime"
			desc := prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace, 
					"kstat_" + cfgModule.ID + "_" + cfgName.ID,
					"snaptime"),
					cfgModule.ID + "::" + cfgName.ID + ":" + "snaptime", 
					[]string{"inst"}, nil, )
			stat.desc = typedDesc{ desc, prometheus.CounterValue }
			name.stats = append(name.stats, stat)
			module.names = append(module.names, name)
		}
		c.modules = append(c.modules, module)
	}

	c.logger = logger
	
	return &c, nil
}

func (c *kstatCollector) Update(ch chan<- prometheus.Metric) error {
	var (	kstatValue *kstat.Named
		err error
		metricValue float64
	)

	tok, err := kstat.Open()
	if err != nil { 
		return err 
	}

	defer tok.Close()

	for _,module := range c.modules {
		for _,name := range module.names {
			inst := 0
			for {
				ksName, err := tok.Lookup(module.ID, inst, name.ID)
				if err != nil {
					//Handle the instance number out-of-bound error
					break 
				}
				for _,stat := range name.stats {
					if strings.HasSuffix(stat.ID, "snaptime") {
						metricValue = float64(ksName.Snaptime)
					} else {
						kstatValue, err = ksName.GetNamed(stat.ID)
						if (err != nil) {
							level.Error(c.logger).Log(module.ID + ":" + 
							strconv.Itoa(inst) + ":" + name.ID + ":" + stat.ID, err)
							continue
						}
						metricValue = float64(kstatValue.UintVal) * stat.scaleFactor
					}

					//Round the value down to the number integer value 
					//like 2.45 to 2.0. At the same time we have 
					//to stick to float64 type.
					ch <- stat.desc.mustNewConstMetric(float64(int(metricValue)), 
						strconv.Itoa(inst))
				}
				inst++	
			}
		}
	}
	return nil
}
