package interserviceclient_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"github.com/sirupsen/logrus"
)

const (
	// OnboardingRootDomain represents onboarding ISC URL
	OnboardingRootDomain = "https://profile-staging.healthcloud.co.ke"

	// OnboardingName represents the onboarding service ISC name
	OnboardingName = "onboarding"
)

func TestGetInterserviceClient(t *testing.T) {

	type args struct {
		OnboardingRootDomain string
		OnboardingName       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: successfully got interservice client",
			args: args{
				OnboardingRootDomain: "https://profile-staging.healthcloud.co.ke",
				OnboardingName:       "onboarding",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userResponse := interserviceclient.GetInterserviceClient(t, tt.args.OnboardingRootDomain, tt.args.OnboardingName)
			if tt.wantErr && userResponse != nil {
				t.Errorf("expected nil auth response but got %v",
					userResponse,
				)
				return
			}
			if !tt.wantErr && userResponse == nil {
				t.Errorf("expected an auth response but got nil, since no error occurred")
				return
			}

		})
	}
}

func TestGetInterserviceBearerTokenHeader(t *testing.T) {
	ctx := context.Background()
	service, _ := interserviceclient.NewInterserviceClient(interserviceclient.ISCService{Name: "otp", RootDomain: "https://example.com"})

	type args struct {
		ctx                  context.Context
		OnboardingRootDomain string
		OnboardingName       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: successfully got interservice client",
			args: args{
				ctx:                  ctx,
				OnboardingRootDomain: "https://example.com",
				OnboardingName:       "otp",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userResponse := interserviceclient.GetInterserviceBearerTokenHeader(t, tt.args.OnboardingRootDomain, tt.args.OnboardingName)
			c := service
			got, err := c.CreateAuthToken(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("InterServiceClient.CreateAuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && userResponse != "Bearer "+got {
				t.Errorf("expected an auth response but got nil, since no error occurred")
				return
			}

		})
	}
}
func onboardingISCClient() (*interserviceclient.InterServiceClient, error) {
	onboardingClient, err := interserviceclient.NewInterserviceClient(
		interserviceclient.ISCService{
			Name:       OnboardingName,
			RootDomain: OnboardingRootDomain,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize onboarding ISC client: %v", err)
	}
	return onboardingClient, nil
}

func TestVerifyTestPhoneNumber(t *testing.T) {
	onboardingClient, err := onboardingISCClient()
	if err != nil {
		t.Errorf("failed to initialize onboarding test ISC client")
	}
	type args struct {
		phone            string
		onboardingClient *interserviceclient.InterServiceClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: verify a phone number does not exist",
			args: args{
				phone:            "+254711999888",
				onboardingClient: onboardingClient,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otp, err := interserviceclient.VerifyTestPhoneNumber(
				t,
				tt.args.phone,
				tt.args.onboardingClient,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyTestPhoneNumber() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if tt.wantErr && otp != "" {
				t.Errorf("expected no otp to be sent but got %v, since the error %v occurred",
					otp,
					err,
				)
				return
			}

			if !tt.wantErr && otp == "" {
				t.Errorf("expected an otp to be sent, since no error occurred")
				return
			}
		})
	}
}

func TestCreateOrLoginTestPhoneNumberUser(t *testing.T) {
	onboardingClient, err := onboardingISCClient()
	if err != nil {
		t.Errorf("failed to initialize onboarding test ISC client")
	}
	type args struct {
		onboardingClient *interserviceclient.InterServiceClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: create a test user successfully",
			args: args{
				onboardingClient: onboardingClient,
			},
			wantErr: false,
		},
		{
			name: "failure: failed to create a test user successfully",
			args: args{
				onboardingClient: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userResponse, err := interserviceclient.CreateOrLoginTestPhoneNumberUser(t, tt.args.onboardingClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrLoginTestPhoneNumberUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && userResponse != nil {
				t.Errorf("expected nil auth response but got %v, since the error %v occurred",
					userResponse,
					err,
				)
				return
			}

			if !tt.wantErr && userResponse == nil {
				t.Errorf("expected an auth response but got nil, since no error occurred")
				return
			}
			if userResponse != nil {
				perms := userResponse.Profile.Permissions
				logrus.Print(perms)
				logrus.Print(userResponse.Profile.UserBioData)
			}
		})
	}

	// clean up
	err = interserviceclient.RemoveTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		t.Errorf("failed to remove test user: %v", err)
		return
	}
}

