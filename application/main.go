package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

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

func main() {

	// Create a log file, overwrites if it exists
	logFile, err := os.OpenFile("VendorStockUpload.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err) //  Handle the error, potentially panic
	}
	defer logFile.Close() // Ensure the file is closed eventually

	logger := log.New(logFile, "[Stock Upload] Error: ", log.Lmsgprefix|log.LstdFlags)
	var errors []Error

	// Get the CSV file path from command-line arguments
	csvFilePath := "sample3.csv"

	// Open the CSV file
	file, _ := os.Open(csvFilePath)
	if err != nil {
		errors = append(errors, Error{RowNo: 0, Err: err})
	}
	defer file.Close()

	// Parse the CSV file
	records, errorSlice := parseCSV(file)
	errors = append(errors, errorSlice...)

	// marshal the records to JSON
	jsonData, err := recordsToJSON(records)
	if err != nil {
		errors = append(errors, Error{RowNo: 0, Err: err})
	}

	jsonFile, err := os.Create("output.json") // or os.OpenFile for more configuration
	if err != nil {
		errors = append(errors, Error{RowNo: 0, Err: err}) // Handle the error
	}
	defer jsonFile.Close()

	// Write the JSON data to the file
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		errors = append(errors, Error{RowNo: 0, Err: err}) // Handle the error
	}

	for _, e := range errors {
		logger.Printf("Row %d: %s", e.RowNo, e.Err)
	}
}

