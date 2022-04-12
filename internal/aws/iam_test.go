package aws

import (
	"reflect"
	"testing"

	"github.com/hown3d/terraform-iam-generator/internal/metrics"
)

func TestGenerateIamPolicy(t *testing.T) {
	type args struct {
		msgs []metrics.CsmMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateIamPolicy(tt.args.msgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateIamPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateIamPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}
