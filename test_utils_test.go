package interserviceclient_test

import (
	"context"
	"testing"

	"github.com/savannahghi/interserviceclient"
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