func TestRemoveTestPhoneNumberUser(t *testing.T) {
	onboardingClient, err := onboardingISCClient()
	if err != nil {
		t.Errorf("failed to initialize onboarding test ISC client")
		return
	}
	_, err = interserviceclient.CreateOrLoginTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		t.Errorf("unable to create user %v", err)
		return
	}

	type args struct {
		onboardingClient *interserviceclient.InterServiceClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: remove the created test user",
			args: args{
				onboardingClient: onboardingClient,
			},
			wantErr: false,
		},
		{
			name: "failure: failed to remove the created test user",
			args: args{
				onboardingClient: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := interserviceclient.RemoveTestPhoneNumberUser(
				t,
				tt.args.onboardingClient,
			); (err != nil) != tt.wantErr {
				t.Errorf("RemoveTestPhoneNumberUser() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestUpdateBioData(t *testing.T) {
	onboardingClient, err := onboardingISCClient()
	if err != nil {
		t.Errorf("failed to initialize onboarding test ISC client")
	}

	response, err := interserviceclient.CreateOrLoginTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		t.Errorf("unable to create user %v", err)
		return
	}

	type args struct {
		t                *testing.T
		onboardingClient *interserviceclient.InterServiceClient
		UID              string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy :) update bio data",
			args: args{
				t:                t,
				onboardingClient: onboardingClient,
				UID:              response.Auth.UID,
			},
			wantErr: false,
		},
		{
			name: "sad :( unable to update bio data",
			args: args{
				t:                t,
				onboardingClient: onboardingClient,
				UID:              "not-a-uid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := interserviceclient.UpdateBioData(tt.args.t, tt.args.onboardingClient, tt.args.UID); (err != nil) != tt.wantErr {
				t.Errorf("UpdateBioData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// clean up
	err = interserviceclient.RemoveTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		t.Errorf("failed to remove test user: %v", err)
		return
	}
}

func TestCreateTestRole(t *testing.T) {
	onboardingClient, err := onboardingISCClient()
	if err != nil {
		t.Errorf("failed to initialize onboarding test ISC client")
		return
	}

	response, err := interserviceclient.CreateOrLoginTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		t.Errorf("unable to create user %v", err)
		return
	}

	// used in test clean up
	var userID, roleID, roleName string
	userID = response.Profile.ID
	roleName = "Test Create Role"

	type args struct {
		t          *testing.T
		user       profileutils.UserResponse
		rootDomain string
		roleName   string
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.Role
		wantErr bool
	}{
		{
			name: "success: create test role",
			args: args{
				t:          t,
				user:       *response,
				rootDomain: OnboardingRootDomain,
				roleName:   roleName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := interserviceclient.CreateTestRole(tt.args.t, tt.args.user, tt.args.rootDomain, tt.args.roleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTestRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("CreateTestRole() = %v, want %v", got, tt.want)
			}
			roleID = got.ID
		})
	}

	// clean up
	interserviceclient.AssignTestRole(t, *response, OnboardingRootDomain, userID, roleID)
	interserviceclient.RemoveTestRole(t, *response, OnboardingRootDomain, roleName)
	err = interserviceclient.RemoveTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		t.Errorf("failed to remove test user: %v", err)
		return
	}
}

func TestAssignTestRole(t *testing.T) {
	onboardingClient, err := onboardingISCClient()
	if err != nil {
		t.Errorf("failed to initialize onboarding test ISC client")
		return
	}

	response, err := interserviceclient.CreateOrLoginTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		t.Errorf("unable to create user %v", err)
		return
	}

	// used in test clean up
	roleName := "Test Assign Role"

	role, err := interserviceclient.CreateTestRole(t, *response, OnboardingRootDomain, roleName)
	if err != nil {
		t.Errorf("unable to create test role %v", err)
		return
	}

	type args struct {
		t          *testing.T
		user       profileutils.UserResponse
		rootDomain string
		userID     string
		roleID     string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "success: assign test role",
			args: args{
				t:          t,
				user:       *response,
				rootDomain: OnboardingRootDomain,
				userID:     response.Profile.ID,
				roleID:     role.ID,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := interserviceclient.AssignTestRole(tt.args.t, tt.args.user, tt.args.rootDomain, tt.args.userID, tt.args.roleID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignTestRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AssignTestRole() = %v, want %v", got, tt.want)
			}
		})
	}

	// clean up
	_, err = interserviceclient.RemoveTestRole(t, *response, OnboardingRootDomain, roleName)
	if err != nil {
		t.Errorf("failed to remove test role: %v", err)
		return
	}
	err = interserviceclient.RemoveTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		t.Errorf("failed to remove test user: %v", err)
		return
	}
}

