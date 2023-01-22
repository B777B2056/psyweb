package models

/* 工作人员 */
type StaffUser struct {
	Id       string `json:"StaffId"`
	Password string `json:"Password"`
}

/* 从数据库查找当前工作人员的账号 */
func (staff *StaffUser) isIdExist() bool {
	return true
}

/* 从数据库获取当前工作人员的密码 */
func (staff *StaffUser) getPasswordFromDB() string {
	return "000000"
}

/* 验证当前登陆账号密码是否与已知工作人员账号密码匹配 */
func (staff *StaffUser) IsPassVerification() bool {
	return staff.isIdExist() && (staff.Password == staff.getPasswordFromDB())
}
