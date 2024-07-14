// zpool collector
// this will :
//  - call zpool list
//  - gather ZPOOL metrics
//  - feed the collector

package collector

import (
	"os/exec"
	"strconv"
	"strings"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	// Prometheus Go toolset
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("zpool", defaultEnabled, NewGZZpoolListExporter)
}

// GZZpoolListCollector declares the data type within the prometheus metrics package.
type GZZpoolListCollector struct {
	gzZpoolListAlloc	*prometheus.GaugeVec
	gzZpoolListFrag		*prometheus.GaugeVec
	gzZpoolListFree		*prometheus.GaugeVec
	gzZpoolListSize		*prometheus.GaugeVec
	gzZpoolListFreeing	*prometheus.GaugeVec
	gzZpoolListHealth	*prometheus.GaugeVec
	gzZpoolListLeaked	*prometheus.GaugeVec
	gzZpoolListGuid		*prometheus.GaugeVec
	gzZfsLogicalUsed	*prometheus.GaugeVec
	logger	log.Logger
}

// NewGZZpoolListExporter returns a newly allocated exporter GZZpoolListCollector.
// It exposes the zpool list command result.
func NewGZZpoolListExporter(logger log.Logger) (Collector, error) {

	return &GZZpoolListCollector{

		gzZpoolListAlloc: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "zpool_alloc_mbytes",
			Help: "zpool allocated, megabytes.",
		}, []string{"zpool"}),

		gzZpoolListFrag: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "zpool_frag_percents",
			Help: "zpool fragmentation, percents.",
		}, []string{"zpool"}),

		gzZpoolListFree: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "zpool_free_mbytes",
			Help: "zpool free, megabytes.",
		}, []string{"zpool"}),

		gzZpoolListSize: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "zpool_size_mbytes",
			Help: "zpool size, megabytes.",
		}, []string{"zpool"}),

		gzZpoolListFreeing: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "zpool_freeing_mbytes",
			Help: "zpool freeing, megabytes.",
		}, []string{"zpool"}),

		gzZpoolListHealth: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "zpool_health",
			Help: "zpool health status (0: OFFLINE, 1: ONLINE)",
		}, []string{"zpool"}),

		gzZpoolListLeaked: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "zpool_leaked_mbytes",
			Help: "zpool leaked mbytes.",
		}, []string{"zpool"}),

		gzZpoolListGuid: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "zpool_guid",
			Help: "zpool guid.",
		}, []string{"zpool"}),

		gzZfsLogicalUsed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "zpool_zfs_logicalused_mbytes",
			Help: "zfs logicalused.",
		}, []string{"zpool"}),
		logger: logger,

	}, nil
}

// Describe describes all the metrics.
func (e *GZZpoolListCollector) Describe(ch chan<- *prometheus.Desc) {
	e.gzZpoolListAlloc.Describe(ch)
	e.gzZpoolListFrag.Describe(ch)
	e.gzZpoolListFree.Describe(ch)
	e.gzZpoolListSize.Describe(ch)		
	e.gzZpoolListFreeing.Describe(ch)
	e.gzZpoolListHealth.Describe(ch)
	e.gzZpoolListLeaked.Describe(ch)
	e.gzZpoolListGuid.Describe(ch)
	e.gzZfsLogicalUsed.Describe(ch)
}

// Collect fetches the stats.
func (e *GZZpoolListCollector) Update(ch chan<- prometheus.Metric) error {
	e.zpoolGet()
	e.zfsGet()
	e.gzZpoolListAlloc.Collect(ch)
	e.gzZpoolListFrag.Collect(ch)
	e.gzZpoolListFree.Collect(ch)
	e.gzZpoolListSize.Collect(ch)		
	e.gzZpoolListFreeing.Collect(ch)
	e.gzZpoolListHealth.Collect(ch)
	e.gzZpoolListLeaked.Collect(ch)
	e.gzZpoolListGuid.Collect(ch)
	e.gzZfsLogicalUsed.Collect(ch)
	return nil;
}

func (e *GZZpoolListCollector) zpoolGet() error {
	out, eerr := exec.Command("zpool", "get", "-Hp", 
		"size,free,allocated,fragmentation,freeing,health,allocated,leaked,guid").Output()
	if eerr != nil {
		level.Error(e.logger).Log("error on executing zpool: %v", eerr)
	} else {
		perr := e.parseZpoolGetOutput(string(out))
		if perr != nil {
			level.Error(e.logger).Log("error on parsing zpool: %v", perr)
		}
	}
	return nil
}

//Yes, zfs get. Though we already have a dedicated collector for zfs,
//but here we need to retrieve only some pool-related statistics
func (e *GZZpoolListCollector) zfsGet() error {
	out, eerr := exec.Command("zfs", "get", "-Hp", "logicalused").Output()
	if eerr != nil {
		level.Error(e.logger).Log("error on executing zfs: %v", eerr)
	} else {
		perr := e.parseZfsGetOutput(string(out))
		if perr != nil {
			level.Error(e.logger).Log("error on parsing zpool: %v", perr)
		}
	}
	return nil
}


func (e *GZZpoolListCollector) parseZpoolGetOutput(out string) error {
	outlines := strings.Split(out, "\n")
	l := len(outlines)

	for _, line := range outlines[0 : l-1] {
		parsed_line := strings.Fields(line)
		pool_name := parsed_line[0]
		val := parsed_line[2]

		switch parsed_line[1] {
		case "size":
			pval, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			e.gzZpoolListSize.With(prometheus.Labels{"zpool": pool_name}).Set(pval / 1024 / 1024)
		case "free":
			pval, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			e.gzZpoolListFree.With(prometheus.Labels{"zpool": pool_name}).Set(pval / 1024 / 1024)
		case "allocated":
			pval, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			e.gzZpoolListAlloc.With(prometheus.Labels{"zpool": pool_name}).Set(pval / 1024 / 1024)
		case "fragmentation":
			pval, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			e.gzZpoolListFrag.With(prometheus.Labels{"zpool": pool_name}).Set(pval)
		case "freeing":
			pval, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			e.gzZpoolListFreeing.With(prometheus.Labels{"zpool": pool_name}).Set(pval)
		case "health":
			if (strings.Contains(val, "ONLINE")) == true {
				e.gzZpoolListHealth.With(prometheus.Labels{"zpool": pool_name}).Set(1)
			} else {
				e.gzZpoolListHealth.With(prometheus.Labels{"zpool": pool_name}).Set(0)
			}
		case "leaked":
			pval, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			e.gzZpoolListFrag.With(prometheus.Labels{"zpool": pool_name}).Set(pval / 1024 / 1024)
/*
		case "guid":
			pval, err := strconv.ParseFloat(val, 64)
			level.Error(e.logger).Log(pval)
			if err != nil {
				return err
			}
			e.gzZpoolListGuid.With(prometheus.Labels{"zpool": pool_name}).Set(pval)
*/
		}
	}
	return nil
}

func (e *GZZpoolListCollector) parseZfsGetOutput(out string) error {
	outlines := strings.Split(out, "\n")
	l := len(outlines)
	for _, line := range outlines[0 : l-1] {
		parsed_line := strings.Fields(line)
		name := parsed_line[0]
		pval, err := strconv.ParseFloat(parsed_line[2], 64)
		if err != nil {
			return err
		}

		if strings.Contains(name, "/") {
			continue
		}
		e.gzZfsLogicalUsed.With(prometheus.Labels{"zpool": name}).Set(float64(pval / 1024 / 1024))
	}
	return nil
}
