package main

import "testing"

func TestGetPackageGitHubRepo(t *testing.T) {
	type args struct {
		packageName string
		version     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "success", args: args{packageName: "vstest-oct1", version: "1.49.0"}, want: "https://api.github.com/repos/varunsh-coder/octokit.js"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPackageGitHubRepo(tt.args.packageName, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackageGitHubRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPackageGitHubRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}
