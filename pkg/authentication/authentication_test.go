package authentication

import (
	"context"
	"reflect"
	"testing"
)

func TestEnv_AuthFunc(t *testing.T) {
	type fields struct {
		SecretSigningKey []byte
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    context.Context
		wantErr bool
	}{
		//{ name: "alice", fields: fields{SecretSigningKey: []byte("testtest")}, args: args{ctx: nil}, want: nil, wantErr: nil},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := &Env{
				SecretSigningKey: tt.fields.SecretSigningKey,
			}
			got, err := env.AuthFunc(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthFunc() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	type args struct {
		username      string
		signingSecret []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateToken(tt.args.username, tt.args.signingSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateToken(t *testing.T) {
	type args struct {
		tokenString   string
		signingSecret []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateToken(tt.args.tokenString, tt.args.signingSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}