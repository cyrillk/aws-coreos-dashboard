package main

import (
	"fmt"
	"net/http"
	"time"
)

type wrappedInstanceInfo struct {
	instanceInfo InstanceInfo
	passed       bool
}

const requestTimeout = 300

func worker(jobs <-chan InstanceInfo, results chan<- wrappedInstanceInfo, appConfig *ApplicationConfig) {
	for instance := range jobs {

		timeout := time.Duration(requestTimeout * time.Millisecond)
		client := http.Client{
			Timeout: timeout,
		}

		ip := pickIP(instance, appConfig.IPAddresses)
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

// FilterInstances should filter instances to use CoreOS
func FilterInstances(instances []InstanceInfo, appConfig *ApplicationConfig) []InstanceInfo {

	jobs := make(chan InstanceInfo, len(instances))
	results := make(chan wrappedInstanceInfo, len(instances))

	// creates workers
	for _ = range instances {
		go worker(jobs, results, appConfig)
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

func pickIP(instance InstanceInfo, ipAddresses string) string {
	if ipAddresses == Private {
		return instance.PrivateIP
	} else if ipAddresses == Public {
		return instance.PublicIP
	}
	return ""
}
