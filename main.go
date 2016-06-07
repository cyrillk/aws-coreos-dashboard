package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gorilla/context"
)

// https://github.com/bndr/gotabulate
// https://github.com/aws/aws-sdk-go

const (
	// Private private
	Private = "private"
	// Public public
	Public = "public"
)

// ApplicationConfig application configuration
type ApplicationConfig struct {
	FleetPort   int
	EtcdPort    int
	IPAddresses string
}

func main() {
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

	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}

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

	appConfig := &ApplicationConfig{
		FleetPort:   49153,
		EtcdPort:    2379,
		IPAddresses: ipAddresses,
	}

	setup(awsConfig, appConfig)
}

func setup(awsConfig *aws.Config, appConfig *ApplicationConfig) {
	http.HandleFunc("/", wrappedHandler(instancesHandler, awsConfig, appConfig))
	// http.HandleFunc("/machines", machinesHandler)
	// http.HandleFunc("/units", servicesHandler)
	// http.HandleFunc("/dockers", dockersHandler)
	http.ListenAndServe(":8080", nil) // TODO configurable port
}

func wrappedHandler(fn http.HandlerFunc, awsConfig *aws.Config, appConfig *ApplicationConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "awsConfig", awsConfig)
		context.Set(r, "appConfig", appConfig)
		fn(w, r)
	}
}

func instancesHandler(w http.ResponseWriter, r *http.Request) {
	awsConfig := context.Get(r, "awsConfig").(*aws.Config)
	appConfig := context.Get(r, "appConfig").(*ApplicationConfig)

	instances := Instances(awsConfig)
	filtered := FilterInstances(instances, appConfig)
	outputString := PrintInstances(filtered)
	outputBytes := []byte(outputString)

	w.Write(outputBytes)
}
