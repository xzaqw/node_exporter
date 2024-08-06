package collector

import (
    "io/ioutil"
    "path/filepath"
    "gopkg.in/yaml.v2"
)

type kstatConfig struct {
	KstatModules []KstatModule `json:"kstat_modules"`
}

type KstatModule struct {
	ID         string      `json:"id"`
	KstatNames []KstatName `json:"kstat_names"`
}

type KstatName struct {
	ID         string      `json:"id"`
	KstatStats []KstatStat `json:"kstat_stats"`
}

type KstatStat struct {
	ID          string      `json:"id"`
	Suffix      string      `json:"suffix"`
	Type        string      `json:"type"`
	Help        string `json:"help"`
	ScaleFactor float64      `json:"scale_factor"`
	LabelString string `json:"label_string"`
}

var kstatConfigInstance = kstatConfig {
        KstatModules: []KstatModule {
		{ 
			ID: "cpu",
                  	KstatNames: []KstatName {
				{
					ID: "sys",
					KstatStats: []KstatStat {
						{
						ID: "bawrite",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "bread",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "dtrace_probes",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "intr",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "intrblk",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "intrthread",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "lread",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "lwrite",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "modload",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "cpu_ticks_idle",
						Suffix: "total",
						Type: "counter",
						Help: "Ticks the CPU spent in idle mode",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "cpu_ticks_kernel",
						Suffix: "total",
						Type: "counter",
						Help: "Ticks the CPU spent in kernel mode",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "cpu_ticks_user",
						Suffix: "total",
						Type: "counter",
						Help: "Ticks the CPU spent in user mode",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "cpu_ticks_wait",
						Suffix: "total",
						Type: "counter",
						Help: "Ticks the CPU spent in wait  mode",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "cpu_nsec_dtrace",
						Suffix: "seconds_total",
						Type: "counter",
						Help: "Seconds the CPU spent in dtrace mode",
						ScaleFactor: 1e-9,
						LabelString: "cpu",
						},
						{
						ID: "cpu_nsec_idle",
						Suffix: "seconds_total",
						Type: "counter",
						Help: "Seconds the CPU spent in idle mode",
						ScaleFactor: 1e-9,
						LabelString: "cpu",
						},
						{
						ID: "cpu_nsec_intr",
						Suffix: "seconds_total",
						Type: "counter",
						Help: "Seconds the CPU spent in interrupt mode",
						ScaleFactor: 1e-9,
						LabelString: "cpu",
						},
						{
						ID: "cpu_nsec_kernel",
						Suffix: "seconds_total",
						Type: "counter",
						Help: "Seconds the CPU spent in kernel mode",
						ScaleFactor: 1e-9,
						LabelString: "cpu",
						},
						{
						ID: "cpu_nsec_user",
						Suffix: "seconds_total",
						Type: "counter",
						Help: "Seconds the CPU spent in user mode",
						ScaleFactor: 1e-9,
						LabelString: "cpu",
						},
						{
						ID: "cpu_load_intr",
						Suffix: "percents",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "cpumigrate",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "iowait",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "nthreads",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "syscall",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "sysexec",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "sysfork",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "sysread",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "sysvfork",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "syswrite",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "trap",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "idlethread",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "inv_swtch",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "mutex_adenters",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
						{
						ID: "xcalls",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						LabelString: "cpu",
						},
					},
				},
	                        { 
					ID: "vm",
					KstatStats: []KstatStat{
						{
						ID: "pgswapin",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						},
						{
						ID: "pgswapout",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						},
                        		},
				},
			},
		},
                { 
			ID: "unix",
                  	KstatNames: []KstatName {
                        	{
					ID: "system_pages",
					KstatStats: []KstatStat{
                                		{
                                		ID: "pagesfree",
                                		Suffix: "total",
                                		Type: "counter",
                                		Help: "",
                                		ScaleFactor: 1,
                                		},
						{
						ID: "pageslocked",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						},
						{
						ID: "pagestotal",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						},
						{
						ID: "physmem",
						Suffix: "total",
						Type: "counter",
						Help: "",
						ScaleFactor: 1,
						},
					},
				},
                        	{ 
					ID: "pset",
                          		KstatStats: []KstatStat{
                                		{
						ID: "avenrun_15min",
						Suffix: "percents",
						Type: "counter",
						Help: "15 min CPU load average",
						ScaleFactor: 1,
						},
						{
						ID: "avenrun_5min",
						Suffix: "percents",
						Type: "counter",
						Help: "5 min CPU load average",
						ScaleFactor: 1,
						},
						{
						ID: "avenrun_1min",
						Suffix: "percents",
						Type: "counter",
						Help: "1 min CPU load average",
						ScaleFactor: 1,
						},
                        		},
				},
			},
		},
	},
}

func (config *kstatConfig) init() (error) {
	var (
		err error
	)
	filename, _ := filepath.Abs("./collector/kstat_config.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil { return err }

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil { return err }

	return nil
}