func (row *CSVRow) UnmarshalCSV(csv []string, rowIndex int) []Error {
	errors := make([]Error, 0)

	const (
		VendorColumnIndex                  = 0
		MaterialCodeColumnIndex            = 1
		MaterialDescColumnIndex            = 2
		ContractNoColumnIndex              = 3
		PONumberColumnIndex                = 4
		POLineItemColumnIndex              = 5
		LINameColumnIndex                  = 6
		LIDateColumnIndex                  = 7
		BatchNoColumnIndex                 = 8
		BatchDueDateColumnIndex            = 9
		DrumSizeColumnIndex                = 10
		TotalNoOfDrumsColumnIndex          = 11
		AvailableDrumNosColumnIndex        = 12
		AvailableFullDrumsColumnIndex      = 13
		FullDrumTotalQuantityColumnIndex   = 14
		BufferDrumNoColumnIndex            = 15
		BufferNoOfDrumsColumnIndex         = 16
		BufferQuantityColumnIndex          = 17
		SampleDrumColumnIndex              = 18
		SampleDrumNoColumnIndex            = 19
		SampleLengthColumnIndex            = 20
		NoOfShortLengthDrumsColumnIndex    = 21
		ShortLengthTotalQtyColumnIndex     = 22
		BatchTestReportDateColumnIndex     = 23
		RemarksColumnIndex                 = 24
		BatchTestReportFileNameColumnIndex = 25
	)

	// Parse Vendor
	row.Vendor = strings.TrimSpace(csv[VendorColumnIndex])

	// Parse MaterialCode
	row.MaterialCode = strings.TrimSpace(csv[MaterialCodeColumnIndex])

	// Parse MaterialDesc
	row.MaterialDesc = strings.TrimSpace(csv[MaterialDescColumnIndex])

	// Parse ContractNo
	row.ContractNo = strings.TrimSpace(csv[ContractNoColumnIndex])

	// Parse PONumber
	row.PONumber = strings.TrimSpace(csv[PONumberColumnIndex])

	// Parse POLineItem
	row.POLineItem = strings.TrimSpace(csv[POLineItemColumnIndex])

	// Parse LIName
	rawLIName := strings.TrimSpace(csv[LINameColumnIndex])
	liNameParts := strings.Split(rawLIName, "-")
	if len(liNameParts) != 2 {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("invalid LI Name format")})
	}

	if len(liNameParts) == 2 {
		row.LIName.LICode = strings.TrimSpace(liNameParts[0])
		row.LIName.LINumber = strings.TrimSpace(liNameParts[1])
	}
	// Parse LIDate
	row.LIDate = strings.TrimSpace(csv[LIDateColumnIndex])

	// Parse BatchNo
	row.BatchNo = strings.TrimSpace(csv[BatchNoColumnIndex])

	// Parse BatchDueDate
	row.BatchDueDate = strings.TrimSpace(csv[BatchDueDateColumnIndex])

	// Parse DrumSize
	rawDrumSize := strings.TrimSpace(csv[DrumSizeColumnIndex])
	drumSize, err := strconv.Atoi(rawDrumSize)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to parse drum size: %w", err)})
	}
	row.DrumSize = drumSize

	// Parse TotalNoOfDrums
	rawTotalNoOfDrums := strings.TrimSpace(csv[TotalNoOfDrumsColumnIndex])
	totalNoOfDrums, err := strconv.Atoi(rawTotalNoOfDrums)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to parse total no of drums: %w", err)})
	}
	row.TotalNoOfDrums = totalNoOfDrums

	// Parse TotalQty
	row.TotalQty = row.DrumSize * row.TotalNoOfDrums

	// Parse Available Drum Nos
	rawDrumNos := strings.TrimSpace(csv[AvailableDrumNosColumnIndex])
	drumNumbers, err := unpackDrumNoRange(rawDrumNos)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to parse available drum numbers: %w", err)})
	}
	row.AvailableDrumNos = drumNumbers

	// Parse Available Full Drums
	rawAvailableFullDrums := strings.TrimSpace(csv[AvailableFullDrumsColumnIndex])
	availableFullDrums, err := strconv.Atoi(rawAvailableFullDrums)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to parse available full drums: %w", err)})
	}
	row.AvailableFullDrums = availableFullDrums

	// Parse Full Drum Total Quantity
	rawFullDrumTotalQuantity := strings.TrimSpace(csv[FullDrumTotalQuantityColumnIndex])
	fullDrumTotalQuantity, err := strconv.Atoi(rawFullDrumTotalQuantity)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to parse full drum total quantity: %w", err)})
	}
	row.FullDrumTotalQuantity = fullDrumTotalQuantity

	// Parse Buffer Drum Nos
	rawBufferDrumNos := strings.TrimSpace(csv[BufferDrumNoColumnIndex])
	bufferDrumNumbers, err := unpackDrumNoRange(rawBufferDrumNos)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to parse buffer drum numbers: %w", err)})
	}
	row.BufferDrumNo = bufferDrumNumbers

	// Parse Buffer No Of Drums
	rawBufferNoOfDrums := strings.TrimSpace(csv[BufferNoOfDrumsColumnIndex])
	bufferNoOfDrums, err := strconv.Atoi(rawBufferNoOfDrums)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to parse buffer no of drums: %w", err)})
	}
	row.BufferNoOfDrums = bufferNoOfDrums

	// Parse Buffer Quantity
	rawBufferQuantity := strings.TrimSpace(csv[BufferQuantityColumnIndex])
	bufferQuantity, err := strconv.Atoi(rawBufferQuantity)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to parse buffer quantity: %w", err)})
	}
	row.BufferQuantity = bufferQuantity

	// Parse Sample Drum
	row.SampleDrum = strings.TrimSpace(csv[SampleDrumColumnIndex])

	// Parse Sample Drum Nos
	rawSampleDrumNos := strings.TrimSpace(csv[SampleDrumNoColumnIndex])
	sampleDrumNumbers, err := unpackDrumNoRange(rawSampleDrumNos)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to parse sample drum numbers: %w", err)})
	}
	row.SampleDrumNo = sampleDrumNumbers

	// Parse Sample Lengths
	rawSampleLengths := strings.TrimSpace(csv[SampleLengthColumnIndex])
	sampleLengths, err := stringToFloat64Slice(rawSampleLengths)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to parse sample lengths: %w", err)})
	}
	row.SampleLength = sampleLengths

	// Parse No Of Short Length Drums
	rawNoOfShortLengthDrums := strings.TrimSpace(csv[NoOfShortLengthDrumsColumnIndex])
	noOfShortLengthDrums, err := strconv.Atoi(rawNoOfShortLengthDrums)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to parse no of short length drums: %w", err)})
	}
	row.NoOfShortLengthDrums = noOfShortLengthDrums

	// Parse Short Length Total Qty
	rawShortLengthTotalQty := strings.TrimSpace(csv[ShortLengthTotalQtyColumnIndex])
	shortLengthTotalQty, err := strconv.ParseFloat(rawShortLengthTotalQty, 64)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to parse short length total quantity: %w", err)})
	}
	row.ShortLengthTotalQty = shortLengthTotalQty

	// Parse approved drum numbers
	approvedDrumNumbers, err := combineSortAndCheckDuplicates(row.AvailableDrumNos, row.BufferDrumNo, row.SampleDrumNo)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to combine, sort and check duplicates: %w", err)})
	}
	row.ApprovedDrumNumbers = approvedDrumNumbers

	// Parse Batch Test Report Date
	row.BatchTestReportDate = strings.TrimSpace(csv[BatchTestReportDateColumnIndex])

	// Parse Remarks
	row.Remarks = strings.TrimSpace(csv[RemarksColumnIndex])

	// Parse Batch Test Report File Name
	row.BatchTestReportFileName = strings.TrimSpace(csv[BatchTestReportFileNameColumnIndex])

	return errors
}

