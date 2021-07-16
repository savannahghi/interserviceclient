package interserviceclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/imroc/req"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/profileutils"
	"github.com/stretchr/testify/assert"
)

const (
	anonymousUserUID  = "AgkGYKUsRifO2O9fTLDuVCMr2hb2" // This is an anonymous user
	verifyPhone       = "testing/verify_phone"
	createUserByPhone = "testing/create_user_by_phone"
	loginByPhone      = "testing/login_by_phone"
	removeUserByPhone = "testing/remove_user"
	addAdmin          = "testing/add_admin_permissions"
	updateBioData     = "testing/update_user_profile"
)

// ContextKey is used as a type for the UID key for the Firebase *auth.Token on context.Context.
// It is a custom type in order to minimize context key collissions on the context
// (.and to shut up golint).
type ContextKey string

// GetInterserviceBearerTokenHeader returns a valid isc bearer token header
func GetInterserviceBearerTokenHeader(t *testing.T, rootDomain string, serviceName string) string {
	ctx := context.Background()
	isc := GetInterserviceClient(t, rootDomain, serviceName)
	authToken, err := isc.CreateAuthToken(ctx)
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

// GetPhoneNumberAuthenticatedContextAndToken returns a phone number logged in context
// and an auth Token that contains the the test user UID useful for test purposes
func GetPhoneNumberAuthenticatedContextAndToken(
	t *testing.T,
	onboardingClient *InterServiceClient,
) (context.Context, *auth.Token, error) {
	ctx := context.Background()
	userResponse, err := CreateOrLoginTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		return nil, nil, err
	}
	authToken := &auth.Token{
		UID: userResponse.Auth.UID,
	}
	authenticatedContext := context.WithValue(ctx, AuthTokenContextKey, authToken)
	return authenticatedContext, authToken, nil
}

// GetDefaultHeaders returns headers used in inter service communication acceptance tests
func GetDefaultHeaders(t *testing.T, rootDomain string, serviceName string) map[string]string {
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": GetInterserviceBearerTokenHeader(t, rootDomain, serviceName),
	}
}

// VerifyTestPhoneNumber checks if the test `Phone Number` exists as a primary
// phone number in any user profile record
func VerifyTestPhoneNumber(
	t *testing.T,
	phone string,
	onboardingClient *InterServiceClient,
) (string, error) {
	ctx := context.Background()

	verifyPhonePayload := map[string]interface{}{
		"phoneNumber": phone,
	}

	resp, err := onboardingClient.MakeRequest(
		ctx,
		http.MethodPost,
		verifyPhone,
		verifyPhonePayload,
	)

	if err != nil {
		return "", fmt.Errorf("unable to make a verify phone number request: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to convert response to string: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s", string(body))
	}

	var otp profileutils.OtpResponse
	err = json.Unmarshal(body, &otp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal OTP: %v", err)
	}

	return otp.OTP, nil
}

// LoginTestPhoneUser returns user response data for a created test user allowing
// them to run and access test resources
func LoginTestPhoneUser(
	t *testing.T,
	phone string,
	PIN string,
	flavour feedlib.Flavour,
	onboardingClient *InterServiceClient,
) (*profileutils.UserResponse, error) {
	ctx := context.Background()

	loginPayload := map[string]interface{}{
		"phoneNumber": phone,
		"pin":         PIN,
		"flavour":     flavour,
	}

	resp, err := onboardingClient.MakeRequest(
		ctx,
		http.MethodPost,
		loginByPhone,
		loginPayload,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to make a login request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to login : %s, with status code %v",
			phone,
			resp.StatusCode,
		)
	}
	code, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to convert response to string: %v", err)
	}

	var response *profileutils.UserResponse
	err = json.Unmarshal(code, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal OTP: %v", err)
	}

	return response, nil
}

