package interserviceclient_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/interserviceclient"
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
}

func TestAddAdminPermissions(t *testing.T) {
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

	type args struct {
		t                *testing.T
		onboardingClient *interserviceclient.InterServiceClient
		phone            string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case :)",
			args: args{
				t:                t,
				onboardingClient: onboardingClient,
				phone:            *response.Profile.PrimaryPhone,
			},
			wantErr: false,
		},
		{
			name: "sad case :(",
			args: args{
				t:                t,
				onboardingClient: onboardingClient,
				phone:            "not-a-phone-number",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := interserviceclient.AddAdminPermissions(tt.args.t, tt.args.onboardingClient, tt.args.phone); (err != nil) != tt.wantErr {
				t.Errorf("AddAdminPermissions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
