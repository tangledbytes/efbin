package ef

import (
	"reflect"
	"testing"
)

func TestConvertStringSetToFlags(t *testing.T) {
	type args struct {
		set []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "When slice is empty",
			args: args{
				set: []string{},
			},
			want: []string{},
		},
		{
			name: "When slice has set strings with len(key) > 1",
			args: args{
				set: []string{"IJK=1", "hello=world", "platform=linux/amd64,linux/arm64", "build_arg=uid=501,gid=20,username=utkarsh"},
			},
			want: []string{"--IJK", "1", "--hello", "world", "--platform", "linux/amd64,linux/arm64", "--build-arg", "uid=501,gid=20,username=utkarsh"},
		},
		{
			name: "When slice has set strings with len(key) = 1",
			args: args{
				set: []string{"i=", "t=", "p=access"},
			},
			want: []string{"-i", "-t", "-p", "access"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertStringSetToFlags(tt.args.set); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertStringSetToFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransformArgsWithSet(t *testing.T) {
	type args struct {
		args   []string
		set    []string
		target string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "When everything is empty",
			args: args{
				args:   []string{},
				set:    []string{},
				target: "",
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "When target and set are empty",
			args: args{
				args:   []string{"arg1", "arg2"},
				set:    []string{},
				target: "",
			},
			want:    []string{"arg1", "arg2"},
			wantErr: false,
		},
		{
			name: "When only set is empty",
			args: args{
				args:   []string{"arg1", "arg2"},
				set:    []string{},
				target: "arg1",
			},
			want:    []string{"arg1", "arg2"},
			wantErr: false,
		},
		{
			name: "When only target is empty",
			args: args{
				args:   []string{"arg1", "arg2"},
				set:    []string{"-i", "-t", "-p", "access"},
				target: "",
			},
			want:    []string{"-i", "-t", "-p", "access", "arg1", "arg2"},
			wantErr: false,
		},
		{
			name: "When target is the first argument with no skips",
			args: args{
				args:   []string{"arg1", "arg2"},
				set:    []string{"-i", "-t", "-p", "access"},
				target: "arg1",
			},
			want:    []string{"arg1", "-i", "-t", "-p", "access", "arg2"},
			wantErr: false,
		},
		{
			name: "When target is the last argument with no skips",
			args: args{
				args:   []string{"arg1", "arg2"},
				set:    []string{"-i", "-t", "-p", "access"},
				target: "arg2",
			},
			want:    []string{"arg1", "arg2", "-i", "-t", "-p", "access"},
			wantErr: false,
		},
		{
			name: "When target is the first argument with 0 skips",
			args: args{
				args:   []string{"arg1", "arg2"},
				set:    []string{"-i", "-t", "-p", "access"},
				target: "arg1:0",
			},
			want:    []string{"arg1", "-i", "-t", "-p", "access", "arg2"},
			wantErr: false,
		},
		{
			name: "When target is the first argument's duplicate with 1 skips",
			args: args{
				args:   []string{"arg1", "arg1", "arg2"},
				set:    []string{"-i", "-t", "-p", "access"},
				target: "arg1:1",
			},
			want:    []string{"arg1", "arg1", "-i", "-t", "-p", "access", "arg2"},
			wantErr: false,
		},
		{
			name: "When target is the first argument's duplicate with 100 (invalid) skips",
			args: args{
				args:   []string{"arg1", "arg1", "arg2"},
				set:    []string{"-i", "-t", "-p", "access"},
				target: "arg1:100",
			},
			want:    []string{"arg1", "arg1", "arg2"},
			wantErr: false,
		},
		{
			name: "When target is invalid with skips",
			args: args{
				args:   []string{"arg1", "arg1", "arg2"},
				set:    []string{},
				target: "arg1:100AB",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "When target is invalid with incorrect structure",
			args: args{
				args:   []string{"arg1", "arg1", "arg2"},
				set:    []string{},
				target: "arg1:100AB:abcd",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransformArgsWithSet(tt.args.args, tt.args.set, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformArgsWithSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransformArgsWithSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
