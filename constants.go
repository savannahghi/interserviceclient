package interserviceclient

/* #nosec */
const (
	// Secret Key for signing json web tokens
	JWTSecretKey = "JWT_KEY"

	// TestUserPin used for testing purposes
	TestUserPin = "1234"

	// AuthTokenContextKey is used to add/retrieve the Firebase UID on the context
	AuthTokenContextKey = ContextKey("UID")

	// The file that contains dependency definition. Each service which depends on other service
	// via REST, need to have this file in their root
	DepsFileName = "deps.yaml"

	// running the service under e2e
	E2eEnv = "e2e"

	// TestUserPhoneNumber is used by integration tests
	TestUserPhoneNumber = "+254711223344"

	// TestUserPhoneNumberWithPin is used by integration tests
	TestUserPhoneNumberWithPin = "+254778990088"
)