// AddAdminPermissions adds ADMIN permissions to our test user
func AddAdminPermissions(
	t *testing.T,
	onboardingClient *InterServiceClient,
	phone string,
) error {
	ctx := context.Background()

	phonePayload := map[string]interface{}{
		"phoneNumber": phone,
	}

	resp, err := onboardingClient.MakeRequest(
		ctx,
		http.MethodPost,
		addAdmin,
		phonePayload,
	)

	if err != nil {
		return fmt.Errorf("unable to make add admin request: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to convert response to string: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got status code %v with resp body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// UpdateBioData adds Bio Data to our test user
func UpdateBioData(
	t *testing.T,
	onboardingClient *InterServiceClient,
	UID string,
) error {
	ctx := context.Background()

	bioDataPayload := map[string]interface{}{
		"uid":         UID,
		"firstName":   "Dumbledore 'the'",
		"lastName":    "Greatest Test User",
		"gender":      enumutils.GenderMale,
		"dateOfBirth": "2000-01-01",
	}

	resp, err := onboardingClient.MakeRequest(
		ctx,
		http.MethodPost,
		updateBioData,
		bioDataPayload,
	)

	if err != nil {
		return fmt.Errorf("unable to make update user profile request: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to convert response to string: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got status code %v with resp body: %s", resp.StatusCode, string(body))
	}
	return nil
}

// CreateOrLoginTestPhoneNumberUser creates an phone number test user if they
// do not exist or `Logs them in` if the test user exists to retrieve
// authenticated user response
// For documentation and test purposes only
func CreateOrLoginTestPhoneNumberUser(t *testing.T, onboardingClient *InterServiceClient) (*profileutils.UserResponse, error) {
	ctx := context.Background()

	phone := TestUserPhoneNumber
	PIN := TestUserPin
	flavour := feedlib.FlavourConsumer

	if onboardingClient == nil {
		return nil, fmt.Errorf("nil ISC client")
	}

	otp, err := VerifyTestPhoneNumber(t, phone, onboardingClient)
	if err != nil {
		if strings.Contains(
			err.Error(),
			strconv.Itoa(int(errorcodeutil.PhoneNumberInUse)),
		) {
			userResponse, err := LoginTestPhoneUser(
				t,
				phone,
				PIN,
				flavour,
				onboardingClient,
			)
			if err != nil {
				return nil, fmt.Errorf("unable to log in the test user: %v", err)
			}

			perms := userResponse.Profile.Permissions
			if len(perms) == 0 {
				err = AddAdminPermissions(t, onboardingClient, phone)
				if err != nil {
					return nil, fmt.Errorf("unable to add admin permissions: %v", err)
				}
			}

			if userResponse.Profile.UserBioData.FirstName == nil {
				err = UpdateBioData(t, onboardingClient, userResponse.Auth.UID)
				if err != nil {
					return nil, fmt.Errorf("unable to update user profile: %v", err)
				}
			}

			return userResponse, nil
		}

		return nil, fmt.Errorf("failed to verify test phone number: %v", err)
	}
	createUserPayload := map[string]interface{}{
		"phoneNumber": phone,
		"pin":         PIN,
		"flavour":     flavour,
		"otp":         otp,
	}

	resp, err := onboardingClient.MakeRequest(
		ctx,
		http.MethodPost,
		createUserByPhone,
		createUserPayload,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to make a sign up request: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unable to sign up : %s, with status code %v",
			phone,
			resp.StatusCode,
		)
	}
	signUpResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to convert response to string: %v", err)
	}

	var response *profileutils.UserResponse
	err = json.Unmarshal(signUpResp, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal OTP: %v", err)
	}

	perms := response.Profile.Permissions
	if len(perms) == 0 {
		err = AddAdminPermissions(t, onboardingClient, phone)
		if err != nil {
			return nil, fmt.Errorf("unable to add admin permissions: %v", err)
		}
	}

	if response.Profile.UserBioData.FirstName == nil {
		err = UpdateBioData(t, onboardingClient, response.Auth.UID)
		if err != nil {
			return nil, fmt.Errorf("unable to update user profile: %v", err)
		}
	}

	return response, nil
}
