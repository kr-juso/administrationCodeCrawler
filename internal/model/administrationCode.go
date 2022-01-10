package model

import (
	"time"
)

type AdministrationCode struct {
	Code        string     // 행정동코드
	City        string     // 시도명
	State       string     // 시군구명(optional)
	Town        string     // 읍면동명(optional)
	CreateDate  time.Time  // 생성일자
	DestroyDate *time.Time // 말소일자
}

func (administrationCode AdministrationCode) ToCsvRow() []string {
	row := make([]string, 6)
	row[0] = administrationCode.Code
	row[1] = administrationCode.City
	row[2] = administrationCode.State
	row[3] = administrationCode.Town
	row[4] = administrationCode.CreateDate.Format("20060102") // yyyyMMdd
	if administrationCode.DestroyDate != nil {
		row[5] = administrationCode.DestroyDate.Format("20060102") // yyyyMMdd
	} else {
		row[5] = ""
	}

	return row
}