func (row *CSVRow) validateRow(rowIndex int) []Error {
	errors := make([]Error, 0)

	// Validate Vendor
	if row.Vendor == "" {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("vendor is required")})
	}

	// Validate MaterialCode
	if row.MaterialCode == "" {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("material code is required")})
	}

	// Validate MaterialDesc
	if row.MaterialDesc == "" {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("material description is required")})
	}

	// Validate ContractNo
	if row.ContractNo == "" {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("contract no. is required")})
	}

	// Validate LiName
	if row.LIName.LICode == "" || row.LIName.LINumber == "" {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("LI No. is required")})
	}

	// Validate LiDate
	if row.LIDate == "" || !validateDateFormat(row.LIDate) {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("invalid LI date format")})
	}

	// Validate BatchNo
	if row.BatchNo == "" || !validateBatchNoFormat(row.BatchNo) {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("invalid batch no. format")})
	}

	// Validate BatchDueDate
	if row.BatchDueDate == "" || !validateDateFormat(row.BatchDueDate) {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("invalid batch due date format")})
	}

	// Validate DrumSize
	if !validateDrumSize(row.DrumSize) {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("invalid drum size")})
	}

	// Validate TotalNoOfDrums
	if row.TotalNoOfDrums <= 0 {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("total no of drums must be greater than 0")})
	}

	// Validate TotalQty
	if row.TotalQty <= 0 {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("total quantity must be greater than 0")})
	}

	// Validate AvailableDrumNos
	if len(row.AvailableDrumNos) != row.AvailableFullDrums {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("available drum nos. does not match available full drums")})
	}

	// Validate FullDrumTotalQuantity
	if row.FullDrumTotalQuantity != row.AvailableFullDrums*row.DrumSize {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("full drum total quantity does not match available full drums")})
	}

	// Validate BufferDrumNo
	if len(row.BufferDrumNo) != row.BufferNoOfDrums {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("buffer drum nos. does not match buffer no. of drums")})
	}

	// Validate BufferQuantity
	if row.BufferQuantity != row.BufferNoOfDrums*row.DrumSize {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("buffer quantity does not match buffer no. of drums")})
	}

	// Validate SampleDrumNo
	if len(row.SampleDrumNo) != row.NoOfShortLengthDrums {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("sample drum nos. does not match no of short length drums")})
	}

	// Validate SampleLength
	if len(row.SampleLength) != len(row.SampleDrumNo) {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("sample length does not match sample drum nos")})
	}

	// Validate ShortLengthTotalQty
	if row.ShortLengthTotalQty != float64(row.DrumSize*row.NoOfShortLengthDrums)-sumFloat64Slice(row.SampleLength) {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("short length total qty does not match with sample length total qty - short length total qty")})
	}

	// Validate ApprovedDrumNumbers
	if len(row.ApprovedDrumNumbers) != row.AvailableFullDrums+row.BufferNoOfDrums+row.NoOfShortLengthDrums {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("approved drum numbers does not match with available drum numbers, buffer drum numbers, sample drum numbers")})
	}

	// Validate BatchTestReportDate
	if row.BatchTestReportDate != "" && !validateDateFormat(row.BatchTestReportDate) {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("invalid batch test report date format")})
	}

	// validate BatchTestReportFileName
	if row.BatchTestReportDate != "" {
		if row.BatchTestReportFileName == "" {
			errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("batch test report file name is required")})
		}
	}

	// Validate TotalNoOfDrums
	if row.BatchTestReportDate != "" && row.TotalNoOfDrums != row.AvailableFullDrums+row.BufferNoOfDrums+row.NoOfShortLengthDrums {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("total no of drums does not match with no of available drums, no of buffer drums, no of short drums")})
	}

	// Validate total quantity
	if row.BatchTestReportDate != "" && float64(row.TotalQty) != float64(row.FullDrumTotalQuantity)+float64(row.BufferQuantity)+row.ShortLengthTotalQty+sumFloat64Slice(row.SampleLength) {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("total quantity does not match with full drum total qty, buffer drum total qty, short length total qty, sample length total qty")})
	}

	return errors
}

func parseCSV(reader io.Reader) (UploadInventoryInput, []Error) {
	var errors []Error
	csvReader := csv.NewReader(reader)
	rows, err := csvReader.ReadAll()
	if err != nil {
		errors = append(errors, Error{RowNo: 0, Err: err})
	}
	fmt.Println("rows read:", len(rows)-1)

	var records []CSVRow
	for i, row := range rows[1:] { // Skip the header row
		var record CSVRow
		errorSlice := record.UnmarshalCSV(row, i)
		errors = append(errors, errorSlice...)
		errorSlice = record.validateRow(i)
		errors = append(errors, errorSlice...)
		records = append(records, record)
	}

	res, errorSlice := processRows(records)
	errors = append(errors, errorSlice...)

	// check for overlapping drum numbers
	errorSlice = res.validateOverlappingDrumNumbers()
	errors = append(errors, errorSlice...)

	return res, errors
}

