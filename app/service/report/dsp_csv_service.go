package report

import (
	"encoding/csv"
	"github.com/google/wire"
	"os"
	"strconv"
	"td_report/pkg/logger"
)

var DspCsvServiceSet = wire.NewSet(wire.Struct(new(DspCsvService), "*"))

type DspCsvService struct {
}

func NewDspCsvService() *DspCsvService {
	return &DspCsvService{}
}

type CsvTable struct {
	FileName string
	Records  []CsvRecord
}

type CsvRecord struct {
	Record map[string]string
}

func (c *CsvRecord) GetInt(field string) int {
	var r int
	var err error
	if r, err = strconv.Atoi(c.Record[field]); err != nil {
		logger.Error(err)
		panic(err)
	}

	return r
}

func (c *CsvRecord) GetString(field string) string {
	data, ok := c.Record[field]
	if ok {
		return data
	} else {
		logger.Info("Get fileld failed! fileld:", field)
		return ""
	}
}

func (t *DspCsvService) LoadCsvCfg(filename string, row int) *CsvTable {
	file, err := os.Open(filename)
	if err != nil {
		logger.Error(err)
		return nil
	}

	defer file.Close()

	reader := csv.NewReader(file)
	if reader == nil {
		logger.Error("NewReader return nil, file:", file)
		return nil
	}
	records, err := reader.ReadAll()
	if err != nil {
		logger.Error(err)
		return nil
	}
	if len(records) < row {
		logger.Info(filename, " is empty")
		return nil
	}
	colNum := len(records[0])
	recordNum := len(records)
	var allRecords []CsvRecord
	for i := row; i < recordNum; i++ {
		record := &CsvRecord{make(map[string]string)}
		for k := 0; k < colNum; k++ {
			record.Record[records[0][k]] = records[i][k]
		}
		allRecords = append(allRecords, *record)
	}
	var result = &CsvTable{
		filename,
		allRecords,
	}
	return result
}
