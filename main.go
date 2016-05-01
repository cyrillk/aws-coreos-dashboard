package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gorilla/context"
)

// https://github.com/bndr/gotabulate
// https://github.com/aws/aws-sdk-go

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

	vpcID := os.Getenv("VPC_ID")
	if vpcID == "" {
		panic("Missing VPC_ID environment variable")
	}

	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")

	config := &aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}

	setup(config, vpcID)
}

func setup(config *aws.Config, vpcID string) {
	http.HandleFunc("/", withParameters(handler, vpcID, config))
	// http.HandleFunc("/services", servicesHandler)
	// http.HandleFunc("/dockers", dockersHandler)
	http.ListenAndServe(":8080", nil)
}

func withParameters(fn http.HandlerFunc, vpcID string, config *aws.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "vpcID", vpcID)
		context.Set(r, "config", config)
		fn(w, r)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	vpcID := context.Get(r, "vpcID").(string)
	config := context.Get(r, "config").(*aws.Config)

	fmt.Fprintf(w, BuildTableOfInstances(vpcID, config))
}
