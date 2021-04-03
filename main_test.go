package main

import (
	"reflect"
	"testing"
)

func Test_getGitIgnoreFiles(t *testing.T) {
	type args struct {
		root string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"Case 1",
			args{"testdata"},
			[]string{
				"testdata/.gitignore",
				"testdata/nested-gitignore/.gitignore",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getGitIgnoreFiles(tt.args.root); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getGitIgnoreFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getItems(t *testing.T) {
	type args struct {
		gitIgnoreFiles []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"Case 1",
			args{
				[]string{
					"testdata/.gitignore", "testdata/nested-gitignore/.gitignore",
				},
			},
			[]string{
				"testdata/nested-gitignore/should-be-ignored",
				"testdata/nested-gitignore/should-be-ignored.txt",
				"testdata/should be ignored",
				"testdata/should-be-ignored",
				"testdata/should-be-ignored.txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getItems(tt.args.gitIgnoreFiles); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter(t *testing.T) {
	type args struct {
		allItems []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"Case 1",
			args{
				[]string{
					"testdata/nested-gitignore",
					"testdata/nested-gitignore/should-be-ignored.txt",
					"testdata/should be ignored",
					"testdata/should-be-ignored",
					"testdata/should-be-ignored.txt",
				},
			},
			[]string{
				"testdata/nested-gitignore",
				"testdata/should be ignored",
				"testdata/should-be-ignored",
				"testdata/should-be-ignored.txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filter(tt.args.allItems); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readGitIgnoreFileContent(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"Case 1",
			args{"testdata/.gitignore"},
			[]string{
				"should-be-ignored/",
				"should-be-ignored.txt",
				"should\\ be\\ ignored/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readGitIgnoreFileContent(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readGitIgnoreFileContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
