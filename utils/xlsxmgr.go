package utils

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var XlsxMgr xlsxManager

func init() {
	XlsxMgr = xlsxManager{}
	XlsxMgr.sheetMap = make(map[string]*Table)
}

type dataRow struct {
	datas []string
}

func (s* dataRow) toInt(index int) (int, error) {
	if index >= len(s.datas){
		return 0, errors.New("out of range")
	}

	d := s.datas[index]
	int, err := strconv.Atoi(d)
	return int, err
}

func (s* dataRow) toString(index int) (string, error) {
	if index >= len(s.datas){
		return "", errors.New("out of range")
	}
	return s.datas[index], nil
}

func (s* dataRow) toFloat32(index int) (float32, error) {
	if index >= len(s.datas){
		return 0.0, errors.New("out of range")
	}

	d := s.datas[index]
	f, err := strconv.ParseFloat(d, 32)
	return float32(f), err
}

func (s* dataRow) toFloat64(index int) (float64, error) {
	if index >= len(s.datas){
		return 0.0, errors.New("out of range")
	}

	d := s.datas[index]
	f, err := strconv.ParseFloat(d, 32)
	return f, err
}

type Table struct {
	sheets        	[] *dataRow
	types 			dataRow
	keyToIndexMap	map[string] int
}

func (s*Table) GetInt(key string, idx int)(int, error){
	if idx < len(s.sheets){
		i,ok:= s.keyToIndexMap[key]
		if ok {
			row := s.sheets[idx]
			return row.toInt(i)
		}
	}
	return 0, errors.New("GetInt error")
}

func (s*Table) GetString(key string, idx int)(string, error){
	if idx < len(s.sheets){
		i,ok:= s.keyToIndexMap[key]
		if ok {
			row := s.sheets[idx]
			return row.toString(i)
		}
	}
	return "", errors.New("GetString error")
}

func (s*Table) GetFloat32(key string, idx int)(float32, error){
	if idx < len(s.sheets){
		i,ok:= s.keyToIndexMap[key]
		if ok {
			row := s.sheets[idx]
			return row.toFloat32(i)
		}
	}
	return 0.0, errors.New("GetFloat32 error")
}

func (s*Table) GetFloat64(key string, idx int)(float64, error){
	if idx < len(s.sheets){
		i,ok:= s.keyToIndexMap[key]
		if ok {
			row := s.sheets[idx]
			return row.toFloat64(i)
		}
	}
	return 0.0, errors.New("GetFloat64 error")
}

func (s*Table) GetCnt() int {
	return len(s.sheets)
}

func (s *Table) ToString() string {

	indexToKey := make(map[int] string)
	for k, i := range s.keyToIndexMap{
		indexToKey[i] = k
	}

	objStr := ""
	for i, v := range s.sheets{
		objStr += "{"
		for j, d := range v.datas{
			value := d
			key := indexToKey[j]
			vtype := s.types.datas[j]
			vtype = strings.ToLower(vtype)

			if j != 0{
				if vtype == "string"{
					objStr += fmt.Sprintf(",\"%s\":\"%s\"",key, value)
				}else if vtype == "int"{
					int, _ := strconv.Atoi(value)
					objStr += fmt.Sprintf(",\"%s\":%d",key, int)
				}else{
					objStr += fmt.Sprintf(",\"%s\":%s",key, value)
				}
			}else{
				if vtype == "string"{
					objStr += fmt.Sprintf("\"%s\":\"%s\"",key, value)
				}else if vtype == "int"{
					int, _ := strconv.Atoi(value)
					objStr += fmt.Sprintf("\"%s\":%d",key, int)
				}else{
					objStr += fmt.Sprintf("\"%s\":%s",key, value)
				}
			}
		}
		if i != len(s.sheets) - 1{
			objStr+="},"
		}else{
			objStr+="}"
		}
	}
	ret := fmt.Sprintf("[%s]", objStr)
	return ret
}

func newDataRow(n int) dataRow{
	return dataRow{datas: make([]string, n)}
}

type xlsxManager struct {
	sheetMap        map[string]*Table
	mutex 			sync.RWMutex
	rootDir         string
}

func (s* xlsxManager) SetRootDir(rootDir string)  {
	s.rootDir = rootDir
}
func (s* xlsxManager) Load(xlsx string){
	xlsx = filepath.Join(s.rootDir, xlsx)
	f, err := excelize.OpenFile(xlsx)
	if err != nil{
		Log.Error("Load xlsx %s error:%s", xlsx, err.Error())
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, name := range f.GetSheetMap() {
		rows := f.GetRows(name)
		n := len(rows)
		if n >0 {
			key := xlsx+"/"+name
			t := Table{}
			t.sheets = make([] *dataRow, n-2)
			for idx, row := range rows {

				if idx == 0{
					//类型
					ncol := len(row)
					r := newDataRow(ncol)
					for index, colCell := range row {
						r.datas[index] = colCell
					}
					t.types = r

				}else if idx == 1{
					//描述
					m := make(map[string] int)
					for index, colCell := range row {
						m[colCell] = index
					}
					t.keyToIndexMap = m

				}else{
					//数据
					ncol := len(row)
					r := newDataRow(ncol)
					for index, colCell := range row {
						r.datas[index] = colCell
					}
					t.sheets[idx-2] = &r
				}
			}
			s.sheetMap[key] = &t
			t.ToString()
		}
	}

	Log.Info("Load xlsx %s finish", xlsx)
}

func (s* xlsxManager) Get(xlsx string, sheet string) *Table {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	xlsx = filepath.Join(s.rootDir, xlsx)
	key := xlsx+"/"+sheet
	return s.sheetMap[key]
}