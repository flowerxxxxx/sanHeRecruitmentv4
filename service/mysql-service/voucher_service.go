package mysql_service

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/formatUtil"
)

type VoucherService struct {
}

// SaveUpgradeVoucher 存储公司凭证
func (vs *VoucherService) SaveUpgradeVoucher(username, voucherUrl string, companyId int, TimeId int64) (err error) {
	var voucherInfo = mysqlModel.Voucher{
		Username:   username,
		VoucherPic: voucherUrl,
		CompanyId:  companyId,
		TimeId:     TimeId,
	}
	if err = dao.DB.Save(&voucherInfo).Error; false {
		return
	}
	return
}

func (vs *VoucherService) QueryUpgradeVouchers(companyId int, username, host string, TimeId int64) []mysqlModel.Voucher {
	var voucherInfo []mysqlModel.Voucher
	dao.DB.Where("company_id=?", companyId).Where("username=?", username).Where("time_id=?", TimeId).Find(&voucherInfo)
	for i, m := 0, len(voucherInfo); i < m; i++ {
		voucherInfo[i].VoucherPic = formatUtil.GetPicHeaderBody(host, voucherInfo[i].VoucherPic)
	}
	return voucherInfo
}

func (vs *VoucherService) AdminQueryUpgradeVouchers(timeId int, fromUsername, host string) []mysqlModel.Voucher {
	var voucherInfo []mysqlModel.Voucher
	dao.DB.Where("username =?", fromUsername).Where("time_id =?", timeId).Find(&voucherInfo)
	for i, m := 0, len(voucherInfo); i < m; i++ {
		voucherInfo[i].VoucherPic = formatUtil.GetPicHeaderBody(host, voucherInfo[i].VoucherPic)
	}
	return voucherInfo
}
