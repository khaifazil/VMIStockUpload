package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVRow_UnmarshalCSV(t *testing.T) {
	type fields struct {
		Vendor                  string
		MaterialCode            string
		MaterialDesc            string
		ContractNo              string
		PONumber                string
		POLineItem              string
		LIName                  LIName
		LIDate                  string
		BatchNo                 string
		BatchDueDate            string
		DrumSize                int
		TotalNoOfDrums          int
		TotalQty                int
		AvailableDrumNos        []int
		AvailableFullDrums      int
		FullDrumTotalQuantity   int
		BufferDrumNo            []int
		BufferNoOfDrums         int
		BufferQuantity          int
		SampleDrum              string
		SampleDrumNo            []int
		SampleLength            []float64
		NoOfShortLengthDrums    int
		ShortLengthTotalQty     float64
		ApprovedDrumNumbers     []int
		BatchTestReportDate     string
		Remarks                 string
		BatchTestReportFileName string
	}
	type args struct {
		csv      []string
		rowIndex int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Error
	}{
		{
			name: "Test Case 1: Valid CSV Row",
			fields: fields{
				Vendor:                  "Vendor1",
				MaterialCode:            "Material1",
				MaterialDesc:            "Description1",
				ContractNo:              "Contract1",
				PONumber:                "PONumber1",
				POLineItem:              "POLineItem1",
				LIName:                  LIName{LICode: "LICode1", LINumber: "LINumber1"},
				LIDate:                  "01-01-2022",
				BatchNo:                 "BatchNo1",
				BatchDueDate:            "01-01-2022",
				DrumSize:                500,
				TotalNoOfDrums:          10,
				TotalQty:                5000,
				AvailableDrumNos:        []int{1, 2, 3, 4, 5},
				AvailableFullDrums:      5,
				FullDrumTotalQuantity:   2500,
				BufferDrumNo:            []int{6, 7},
				BufferNoOfDrums:         2,
				BufferQuantity:          1000,
				SampleDrum:              "Yes",
				SampleDrumNo:            []int{8, 9},
				SampleLength:            []float64{100.0, 200.0},
				NoOfShortLengthDrums:    2,
				ShortLengthTotalQty:     300.0,
				ApprovedDrumNumbers:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				BatchTestReportDate:     "01-01-2022",
				Remarks:                 "Remarks1",
				BatchTestReportFileName: "FileName1",
			},
			args: args{
				csv:      []string{"Vendor1", "Material1", "Description1", "Contract1", "PONumber1", "POLineItem1", "LICode1-LINumber1", "01-01-2022", "BatchNo1", "01-01-2022", "500", "10", "1-5", "5", "2500", "6-7", "2", "1000", "Yes", "8-9", "100,200", "2", "300", "01-01-2022", "Remarks1", "FileName1"},
				rowIndex: 1,
			},
			want: []Error{},
		},
		{
			name: "Test Case 2: Invalid CSV Row",
			fields: fields{
				Vendor:                  "Vendor2",
				MaterialCode:            "Material2",
				MaterialDesc:            "Description2",
				ContractNo:              "Contract2",
				PONumber:                "PONumber2",
				POLineItem:              "POLineItem2",
				LIName:                  LIName{LICode: "LICode2", LINumber: "LINumber2"},
				LIDate:                  "01-01-2022",
				BatchNo:                 "BatchNo2",
				BatchDueDate:            "01-01-2022",
				DrumSize:                500,
				TotalNoOfDrums:          10,
				TotalQty:                5000,
				AvailableDrumNos:        []int{1, 2, 3, 4, 5},
				AvailableFullDrums:      5,
				FullDrumTotalQuantity:   2500,
				BufferDrumNo:            []int{6, 7},
				BufferNoOfDrums:         2,
				BufferQuantity:          1000,
				SampleDrum:              "Yes",
				SampleDrumNo:            []int{8, 9},
				SampleLength:            []float64{100.0, 200.0},
				NoOfShortLengthDrums:    2,
				ShortLengthTotalQty:     300.0,
				ApprovedDrumNumbers:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				BatchTestReportDate:     "01-01-2022",
				Remarks:                 "Remarks2",
				BatchTestReportFileName: "FileName2",
			},
			args: args{
				csv:      []string{"Vendor2", "Material2", "Description2", "Contract2", "PONumber2", "POLineItem2", "LICode2", "01-01-2022", "BatchNo2", "01-01-2022", "500", "10", "1-5", "5", "2500", "6-7", "2", "1000", "Yes", "8-9", "100,200", "2", "300", "01-01-2022", "Remarks2", "FileName2"},
				rowIndex: 1,
			},
			want: []Error{
				{RowNo: 2, Err: fmt.Errorf("invalid LI Name format")},
			},
		},
		{
			name: "Test Case 3: Valid CSV Row with Different Values",
			fields: fields{
				Vendor:                  "Vendor4",
				MaterialCode:            "Material4",
				MaterialDesc:            "Description4",
				ContractNo:              "Contract4",
				PONumber:                "PONumber4",
				POLineItem:              "POLineItem4",
				LIName:                  LIName{LICode: "LICode4", LINumber: "LINumber4"},
				LIDate:                  "04-04-2024",
				BatchNo:                 "4/4",
				BatchDueDate:            "04-04-2024",
				DrumSize:                1000,
				TotalNoOfDrums:          4,
				TotalQty:                4000,
				AvailableDrumNos:        []int{1, 2},
				AvailableFullDrums:      2,
				FullDrumTotalQuantity:   2000,
				BufferDrumNo:            []int{3},
				BufferNoOfDrums:         1,
				BufferQuantity:          1000,
				SampleDrum:              "Yes",
				SampleDrumNo:            []int{4},
				SampleLength:            []float64{500.0},
				NoOfShortLengthDrums:    1,
				ShortLengthTotalQty:     500.0,
				ApprovedDrumNumbers:     []int{1, 2, 3, 4},
				BatchTestReportDate:     "04-04-2024",
				Remarks:                 "Remarks4",
				BatchTestReportFileName: "FileName4",
			},
			args: args{
				csv:      []string{"Vendor4", "Material4", "Description4", "Contract4", "PONumber4", "POLineItem4", "LICode4-LINumber4", "04-04-2024", "4/4", "04-04-2024", "1000", "4", "1-2", "2", "2000", "3", "1", "1000", "Yes", "4", "500", "1", "500", "04-04-2024", "Remarks4", "FileName4"},
				rowIndex: 4,
			},
			want: []Error{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := &CSVRow{
				Vendor:                  tt.fields.Vendor,
				MaterialCode:            tt.fields.MaterialCode,
				MaterialDesc:            tt.fields.MaterialDesc,
				ContractNo:              tt.fields.ContractNo,
				PONumber:                tt.fields.PONumber,
				POLineItem:              tt.fields.POLineItem,
				LIName:                  tt.fields.LIName,
				LIDate:                  tt.fields.LIDate,
				BatchNo:                 tt.fields.BatchNo,
				BatchDueDate:            tt.fields.BatchDueDate,
				DrumSize:                tt.fields.DrumSize,
				TotalNoOfDrums:          tt.fields.TotalNoOfDrums,
				TotalQty:                tt.fields.TotalQty,
				AvailableDrumNos:        tt.fields.AvailableDrumNos,
				AvailableFullDrums:      tt.fields.AvailableFullDrums,
				FullDrumTotalQuantity:   tt.fields.FullDrumTotalQuantity,
				BufferDrumNo:            tt.fields.BufferDrumNo,
				BufferNoOfDrums:         tt.fields.BufferNoOfDrums,
				BufferQuantity:          tt.fields.BufferQuantity,
				SampleDrum:              tt.fields.SampleDrum,
				SampleDrumNo:            tt.fields.SampleDrumNo,
				SampleLength:            tt.fields.SampleLength,
				NoOfShortLengthDrums:    tt.fields.NoOfShortLengthDrums,
				ShortLengthTotalQty:     tt.fields.ShortLengthTotalQty,
				ApprovedDrumNumbers:     tt.fields.ApprovedDrumNumbers,
				BatchTestReportDate:     tt.fields.BatchTestReportDate,
				Remarks:                 tt.fields.Remarks,
				BatchTestReportFileName: tt.fields.BatchTestReportFileName,
			}
			got := row.UnmarshalCSV(tt.args.csv, tt.args.rowIndex)
			assert.Equal(t, tt.want, got, "UnmarshalCSV() = %v, want %v", got, tt.want)
		})
	}
}
