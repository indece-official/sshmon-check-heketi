package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/heketi/heketi/client/api/go-client"
	"github.com/heketi/heketi/pkg/glusterfs/api"
	"github.com/miekg/dns"
)

// Variables set during build
var (
	ProjectName  string
	BuildVersion string
	BuildDate    string
)

var statusMap = []string{
	"OK",
	"WARN",
	"CRIT",
	"UNKNOWN",
}

var (
	flagVersion = flag.Bool("v", false, "Print the version info and exit")
	flagService = flag.String("service", "", "Service name (defaults to Heketi_<host>)")
	flagHost    = flag.String("host", "", "Host")
	flagUser    = flag.String("user", "", "Username")
	flagKey     = flag.String("key", "", "Key")
	flagPort    = flag.Int("port", 5080, "Port")
	flagDNS     = flag.String("dns", "", "Use other dns server")
)

func resolveDNS(host string) (string, error) {
	c := dns.Client{}
	m := dns.Msg{}

	m.SetQuestion(host+".", dns.TypeA)

	r, _, err := c.Exchange(&m, *flagDNS)
	if err != nil {
		return "", fmt.Errorf("Can't resolve '%s' on %s: %s", host, *flagDNS, err)
	}

	if len(r.Answer) == 0 {
		return "", fmt.Errorf("Can't resolve '%s' on %s: No results", host, *flagDNS)
	}

	aRecord := r.Answer[0].(*dns.A)

	return aRecord.A.String(), nil
}

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Printf("%s %s (Build %s)\n", ProjectName, BuildVersion, BuildDate)
		fmt.Printf("\n")
		fmt.Printf("https://github.com/indece-official/sshmon-check-heketi\n")
		fmt.Printf("\n")
		fmt.Printf("Copyright 2020 by indece UG (haftungsbeschränkt)\n")

		os.Exit(0)

		return
	}

	serviceName := *flagService
	if serviceName == "" {
		serviceName = fmt.Sprintf("Heketi_%s", *flagHost)
	}

	url := fmt.Sprintf("http://%s:%d", *flagHost, *flagPort)

	// Create a client
	heketi := client.NewClient(url, *flagUser, *flagKey)

	// List clusters
	clusterList, err := heketi.ClusterList()
	if err != nil {
		fmt.Printf(
			"2 %s - %s - Can't get cluster list from heketi on %s: %s\n",
			serviceName,
			statusMap[2],
			*flagHost,
			err,
		)
		return
	}

	fmt.Printf(
		"0 %s - %s - Heketi controller on %s is up and running\n",
		serviceName,
		statusMap[0],
		*flagHost,
	)

	for _, clusterID := range clusterList.Clusters {
		clusterInfo, err := heketi.ClusterInfo(clusterID)
		if err != nil {
			fmt.Printf(
				"2 %s - %s - Can't get cluster info for heketi cluster '%s': %s\n",
				fmt.Sprintf("%s_%s", serviceName, clusterID),
				statusMap[2],
				clusterID,
				err,
			)
			continue
		}

		badNodes := []string{}
		failed := false

		for _, nodeID := range clusterInfo.Nodes {
			nodeInfo, err := heketi.NodeInfo(nodeID)
			if err != nil {
				fmt.Printf(
					"2 %s - %s - Can't get node info for node '%s' in heketi cluster '%s': %s\n",
					fmt.Sprintf("%s_%s", serviceName, clusterID),
					statusMap[2],
					nodeID,
					clusterID,
					err,
				)
				break
			}

			if nodeInfo.State != api.EntryStateOnline {
				badNodes = append(badNodes, fmt.Sprintf(
					"Node %s[%s] (%s)",
					nodeID,
					strings.Join(nodeInfo.Hostnames.Storage, ";"),
					nodeInfo.State,
				))
			}
		}

		if failed {
			continue
		}

		if len(badNodes) > 0 {
			fmt.Printf(
				"2 %s - %s - %d of %d nodes of heketi cluster '%s' are unhealthy: %s\n",
				fmt.Sprintf("%s_%s", serviceName, clusterID),
				statusMap[2],
				len(badNodes),
				len(clusterInfo.Nodes),
				clusterID,
				strings.Join(badNodes, ", "),
			)
		} else {
			fmt.Printf(
				"0 %s - %s - All %d nodes of heketi cluster '%s' are healthy\n",
				fmt.Sprintf("%s_%s", serviceName, clusterID),
				statusMap[0],
				len(clusterInfo.Nodes),
				clusterID,
			)
		}
	}
}
