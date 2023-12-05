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
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
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
	createRole        = "roles/create_role"
	assignRole        = "roles/assign_role"
	removeRoleByName  = "roles/remove_role"
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
	authenticatedContext := context.WithValue(ctx, firebasetools.AuthTokenContextKey, authToken)
	return authenticatedContext, authToken, nil
}

// GetTestAuthorizedContextAndToken returns an authorized phone number with permissions logged in context
// and an auth Token that contains the the test user UID useful for test purposes
func GetTestAuthorizedContextAndToken(
	t *testing.T,
	onboardingClient *InterServiceClient,
) (context.Context, *auth.Token, error) {
	ctx := context.Background()
	userResponse, err := CreateOrLoginTestPhoneNumberAuthorizedUser(t, onboardingClient)
	if err != nil {
		return nil, nil, err
	}
	authToken := &auth.Token{
		UID: userResponse.Auth.UID,
	}
	authenticatedContext := context.WithValue(ctx, firebasetools.AuthTokenContextKey, authToken)
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
		"appID":       uuid.NewString(),
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

			return userResponse, nil
		}

		return nil, fmt.Errorf("failed to verify test phone number: %v", err)
	}

	response, err := CreateTestPhoneNumberUser(t, onboardingClient, otp)
	if err != nil {
		return nil, fmt.Errorf("unable to create test user:%v", err)
	}

	return response, nil
}

// CreateOrLoginTestPhoneNumberAuthorizedUser creates an phone number test user if they
// do not exist or `Logs them in` if the test user exists to retrieve
// authenticated user response
// For documentation and test purposes only
func CreateOrLoginTestPhoneNumberAuthorizedUser(t *testing.T, onboardingClient *InterServiceClient) (*profileutils.UserResponse, error) {
	userResponse, err := CreateOrLoginTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		return nil, err
	}

	return userResponse, nil
}

// CreateTestPhoneNumberUser creates the user for test phone number
func CreateTestPhoneNumberUser(t *testing.T, onboardingClient *InterServiceClient, otp string) (*profileutils.UserResponse, error) {
	ctx := context.Background()

	phone := TestUserPhoneNumber
	PIN := TestUserPin
	flavour := feedlib.FlavourConsumer

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

	return response, nil
}

// RemoveTestPhoneNumberUser removes the records created by the
// test phonenumber user
func RemoveTestPhoneNumberUser(
	t *testing.T,
	onboardingClient *InterServiceClient,
) error {
	ctx := context.Background()

	if onboardingClient == nil {
		return fmt.Errorf("nil ISC client")
	}

	payload := map[string]interface{}{
		"phoneNumber": TestUserPhoneNumber,
	}
	resp, err := onboardingClient.MakeRequest(
		ctx,
		http.MethodPost,
		removeUserByPhone,
		payload,
	)
	if err != nil {
		return fmt.Errorf("unable to make a request to remove test user: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil // This is a test utility. Do not block if the user is not found
	}

	return nil
}

// RemoveTestPhoneNumberAuthorizedUser removes the records created by the
// test phonenumber user
func RemoveTestPhoneNumberAuthorizedUser(
	t *testing.T,
	onboardingClient *InterServiceClient,
) error {
	if onboardingClient == nil {
		return fmt.Errorf("nil ISC client")
	}

	return RemoveTestPhoneNumberUser(t, onboardingClient)
}

// GetTestGraphQLHeaders gets relevant GraphQLHeaders for running
// GraphQL acceptance tests
func GetTestGraphQLHeaders(
	t *testing.T,
	onboardingClient *InterServiceClient,
) (map[string]string, error) {
	authorization, err := GetTestBearerTokenHeader(t, onboardingClient)
	if err != nil {
		return nil, fmt.Errorf("can't Generate Bearer Token: %s", err)
	}
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": authorization,
	}, nil
}

// GetTestBearerTokenHeader gets bearer Token Header for running
// GraphQL acceptance tests
func GetTestBearerTokenHeader(
	t *testing.T,
	onboardingClient *InterServiceClient,
) (string, error) {
	user, err := CreateOrLoginTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Bearer %s", *user.Auth.IDToken), nil
}
