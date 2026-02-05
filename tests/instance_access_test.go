package tests

import (
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestSSMReachability(t *testing.T) {
	if os.Getenv("TEST_AWS") == "" {
		t.Skip("Skipping: set TEST_AWS=1 to run AWS SSM tests")
	}

	t.Parallel()

	uniqueID := random.UniqueId()
	exampleDir := "../examples/ssm-instance"

	tfOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: exampleDir,
		Vars: map[string]interface{}{
			"run_id": uniqueID,
			"region": "us-east-1",
		},
	})

	// Teardown guarantee
	defer terraform.Destroy(t, tfOptions)

	terraform.InitAndApply(t, tfOptions)

	instanceID := terraform.Output(t, tfOptions, "instance_id")
	assert.NotEmpty(t, instanceID)

	// Verify SSM connectivity using AWS SDK
	cfg, err := config.LoadDefaultConfig(t.Context(), config.WithRegion("us-east-1"))
	if err != nil {
		t.Fatalf("unable to load SDK config, %v", err)
	}

	ssmClient := ssm.NewFromConfig(cfg)

	// Wait for instance to come online in SSM
	// This can take a minute or two after boot
	maxRetries := 30
	retryInterval := 10 * time.Second

	t.Logf("Waiting for instance %s to register with SSM...", instanceID)

	var isOnline bool
	for i := 0; i < maxRetries; i++ {
		output, err := ssmClient.DescribeInstanceInformation(t.Context(), &ssm.DescribeInstanceInformationInput{
			Filters: []types.InstanceInformationStringFilter{
				{
					Key:    aws.String("InstanceIds"),
					Values: []string{instanceID},
				},
			},
		})

		if err == nil && len(output.InstanceInformationList) > 0 {
			status := output.InstanceInformationList[0].PingStatus
			t.Logf("Instance status: %s", status)
			if status == types.PingStatusOnline {
				isOnline = true
				break
			}
		}

		time.Sleep(retryInterval)
	}

	assert.True(t, isOnline, "Instance failed to register with SSM and become Online")
}
