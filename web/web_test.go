package web

import "testing"

func TestGetBuild(t *testing.T) {
	type args struct {
		build Build
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"download",
			args{Build{Version: "2021-07-03-0750", Graphic: "tiles"}},
			"cataclysm-tiles-2021-07-03-0750.tar.gz",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBuild(tt.args.build)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBuild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}
