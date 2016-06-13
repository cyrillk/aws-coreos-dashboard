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
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")

	// http.HandleFunc("/units", servicesHandler)
	// http.HandleFunc("/dockers", dockersHandler)

	router.GET("/", func(c *gin.Context) {
		handleMachines(c, awsConfig, appConfig)
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

func handleMachines(c *gin.Context, awsConfig *aws.Config, appConfig *ApplicationConfig) {
	instances := Instances(awsConfig)
	filtered := FilterInstances(instances, appConfig)
	grouped := GroupInstances(filtered, appConfig)

	var body string

	for _, group := range grouped {
		body = body + "\n" + PrintInstances(group)
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"body": body,
	})
}

func handleInstances(c *gin.Context, awsConfig *aws.Config, appConfig *ApplicationConfig) {
	instances := Instances(awsConfig)
	filtered := FilterInstances(instances, appConfig)
	sorted := SortInstances(filtered)
	body := PrintInstances(sorted)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
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
