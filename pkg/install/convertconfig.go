/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package install

import (
	"fmt"
	"net"

	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
)

// generateInstallConfig builds an InstallConfig for the installer from our ClusterDeploymentSpec.
// The two types are extremely similar, but have different goals and in some cases deviation was required
// as ClusterDeployment is used as a CRD API.
//
// It is assumed the caller will lookup the admin password and ssh key from their respective secrets.
func generateInstallConfig(cd *hivev1.ClusterDeployment, adminPassword, sshKey, pullSecret string) (*types.InstallConfig, error) {
	/*
		networkType, err := convertNetworkingType(spec.Config.Networking.Type)
		if err != nil {
			return nil, err
		}
	*/

	spec := cd.Spec
	platform := types.Platform{}
	if spec.Config.Platform.AWS != nil {
		aws := spec.Config.Platform.AWS
		platform.AWS = &types.AWSPlatform{
			Region:       aws.Region,
			UserTags:     aws.UserTags,
			VPCID:        aws.VPCID,
			VPCCIDRBlock: aws.VPCCIDRBlock,
		}
		if aws.DefaultMachinePlatform != nil {
			platform.AWS.DefaultMachinePlatform = &types.AWSMachinePoolPlatform{
				InstanceType: aws.DefaultMachinePlatform.InstanceType,
				IAMRoleName:  aws.DefaultMachinePlatform.IAMRoleName,
				EC2RootVolume: types.EC2RootVolume{
					IOPS: aws.DefaultMachinePlatform.EC2RootVolume.IOPS,
					Size: aws.DefaultMachinePlatform.EC2RootVolume.Size,
					Type: aws.DefaultMachinePlatform.EC2RootVolume.Type,
				},
			}
		}
	}

	machinePools := []types.MachinePool{}
	for _, mp := range spec.Config.Machines {
		newMP := types.MachinePool{
			Name:     mp.Name,
			Replicas: mp.Replicas,
		}
		if mp.Platform.AWS != nil {
			newMP.Platform.AWS = &types.AWSMachinePoolPlatform{
				InstanceType: mp.Platform.AWS.InstanceType,
				IAMRoleName:  mp.Platform.AWS.IAMRoleName,
				EC2RootVolume: types.EC2RootVolume{
					IOPS: mp.Platform.AWS.EC2RootVolume.IOPS,
					Size: mp.Platform.AWS.EC2RootVolume.Size,
					Type: mp.Platform.AWS.EC2RootVolume.Type,
				},
			}
		}
		machinePools = append(machinePools, newMP)
	}

	ic := &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: spec.Config.ClusterID,
		},
		ClusterID: cd.Status.ClusterUUID,
		Admin: types.Admin{
			Email:    spec.Config.Admin.Email,
			Password: adminPassword,
			SSHKey:   sshKey,
		},
		BaseDomain: spec.Config.BaseDomain,
		Networking: types.Networking{
			// TODO: installer currently only works with flannel, flagged for TODO on their end here: https://github.com/openshift/installer/blob/master/pkg/asset/installconfig/installconfig.go#L82
			Type: "flannel",
			ServiceCIDR: ipnet.IPNet{
				IPNet: parseCIDR(spec.Config.Networking.ServiceCIDR),
			},
			PodCIDR: ipnet.IPNet{
				IPNet: parseCIDR(spec.Config.Networking.PodCIDR),
			},
		},
		PullSecret: pullSecret,
		Platform:   platform,
		Machines:   machinePools,
	}
	return ic, nil
}

func parseCIDR(s string) net.IPNet {
	if s == "" {
		return net.IPNet{}
	}
	_, cidr, _ := net.ParseCIDR(s)
	return *cidr
}

func convertNetworkingType(hnt hivev1.NetworkType) (types.NetworkType, error) {
	switch hnt {
	case hivev1.NetworkTypeOpenshiftOVN:
		return types.NetworkTypeOpenshiftOVN, nil
	case hivev1.NetworkTypeOpenshiftSDN:
		return types.NetworkTypeOpenshiftSDN, nil
	default:
		return "", fmt.Errorf("unknown NetworkType: %s", hnt)
	}
}
