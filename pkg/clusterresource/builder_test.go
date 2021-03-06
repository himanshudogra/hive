package clusterresource

import (
	"fmt"
	"github.com/openshift/hive/pkg/apis"
	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"testing"
)

const (
	clusterName                      = "mycluster"
	baseDomain                       = "example.com"
	namespace                        = "mynamespace"
	workerNodeCount                  = 3
	pullSecret                       = "fakepullsecret"
	deleteAfter                      = "8h"
	imageSetName                     = "fake-image-set"
	sshPrivateKey                    = "fakeprivatekey"
	sshPublicKey                     = "fakepublickey"
	fakeManifestFile                 = "my.yaml"
	fakeManifestFileContents         = "fakemanifest"
	fakeAWSAccessKeyID               = "fakeaccesskeyid"
	fakeAWSSecretAccessKey           = "fakesecretAccessKey"
	fakeAzureServicePrincipal        = "fakeSP"
	fakeAzureBaseDomainResourceGroup = "azure-resource-group"
	fakeGCPServiceAccount            = "fakeSA"
	fakeGCPProjectID                 = "gcp-project-id"
	adoptAdminKubeconfig             = "adopted-admin-kubeconfig"
	adoptClusterID                   = "adopted-cluster-id"
	adoptInfraID                     = "adopted-infra-id"
	machineNetwork                   = "10.0.0.0/16"
	fakeOpenStackCloudsYAML          = "fakeYAML"
)

func createTestBuilder() *Builder {
	b := &Builder{
		Name:             clusterName,
		Namespace:        namespace,
		WorkerNodesCount: workerNodeCount,
		PullSecret:       pullSecret,
		SSHPrivateKey:    sshPrivateKey,
		SSHPublicKey:     sshPublicKey,
		BaseDomain:       baseDomain,
		Labels: map[string]string{
			"foo": "bar",
		},
		InstallerManifests: map[string][]byte{
			fakeManifestFile: []byte(fakeManifestFileContents),
		},
		DeleteAfter:    deleteAfter,
		ImageSet:       imageSetName,
		MachineNetwork: machineNetwork,
	}
	return b
}

func createAWSClusterBuilder() *Builder {
	b := createTestBuilder()
	b.CloudBuilder = &AWSCloudBuilder{
		AccessKeyID:     fakeAWSAccessKeyID,
		SecretAccessKey: fakeAWSSecretAccessKey,
	}
	return b
}

func createAzureClusterBuilder() *Builder {
	b := createTestBuilder()
	b.CloudBuilder = &AzureCloudBuilder{
		ServicePrincipal:            []byte(fakeAzureServicePrincipal),
		BaseDomainResourceGroupName: fakeAzureBaseDomainResourceGroup,
	}
	return b
}

func createGCPClusterBuilder() *Builder {
	b := createTestBuilder()
	b.CloudBuilder = &GCPCloudBuilder{
		ServiceAccount: []byte(fakeGCPServiceAccount),
		ProjectID:      fakeGCPProjectID,
	}
	return b
}

func createOpenStackClusterBuilder() *Builder {
	b := createTestBuilder()
	b.CloudBuilder = &OpenStackCloudBuilder{
		CloudsYAMLContent: []byte(fakeOpenStackCloudsYAML),
	}
	return b
}

