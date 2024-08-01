package collector

import (
	"os/exec"
	"strings"
	"strconv"
	"fmt"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)


type netCollector struct {
	iPackets	*prometheus.GaugeVec
	oPackets	*prometheus.GaugeVec
	rBytes		*prometheus.GaugeVec
	oBytes		*prometheus.GaugeVec
	iErrors		*prometheus.GaugeVec
	oErrors 	*prometheus.GaugeVec
	class		*prometheus.GaugeVec
	mtu		*prometheus.GaugeVec
	state		*prometheus.GaugeVec
	bridge		*prometheus.GaugeVec
	over		*prometheus.GaugeVec
}

const (
	netCollectorSubsystem = "net"
)

func init() {
	registerCollector(netCollectorSubsystem , defaultEnabled, NewNetCollector)
}

func NewNetCollector(logger log.Logger) (Collector, error) {
	return &netCollector {
		iPackets: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "net_link_ipackets",
			Help: "Link input packets",
		}, []string{"link"}),
		oPackets: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "net_link_opackets",
			Help: "Link output packets",
		}, []string{"link"}),
		rBytes: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "net_link_rbytes",
			Help: "Link received bytes",
		}, []string{"link"}),
		oBytes: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "net_link_obytes",
			Help: "Link transmitted bytes",
		}, []string{"link"}),
		iErrors: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "net_link_ierrors",
			Help: "Link receive errors",
		}, []string{"link"}),
		oErrors: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "net_link_oerrors",
			Help: "Link output errors",
		}, []string{"link"}),
		class: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "net_link_class",
			Help: "Link class",
		}, []string{"link", "class"}),
		mtu: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "net_link_mtu",
			Help: "Link MTU",
		}, []string{"link", "mtu"}),
		state: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "net_link_state",
			Help: "Link state",
		}, []string{"link", "state"}),
		bridge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "net_link_bridge",
			Help: "Link bridge",
		}, []string{"link", "bridge"}),
		over: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "net_link_over",
			Help: "Link over",
		}, []string{"link", "over"}),
	}, nil
}


func (c *netCollector) dladmConfGet() error {
//dladm show-link -po link,class,mtu,state,bridge,over
	var (
		link string
		class string
		mtu string
		state string
		bridge string
		over string
	)
	out, err := exec.Command("dladm", "show-link", "-po", 
		"link,class,mtu,state,bridge,over").Output()
	if (err != nil) {
		return err
	}
	outlines := strings.Split(string(out), "\n")
	for _,l := range(outlines) {
		values := strings.Split(l, ":")
		if values[0] == "" {
			continue
		}
		link 	= values[0]
		class 	= values[1]
		mtu 	= values[2]
		state 	= values[3]
		bridge 	= values[4]
		over 	= values[5]

		c.class.With(prometheus.Labels{"link": link, "class": class}).Set(0)
		c.mtu.With(prometheus.Labels{"link": link, "mtu": mtu}).Set(0)
		c.state.With(prometheus.Labels{"link": link, "state": state}).Set(0)
		c.bridge.With(prometheus.Labels{"link": link, "bridge": bridge}).Set(0)
		c.over.With(prometheus.Labels{"link": link, "over": over}).Set(0)
	}
	return nil
}

func (c *netCollector) dladmStatsGet() error {
//dladm show-link -pso link,ipackets,opackets,rbytes,obytes,ierrors,oerrors
	var (
		link string
		ipackets, opackets, rbytes,
		obytes,ierrors,oerrors uint64
		err error
	)
	out, err := exec.Command("dladm", "show-link", "-pso", 
		"link,ipackets,opackets,rbytes,obytes,ierrors,oerrors").Output()
	if (err != nil) {
		return err
	}

	outlines := strings.Split(string(out), "\n")
	for _,l := range(outlines) {
		values := strings.Split(l, ":")
		if values[0] == "" {
			continue
		}
		link 		= values[0]
		ipackets, err = strconv.ParseUint(values[1], 10, 64)
		if (err != nil) { return err }
		opackets, err = strconv.ParseUint(values[2], 10, 64)
		if (err != nil) { return err }
		rbytes, err = strconv.ParseUint(values[3], 10, 64)
		if (err != nil) { return err }
		obytes, err = strconv.ParseUint(values[4], 10, 64)
		if (err != nil) { return err }
		ierrors, err = strconv.ParseUint(values[5], 10, 64)
		if (err != nil) { return err }
		oerrors, err = strconv.ParseUint(values[6], 10, 64)
		if (err != nil) { return err }

		c.iPackets.With(prometheus.Labels{"link": link}).Set(float64(ipackets))
		c.oPackets.With(prometheus.Labels{"link": link}).Set(float64(opackets))
		c.rBytes.With(prometheus.Labels{"link": link}).Set(float64(rbytes))
		c.oBytes.With(prometheus.Labels{"link": link}).Set(float64(obytes))
		c.iErrors.With(prometheus.Labels{"link": link}).Set(float64(ierrors))
		c.oErrors.With(prometheus.Labels{"link": link}).Set(float64(oerrors))
		ipackets = ipackets
		opackets = opackets
		rbytes = rbytes
		obytes = obytes
		ierrors = ierrors
		oerrors = oerrors
	}
	return nil
}

func (c *netCollector) Update(ch chan<- prometheus.Metric) error {
	c.dladmConfGet()
	c.dladmStatsGet()
	c.iPackets.Collect(ch)
	c.oPackets.Collect(ch)
	c.rBytes.Collect(ch)
	c.oBytes.Collect(ch)
	c.iErrors.Collect(ch)
	c.oErrors.Collect(ch)
	c.class.Collect(ch)
	c.mtu.Collect(ch)
	c.state.Collect(ch)
	c.bridge.Collect(ch)
	c.over.Collect(ch)
	return nil;
}

func parseDladmOutput(out string) error {
	//var err error
	outlines := strings.Split(out, "\n")
	for _,l := range(outlines) {
		values := strings.Split(l, ":")
		if values[0] == "" {
			continue
		}
		fmt.Print(values)
		fmt.Print("\n")
	}
	return nil
}