func processRows(rows []CSVRow) (UploadInventoryInput, []Error) {
	var res UploadInventoryInput
	var errors []Error
	for rowIndex, row := range rows {

		if contractIndex := findContractIndex(res.Contracts, row.ContractNo); contractIndex == -1 {
			// Add a new contract to the res
			newContract := Contracts{
				ContractNo: row.ContractNo,
			}

			newLi, err := createNewLI(row, rowIndex)
			if err != nil {
				errors = append(errors, err...)
			}

			newContract.LIs = append(newContract.LIs, newLi)
			res.Contracts = append(res.Contracts, newContract)
		} else {
			// get the existing contract to update
			contract := res.Contracts[contractIndex]

			// Check if LI exists, if not, add a new LI to the existing contract, else add new batch to existing LI
			if liIndex := findLiIndex(contract.LIs, row.LIName); liIndex == -1 {
				// Add a new LI to the existing contract
				newLi, err := createNewLI(row, rowIndex)
				if err != nil {
					errors = append(errors, err...)
				}

				contract.LIs = append(contract.LIs, newLi)
			} else {
				// get the existing LI to update
				li := contract.LIs[liIndex]

				// verify the hos approval date
				if li.HosApprovalDate != row.LIDate {
					errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("hos approval date does not match")})
				}

				// verify material
				if li.MaterialCode != row.MaterialCode {
					errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("material code does not match")})
				}

				// verify material description
				if li.Description != row.MaterialDesc {
					errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("material description does not match")})
				}

				// Check if batch exists, if not, add a new batch to the existing LI
				if batchIndex := findBatchIndex(li.Batches, row.BatchNo); batchIndex == -1 {
					newBatch, err := createNewBatch(row, rowIndex)
					if err != nil {
						errors = append(errors, err...)
					}
					li.Batches = append(li.Batches, newBatch)
				} else { // update the existing batch

					// Testcases:
					// 1. same drum size, no batch test approval date
					// 2. same drum size, same batch test approval date
					// 3. same drum size, different batch test approval date
					// 4. different drum size, same batch test approval date
					// 5. different drum size, different batch test approval date
					// 6. different drum size, no batch test approval date

					// get the existing batch to update
					batch := li.Batches[batchIndex]

					// verify batch due date
					if batch.SubmissionDate != row.BatchDueDate {
						errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("batch due date does not match")})
					}

					// check if current row is new drum size, if yes, add a new drum partition to the existing batch
					var dpMap = make(map[int]DrumPartition)
					for _, dp := range batch.DrumPartitions {
						dpMap[dp.DrumSize] = dp
					}
					dp, drumSizeExists := dpMap[row.DrumSize]

					// check if batch test approval date is duplicated
					var btaMap = make(map[string]BatchTestApproval)
					for _, bta := range batch.BatchTestApprovals {
						btaMap[bta.ApprovalDate] = bta
					}
					bta, batchTestApprovalExists := btaMap[row.BatchTestReportDate]

					switch {
					case drumSizeExists:

						// 1. same drum size, no batch test approval date
						//		- update the existing drum partition
						//		- skip creating a new batch test approval
						// 2. same drum size, same batch test approval date
						//		- update the existing drum partition
						//		- update the existing batch test approval
						// 3. same drum size, different batch test approval date
						//		- update the existing drum partition
						//		- create new batch test approval and append to newBatch

						// update the existing drum partition
						updatedDrumPartition, err := updateDrumPartition(dp, row, rowIndex)
						if err != nil {
							errors = append(errors, err...)
						}

						// update the existing drum partition
						dpMap[row.DrumSize] = updatedDrumPartition

						switch {
						case row.BatchTestReportDate == "": // Case1: same drum size, no batch test approval date

							// skip creating a new batch test approval

						case batchTestApprovalExists: // Case2: same drum size, same batch test approval date

							// update the existing batch test approval's test drum numbers
							for testDrumNumberIndex, testDrumNumber := range bta.TestDrumNumbers {
								if testDrumNumber.DrumSize == row.DrumSize {
									if len(row.SampleDrumNo) > 0 {
										for i, drumNo := range row.SampleDrumNo {
											testDrumNumber.DrumNumbers = append(testDrumNumber.DrumNumbers, DrumDetails{
												DrumNumber: drumNo,
												Quantity:   row.SampleLength[i],
											})
										}
									}
									bta.TestDrumNumbers[testDrumNumberIndex] = testDrumNumber
									break
								}
							}

							// update the existing batch test approval's approval drum numbers
							for approvalDrumNumberIndex, approvalDrumNumber := range bta.ApprovalDrumNumbers {
								if approvalDrumNumber.DrumSize == row.DrumSize {
									drumNumbers, err2 := combineSortAndCheckDuplicates(approvalDrumNumber.DrumNumbers, row.ApprovedDrumNumbers)
									if err2 != nil {
										errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to combine, sort and check drum no. duplicates : %s", err2)})
									}
									approvalDrumNumber.DrumNumbers = drumNumbers
									bta.ApprovalDrumNumbers[approvalDrumNumberIndex] = approvalDrumNumber
									break
								}
							}

							btaMap[row.BatchTestReportDate] = bta

						case !batchTestApprovalExists: // Case3: same drum size, different batch test approval date

							// create new batch test approval
							newBatchTestApproval := createBatchTestApproval(row)

							// add to btaMap
							btaMap[row.BatchTestReportDate] = newBatchTestApproval
						}

					case !drumSizeExists:
						// 4. different drum size, no batch test approval date
						//		- create new drum partition and append to the existing batch
						//		- skip creating a new batch test approval
						// 5. different drum size, same batch test approval date
						//		- create new drum partition and append to the existing batch
						//		- update the existing batch test approval
						// 6. different drum size, different batch test approval date
						//		- create new drum partition and append to the existing batch
						//		- create new batch test approval and append to the existing batch

						newDrumPartition, err := createDrumPartition(row, rowIndex)
						if err != nil {
							errors = append(errors, err...)
						}

						// add to dpMap
						dpMap[row.DrumSize] = newDrumPartition

						switch {
						case row.BatchTestReportDate == "": // Case4: different drum size, no batch test approval date

							// skip creating a new batch test approval

						case batchTestApprovalExists: // Case5: different drum size, same batch test approval date

							// update the existing batch test approval's test drum numbers
							testDrumNumberDetails, _, _, _ := unpackSampleDrumNos(row.SampleDrumNo, row.SampleLength, row.DrumSize)

							bta.TestDrumNumbers = append(bta.TestDrumNumbers, BatchTestDrumNumbers{
								DrumSize:    row.DrumSize,
								DrumNumbers: testDrumNumberDetails,
							})

							// update the existing batch test approval's approval drum numbers
							bta.ApprovalDrumNumbers = append(bta.ApprovalDrumNumbers, ApprovalDrumNumber{
								DrumSize:    row.DrumSize,
								DrumNumbers: row.ApprovedDrumNumbers,
							})

							btaMap[row.BatchTestReportDate] = bta

						case !batchTestApprovalExists: // Case6: different drum size, different batch test approval date

							// create new batch test approval
							newBatchTestApproval := createBatchTestApproval(row)

							// add to btaMap
							btaMap[row.BatchTestReportDate] = newBatchTestApproval
						}

					}

					// update the existing batch with the new drum partition map
					var updatedDrumPartitions []DrumPartition
					for _, drumPartition := range dpMap {
						updatedDrumPartitions = append(updatedDrumPartitions, drumPartition)
						// sort drum partitions by drum size
						sort.Slice(updatedDrumPartitions, func(i, j int) bool {
							return updatedDrumPartitions[i].DrumSize < updatedDrumPartitions[j].DrumSize
						})
					}
					batch.DrumPartitions = updatedDrumPartitions

					// update the existing batch with the new batch test approval map
					var updatedBatchTestApprovals []BatchTestApproval
					for _, batchTestApproval := range btaMap {
						updatedBatchTestApprovals = append(updatedBatchTestApprovals, batchTestApproval)
						// sort batch test approvals by approval date
						sort.Slice(updatedBatchTestApprovals, func(i, j int) bool {
							return updatedBatchTestApprovals[i].ApprovalDate < updatedBatchTestApprovals[j].ApprovalDate
						})
					}
					batch.BatchTestApprovals = updatedBatchTestApprovals

					// update batch total qty
					batch.TotalQuantity += row.TotalQty

					// update batch status
					batch.Status = determineBatchStatus(batch)

					if batch.BatchTestApprovals == nil {
						batch.BatchTestApprovals = []BatchTestApproval{}
					}

					li.Batches[batchIndex] = batch
				}

				// update the existing LI
				contract.LIs[liIndex] = li
			}

			res.Contracts[contractIndex] = contract
		}
	}
	return res, errors
}

