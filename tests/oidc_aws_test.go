package tests

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestOIDCModuleTrustPolicy(t *testing.T) {
	// Soft guard so devs without AWS/OIDC don't fail the suite accidentally.
	if os.Getenv("TEST_AWS") == "" {
		t.Skip("Skipping: set TEST_AWS=1 to run AWS OIDC tests")
	}

	// Path to your example
	exampleDir := "../examples/oidc-aws-setup"

	// Standard Terratest options
	tf := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: exampleDir,
		// Supply vars here if you don't use tfvars files
		// Vars: map[string]interface{}{
		// 	"region":       "us-east-1",
		// 	"github_owner": "pwjohnston2002",
		// 	"github_repo":  "secrets-workflow",
		// 	"branch_name":  "main",
		// },
	})

	var roleName string // Declare roleName here so it's accessible in the deferred function

	// Always teardown, and verify teardown succeeded
	defer func() {
		terraform.Destroy(t, tf)

		// Post-destroy verification: IAM role should be gone
		cfg, err := config.LoadDefaultConfig(t.Context())
		if err != nil {
			t.Logf("Skipping teardown verification: loading AWS config failed: %v", err)
			return
		}
		iamc := iam.NewFromConfig(cfg)

		// Expect an error of type NotFoundException
		_, err = iamc.GetRole(t.Context(), &iam.GetRoleInput{RoleName: aws.String(roleName)})
		if err == nil {
			t.Errorf("expected IAM role %s to be deleted, but GetRole succeeded", roleName)
		} else {
			var nf *types.NoSuchEntityException
			assert.ErrorAs(t, err, &nf, "expected NoSuchEntityException when getting deleted role")
			t.Logf("IAM role %s successfully deleted (verified by NoSuchEntityException)", roleName)
		}
	}()

	terraform.InitAndApply(t, tf)

	// Grab output (role ARN) so we can query IAM directly
	roleArn := terraform.Output(t, tf, "role_arn")
	assert.NotEmpty(t, roleArn, "role_arn output should not be empty")

	// Derive role name from ARN (arn:aws:iam::<acct>:role/<name>)
	parts := strings.Split(roleArn, "/")
	if len(parts) < 2 {
		t.Fatalf("unexpected role ARN format: %s", roleArn)
	}
	roleName = parts[len(parts)-1]

	// AWS SDK client
	cfg, err := config.LoadDefaultConfig(t.Context())
	if err != nil {
		t.Fatalf("loading AWS config: %v", err)
	}
	iamc := iam.NewFromConfig(cfg)

	// Fetch the role and its trust policy (AssumeRolePolicyDocument)
	getOut, err := iamc.GetRole(t.Context(), &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		t.Fatalf("GetRole failed: %v", err)
	}

	// Policy doc is URL-encoded JSON per AWS; decode automatically by SDK v2? (It returns string)
	// For robust substring checks, we just look for key substrings that must exist.

	policy := aws.ToString(getOut.Role.AssumeRolePolicyDocument)
	assert.NotEmpty(t, policy, "trust policy should not be empty")

	// --- Assertions you should tailor to your module’s policy ---
	// 1) Federated principal is the GitHub OIDC provider
	//    e.g., "arn:aws:iam::<acct>:oidc-provider/token.actions.githubusercontent.com"
	requireContains(t, policy, "oidc-provider/token.actions.githubusercontent.com")

	// 2) Audience condition for STS
	requireContains(t, policy, `"token.actions.githubusercontent.com:aud"`)
	requireContains(t, policy, `"sts.amazonaws.com"`)

	// 3) Subject condition scoping to repo/branch
	//    Typical pattern: repo:<owner>/<repo>:ref:refs/heads/<branch>
	requireContains(t, policy, "repo:pwjohnston2002/secrets-workflow")
	requireContains(t, policy, "refs/heads/main")

	// Optional small wait then verify role still present (exercise IAM eventual consistency)
	time.Sleep(2 * time.Second)
	_, err = iamc.GetRole(t.Context(), &iam.GetRoleInput{RoleName: aws.String(roleName)})
	assert.NoError(t, err, "role should still be retrievable after small delay")
}

// tiny helper for readable failures
func requireContains(t *testing.T, haystack, needle string) {
	if !strings.Contains(haystack, needle) {
		t.Fatalf("expected trust policy to contain %q; got:\n%s\n", needle, haystack)
	}
}
