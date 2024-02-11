package main

import (
	"io"
	"reflect"
	"testing"
)

func Test_createBatchTestApproval(t *testing.T) {
	type args struct {
		row CSVRow
	}
	tests := []struct {
		name string
		args args
		want BatchTestApproval
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createBatchTestApproval(tt.args.row); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createBatchTestApproval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createDrumPartition(t *testing.T) {
	type args struct {
		row      CSVRow
		rowIndex int
	}
	tests := []struct {
		name  string
		args  args
		want  DrumPartition
		want1 []Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := createDrumPartition(tt.args.row, tt.args.rowIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDrumPartition() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("createDrumPartition() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_createNewBatch(t *testing.T) {
	type args struct {
		row      CSVRow
		rowIndex int
	}
	tests := []struct {
		name  string
		args  args
		want  Batch
		want1 []Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := createNewBatch(tt.args.row, tt.args.rowIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createNewBatch() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("createNewBatch() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_createNewLI(t *testing.T) {
	type args struct {
		row      CSVRow
		rowIndex int
	}
	tests := []struct {
		name  string
		args  args
		want  LI
		want1 []Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := createNewLI(tt.args.row, tt.args.rowIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createNewLI() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("createNewLI() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_parseCSV(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name  string
		args  args
		want  UploadInventoryInput
		want1 []Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseCSV(tt.args.reader)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCSV() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseCSV() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_processRows(t *testing.T) {
	type args struct {
		rows []CSVRow
	}
	tests := []struct {
		name  string
		args  args
		want  UploadInventoryInput
		want1 []Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := processRows(tt.args.rows)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processRows() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("processRows() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_recordsToJSON(t *testing.T) {
	type args struct {
		records UploadInventoryInput
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := recordsToJSON(tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("recordsToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("recordsToJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateDrumPartition(t *testing.T) {
	type args struct {
		dp       DrumPartition
		row      CSVRow
		rowIndex int
	}
	tests := []struct {
		name  string
		args  args
		want  DrumPartition
		want1 []Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := updateDrumPartition(tt.args.dp, tt.args.row, tt.args.rowIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateDrumPartition() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("updateDrumPartition() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