func determineBatchStatus(batch Batch) string {
	var totalBatchBufferQty int
	var totalBatchTestQty float64
	var totalBatchShortQty float64
	var totalAvailableQty int

	for _, drumPartition := range batch.DrumPartitions {
		totalBatchBufferQty += drumPartition.BufferQuantity
		totalBatchTestQty += drumPartition.TestQuantity
		totalBatchShortQty += drumPartition.ShortQuantity
		totalAvailableQty += drumPartition.AvailableQuantity
	}

	if float64(batch.TotalQuantity) == float64(totalBatchBufferQty)+totalBatchTestQty+totalBatchShortQty {
		return "BUFFER"
	}

	if totalBatchBufferQty > 0 && totalBatchBufferQty < batch.TotalQuantity && totalAvailableQty > 0 {
		return "PARTIAL_BUFFER"
	}

	if totalAvailableQty > 0 {
		return "AVAILABLE"
	}

	return "DOCS_PENDING_UPLOAD"
}

func updateDrumPartition(dp DrumPartition, row CSVRow, rowIndex int) (DrumPartition, []Error) {
	// Drum size partition exists, update the existing drum partition
	var errors []Error
	var err error

	dp.Quantity += row.TotalQty

	dp.AvailableQuantity += row.FullDrumTotalQuantity
	dp.AvailableDrumNumbers, err = combineSortAndCheckDuplicates(dp.AvailableDrumNumbers, row.AvailableDrumNos)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to combine, sort and check duplicates: %s", err)})
	}
	dp.BufferQuantity += row.BufferQuantity
	dp.BufferDrumNumbers, err = combineSortAndCheckDuplicates(dp.BufferDrumNumbers, row.BufferDrumNo)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to combine, sort and check duplicates: %s", err)})
	}

	TestDrumNumbers, TestQuantity, ShortDrumNumbers, ShortQuantity := unpackSampleDrumNos(row.SampleDrumNo, row.SampleLength, dp.DrumSize)
	if err != nil {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("failed to unpack sample drum nos: %s", err)})
	}

	dp.TestQuantity += TestQuantity
	dp.TestDrumNumbers = append(dp.TestDrumNumbers, TestDrumNumbers...)
	dp.ShortQuantity += ShortQuantity
	dp.ShortDrumNumbers = append(dp.ShortDrumNumbers, ShortDrumNumbers...)

	// recalculate unapproved quantity after updating the drum partition
	unapprovedQty := float64(dp.Quantity) - float64(dp.AvailableQuantity) - float64(dp.BufferQuantity) - dp.TestQuantity - dp.ShortQuantity
	dp.UnapprovedQuantity = int(unapprovedQty)

	// validate the updated drum partition
	if float64(dp.Quantity) != float64(dp.AvailableQuantity)+float64(dp.BufferQuantity)+dp.TestQuantity+dp.ShortQuantity+float64(dp.UnapprovedQuantity) {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("total quantity does not match with available quantity, buffer quantity, test quantity, short quantity, unapproved quantity")})
	}

	return dp, errors
}

