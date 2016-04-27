package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gorilla/context"
)

// https://github.com/bndr/gotabulate
// https://github.com/aws/aws-sdk-go

func main() {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		panic("Missing AWS_REGION environment variable")
	}

	vpcID := os.Getenv("VPC_ID")
	if vpcID == "" {
		panic("Missing VPC_ID environment variable")
	}

	config := &aws.Config{Region: aws.String(region)}

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
