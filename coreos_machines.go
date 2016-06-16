package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

type fleetMachines struct {
	Machines []fleetMachine `json:"machines"`
}

type fleetMachine struct {
	ID        string `json:"id"`
	PrimaryIP string `json:"primaryIP"`
}

// CoreOSMachine CoreOS machine
type CoreOSMachine struct {
	InstanceID    string
	Name          string
	PrivateIP     string
	PublicIP      string
	InstanceType  string
	InstanceState string
	MachineID     string
}

// http://52.48.119.163:49153/fleet/v1/units
// http://52.48.119.163:49153/fleet/v1/machines
// http://52.48.119.163:49153/fleet/v1/states

const machinesRequestTimeout = 300

func parseResponse(body []byte) fleetMachines {
	var m fleetMachines
	err := json.Unmarshal(body, &m)

	if err != nil {
		panic(err)
	}

	return m
}

// GroupInstances groups instances by cluster
func GroupInstances(instances []InstanceInfo, appConfig *ApplicationConfig) [][]CoreOSMachine {

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

		for _, m := range sortFleetMachines(machinesSlice) {
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

	var groupedInstances [][]CoreOSMachine

	for k := range groupedMachines {
		var group []CoreOSMachine

		for _, m := range groupedMachines[k].Machines {
			for _, f := range instances {
				if m.PrimaryIP == f.PrivateIP {
					group = append(group, CoreOSMachine{
						InstanceID:    f.InstanceID,
						Name:          f.Name,
						PrivateIP:     f.PrivateIP,
						PublicIP:      f.PublicIP,
						InstanceType:  f.InstanceType,
						InstanceState: f.InstanceState,
						MachineID:     m.ID[:8],
					})
				}
			}
		}

		groupedInstances = append(groupedInstances, group)
	}

	return groupedInstances
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

		timeout := time.Duration(machinesRequestTimeout * time.Millisecond)
		client := http.Client{
			Timeout: timeout,
		}

		ip := pickIP(instance.PrivateIP, instance.PublicIP, appConfig.IPAddresses)
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

type fleetMachineSorter []fleetMachine

func (a fleetMachineSorter) Len() int           { return len(a) }
func (a fleetMachineSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a fleetMachineSorter) Less(i, j int) bool { return a[i].ID < a[j].ID }

func sortFleetMachines(m []fleetMachine) []fleetMachine {
	sort.Sort(fleetMachineSorter(m))
	return m
}

type groupSorter [][]CoreOSMachine

func (a groupSorter) Len() int           { return len(a) }
func (a groupSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a groupSorter) Less(i, j int) bool { return a[i][0].PrivateIP < a[j][0].PrivateIP }

// SortGroups sorts groups of instances
func SortGroups(m [][]CoreOSMachine) [][]CoreOSMachine {
	sort.Sort(groupSorter(m))
	return m
}

type machineSorter []CoreOSMachine

func (a machineSorter) Len() int           { return len(a) }
func (a machineSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a machineSorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

// SortMachines sorts machines
func SortMachines(m []CoreOSMachine) []CoreOSMachine {
	sort.Sort(machineSorter(m))
	return m
}