func createNewLI(row CSVRow, rowIndex int) (LI, []Error) {
	// Add a new LI to the res
	newLi := LI{
		MaterialCode:    row.MaterialCode,
		LiCode:          row.LIName.LICode,
		LiNumber:        row.LIName.LINumber,
		Description:     row.MaterialDesc,
		HosApprovalDate: row.LIDate,
		Status:          "VENDOR_ACKNOWLEDGED",
	}

	newBatch, errs := createNewBatch(row, rowIndex)
	newLi.Batches = append(newLi.Batches, newBatch)

	return newLi, errs
}

func createNewBatch(row CSVRow, rowIndex int) (Batch, []Error) {
	var errors []Error
	// Add the batch to the new LI
	newBatch := Batch{
		BatchNo:        row.BatchNo,
		SubmissionDate: row.BatchDueDate,
		Remarks:        row.Remarks,
	}

	newBatch.TotalQuantity = row.TotalQty

	// add row's drum size to the new li and batch
	newDrumPartition, err := createDrumPartition(row, rowIndex)
	if err != nil {
		errors = append(errors, err...)
	}

	// if batch test approval date is empty, skip creating a new batch test approval and set the batch test status to "DOCS_PENDING_UPLOAD"
	if row.BatchTestReportDate == "" {
		newBatch.BatchTestApprovals = []BatchTestApproval{}
		newBatch.Status = "DOCS_PENDING_UPLOAD"
		newBatch.DrumPartitions = append(newBatch.DrumPartitions, newDrumPartition)
		return newBatch, errors
	}

	// create new batch test approval and append to newBatch
	newBatchTestApproval := createBatchTestApproval(row)

	// determine batch status
	if newDrumPartition.BufferQuantity == 0 && newDrumPartition.AvailableQuantity > 0 {
		newBatch.Status = "AVAILABLE"
	} else if float64(newBatch.TotalQuantity)-(float64(newDrumPartition.BufferQuantity)+newDrumPartition.TestQuantity+newDrumPartition.ShortQuantity) == 0 {
		newBatch.Status = "BUFFER"
	} else {
		newBatch.Status = "PARTIAL_BUFFER"
	}
	if newDrumPartition.UnapprovedQuantity == newBatch.TotalQuantity {
		newBatch.Status = "DOCS_PENDING_UPLOAD"
	}

	newBatch.BatchTestApprovals = append(newBatch.BatchTestApprovals, newBatchTestApproval)
	newBatch.DrumPartitions = append(newBatch.DrumPartitions, newDrumPartition)
	return newBatch, errors
}

