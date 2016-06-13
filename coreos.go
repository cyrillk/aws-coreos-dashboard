package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

type wrappedInstanceInfo struct {
	instanceInfo InstanceInfo
	passed       bool
}

// type localInstanceInfo interface {
// }
// type emptyInstanceInfo struct {
// 	localInstanceInfo
// }
// type clusteredInstanceInfo struct {
// 	// localInstanceInfo
// 	InstanceInfo  InstanceInfo
// 	FleetMachines []fleetMachine
// }

type fleetMachines struct {
	Machines []fleetMachine `json:"machines"`
}

type fleetMachine struct {
	ID        string `json:"id"`
	PrimaryIP string `json:"primaryIP"`
}

// http://52.48.119.163:49153/fleet/v1/units
// http://52.48.119.163:49153/fleet/v1/machines
// http://52.48.119.163:49153/fleet/v1/states

const requestTimeout = 300

func parseResponse(body []byte) fleetMachines {
	var m fleetMachines
	err := json.Unmarshal(body, &m)

	if err != nil {
		panic(err)
	}

	return m
}

// GroupInstances groups instances by cluster
func GroupInstances(instances []InstanceInfo, appConfig *ApplicationConfig) [][]InstanceInfo {

	jobs := make(chan InstanceInfo, len(instances))
	results := make(chan []fleetMachine, len(instances))

	// creates workers
	for _ = range instances {
		go groupWorker(jobs, results, appConfig)
	}

	// run jobs
	for _, instance := range instances {
		jobs <- instance
	}

	close(jobs)

	var resultMachines [][]fleetMachine

	for _ = range instances {
		result := <-results
		resultMachines = append(resultMachines, result)
	}

	groupedMachines := make(map[string]fleetMachines)

	for _, machinesSlice := range resultMachines {
		var key string

		for _, m := range sortMachines(machinesSlice) {
			if key == "" {
				key = m.ID
			} else {
				key = fmt.Sprintf("%s|%s", key, m.ID)
			}
		}

		groupedMachines[key] = fleetMachines{
			Machines: machinesSlice,
		}
	}

	var groupedInstances [][]InstanceInfo

	for k := range groupedMachines {
		var group []InstanceInfo

		for _, m := range groupedMachines[k].Machines {
			for _, f := range instances {
				if m.PrimaryIP == f.PrivateIP {
					group = append(group, f)
				}
			}
		}

		groupedInstances = append(groupedInstances, group)
	}

	return groupedInstances
}

type machineSorter []fleetMachine

func (a machineSorter) Len() int           { return len(a) }
func (a machineSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a machineSorter) Less(i, j int) bool { return a[i].ID < a[j].ID }

func sortMachines(m []fleetMachine) []fleetMachine {
	sort.Sort(machineSorter(m))
	return m
}

func contains(s []fleetMachine, machine fleetMachine) bool {
	for _, m := range s {
		if m.PrimaryIP == machine.PrimaryIP {
			return true
		}
	}
	return false
}

func groupWorker(jobs <-chan InstanceInfo, results chan<- []fleetMachine, appConfig *ApplicationConfig) {
	for instance := range jobs {

		timeout := time.Duration(requestTimeout * time.Millisecond)
		client := http.Client{
			Timeout: timeout,
		}

		ip := pickIP(instance, appConfig.IPAddresses)
		address := fmt.Sprintf("%s:%d", ip, appConfig.FleetPort)

		resp, err := client.Get(fmt.Sprintf("http://%s/fleet/v1/machines", address))
		defer resp.Body.Close()

		if err != nil {
			// panic(err)
		}

		if resp.StatusCode != 200 {
			// panic(fmt.Sprintf("Response code is %d", resp.StatusCode))
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			// panic(err)
		}

		fleetMachines := parseResponse(body)
		results <- fleetMachines.Machines
	}
}

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

func pickIP(instance InstanceInfo, ipAddresses string) string {
	if ipAddresses == Private {
		return instance.PrivateIP
	} else if ipAddresses == Public {
		return instance.PublicIP
	}
	return ""
}
