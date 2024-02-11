package model

type UploadInventoryInput struct {
	Contracts []Contracts `json:"contracts"`
}

type Contracts struct {
	ContractNo string `json:"contract_no"`
	LIs        []LI   `json:"lis"`
}

type LI struct {
	MaterialCode    string  `json:"material_code"`
	LiCode          string  `json:"li_code"`
	LiNumber        string  `json:"li_number"`
	Description     string  `json:"description"`
	Batches         []Batch `json:"batches"`
	HosApprovalDate string  `json:"hos_approval_date"`
	Status          string  `json:"status"`
}

type Batch struct {
	BatchNo            string              `json:"batch_no"`
	TotalQuantity      int                 `json:"total_quantity"`
	SubmissionDate     string              `json:"submission_date"`
	DrumPartitions     []DrumPartition     `json:"drum_partition"`
	BatchTestApprovals []BatchTestApproval `json:"batch_test_approvals"`
	Remarks            string              `json:"remarks"`
	Status             string              `json:"status"`
}

type DrumPartition struct {
	DrumSize             int           `json:"drum_size"`
	Quantity             int           `json:"quantity"`
	UnapprovedQuantity   int           `json:"unapproved_quantity"`
	AvailableQuantity    int           `json:"available_quantity"`
	AvailableDrumNumbers []int         `json:"available_drum_numbers"`
	BufferQuantity       int           `json:"buffer_quantity"`
	BufferDrumNumbers    []int         `json:"buffer_drum_numbers"`
	TestQuantity         float64       `json:"test_quantity"`
	TestDrumNumbers      []DrumDetails `json:"test_drum_numbers"`
	ShortQuantity        float64       `json:"short_quantity"`
	ShortDrumNumbers     []DrumDetails `json:"short_drum_numbers"`
}

type DrumDetails struct {
	DrumNumber int     `json:"number"`
	Quantity   float64 `json:"quantity"`
}

type BatchTestApproval struct {
	ApprovalDate        string                 `json:"approval_date"`
	TestDrumNumbers     []BatchTestDrumNumbers `json:"test_drum_numbers"`
	ApprovalDrumNumbers []ApprovalDrumNumber   `json:"approval_drum_numbers"`
	Status              string                 `json:"status"`
	ApprovalComment     string                 `json:"approval_comment"`
}

type BatchTestDrumNumbers struct {
	DrumSize    int           `json:"drum_size"`
	DrumNumbers []DrumDetails `json:"drum_numbers"`
}

type ApprovalDrumNumber struct {
	DrumSize    int   `json:"drum_size"`
	DrumNumbers []int `json:"drum_numbers"`
}

type CSVRow struct {
	Vendor                  string    `csv:"Vendor"`
	MaterialCode            string    `csv:"Material"`
	MaterialDesc            string    `csv:"Description"`
	ContractNo              string    `csv:"Contract"`
	PONumber                string    `csv:"PO Number"`
	POLineItem              string    `csv:"PO line item"`
	LIName                  LIName    `csv:"Li No"`
	LIDate                  string    `csv:"LI Date"`
	BatchNo                 string    `csv:"Batch No."`
	BatchDueDate            string    `csv:"Batch Due date"`
	DrumSize                int       `csv:"Drum Size"`
	TotalNoOfDrums          int       `csv:"Total nos. of Drum"`
	TotalQty                int       `csv:"-"`
	AvailableDrumNos        []int     `csv:"Available Drum Nos."`
	AvailableFullDrums      int       `csv:"Available Full Drums"`
	FullDrumTotalQuantity   int       `csv:"Full Drum Total Quantity"`
	BufferDrumNo            []int     `csv:"Buffer Drum No."`
	BufferNoOfDrums         int       `csv:"Buffer No. of Drum"`
	BufferQuantity          int       `csv:"Buffer Quantity"`
	SampleDrum              string    `csv:"Sample Drum (Yes/No)"`
	SampleDrumNo            []int     `csv:"Sample Drum No."`
	SampleLength            []float64 `csv:"Sample Length (m)"`
	NoOfShortLengthDrums    int       `csv:"No of Short length Drums"`
	ShortLengthTotalQty     float64   `csv:"Short Length total Quantity"`
	ApprovedDrumNumbers     []int     `csv:"-"`
	BatchTestReportDate     string    `csv:"Batch Test Report Date"`
	Remarks                 string    `csv:"Remarks"`
	BatchTestReportFileName string    `csv:"Batch Test Report File Name"`
}

type LIName struct {
	LICode   string
	LINumber string
}

type Error struct {
	RowNo int
	Err   error
}
