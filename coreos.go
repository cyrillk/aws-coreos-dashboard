package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func pickIP(instance InstanceInfo, appConfig *ApplicationConfig) string {
	if appConfig.IPAddresses == Private {
		return instance.PrivateIP
	} else if appConfig.IPAddresses == Public {
		return instance.PublicIP
	}
	return ""
}

// FilterInstances should filter instances to use CoreOS
func FilterInstances(instances []InstanceInfo, appConfig *ApplicationConfig) []InstanceInfo {

	var fleetInstances []InstanceInfo

	for _, instance := range instances {

		ip := pickIP(instance, appConfig)

		if ip != "" {
			address := fmt.Sprintf("%s:%d", ip, appConfig.FleetPort)

			log.Println(address)

			timeout := time.Duration(1 * time.Second)
			client := http.Client{
				Timeout: timeout,
			}

			_, err := client.Get(fmt.Sprintf("http://%s/", address))
			// conn, err := net.Dial("tcp", address)

			if err == nil {
				fleetInstances = append(fleetInstances, instance)
			}
			// if err != nil {
			// 	log.Println("Connection error:", err)
			// } else {
			// 	// defer conn.Close()
			// 	log.Println("OK")
			// }
		}
	}

	return fleetInstances
}
