package interserviceclient

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// GetInterserviceBearerTokenHeader returns a valid isc bearer token header
func GetInterserviceBearerTokenHeader(t *testing.T, rootDomain string, serviceName string) string {
	isc := GetInterserviceClient(t, rootDomain, serviceName)
	authToken, err := isc.CreateAuthToken()
	assert.Nil(t, err)
	assert.NotZero(t, authToken)
	bearerHeader := fmt.Sprintf("Bearer %s", authToken)
	return bearerHeader
}

// GetInterserviceClient returns an isc client used in acceptance testing
func GetInterserviceClient(t *testing.T, rootDomain string, serviceName string) *InterServiceClient {
	service := ISCService{
		Name:       serviceName,
		RootDomain: rootDomain,
	}
	isc, err := NewInterserviceClient(service)
	assert.Nil(t, err)
	assert.NotNil(t, isc)
	return isc
}