func TestRemoveTestRole(t *testing.T) {
	onboardingClient, err := onboardingISCClient()
	if err != nil {
		t.Errorf("failed to initialize onboarding test ISC client")
		return
	}

	response, err := interserviceclient.CreateOrLoginTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		t.Errorf("unable to create user %v", err)
		return
	}

	roleName := "Test Delete Role"

	role, err := interserviceclient.CreateTestRole(t, *response, OnboardingRootDomain, roleName)
	if err != nil {
		t.Errorf("unable to create test role %v", err)
		return
	}

	_, err = interserviceclient.AssignTestRole(t, *response, OnboardingRootDomain, response.Profile.ID, role.ID)
	if err != nil {
		t.Errorf("unable to assign test role %v", err)
		return
	}

	type args struct {
		t          *testing.T
		user       profileutils.UserResponse
		rootDomain string
		roleName   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "success: remove test role",
			args: args{
				t:          t,
				user:       *response,
				rootDomain: OnboardingRootDomain,
				roleName:   roleName,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := interserviceclient.RemoveTestRole(tt.args.t, tt.args.user, tt.args.rootDomain, tt.args.roleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveTestRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RemoveTestRole() = %v, want %v", got, tt.want)
			}
		})
	}

	// clean up
	err = interserviceclient.RemoveTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		t.Errorf("failed to remove test user: %v", err)
		return
	}
}

func TestCreateOrLoginTestPhoneNumberAuthorizedUser(t *testing.T) {
	onboardingClient, err := onboardingISCClient()
	if err != nil {
		t.Errorf("failed to initialize onboarding test ISC client")
	}

	type args struct {
		t                *testing.T
		onboardingClient *interserviceclient.InterServiceClient
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.UserResponse
		wantErr bool
	}{
		{
			name: "success: create a test user successfully",
			args: args{
				onboardingClient: onboardingClient,
			},
			wantErr: false,
		},
		{
			name: "failure: failed to create a test user successfully",
			args: args{
				onboardingClient: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := interserviceclient.CreateOrLoginTestPhoneNumberAuthorizedUser(tt.args.t, tt.args.onboardingClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrLoginTestPhoneNumberAuthorizedUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("CreateOrLoginTestPhoneNumberAuthorizedUser() = %v, want %v", got, tt.want)
			}
		})
	}

	// clean up
	err = interserviceclient.RemoveTestPhoneNumberAuthorizedUser(t, onboardingClient)
	if err != nil {
		t.Errorf("failed to remove test user: %v", err)
		return
	}
}
