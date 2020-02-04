package files

import (
	"fmt"
	"github.com/shopspring/decimal"
	"path"

	"github.com/ggoop/goutils/configs"
	"github.com/ggoop/goutils/md"
	"github.com/ggoop/goutils/utils"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/ggoop/goutils/glog"
)

type FileData struct {
	FileName string
	Dir      string
	FullPath string
}
type ImportData struct {
	Name    string
	Columns map[string]string
	Datas   []map[string]interface{}
}

type ExcelColumn struct {
	Name      string
	Title     string
	Hidden    bool
	excelName string
}
type ExcelCell struct {
	Name  string
	Value interface{}
}
type ToExcel struct {
	Columns []ExcelColumn
	Datas   [][]ExcelCell
}
type ExcelSv struct {
}

/**
* 创建服务实例
 */
func NewExcelSv() *ExcelSv {
	return &ExcelSv{}
}
func GetMapStringValue(key string, row map[string]interface{}) string {
	v := GetMapValue(key, row)
	if v == nil {
		return ""
	}
	if v, ok := v.(string); ok {
		return v
	}
	return ""
}
func GetMapIntValue(key string, row map[string]interface{}) int {
	v := GetMapValue(key, row)
	if v == nil {
		return 0
	}
	return utils.ToInt(v)
}

func GetMapSBoolValue(key string, row map[string]interface{}) md.SBool {
	return md.SBool_Parse(GetMapValue(key, row))
}
func GetMapTimeValue(key string, row map[string]interface{}) *md.Time {
	v := GetMapValue(key, row)
	if v == nil {
		return nil
	}
	if vv, ok := v.(string); ok {
		return md.CreateTimePtr(vv)
	}
	return nil
}
func GetMapBoolValue(key string, row map[string]interface{}, defaultValue bool) bool {
	v := GetMapValue(key, row)
	if v == nil {
		return defaultValue
	}
	if vv, ok := v.(string); ok {
		if vv == "true" || vv == "1" || vv == "T" {
			return true
		} else {
			return false
		}
	} else if vv, ok := v.(int); ok {
		if vv > 0 {
			return true
		} else {
			return false
		}
	}
	return utils.ToBool(v)
}
func GetMapDecimalValue(key string, row map[string]interface{}) decimal.Decimal {
	v := GetMapValue(key, row)
	if v == nil {
		return decimal.Zero
	}
	if vv, ok := v.(decimal.Decimal); ok {
		return vv
	} else if vv, ok := v.(string); ok {
		rv, _ := decimal.NewFromString(vv)
		return rv
	} else if vv, ok := v.(float64); ok {
		return decimal.NewFromFloat(vv)
	} else if vv, ok := v.(float32); ok {
		return decimal.NewFromFloat32(vv)
	}
	return decimal.Zero
}
func GetMapValue(key string, row map[string]interface{}) interface{} {
	if v, ok := row[key]; ok {
		return v
	}
	if v, ok := row[utils.SnakeString(key)]; ok {
		return v
	}
	return nil
}
func (s *ExcelSv) GetExcelDatas(filePath string, sheetNames ...string) ([]ImportData, error) {
	xlsx, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	rtnDatas := make([]ImportData, 0)
	if len(sheetNames) == 0 || sheetNames[0] == "*" {
		for _, sheetName := range xlsx.GetSheetMap() {
			if data, err := s.getSheetData(xlsx, sheetName); err != nil {
				return nil, err
			} else if len(data.Columns) > 0 && len(data.Datas) > 0 {
				rtnDatas = append(rtnDatas, data)
			}
		}
	} else {
		for _, sheetName := range sheetNames {
			if data, err := s.getSheetData(xlsx, sheetName); err != nil {
				return nil, err
			} else if len(data.Columns) > 0 && len(data.Datas) > 0 {
				rtnDatas = append(rtnDatas, data)
			}
		}
	}
	return rtnDatas, nil
}
func (s *ExcelSv) GetExcelData(filePath string) (data ImportData, err error) {
	xlsx, err := excelize.OpenFile(filePath)
	if err != nil {
		return data, err
	}
	return s.getSheetData(xlsx, xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))
}
func (s *ExcelSv) getSheetData(xlsx *excelize.File, sheetName string) (ImportData, error) {
	importData := ImportData{Columns: make(map[string]string), Name: sheetName}
	dataFromRow := 2
	//第一行为字段标识
	//条二行为行名称
	// 获取 Sheet1 上所有单元格,模板，需要预制一个列标识，_state,为空后，后边的行将都不会导入
	if sheetName == "" {
		return importData, nil
	}
	allRows := xlsx.GetRows(sheetName)
	if len(allRows) <= dataFromRow {
		return importData, nil
	}
	//取列数
	colsMap := make(map[int]string)
	cols := allRows[0]
	titles := allRows[1]
	for c, name := range cols {
		if name != "" {
			colsMap[c] = name
			importData.Columns[name] = titles[c]
		}
	}
	if len(colsMap) == 0 {
		return importData, nil
	}
	dataRows := allRows[dataFromRow:]

	datas := make([]map[string]interface{}, 0)
	isData := false
	for i, row := range dataRows {
		isData = false
		values := make(map[string]interface{}, 0)
		for c, value := range row {
			if cName, ok := colsMap[c]; ok {
				if c == 0 && value != "" {
					isData = true
				}
				//处理合并单元格时，取出来的空值时
				if isData && c > 0 && value == "" {
					value = xlsx.GetCellValue(sheetName, fmt.Sprintf("%v%v", excelize.ToAlphaString(c), i+1+dataFromRow))
				}
				values[cName] = value
			}
		}
		if isData {
			datas = append(datas, values)
		} else {
			break
		}
	}
	importData.Datas = datas
	return importData, nil
}

