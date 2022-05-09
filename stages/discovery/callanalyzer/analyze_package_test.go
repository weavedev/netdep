package callanalyzer

import (
	"reflect"
	"testing"
)

func TestAnalyzePackageCalls(t *testing.T) {
	type args struct {
		pkg *ssa.Package
	}
	tests := []struct {
		name string
		args args
		want []*Caller
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnalyzePackageCalls(tt.args.pkg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnalyzePackageCalls() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoverBlock(t *testing.T) {
	type args struct {
		block *ssa.BasicBlock
	}
	tests := []struct {
		name string
		args args
		want []*Caller
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := discoverBlock(tt.args.block); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("discoverBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoverBlocks(t *testing.T) {
	type args struct {
		blocks []*ssa.BasicBlock
	}
	tests := []struct {
		name string
		args args
		want []*Caller
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := discoverBlocks(tt.args.blocks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("discoverBlocks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoverCall(t *testing.T) {
	type args struct {
		call *ssa.Call
	}
	tests := []struct {
		name string
		args args
		want *Caller
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := discoverCall(tt.args.call); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("discoverCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMainFunction(t *testing.T) {
	type args struct {
		pkg *ssa.Package
	}
	tests := []struct {
		name string
		args args
		want *ssa.Function
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMainFunction(tt.args.pkg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMainFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}
