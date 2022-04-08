package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_hasTerraformUserAgent(t *testing.T) {
	type args struct {
		userAgent string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no terraform useragent",
			args: args{
				userAgent: "aws-sdk-go/1.41.4 (go1.16.8; linux; amd64) amazon-ssm-agent/",
			},
			want: false,
		},
		{
			name: "has terraform useragent",
			args: args{
				userAgent: "APN/1.0 HashiCorp/1.0 Terraform",
			},
			want: true,
		},
		{
			name: "has terraform useragent with suffix",
			args: args{
				userAgent: "APN/1.0 HashiCorp/1.0 Terraform/blablabla",
			},
			want: true,
		},
		{
			name: "has surrounding parenthesis",
			args: args{
				userAgent: "[APN/1.0 HashiCorp/1.0 Terraform/blablabla]",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, hasTerraformUserAgent(tt.args.userAgent))
		})
	}
}
