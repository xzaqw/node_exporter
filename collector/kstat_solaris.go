package collector

import (
	"strconv"
	"github.com/go-kit/log"
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
		err error
	)

	err = cfg.init()
	if err != nil {
		return nil, err 
	}


	for _, cfgModule := range cfg.KstatModules {
		module := kstatModule{}
		module.ID = cfgModule.ID
		for _,cfgName := range cfgModule.KstatNames {
			name := kstatName{}
			name.ID = cfgName.ID
			for _, cfgStat := range cfgName.KstatStats {
				stat := kstatStat{}
				stat.ID = cfgStat.ID
				stat.desc = typedDesc{prometheus.NewDesc(
					prometheus.BuildFQName(
						namespace, 
						"kstat_" + cfgModule.ID + "_" + cfgName.ID,
						cfgStat.ID + "_" + cfgStat.Suffix),
					cfgStat.Help, []string{"inst"}, nil,
					), 
					prometheus.CounterValue}
				stat.scaleFactor = cfgStat.ScaleFactor
				name.stats = append(name.stats, stat)
			}
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
	)

	tok, err := kstat.Open()
	if err != nil { goto exit }

	defer tok.Close()

	for _,module := range c.modules {
		for _,name := range module.names {
			//Walk through all instances
			inst := 0
			for {
				ksName, err := tok.Lookup(module.ID, inst, name.ID)
				if err != nil { break }

				for _,stat := range name.stats {
					kstatValue, err = ksName.GetNamed(stat.ID)
					if (err != nil) {
						goto exit 
					}
					v := float64(kstatValue.UintVal) * stat.scaleFactor

					//Round the value down to the number integer value 
					//like 2.45 to 2.0. At the same time we have 
					//to stick to float64 type.
					ch <- stat.desc.mustNewConstMetric(float64(int(v)), strconv.Itoa(inst))
				}
				inst++
			}
		}
	}

exit:
	if err != nil {
		return err
	}
	return nil
}
