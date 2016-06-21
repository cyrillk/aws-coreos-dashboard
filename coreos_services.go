package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"time"
)

type fleetStates struct {
	States []fleetState `json:"states"`
}

type fleetState struct {
	MachineID   string `json:"machineID"`
	Name        string `json:"name"`
	ActiveState string `json:"systemdActiveState"`
	LoadState   string `json:"systemdLoadState"`
	SubState    string `json:"systemdSubState"`
}

// CoreOSService CoreOS unit service
type CoreOSService struct {
	MachineID   string
	Name        string
	ActiveState string
	LoadState   string
	SubState    string
	PrivateIP   string
}

// CoreOSCluster CoreOS cluster
type CoreOSCluster struct {
	Services []CoreOSService
	Machines []CoreOSMachine
}

// http://52.48.119.163:49153/fleet/v1/units
// http://52.48.119.163:49153/fleet/v1/machines
// http://52.48.119.163:49153/fleet/v1/state

// CoreOSServicesModule CoreOS unit services
type CoreOSServicesModule struct {
}

const statesRequestTimeout = 300

type servicesSorter []CoreOSService

func (a servicesSorter) Len() int           { return len(a) }
func (a servicesSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a servicesSorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

// SortServices sorts services
func SortServices(m []CoreOSService) []CoreOSService {
	sort.Sort(servicesSorter(m))
	return m
}

type clusterSorter []CoreOSCluster

func (a clusterSorter) Len() int      { return len(a) }
func (a clusterSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a clusterSorter) Less(i, j int) bool {
	return a[i].Machines[0].PrivateIP < a[j].Machines[0].PrivateIP
}

// SortClusters sorts clusteres
func SortClusters(m []CoreOSCluster) []CoreOSCluster {
	sort.Sort(clusterSorter(m))
	return m
}

// Sort sorts
func (coreOSServicesModule CoreOSServicesModule) Sort(clusteredServices []CoreOSCluster) []CoreOSCluster {
	var sorted []CoreOSCluster

	for _, group := range SortClusters(clusteredServices) {
		sortedMachines := SortMachines(group.Machines)
		sortedServices := SortServices(group.Services)

		sorted = append(sorted, CoreOSCluster{
			Machines: sortedMachines,
			Services: sortedServices,
		})
	}

	return sorted
}

func parseFleetStatesJSON(body []byte) fleetStates {
	var m fleetStates
	err := json.Unmarshal(body, &m)

	if err != nil {
		panic(err)
	}

	return m
}

// RetrieveServices retrieves unit services
func (coreOSServicesModule CoreOSServicesModule) RetrieveServices(
	groupedMachines [][]CoreOSMachine, appConfig *ApplicationConfig) []CoreOSCluster {

	jobs := make(chan []CoreOSMachine, len(groupedMachines))
	results := make(chan CoreOSCluster, len(groupedMachines))

	// creates workers
	for _ = range groupedMachines {
		go servicesWorker(jobs, results, appConfig)
	}

	// run jobs
	for _, group := range groupedMachines {
		jobs <- group
	}

	close(jobs)

	var resultServices []CoreOSCluster

	for _ = range groupedMachines {
		result := <-results
		resultServices = append(resultServices, result)
	}

	return resultServices
}

func servicesWorker(jobs <-chan []CoreOSMachine, results chan<- CoreOSCluster, appConfig *ApplicationConfig) {
	for machines := range jobs {

		machine := machines[rand.Intn(len(machines))] // TODO make it different possibly

		client := http.Client{
			Timeout: time.Duration(statesRequestTimeout * time.Millisecond),
		}

		ip := pickIP(machine.PrivateIP, machine.PublicIP, appConfig.IPAddresses)
		address := fmt.Sprintf("%s:%d", ip, appConfig.FleetPort)

		resp, err := client.Get(fmt.Sprintf("http://%s/fleet/v1/state", address))
		defer resp.Body.Close()

		if err != nil {
			// TODO panic(err)
		}

		if resp.StatusCode != 200 {
			// TODO panic(fmt.Sprintf("Response code is %d", resp.StatusCode))
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			// TODO panic(err)
		}

		fleetStates := parseFleetStatesJSON(body)

		var services []CoreOSService

		for _, state := range fleetStates.States {
			services = append(services, CoreOSService{
				MachineID:   state.MachineID[:8],
				Name:        state.Name,
				ActiveState: state.ActiveState,
				LoadState:   state.LoadState,
				SubState:    state.SubState,
			})
		}

		clustered := CoreOSCluster{
			Services: services,
			Machines: machines,
		}

		results <- clustered
	}
}
