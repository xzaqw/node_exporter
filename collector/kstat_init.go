package collector

/*
cpu:0:sys:iowait
*/
type Statistic struct {
        name string
        scaling_factor float64
        help string
}

var statistics = map[string] map[string] []map[string] []map[string] Statistic {
        "cpu":{
                "instances": {},
                "names": {
                        {"sys":
                                {

                                        {"iowait": {
                                                name: "wef",
                                                help: "wef",
                                                },
                                        },
                                        {"syscall": {
                                                name: "wef",
                                                help: "wef",
                                                },
                                        },

                                },
                        },
                        {"vm":
                                {

                                        {"anonfree": {
                                                name: "iowait_total",
                                                help: "IOWAIT HELP",
                                                },
                                        },
                                },
                        },

                },
        },
}

/*
var kstatPaths = []kstatPath {
	{
		"cpu",
		"cpu_ticks_idle",
		"cpu_ticks_idle_total",
		"Ticks the CPUs spent in each mode.",
		1,
	},
}
*/