func createBatchTestApproval(row CSVRow) BatchTestApproval {

	res := BatchTestApproval{
		ApprovalDate:    row.BatchTestReportDate,
		Status:          "APPROVED",
		ApprovalComment: "Batch Test report uploaded on Go-live phase 1",
	}

	newBatchTestDrumNumbers := BatchTestDrumNumbers{
		DrumSize: row.DrumSize,
	}

	if len(row.SampleDrumNo) > 0 {
		for i, drumNo := range row.SampleDrumNo {
			newBatchTestDrumNumbers.DrumNumbers = append(newBatchTestDrumNumbers.DrumNumbers, DrumDetails{
				DrumNumber: drumNo,
				Quantity:   row.SampleLength[i],
			})
		}
	}

	res.TestDrumNumbers = append(res.TestDrumNumbers, newBatchTestDrumNumbers)

	newApprovalDrumNumbers := ApprovalDrumNumber{
		DrumSize:    row.DrumSize,
		DrumNumbers: row.ApprovedDrumNumbers,
	}

	if len(row.ApprovedDrumNumbers) == 0 {
		newApprovalDrumNumbers.DrumNumbers = []int{}
	}

	res.ApprovalDrumNumbers = append(res.ApprovalDrumNumbers, newApprovalDrumNumbers)

	return res

}

func createDrumPartition(row CSVRow, rowIndex int) (DrumPartition, []Error) {
	var res DrumPartition
	var errors []Error

	res.DrumSize = row.DrumSize
	res.Quantity = row.TotalQty
	res.AvailableDrumNumbers = row.ApprovedDrumNumbers
	if len(row.ApprovedDrumNumbers) == 0 {
		res.AvailableDrumNumbers = []int{}
	}
	res.BufferDrumNumbers = row.BufferDrumNo
	if len(row.BufferDrumNo) == 0 {
		res.BufferDrumNumbers = []int{}
	}
	res.BufferQuantity = res.DrumSize * len(row.BufferDrumNo)

	// unpack sample drum numbers and sample length into test drum numbers, test quantity, short drum numbers and short quantity
	res.TestDrumNumbers, res.TestQuantity, res.ShortDrumNumbers, res.ShortQuantity = unpackSampleDrumNos(row.SampleDrumNo, row.SampleLength, row.DrumSize)

	res.AvailableDrumNumbers = removeDuplicateDrumNumbers(res.AvailableDrumNumbers, row.SampleDrumNo)
	res.AvailableDrumNumbers = removeDuplicateDrumNumbers(res.AvailableDrumNumbers, row.BufferDrumNo)
	res.AvailableQuantity = res.DrumSize * len(res.AvailableDrumNumbers)
	if len(res.AvailableDrumNumbers) == 0 {
		res.AvailableDrumNumbers = []int{}
	}
	UnapprovedQty := float64(res.Quantity) - float64(res.AvailableQuantity) - float64(res.BufferQuantity) - res.TestQuantity - res.ShortQuantity
	res.UnapprovedQuantity = int(UnapprovedQty)

	if float64(res.Quantity)-float64(res.UnapprovedQuantity)-float64(res.AvailableQuantity)-float64(res.BufferQuantity)-res.TestQuantity-res.ShortQuantity != 0 {
		errors = append(errors, Error{RowNo: rowIndex + 1, Err: fmt.Errorf("drum partition total quantity does not match sum of available quantity, buffer quantity, test quantity, short quantity, unapproved quantity")})
	}
	return res, errors
}

func unpackSampleDrumNos(sampleDrumNumbers []int, sampleLength []float64, drumSize int) ([]DrumDetails, float64, []DrumDetails, float64) {

	testDrumNumbers := make([]DrumDetails, 0)
	var testQuantity float64
	shortDrumNumbers := make([]DrumDetails, 0)
	var shortQuantity float64

	if len(sampleDrumNumbers) > 0 {
		for i, drumNo := range sampleDrumNumbers {
			testDrumNumbers = append(testDrumNumbers, DrumDetails{
				DrumNumber: drumNo,
				Quantity:   sampleLength[i],
			})
			testQuantity += sampleLength[i]

			// unpack sample drum numbers and sample length into short drum numbers and short quantity
			shortDrumNumbers = append(shortDrumNumbers, DrumDetails{
				DrumNumber: drumNo,
				Quantity:   float64(drumSize) - sampleLength[i],
			})
			shortQuantity += float64(drumSize) - sampleLength[i]
		}
	}
	return testDrumNumbers, testQuantity, shortDrumNumbers, shortQuantity
}

func findContractIndex(slice []Contracts, contractNo string) int {
	for i, value := range slice {
		if value.ContractNo == contractNo {
			return i
		}
	}
	return -1
}