func (s *ExcelSv) ToExcel(data *ToExcel) (*FileData, error) {

	xlsx := excelize.NewFile()
	sheetName := xlsx.GetSheetName(xlsx.GetActiveSheetIndex())
	colMap := make(map[string]ExcelColumn)
	//增加系统默认导出列
	columns := make([]ExcelColumn, 0)
	for _, c := range data.Columns {
		columns = append(columns, c)
	}
	startIndex := 2
	for i, c := range columns {
		cName := excelize.ToAlphaString(i)
		columns[i].excelName = cName
		colMap[c.Name] = columns[i]

		xlsx.SetCellValue(sheetName, fmt.Sprintf("%s%d", cName, 1), c.Name)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("%s%d", cName, 2), c.Title)
	}
	xlsx.SetRowVisible(sheetName, 0, false)
	//设置数据列宽度
	xlsx.SetColWidth("Sheet1", "A", columns[len(columns)-1].excelName, 20)
	//border
	if style, err := xlsx.NewStyle(`{"border":[{"type":"left","color":"666666","style":1},{"type":"top","color":"666666","style":1},{"type":"bottom","color":"666666","style":1},{"type":"right","color":"666666","style":1}]}`); err == nil {
		xlsx.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s%d", columns[len(columns)-1].excelName, len(data.Datas)+startIndex), style)
	}
	//header
	if style, err := xlsx.NewStyle(`{"font":{"bold":true},"fill":{"type":"pattern","pattern":1,"color":["#dddddd"]},"border":[{"type":"left","color":"666666","style":1},{"type":"top","color":"666666","style":1},{"type":"bottom","color":"666666","style":1},{"type":"right","color":"666666","style":1}]}`); err == nil {
		xlsx.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s%d", columns[len(columns)-1].excelName, startIndex), style)
	}
	if style, err := xlsx.NewStyle(`{"font":{"bold":true},"fill":{"type":"pattern","pattern":1,"color":["#dddddd"]},"border":[{"type":"left","color":"666666","style":1},{"type":"top","color":"666666","style":1},{"type":"bottom","color":"666666","style":1},{"type":"right","color":"666666","style":1}]}`); err == nil {
		xlsx.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s%d", "A", len(data.Datas)+startIndex), style)
	}

	for r, row := range data.Datas {
		for _, cell := range row {
			if c, ok := colMap[cell.Name]; ok {
				xlsx.SetCellValue(sheetName, fmt.Sprintf("%s%d", c.excelName, r+startIndex+1), cell.Value)
			}
		}
	}
	fileData := FileData{}
	fileData.FileName = fmt.Sprintf("%s.%s", utils.GUID(), "xlsx")
	fileData.Dir = path.Join(configs.Default.App.Storage, "export", md.NewTime().Format("200601"))
	utils.CreatePath(fileData.Dir)
	fileData.FullPath = utils.JoinCurrentPath(path.Join(fileData.Dir, fileData.FileName))
	if err := xlsx.SaveAs(fileData.FullPath); err != nil {
		glog.Error(err)
		return nil, err
	}
	return &fileData, nil
}
