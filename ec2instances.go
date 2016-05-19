package main

import (
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// InstanceInfo ec2 instance info
type InstanceInfo struct {
	InstanceID    string
	Name          string
	PrivateIP     string
	PublicIP      string
	InstanceType  string
	InstanceState string
}

type nameSorter []InstanceInfo

func (a nameSorter) Len() int           { return len(a) }
func (a nameSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a nameSorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Instances retrives a list of EC2 instances
func Instances(awsConfig *aws.Config) []InstanceInfo {
	return sortInstances(getInstances(awsConfig))
}

// SortInstances sorts instances
func sortInstances(instances []InstanceInfo) []InstanceInfo {
	sort.Sort(nameSorter(instances))
	return instances
}

// GetInstances gets ec2 instances info
func getInstances(awsConfig *aws.Config) []InstanceInfo {
	svc := ec2.New(session.New(), awsConfig)

	params := &ec2.DescribeInstancesInput{
	// TODO DryRun: aws.Bool(true),
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
		result.InstanceID = *instance.InstanceId
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
		result.InstanceState = *instance.State.Name
	}

	if instance.InstanceType != nil {
		result.InstanceType = *instance.InstanceType
	}

	return result
}
