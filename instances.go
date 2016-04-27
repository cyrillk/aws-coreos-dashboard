package main

import (
	"fmt"
	"log"
	"os/exec"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// InstanceInfo ec2 instance info
type InstanceInfo struct {
	instanceID   string
	InstanceType string
	PrivateIP    string
	PublicIP     string
	StateName    string
	Name         string
}

type nameSorter []InstanceInfo

func (a nameSorter) Len() int           { return len(a) }
func (a nameSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a nameSorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

// BuildTableOfInstances builds table of instances
func BuildTableOfInstances(vpcID string, config *aws.Config) string {

	instances := sortInstances(getInstances(vpcID, config))

	var data [][]string

	for _, instance := range instances {

		if instance.PublicIP != "" {
			tunnel := fmt.Sprintf("FLEETCTL_TUNNEL=%s", instance.PublicIP)

			fmt.Println(">>> " + tunnel)

			args := []string{tunnel, "fleetctl", "list-machines", "-no-legend", "-fields", "ip"}
			out, err := exec.Command("env", args...).Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("The %s\n", out)
		}

		s := []string{
			instance.Name,
			instance.PrivateIP,
			instance.PublicIP,
			instance.InstanceType,
			instance.StateName,
		}
		data = append(data, s)
	}

	headers := []string{"Name", "Private IP", "Public IP", "Type", "State"}

	return BuildTable(headers, data)
}

// SortInstances sorts instances
func sortInstances(instances []InstanceInfo) []InstanceInfo {
	sort.Sort(nameSorter(instances))
	return instances
}

// GetInstances gets ec2 instances info
func getInstances(vpcID string, config *aws.Config) []InstanceInfo {
	svc := ec2.New(session.New(), config)

	params := &ec2.DescribeInstancesInput{
		// TODO DryRun: aws.Bool(true),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(vpcID),
				},
			},
		},
	}

	resp, err := svc.DescribeInstances(params)
	if err != nil {
		panic(err)
	}

	var result []InstanceInfo

	for idx := range resp.Reservations {
		for _, instance := range resp.Reservations[idx].Instances {
			result = append(result, *parseInstance(instance))
		}
	}

	return result
}

func parseInstance(instance *ec2.Instance) *InstanceInfo {
	var name *string

	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			name = tag.Value
		}
	}

	result := new(InstanceInfo)

	if instance.InstanceId != nil {
		result.instanceID = *instance.InstanceId
	}

	if name != nil {
		result.Name = *name
	}

	if instance.PrivateIpAddress != nil {
		result.PrivateIP = *instance.PrivateIpAddress
	}

	if instance.PublicIpAddress != nil {
		result.PublicIP = *instance.PublicIpAddress
	}

	if instance.State.Name != nil {
		result.StateName = *instance.State.Name
	}

	if instance.InstanceType != nil {
		result.InstanceType = *instance.InstanceType
	}

	return result
}
