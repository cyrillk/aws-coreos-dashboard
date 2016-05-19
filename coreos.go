package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Filter should filter instances to use CoreOS
func Filter(instances []InstanceInfo, appConfig *ApplicationConfig) []InstanceInfo {

	for _, instance := range instances {

		if instance.PublicIP != "" {
			address := fmt.Sprintf("%s:%d", instance.PublicIP, appConfig.FleetPort)

			log.Println(address)

			timeout := time.Duration(1 * time.Second)
			client := http.Client{
				Timeout: timeout,
			}

			_, err := client.Get(fmt.Sprintf("http://%s/", address))
			// conn, err := net.Dial("tcp", address)

			if err != nil {
				log.Println("Connection error:", err)
			} else {
				// defer conn.Close()
				log.Println("OK")
			}
		}
	}

	var data [][]string

	for _, instance := range instances {

		s := []string{
			instance.Name,
			instance.PrivateIP,
			instance.PublicIP,
			instance.InstanceType,
			instance.InstanceState,
		}
		data = append(data, s)
	}

	headers := []string{"Name", "Private IP", "Public IP", "Type", "State"}

	BuildTable(headers, data)

	return instances
}
