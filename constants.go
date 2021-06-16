package interserviceclient

/* #nosec */
const (
	defaultRegion = "KE"
	// Secret Key for signing json web tokens
	JWTSecretKey = "JWT_KEY"

	// The file that contains dependency definition. Each service which depends on other service
	// via REST, need to have this file in their root
	DepsFileName = "deps.yaml"

	// env variable pointing to where this service is running e.g staging, testing, prod
	Environment = "ENVIRONMENT"

	// running the service under staging
	StagingEnv = "staging"

	// running the service under demo
	DemoEnv = "demo"

	// running the service under testing
	TestingEnv = "testing"

	// running the service under production
	ProdEnv = "prod"

	// running the service under e2e
	E2eEnv = "e2e"

	// TestUserPhoneNumber is used by integration tests
	TestUserPhoneNumber = "+254711223344"

	// TestUserPhoneNumberWithPin is used by integration tests
	TestUserPhoneNumberWithPin = "+254778990088"
)
