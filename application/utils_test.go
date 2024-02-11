package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_combineSortAndCheckDuplicates(t *testing.T) {
	type args struct {
		slices [][]int
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "No duplicates",
			args: args{
				slices: [][]int{{1, 2, 3}, {4, 5, 6}},
			},
			want:    []int{1, 2, 3, 4, 5, 6},
			wantErr: false,
		},
		{
			name: "With duplicates",
			args: args{
				slices: [][]int{{1, 2, 3}, {3, 4, 5}},
			},
			want:    []int{1, 2, 3, 3, 4, 5},
			wantErr: true,
		},
		{
			name: "Empty slices",
			args: args{
				slices: [][]int{{}, {}},
			},
			want:    []int{},
			wantErr: false,
		},
		{
			name: "Single slice",
			args: args{
				slices: [][]int{{1, 2, 3}},
			},
			want:    []int{1, 2, 3},
			wantErr: false,
		},
		{
			name: "Single slice with duplicates",
			args: args{
				slices: [][]int{{1, 2, 2, 3}},
			},
			want:    []int{1, 2, 2, 3},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := combineSortAndCheckDuplicates(tt.args.slices...)
			if (err != nil) != tt.wantErr {
				t.Errorf("combineSortAndCheckDuplicates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_findBatchIndex(t *testing.T) {
	type args struct {
		batches []Batch
		no      string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Batch found",
			args: args{
				batches: []Batch{
					{BatchNo: "Batch1"},
					{BatchNo: "Batch2"},
					{BatchNo: "Batch3"},
				},
				no: "Batch2",
			},
			want: 1,
		},
		{
			name: "Batch not found",
			args: args{
				batches: []Batch{
					{BatchNo: "Batch1"},
					{BatchNo: "Batch2"},
					{BatchNo: "Batch3"},
				},
				no: "Batch4",
			},
			want: -1,
		},
		{
			name: "Empty batches",
			args: args{
				batches: []Batch{},
				no:      "Batch1",
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findBatchIndex(tt.args.batches, tt.args.no)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_findBatchTestApprovalIndex(t *testing.T) {
	type args struct {
		bta  []BatchTestApproval
		date string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "BatchTestApproval found",
			args: args{
				bta: []BatchTestApproval{
					{ApprovalDate: "2022-01-01"},
					{ApprovalDate: "2022-01-02"},
					{ApprovalDate: "2022-01-03"},
				},
				date: "2022-01-02",
			},
			want: 1,
		},
		{
			name: "BatchTestApproval not found",
			args: args{
				bta: []BatchTestApproval{
					{ApprovalDate: "2022-01-01"},
					{ApprovalDate: "2022-01-02"},
					{ApprovalDate: "2022-01-03"},
				},
				date: "2022-01-04",
			},
			want: -1,
		},
		{
			name: "Empty BatchTestApproval",
			args: args{
				bta:  []BatchTestApproval{},
				date: "2022-01-01",
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findBatchTestApprovalIndex(tt.args.bta, tt.args.date)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_findContractIndex(t *testing.T) {
	type args struct {
		slice      []Contracts
		contractNo string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Contract found",
			args: args{
				slice: []Contracts{
					{ContractNo: "Contract1"},
					{ContractNo: "Contract2"},
					{ContractNo: "Contract3"},
				},
				contractNo: "Contract2",
			},
			want: 1,
		},
		{
			name: "Contract not found",
			args: args{
				slice: []Contracts{
					{ContractNo: "Contract1"},
					{ContractNo: "Contract2"},
					{ContractNo: "Contract3"},
				},
				contractNo: "Contract4",
			},
			want: -1,
		},
		{
			name: "Empty Contracts",
			args: args{
				slice:      []Contracts{},
				contractNo: "Contract1",
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findContractIndex(tt.args.slice, tt.args.contractNo)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_findLiIndex(t *testing.T) {
	type args struct {
		slice []LI
		liNo  LIName
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "LI found",
			args: args{
				slice: []LI{
					{LiCode: "LI1", LiNumber: "1"},
					{LiCode: "LI2", LiNumber: "2"},
					{LiCode: "LI3", LiNumber: "3"},
				},
				liNo: LIName{
					LICode:   "LI2",
					LINumber: "2",
				},
			},
			want: 1,
		},
		{
			name: "LI not found",
			args: args{
				slice: []LI{
					{LiCode: "LI1", LiNumber: "1"},
					{LiCode: "LI2", LiNumber: "2"},
					{LiCode: "LI3", LiNumber: "3"},
				},
				liNo: LIName{
					LICode:   "LI1",
					LINumber: "4",
				},
			},
			want: -1,
		},
		{
			name: "Empty LI",
			args: args{
				slice: []LI{},
				liNo: LIName{
					LICode:   "LI1",
					LINumber: "1",
				},
			},
			want: -1,
		},
		{
			name: "Empty LI",
			args: args{
				slice: []LI{},
				liNo:  LIName{},
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findLiIndex(tt.args.slice, tt.args.liNo)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_unpackSampleDrumNos(t *testing.T) {
	type args struct {
		sampleDrumNumbers []int
		sampleLength      []float64
		drumSize          int
	}
	tests := []struct {
		name  string
		args  args
		want  []DrumDetails
		want1 float64
		want2 []DrumDetails
		want3 float64
	}{
		{
			name: "valid sample drum numbers and sample length",
			args: args{
				sampleDrumNumbers: []int{1, 2, 3},
				sampleLength:      []float64{2.5, 2.0, 5.0},
				drumSize:          250,
			},
			want:  []DrumDetails{{DrumNumber: 1, Quantity: 2.5}, {DrumNumber: 2, Quantity: 2.0}, {DrumNumber: 3, Quantity: 5.0}},
			want1: 9.5,
			want2: []DrumDetails{{DrumNumber: 1, Quantity: 247.5}, {DrumNumber: 2, Quantity: 248}, {DrumNumber: 3, Quantity: 245}},
			want3: 740.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := unpackSampleDrumNos(tt.args.sampleDrumNumbers, tt.args.sampleLength, tt.args.drumSize)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.want3, got3)
		})
	}
}

func Test_sumFloat64Slice(t *testing.T) {
	tests := []struct {
		name   string
		length []float64
		want   float64
	}{
		{
			name:   "PositiveNumbers",
			length: []float64{1.1, 2.2, 3.3},
			want:   6.6,
		},
		{
			name:   "IncludingZero",
			length: []float64{0, 1.1, 2.2, 3.3},
			want:   6.6,
		},
		{
			name:   "NegativeNumbers",
			length: []float64{-1.1, -2.2, -3.3},
			want:   -6.6,
		},
		{
			name:   "EmptySlice",
			length: []float64{},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sumFloat64Slice(tt.length)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_unpackDrumNoRange(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    []int
		wantErr bool
	}{
		{
			name:    "ValidRange",
			str:     "1-5",
			want:    []int{1, 2, 3, 4, 5},
			wantErr: false,
		},
		{
			name:    "SingleNumber",
			str:     "1",
			want:    []int{1},
			wantErr: false,
		},
		{
			name:    "EmptyString",
			str:     "",
			want:    []int{},
			wantErr: false,
		},
		{
			name:    "NonNumericCharacter",
			str:     "a-b",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Skip number",
			str:     "1-3,5",
			want:    []int{1, 2, 3, 5},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unpackDrumNoRange(tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackDrumNoRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_stringToFloat64Slice(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    []float64
		wantErr bool
	}{
		{
			name:    "ValidFloats",
			str:     "1.1,2.2,3.3",
			want:    []float64{1.1, 2.2, 3.3},
			wantErr: false,
		},
		{
			name:    "InvalidFloat",
			str:     "1.1,a,3.3",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "zero",
			str:     "0",
			want:    []float64{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringToFloat64Slice(tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringToFloat64Slice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_removeDuplicateDrumNumbers(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		dupSlice []int
		want     []int
	}{
		{
			name:     "NoDuplicates",
			slice:    []int{1, 2, 3},
			dupSlice: []int{4, 5, 6},
			want:     []int{1, 2, 3},
		},
		{
			name:     "AllDuplicates",
			slice:    []int{1, 2, 3},
			dupSlice: []int{1, 2, 3},
			want:     []int{},
		},
		{
			name:     "SomeDuplicates",
			slice:    []int{1, 2, 3, 4, 5},
			dupSlice: []int{2, 4},
			want:     []int{1, 3, 5},
		},
		{
			name:     "EmptySlice",
			slice:    []int{},
			dupSlice: []int{1, 2, 3},
			want:     []int{},
		},
		{
			name:     "EmptyDupSlice",
			slice:    []int{1, 2, 3},
			dupSlice: []int{},
			want:     []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeDuplicateDrumNumbers(tt.slice, tt.dupSlice)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_validateBatchNoFormat(t *testing.T) {
	type args struct {
		batchNo string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid batch number",
			args: args{
				batchNo: "10/15",
			},
			want: true,
		},
		{
			name: "Invalid batch number",
			args: args{
				batchNo: "123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateBatchNoFormat(tt.args.batchNo)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_validateDateFormat(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid date",
			args: args{
				date: "01-01-2024",
			},
			want: true,
		},
		{
			name: "Invalid date",
			args: args{
				date: "2022-01-01",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateDateFormat(tt.args.date)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_validateDrumSize(t *testing.T) {
	type args struct {
		drumSize int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid drum size",
			args: args{
				drumSize: 250,
			},
			want: true,
		},
		{
			name: "Invalid drum size",
			args: args{
				drumSize: -1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateDrumSize(tt.args.drumSize)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDetermineBatchStatus(t *testing.T) {
	tests := []struct {
		name     string
		batch    Batch
		expected string
	}{
		{
			name: "Test BUFFER Status",
			batch: Batch{
				TotalQuantity: 10,
				DrumPartitions: []DrumPartition{
					{
						BufferQuantity: 5,
						TestQuantity:   2.5,
						ShortQuantity:  2.5,
					},
				},
			},
			expected: "BUFFER",
		},
		{
			name: "Test PARTIAL_BUFFER Status",
			batch: Batch{
				TotalQuantity: 10,
				DrumPartitions: []DrumPartition{
					{
						BufferQuantity:    4,
						TestQuantity:      2,
						ShortQuantity:     2,
						AvailableQuantity: 2,
					},
				},
			},
			expected: "PARTIAL_BUFFER",
		},
		{
			name: "Test AVAILABLE Status",
			batch: Batch{
				TotalQuantity: 10,
				DrumPartitions: []DrumPartition{
					{
						BufferQuantity:    0,
						TestQuantity:      0,
						ShortQuantity:     0,
						AvailableQuantity: 10,
					},
				},
			},
			expected: "AVAILABLE",
		},
		{
			name: "Test DOCS_PENDING_UPLOAD Status",
			batch: Batch{
				TotalQuantity: 10,
				DrumPartitions: []DrumPartition{
					{
						BufferQuantity:    0,
						TestQuantity:      0,
						ShortQuantity:     0,
						AvailableQuantity: 0,
					},
				},
			},
			expected: "DOCS_PENDING_UPLOAD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determineBatchStatus(tt.batch); got != tt.expected {
				t.Errorf("determineBatchStatus() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_collectApprovedDrumNumbers(t *testing.T) {
	type args struct {
		LI           LI
		materialCode string
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "returns approved drum numbers for given material code",
			args: args{
				LI: LI{
					MaterialCode: "material1",
					Batches: []Batch{
						{
							BatchTestApprovals: []BatchTestApproval{
								{
									ApprovalDrumNumbers: []ApprovalDrumNumber{
										{
											DrumNumbers: []int{1, 2, 3},
										},
									},
								},
							},
						},
						{
							BatchTestApprovals: []BatchTestApproval{
								{
									ApprovalDrumNumbers: []ApprovalDrumNumber{
										{
											DrumNumbers: []int{4, 5, 6},
										},
									},
								},
							},
						},
					},
				},
				materialCode: "material1",
			},

			want: [][]int{{1, 2, 3}, {4, 5, 6}},
		},
		{
			name: "returns empty slice when no drums are approved",
			args: args{
				LI: LI{
					MaterialCode: "material1",
					Batches: []Batch{
						{
							BatchTestApprovals: []BatchTestApproval{
								{
									ApprovalDrumNumbers: []ApprovalDrumNumber{
										{
											DrumNumbers: []int{},
										},
									},
								},
							},
						},
					},
				},
				materialCode: "material1",
			},
			want: [][]int{[]int{}},
		},
		{
			name: "returns empty slice when no line items match material code",
			args: args{
				LI: LI{
					MaterialCode: "material1",
					Batches: []Batch{
						{
							BatchTestApprovals: []BatchTestApproval{
								{
									ApprovalDrumNumbers: []ApprovalDrumNumber{
										{
											DrumNumbers: []int{1, 2, 3},
										},
									},
								},
							},
						},
						{
							BatchTestApprovals: []BatchTestApproval{
								{
									ApprovalDrumNumbers: []ApprovalDrumNumber{
										{
											DrumNumbers: []int{4, 5, 6},
										},
									},
								},
							},
						},
					},
				},
				materialCode: "material2",
			},
			want: [][]int(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectApprovedDrumNumbers(tt.args.LI, tt.args.materialCode)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_validateOverlappingDrumNumbers(t *testing.T) {
	type fields struct {
		Contracts []Contracts
	}
	tests := []struct {
		name   string
		fields fields
		want   []Error
	}{
		{
			name: "returns error when there are overlapping drum numbers",
			fields: fields{
				Contracts: []Contracts{
					{
						ContractNo: "contract1",
						LIs: []LI{
							{
								MaterialCode: "material1",
								Batches: []Batch{
									{
										BatchTestApprovals: []BatchTestApproval{
											{
												ApprovalDrumNumbers: []ApprovalDrumNumber{
													{
														DrumNumbers: []int{1, 2, 3},
													},
												},
											},
										},
									},
									{
										BatchTestApprovals: []BatchTestApproval{
											{
												ApprovalDrumNumbers: []ApprovalDrumNumber{
													{
														DrumNumbers: []int{4, 5},
													},
												},
											},
										},
									},
								},
							},
							{
								MaterialCode: "material1",
								Batches: []Batch{
									{
										BatchTestApprovals: []BatchTestApproval{
											{
												ApprovalDrumNumbers: []ApprovalDrumNumber{
													{
														DrumNumbers: []int{3},
													},
												},
											},
										},
									},
									{
										BatchTestApprovals: []BatchTestApproval{
											{
												ApprovalDrumNumbers: []ApprovalDrumNumber{
													{
														DrumNumbers: []int{4},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: []Error{
				{
					RowNo: 0,
					Err:   fmt.Errorf("overlapping drum numbers found for material code: material1, duplicates found: [3 4]"),
				},
			},
		},
		{
			name: "returns no error when there are no overlapping drum numbers",
			fields: fields{
				Contracts: []Contracts{
					{
						ContractNo: "contract1",
						LIs: []LI{
							{
								MaterialCode: "material1",
								Batches: []Batch{
									{
										BatchTestApprovals: []BatchTestApproval{
											{
												ApprovalDrumNumbers: []ApprovalDrumNumber{
													{
														DrumNumbers: []int{1, 2, 3},
													},
												},
											},
										},
									},
									{
										BatchTestApprovals: []BatchTestApproval{
											{
												ApprovalDrumNumbers: []ApprovalDrumNumber{
													{
														DrumNumbers: []int{4, 5},
													},
												},
											},
										},
									},
								},
							},
							{
								MaterialCode: "material1",
								Batches: []Batch{
									{
										BatchTestApprovals: []BatchTestApproval{
											{
												ApprovalDrumNumbers: []ApprovalDrumNumber{
													{
														DrumNumbers: []int{6},
													},
												},
											},
										},
									},
									{
										BatchTestApprovals: []BatchTestApproval{
											{
												ApprovalDrumNumbers: []ApprovalDrumNumber{
													{
														DrumNumbers: []int{7},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: []Error{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UploadInventoryInput{
				Contracts: tt.fields.Contracts,
			}
			got := u.validateOverlappingDrumNumbers()
			assert.Equal(t, tt.want, got)
		})
	}
}
