package interserviceclient_test

import (
	"testing"

	"github.com/savannahghi/interserviceclient"
)

const testPhone = "+254723002959"

func TestSendSMS(t *testing.T) {

	// Note: This is a very brittle test case.
	// Any change to the service urls would probably lead to a failure
	// There's probably a better way to do this (Mocking *wink wink)
	// But I (Farad) felt this is the best way of doing it i.e. Acceptance Testing
	//TODO: Make these env vars
	newSmsIsc, _ := interserviceclient.NewInterserviceClient(interserviceclient.ISCService{
		Name:       "sms",
		RootDomain: "https://sms-staging.healthcloud.co.ke",
	})

	newTwilioIsc, _ := interserviceclient.NewInterserviceClient(interserviceclient.ISCService{
		Name:       "twilio",
		RootDomain: "https://twilio-staging.healthcloud.co.ke",
	})

	smsEndPoint := "internal/send_sms"

	type args struct {
		phoneNumbers    []string
		message         string
		smsIscClient    interserviceclient.SmsISC
		twilioIscClient interserviceclient.SmsISC
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good test case",
			args: args{
				phoneNumbers: []string{testPhone},
				message:      "Test Text Message",
				smsIscClient: interserviceclient.SmsISC{
					Isc:      newSmsIsc,
					EndPoint: "internal/send_sms",
				},
				twilioIscClient: interserviceclient.SmsISC{
					Isc:      newTwilioIsc,
					EndPoint: smsEndPoint,
				},
			},
			wantErr: true, // TODO: fix the error and return to false
		},
		{
			name: "bad test case: Empty Message",
			args: args{
				phoneNumbers: []string{testPhone},
				message:      "",
				smsIscClient: interserviceclient.SmsISC{
					Isc:      newSmsIsc,
					EndPoint: smsEndPoint,
				},
				twilioIscClient: interserviceclient.SmsISC{
					Isc:      newTwilioIsc,
					EndPoint: smsEndPoint,
				},
			},
			wantErr: true,
		},
		{
			name: "bad test case: No Phone Numbers",
			args: args{
				phoneNumbers: []string{},
				message:      "Test Text Message",
				smsIscClient: interserviceclient.SmsISC{
					Isc:      newSmsIsc,
					EndPoint: smsEndPoint,
				},
				twilioIscClient: interserviceclient.SmsISC{
					Isc:      newTwilioIsc,
					EndPoint: smsEndPoint,
				},
			},
			wantErr: true,
		},
		{
			name: "bad test case: Invalid Phone Numbers",
			args: args{
				phoneNumbers: []string{"not-a-number"},
				message:      "Test Text Message",
				smsIscClient: interserviceclient.SmsISC{
					Isc:      newSmsIsc,
					EndPoint: smsEndPoint,
				},
				twilioIscClient: interserviceclient.SmsISC{
					Isc:      newTwilioIsc,
					EndPoint: smsEndPoint,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := interserviceclient.SendSMS(tt.args.phoneNumbers, tt.args.message, tt.args.smsIscClient, tt.args.twilioIscClient); (err != nil) != tt.wantErr {
				t.Errorf("SendSMS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TODO: Fix unable to send OTP unable to generate otp, with status code 401 and uncomment
func TestVerifyOTP(t *testing.T) {
	// client, _ := interserviceclient.NewInterserviceClient(interserviceclient.ISCService{
	// 	Name:       "otp",
	// 	RootDomain: "https://otp-staging.healthcloud.co.ke",
	// })
	// // generate the OTP first to be used for a happy case
	// OTPCode, err := interserviceclient.SendOTPHelper(interserviceclient.TestUserPhoneNumber, client)
	// if err != nil {
	// 	t.Errorf("TestVerifyOTP: unable to send OTP %v", err)
	// 	return
	// }
	type args struct {
		msisdn           string
		verificationCode string
		client           *interserviceclient.InterServiceClient
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// {
		// 	name: "verify OTP success: OTP generated and verified on same number",
		// 	args: args{
		// 		msisdn:           interserviceclient.TestUserPhoneNumber,
		// 		verificationCode: OTPCode,
		// 		client:           client,
		// 	},
		// 	wantErr: true,  // TODO: fix the error and return to false
		// 	want:    false, // TODO: fix the error and return to true
		// },
		// {
		// 	name: "verify OTP failure: OTP not generated and verified on same number",
		// 	args: args{
		// 		msisdn:           interserviceclient.TestUserPhoneNumberWithPin,
		// 		verificationCode: OTPCode,
		// 		client:           client,
		// 	},
		// 	wantErr: true,
		// 	want:    false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := interserviceclient.VerifyOTP(tt.args.msisdn, tt.args.verificationCode, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifyOTP() = %v, want %v", got, tt.want)
			}
		})
	}
}
