package interserviceclient_test

// TODO enable this test to run once onboarding repository is updated
// func TestOnboardingISCImpl_RegisterUser(t *testing.T) {
// 	ctx := context.Background()
// 	isc, _ := interserviceclient.NewInterserviceClient(interserviceclient.ISCService{
// 		Name:       "onboarding",
// 		RootDomain: "https://profile-staging.healthcloud.co.ke",
// 	})

// 	onboarding := interserviceclient.NewOnboardingISC(isc)

// 	input := struct {
// 		FirstName   string
// 		LastName    string
// 		PhoneNumber string
// 		DateOfBirth scalarutils.Date
// 	}{
// 		FirstName:   "Test",
// 		LastName:    "Test2",
// 		PhoneNumber: interserviceclient.TestUserPhoneNumber,
// 		DateOfBirth: scalarutils.Date{
// 			Day:   1,
// 			Month: 1,
// 			Year:  2020,
// 		},
// 	}

// 	type args struct {
// 		ctx     context.Context
// 		payload interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *profileutils.UserProfile
// 		wantErr bool
// 	}{
// 		{
// 			name: "happy registered a user",
// 			args: args{
// 				ctx:     ctx,
// 				payload: input,
// 			},
// 			want:    &profileutils.UserProfile{},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := onboarding.RegisterUser(tt.args.ctx, tt.args.payload)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("OnboardingISCImpl.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("OnboardingISCImpl.RegisterUser() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
