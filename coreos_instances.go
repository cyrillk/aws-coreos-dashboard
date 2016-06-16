package main

import (
	"fmt"
	"net/http"
	"time"
)

// http://52.48.119.163:49153/fleet/v1/units
// http://52.48.119.163:49153/fleet/v1/machines
// http://52.48.119.163:49153/fleet/v1/states

type wrappedInstanceInfo struct {
	instanceInfo InstanceInfo
	passed       bool
}

const instancesRequestTimeout = 300

// FilterInstances should filter instances to use CoreOS
func FilterInstances(instances []InstanceInfo, appConfig *ApplicationConfig) []InstanceInfo {

	jobs := make(chan InstanceInfo, len(instances))
	results := make(chan wrappedInstanceInfo, len(instances))

	// creates workers
	for _ = range instances {
		go filterWorker(jobs, results, appConfig)
	}

	// run jobs
	for _, instance := range instances {
		jobs <- instance
	}

	close(jobs)

	var fleetInstances []InstanceInfo

	for _ = range instances {
		result := <-results
		if result.passed {
			fleetInstances = append(fleetInstances, result.instanceInfo)
		}
	}

	return fleetInstances
}

func filterWorker(jobs <-chan InstanceInfo, results chan<- wrappedInstanceInfo, appConfig *ApplicationConfig) {
	for instance := range jobs {

		timeout := time.Duration(instancesRequestTimeout * time.Millisecond)
		client := http.Client{
			Timeout: timeout,
		}

		ip := pickIP(instance.PrivateIP, instance.PublicIP, appConfig.IPAddresses)
		address := fmt.Sprintf("%s:%d", ip, appConfig.FleetPort)

		_, err := client.Get(fmt.Sprintf("http://%s/", address))
		// conn, err := net.Dial("tcp", address)

		if err == nil {
			results <- wrappedInstanceInfo{
				instanceInfo: instance,
				passed:       true,
			}
		} else {
			results <- wrappedInstanceInfo{
				instanceInfo: instance,
				passed:       false,
			}
		}
	}
}

func pickIP(privateIP string, publicIP string, ipAddressType string) string {
	if ipAddressType == Private {
		return privateIP
	} else if ipAddressType == Public {
		return publicIP
	}
	return ""
}