func TestBuildClusterResources(t *testing.T) {
	tests := []struct {
		name     string
		builder  *Builder
		validate func(t *testing.T, allObjects []runtime.Object)
	}{
		{
			name:    "AWS cluster",
			builder: createAWSClusterBuilder(),
			validate: func(t *testing.T, allObjects []runtime.Object) {
				cd := findClusterDeployment(allObjects, clusterName)
				workerPool := findMachinePool(allObjects, fmt.Sprintf("%s-%s", clusterName, "worker"))

				credsSecretName := fmt.Sprintf("%s-aws-creds", clusterName)
				credsSecret := findSecret(allObjects, credsSecretName)
				require.NotNil(t, credsSecret)
				assert.Equal(t, credsSecret.Name, cd.Spec.Platform.AWS.CredentialsSecretRef.Name)

				assert.Equal(t, awsInstanceType, workerPool.Spec.Platform.AWS.InstanceType)
			},
		},
		{
			name: "adopt AWS cluster",
			builder: func() *Builder {
				awsBuilder := createAWSClusterBuilder()
				awsBuilder.Adopt = true
				awsBuilder.AdoptInfraID = adoptInfraID
				awsBuilder.AdoptClusterID = adoptClusterID
				awsBuilder.AdoptAdminKubeconfig = []byte(adoptAdminKubeconfig)
				return awsBuilder
			}(),
			validate: func(t *testing.T, allObjects []runtime.Object) {
				cd := findClusterDeployment(allObjects, clusterName)

				assert.Equal(t, true, cd.Spec.Installed)
				assert.Equal(t, adoptInfraID, cd.Spec.ClusterMetadata.InfraID)
				assert.Equal(t, adoptClusterID, cd.Spec.ClusterMetadata.ClusterID)

				adminKubeconfig := findSecret(allObjects, fmt.Sprintf("%s-adopted-admin-kubeconfig", clusterName))
				require.NotNil(t, adminKubeconfig)
				assert.Equal(t, adminKubeconfig.Name, cd.Spec.ClusterMetadata.AdminKubeconfigSecretRef.Name)
			},
		},
		{
			name:    "Azure cluster",
			builder: createAzureClusterBuilder(),
			validate: func(t *testing.T, allObjects []runtime.Object) {
				cd := findClusterDeployment(allObjects, clusterName)
				workerPool := findMachinePool(allObjects, fmt.Sprintf("%s-%s", clusterName, "worker"))

				credsSecretName := fmt.Sprintf("%s-azure-creds", clusterName)
				credsSecret := findSecret(allObjects, credsSecretName)
				require.NotNil(t, credsSecret)
				assert.Equal(t, credsSecret.Name, cd.Spec.Platform.Azure.CredentialsSecretRef.Name)

				assert.Equal(t, azureInstanceType, workerPool.Spec.Platform.Azure.InstanceType)
			},
		},
		{
			name:    "GCP cluster",
			builder: createGCPClusterBuilder(),
			validate: func(t *testing.T, allObjects []runtime.Object) {
				cd := findClusterDeployment(allObjects, clusterName)
				workerPool := findMachinePool(allObjects, fmt.Sprintf("%s-%s", clusterName, "worker"))

				credsSecretName := fmt.Sprintf("%s-gcp-creds", clusterName)
				credsSecret := findSecret(allObjects, credsSecretName)
				require.NotNil(t, credsSecret)
				assert.Equal(t, credsSecret.Name, cd.Spec.Platform.GCP.CredentialsSecretRef.Name)

				assert.Equal(t, gcpInstanceType, workerPool.Spec.Platform.GCP.InstanceType)
			},
		},
		{
			name:    "OpenStack cluster",
			builder: createOpenStackClusterBuilder(),
			validate: func(t *testing.T, allObjects []runtime.Object) {
				cd := findClusterDeployment(allObjects, clusterName)

				credsSecretName := fmt.Sprintf("%s-openstack-creds", clusterName)
				credsSecret := findSecret(allObjects, credsSecretName)
				require.NotNil(t, credsSecret)
				assert.Equal(t, credsSecret.Name, cd.Spec.Platform.OpenStack.CredentialsSecretRef.Name)
			},
		},
	}

	for _, test := range tests {
		apis.AddToScheme(scheme.Scheme)
		t.Run(test.name, func(t *testing.T) {
			require.NoError(t, test.builder.Validate())
			allObjects, err := test.builder.Build()
			assert.NoError(t, err)

			cd := findClusterDeployment(allObjects, clusterName)
			require.NotNil(t, cd)

			assert.Equal(t, clusterName, cd.Name)
			assert.Equal(t, "bar", cd.Labels["foo"])
			assert.Equal(t, baseDomain, cd.Spec.BaseDomain)
			assert.Equal(t, deleteAfter, cd.Annotations[deleteAfterAnnotation])
			assert.Equal(t, imageSetName, cd.Spec.Provisioning.ImageSetRef.Name)

			installConfigSecret := findSecret(allObjects, fmt.Sprintf("%s-install-config", clusterName))
			require.NotNil(t, installConfigSecret)
			assert.Equal(t, installConfigSecret.Name, cd.Spec.Provisioning.InstallConfigSecretRef.Name)

			pullSecretSecret := findSecret(allObjects, fmt.Sprintf("%s-pull-secret", clusterName))
			require.NotNil(t, pullSecretSecret)
			assert.Equal(t, pullSecretSecret.Name, cd.Spec.PullSecretRef.Name)

			sshKeySecret := findSecret(allObjects, fmt.Sprintf("%s-ssh-private-key", clusterName))
			require.NotNil(t, sshKeySecret)
			assert.Equal(t, sshKeySecret.Name, cd.Spec.Provisioning.SSHPrivateKeySecretRef.Name)

			workerPool := findMachinePool(allObjects, fmt.Sprintf("%s-%s", clusterName, "worker"))
			require.NotNil(t, workerPool)
			nc := int64(workerNodeCount)
			assert.Equal(t, &nc, workerPool.Spec.Replicas)

			manifestsConfigMap := findConfigMap(allObjects, fmt.Sprintf("%s-%s", clusterName, "manifests"))
			require.NotNil(t, manifestsConfigMap)
			assert.Equal(t, manifestsConfigMap.Name, cd.Spec.Provisioning.ManifestsConfigMapRef.Name)

			test.validate(t, allObjects)
		})
	}

}

func findSecret(allObjects []runtime.Object, name string) *corev1.Secret {
	for _, ro := range allObjects {
		obj, ok := ro.(*corev1.Secret)
		if !ok {
			continue
		}
		if obj.Name == name {
			return obj
		}
	}
	return nil
}

func findClusterDeployment(allObjects []runtime.Object, name string) *hivev1.ClusterDeployment {
	for _, ro := range allObjects {
		obj, ok := ro.(*hivev1.ClusterDeployment)
		if !ok {
			continue
		}
		if obj.Name == name {
			return obj
		}
	}
	return nil
}

func findMachinePool(allObjects []runtime.Object, name string) *hivev1.MachinePool {
	for _, ro := range allObjects {
		obj, ok := ro.(*hivev1.MachinePool)
		if !ok {
			continue
		}
		if obj.Name == name {
			return obj
		}
	}
	return nil
}

func findConfigMap(allObjects []runtime.Object, name string) *corev1.ConfigMap {
	for _, ro := range allObjects {
		obj, ok := ro.(*corev1.ConfigMap)
		if !ok {
			continue
		}
		if obj.Name == name {
			return obj
		}
	}
	return nil
}
