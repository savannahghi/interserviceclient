package interserviceclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/savannahghi/profileutils"
)

// application endpoints
const (
	registerUser = "internal/register_user"
)

// OnboardingISC is a representation of an ISC client
type OnboardingISC interface {
	RegisterUser(ctx context.Context, payload interface{}) (*profileutils.UserProfile, error)
}

//OnboardingISCImpl represents the implemented methods in this ISC
type OnboardingISCImpl struct {
	isc *InterServiceClient
}

//NewOnboardingISC initializes a new instance of OnboardingISC
func NewOnboardingISC(isc *InterServiceClient) OnboardingISC {
	return &OnboardingISCImpl{
		isc: isc,
	}
}

//RegisterUser makes the request to register a user
func (o *OnboardingISCImpl) RegisterUser(ctx context.Context, payload interface{}) (*profileutils.UserProfile, error) {
	res, err := o.isc.MakeRequest(ctx, http.MethodPost, registerUser, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to send request, error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("register user failed with status code: %v", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	userprofile := profileutils.UserProfile{}

	err = json.Unmarshal(data, &userprofile)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body, %v", err)
	}

	return &userprofile, nil
}
