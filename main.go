package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gin-gonic/gin"
)

// https://github.com/bndr/gotabulate
// https://github.com/aws/aws-sdk-go
// https://github.com/gin-gonic/gin

const (
	// Private private IP address
	Private = "private"
	// Public public IP address
	Public = "public"
)

// ApplicationConfig application configuration
type ApplicationConfig struct {
	FleetPort   int
	EtcdPort    int
	IPAddresses string
}

func main() {
	awsConfig := awsConfig()
	appConfig := appConfig()

	router := gin.Default()
	router.Static("/resources/assets", "./resources/assets")
	router.LoadHTMLGlob("resources/templates/*")

	// http.HandleFunc("/dockers", dockersHandler)

	router.GET("/", func(c *gin.Context) {
		handleMachines(c, awsConfig, appConfig)
	})

	router.GET("/services", func(c *gin.Context) {
		handleServices(c, awsConfig, appConfig)
	})

	router.GET("/machines", func(c *gin.Context) {
		handleMachines(c, awsConfig, appConfig)
	})

	router.GET("/instances", func(c *gin.Context) {
		handleInstances(c, awsConfig, appConfig)
	})

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
}

func handleServices(c *gin.Context, awsConfig *aws.Config, appConfig *ApplicationConfig) {
	instances := Instances(awsConfig)
	filtered := FilterInstances(instances, appConfig)
	groupedInstances := GroupInstances(filtered, appConfig)

	var servicesModule = new(CoreOSServicesModule)
	groupedServices := servicesModule.RetrieveServices(groupedInstances, appConfig)
	sortedServices := servicesModule.Sort(groupedServices)

	c.HTML(http.StatusOK, "services.html", gin.H{
		"groups": sortedServices,
	})
}

func handleMachines(c *gin.Context, awsConfig *aws.Config, appConfig *ApplicationConfig) {
	instances := Instances(awsConfig)
	filtered := FilterInstances(instances, appConfig)
	grouped := GroupInstances(filtered, appConfig)

	var sorted [][]CoreOSMachine

	for _, group := range SortGroups(grouped) {
		sorted = append(sorted, SortMachines(group))
	}

	c.HTML(http.StatusOK, "machines.html", gin.H{
		"groups": sorted,
	})
}

func handleInstances(c *gin.Context, awsConfig *aws.Config, appConfig *ApplicationConfig) {
	instances := Instances(awsConfig)
	filtered := FilterInstances(instances, appConfig)
	sorted := SortInstances(filtered)
	body := PrintInstances(sorted)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"body": body,
	})
}

func awsConfig() *aws.Config {
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if awsAccessKeyID == "" {
		panic("Missing AWS_ACCESS_KEY_ID environment variable")
	}

	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if awsSecretAccessKey == "" {
		panic("Missing AWS_SECRET_ACCESS_KEY environment variable")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		panic("Missing AWS_REGION environment variable")
	}

	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")

	return &aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}
}

func appConfig() *ApplicationConfig {
	ipAddresses := os.Getenv("IP_ADDRESSES")
	if ipAddresses != Private && ipAddresses != Public {
		panic("Invalid IP_ADDRESSES environment variable")
	}

	// fleetPort, err := strconv.Atoi(os.Getenv("FLEET_PORT"))
	// if err != nil {
	// 	panic("Invalid FLEET_PORT environment variable")
	// }
	//
	// etcdPort, err := strconv.Atoi(os.Getenv("ETCD_PORT"))
	// if err != nil {
	// 	panic("Invalid ETCD_PORT environment variable")
	// }

	return &ApplicationConfig{
		FleetPort:   49153,
		EtcdPort:    2379,
		IPAddresses: ipAddresses,
	}
}
