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
	instanceState := "running"

	instance := &ec2.Instance{
		PrivateIpAddress: &privateIP,
		PublicIpAddress:  &publicIP,
		InstanceId:       &instanceID,
		State:            &ec2.InstanceState{Name: &instanceState},
	}

	fmt.Println(instance)

	info := parseInstance(instance)

	if info.PrivateIP != privateIP {
		t.Error("Expected "+privateIP+", got ", info.PrivateIP)
	}

	if info.PublicIP != publicIP {
		t.Error("Expected "+publicIP+", got ", info.PublicIP)
	}

	if info.InstanceID != instanceID {
		t.Error("Expected "+instanceID+", got ", info.InstanceID)
	}

	if info.InstanceState != instanceState {
		t.Error("Expected "+instanceState+", got ", info.InstanceState)
	}
}