func findLiIndex(lis []LI, liNo LIName) int {
	liName := liNo.LICode + liNo.LINumber
	for i, li := range lis {
		if li.LiCode+li.LiNumber == liName {
			return i
		}
	}
	return -1
}

func findBatchIndex(batches []Batch, no string) int {
	for i, value := range batches {
		if value.BatchNo == no {
			return i
		}
	}
	return -1
}

func findBatchTestApprovalIndex(bta []BatchTestApproval, date string) int {
	for i, value := range bta {
		if value.ApprovalDate == date {
			return i
		}
	}
	return -1
}

func unpackDrumNoRange(str string) ([]int, error) {
	result := make([]int, 0)

	if strings.TrimSpace(str) == "" {
		return result, nil
	}

	// Split the string by comma
	parts := strings.Split(str, ",")

	for _, part := range parts {
		// Check if the part contains a dash
		if strings.Contains(part, "-") {
			// Split the part by dash
			rangeParts := strings.Split(part, "-")
			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, err
			}

			// Generate a range of integers and add them to the slice
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else {
			// Convert the part to an integer and add it to the slice
			num, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				return nil, err
			}
			result = append(result, num)
		}
	}

	return result, nil
}

func combineSortAndCheckDuplicates(slices ...[]int) ([]int, error) {
	combined := make([]int, 0)
	for _, slice := range slices {
		combined = append(combined, slice...)
	}
	sort.Ints(combined)

	m := make(map[int]bool)
	var duplicates []int
	for _, item := range combined {
		if m[item] {
			duplicates = append(duplicates, item)
		} else {
			m[item] = true
		}
	}
	if len(duplicates) > 0 {
		return combined, fmt.Errorf("duplicates found: %v", duplicates)
	}
	return combined, nil
}

func removeDuplicateDrumNumbers(slice []int, dupSlice []int) []int {
	dupMap := make(map[int]bool)
	for _, num := range dupSlice {
		dupMap[num] = true
	}

	result := make([]int, 0)
	for _, num := range slice {
		if !dupMap[num] {
			result = append(result, num)
		}
	}

	return result
}

func stringToFloat64Slice(str string) ([]float64, error) {
	// If the string is empty or contains only spaces, return an empty slice and nil error
	if strings.TrimSpace(str) == "0" {
		return []float64{}, nil
	}

	parts := strings.Split(str, ",")
	var result []float64
	for _, part := range parts {
		num, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}
	return result, nil
}

func sumFloat64Slice(length []float64) float64 {
	var sum float64
	for _, l := range length {
		sum += l
	}
	return sum
}

// recordsToJSON converts a slice of records to JSON format
func recordsToJSON(records UploadInventoryInput) ([]byte, error) {
	// Marshal the records to JSON
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal records to JSON: %w", err)
	}

	return jsonData, nil
}

func validateDateFormat(date string) bool {
	_, err := time.Parse("02-01-2006", date)
	return err == nil
}

func validateBatchNoFormat(batchNo string) bool {
	re := regexp.MustCompile(`^\d{1,2}/\d{1,2}$`)
	return re.MatchString(batchNo)
}

func validateDrumSize(drumSize int) bool {
	validSizes := []int{250, 300, 500, 1000}

	// Iterate through valid sizes
	for _, size := range validSizes {
		if drumSize == size {
			return true
		}
	}

	return false
}

func (u UploadInventoryInput) validateOverlappingDrumNumbers() []Error {

	errors := make([]Error, 0)
	for _, contract := range u.Contracts {

		matCodeMap := make(map[string][][]int)

		for _, li := range contract.LIs {
			if _, ok := matCodeMap[li.MaterialCode]; !ok {
				matCodeMap[li.MaterialCode] = collectApprovedDrumNumbers(li, li.MaterialCode)
			} else {
				matCodeMap[li.MaterialCode] = append(matCodeMap[li.MaterialCode], collectApprovedDrumNumbers(li, li.MaterialCode)...)
			}
		}

		for matCode, collatedDrumNos := range matCodeMap {

			_, err := combineSortAndCheckDuplicates(collatedDrumNos...)
			if err != nil {
				errors = append(errors, Error{RowNo: 0, Err: fmt.Errorf("overlapping drum numbers found for material code: %s, %v", matCode, err)})
			}
		}

	}
	return errors
}

func collectApprovedDrumNumbers(li LI, materialCode string) [][]int {
	var approvedDrumNumbers [][]int

	if li.MaterialCode == materialCode {
		for _, batch := range li.Batches {
			for _, batchTestApproval := range batch.BatchTestApprovals {
				for _, approvalDrumNumber := range batchTestApproval.ApprovalDrumNumbers {
					approvedDrumNumbers = append(approvedDrumNumbers, approvalDrumNumber.DrumNumbers)
				}
			}
		}
	}

	return approvedDrumNumbers
}
