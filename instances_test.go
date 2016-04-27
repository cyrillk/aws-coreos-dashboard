package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestParseInstance(t *testing.T) {

	privateIP := "1.1.1.1"
	publicIP := "2.2.2.2"
	instanceID := "XX"
	stateName := "running"

	instance := &ec2.Instance{
		PrivateIpAddress: &privateIP,
		PublicIpAddress:  &publicIP,
		InstanceId:       &instanceID,
		State:            &ec2.InstanceState{Name: &stateName},
	}

	fmt.Println(instance)

	info := parseInstance(instance)

	if info.PrivateIP != privateIP {
		t.Error("Expected "+privateIP+", got ", info.PrivateIP)
	}

	if info.PublicIP != publicIP {
		t.Error("Expected "+publicIP+", got ", info.PublicIP)
	}

	if info.instanceID != instanceID {
		t.Error("Expected "+instanceID+", got ", info.instanceID)
	}

	if info.StateName != stateName {
		t.Error("Expected "+stateName+", got ", info.StateName)
	}

	// info
	// var v float64
	// v = Average([]float64{1, 2})
	// if v != 1.5 {
	// 	t.Error("Expected 1.5, got ", v)
	// }
}
