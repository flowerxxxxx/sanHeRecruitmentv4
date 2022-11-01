package controllerModule

import (
	"sanHeRecruitment/models/exportModel"
	"sanHeRecruitment/models/mysqlModel"
	"strconv"
)

type StatisticsModule struct {
}

func (sm *StatisticsModule) UpgradeSliChanger(input []mysqlModel.WaitingUpgradeXls) []interface{} {
	var xlsExport []interface{}
	for _, x := range input {
		newR := ""
		if x.CompanyExist == 0 {
			newR = "新注册公司"
		} else if x.Qualification == 1 {
			newR = "公司已注册"
		}

		done := ""
		if x.Qualification == 0 {
			done = "未审核"
		} else if x.Qualification == 1 {
			done = "已通过"
		} else if x.Qualification == 2 {
			done = "未通过"
		}
		if x.Phone == "" {
			x.Phone = "未填写"
		}
		xlsExport = append(xlsExport, &exportModel.UpgradeExporter{
			CompanyName: x.CompanyName,
			Phone:       x.Phone,
			Address:     x.Address,
			Description: x.Description,
			NewRegister: newR,
			Done:        done,
			Applicant:   x.Name,
			ApplyTime:   x.ApplyTimeOut,
		})
	}
	return xlsExport
}

func (sm *StatisticsModule) ComPubSliChanger(input []mysqlModel.UserComArticleXls) []interface{} {
	var xlsExport []interface{}
	for _, x := range input {
		smi := ""
		sma := ""
		if x.SalaryMin == -1 && x.SalaryMax == -1 {
			smi = "待商议"
			sma = "待商议"
		} else {
			smi = strconv.Itoa(x.SalaryMin)
			sma = strconv.Itoa(x.SalaryMax)
		}
		artType := ""
		if x.ArtType == "request" {
			artType = "经营需求"
		} else {
			artType = "招聘需求"
		}
		StatusFlag := ""
		if x.Status == 0 {
			StatusFlag = "待审核"
		} else if x.Status == 1 {
			StatusFlag = "通过"
		} else if x.Status == 2 {
			StatusFlag = "未通过"
		}
		xlsExport = append(xlsExport, &exportModel.ComPubExport{
			CompanyName: x.CompanyName,
			PersonScale: x.PersonScale,
			Title:       x.Title,
			View:        x.View,
			Nickname:    x.Nickname,
			Name:        x.Name,
			Tags:        x.JobLabel + "," + x.Tags,
			SalaryMin:   smi,
			SalaryMax:   sma,
			Content:     x.Content,
			ArtType:     artType,
			Status:      StatusFlag,
		})
	}
	return xlsExport
}
