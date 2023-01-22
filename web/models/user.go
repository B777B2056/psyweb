package models

import (
	"database/sql"
	"log"
	"math"
	"math/rand"
	"psyWeb/utils"
	"reflect"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var userMap sync.Map // 模拟数据库

/* 生成6位短信验证码 */
func generateVerificationCode() (code string) {
	code_num := 0
	numeric := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 6; i++ {
		code_num += (numeric[rand.Intn(r)] * int(math.Pow(10, float64(i))))
	}

	if code_num/int(math.Pow(10, float64(5))) == 0 {
		code = "0"
	} else {
		code = ""
	}
	code += strconv.Itoa(code_num)
	return code
}

/* 普通用户 */
type User struct {
	/* 基本信息 */
	PhoneNumber  string           `json:"PhoneNumber"`
	Code         string           `json:"Code"`
	SerialNumber string           `json:"SerialNumber"`
	Name         string           `json:"Name"`
	Gender       utils.UserGender `json:"Gender"`
	Age          uint8            `json:"Age"`
	/* 量表信息 */
	SASScore float32 `json:"SAS"`
	ESSScore float32 `json:"ESS"`
	ISIScore float32 `json:"ISI"`
	SDSScore float32 `json:"SDS"`
	/* 用户状态 */
	Status utils.UserStatus `json:"UserStatus"`
}

/* 查询当前用户是否存在于数据库中 */
func (user *User) IsExist() (bool, error) {
	db := utils.GetPsyWebDataBaseInstance().Db
	var phone_number string
	err := db.QueryRow("SELECT PhoneNumber FROM user WHERE PhoneNumber=?", user.PhoneNumber).Scan(&phone_number)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, err
	}
}

/* 填充所有字段 */
func (user *User) FillFields() {
	db := utils.GetPsyWebDataBaseInstance().Db
	rows, err := db.Query("SELECT * FROM user WHERE PhoneNumber=?", user.PhoneNumber)
	if err != nil {
		log.Println("select *", err)
	}
	for rows.Next() {
		rows.Scan(&user.PhoneNumber, &user.Code, &user.SerialNumber, &user.Name, &user.Gender, &user.Age, &user.SASScore, &user.ESSScore, &user.ISIScore, &user.SDSScore, &user.Status)
	}
}

/* 更新用户状态 */
func (user *User) updateStatus() error {
	db := utils.GetPsyWebDataBaseInstance().Db
	err := db.QueryRow("SELECT Status FROM user WHERE PhoneNumber=?", user.PhoneNumber).Scan(&user.Status)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	// 用户状态进入下一状态
	if user.Status != utils.Done {
		user.Status++
	}
	return nil
}

/* 删除用户 */
func (user *User) Del() error {
	db := utils.GetPsyWebDataBaseInstance().Db
	_, err := db.Exec("DELETE FROM user WHERE PhoneNumber=?", user.PhoneNumber)
	return err
}

/* 新建用户 */
func (user *User) New() (err error) {
	log.Printf("New user: %+v\n", *user)
	user.updateStatus()

	sqlStr := "INSERT INTO user(PhoneNumber, Code, SerialNumber, Name, Gender, Age, SASScore, ESSScore, ISIScore, SDSScore, Status) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	_, err = utils.GetPsyWebDataBaseInstance().Db.Exec(sqlStr, user.PhoneNumber, user.Code, user.SerialNumber, user.Name, user.Gender, user.Age, user.SASScore, user.ESSScore, user.ISIScore, user.SDSScore, user.Status)
	if err != nil {
		user.Del() // 出错则删除已插入的用户信息
		log.Printf("insert %s failed, err:%v\n", "PhoneNumber", err)
	}
	return err
}

/* 更新量表相关的内容 */
func (user *User) UpdateScaleResult() (err error) {
	log.Printf("Update scale user: %+v\n", *user)
	user.updateStatus() // 量表完成后，更新状态为下一状态

	s_type := reflect.TypeOf(user).Elem()
	s_value := reflect.ValueOf(user).Elem()
	for i := 0; i < s_value.NumField(); i++ {
		field_name := s_type.Field(i).Name
		// 不更新与量表无关的变量
		if (field_name == "PhoneNumber") || (field_name == "Code") {
			continue
		}
		field_val := s_value.Field(i).Interface()
		_, err = utils.GetPsyWebDataBaseInstance().Db.Exec("UPDATE user SET "+field_name+"=? WHERE PhoneNumber=?", field_val, user.PhoneNumber)
		if err != nil {
			log.Printf("update %s failed, err:%v\n", field_name, err)
			break
		}
	}
	return err
}

/* 更新脑电诊断相关的内容 */
func (user *User) UpdateEEGResult() (err error) {
	log.Printf("Update eeg user: %s", user.PhoneNumber)
	user.updateStatus() // 诊断完成后，更新状态为下一状态
	_, err = utils.GetPsyWebDataBaseInstance().Db.Exec("UPDATE user SET Status=? WHERE PhoneNumber=?", user.Status, user.PhoneNumber)
	return err
}

/* 发送验证码至用户手机 */
func (user *User) SendVerificationCodeToUserPhone() {
	user.Code = generateVerificationCode()
	log.Printf("Send Verification Code: %s", user.Code)
	userMap.Store(user.PhoneNumber, user.Code)
	utils.SendSMSByTencentCloud(user.PhoneNumber, user.Code)
}

/* 检查给定用户信息是否匹配已有用户信息（手机号，验证码是否匹配），以及查询用户状态 */
type UserVerification struct {
	Result bool             `json:"VerificationResult"`
	Status utils.UserStatus `json:"UserStatus"`
}

func (info *UserVerification) Check(user User) *UserVerification {
	if ptr, _ := userMap.Load(user.PhoneNumber); ptr != nil {
		// 检查验证码是否匹配
		correct_code := ptr.(string)
		info.Result = (correct_code == user.Code)
		// 查询User状态（仅验证码匹配时进行）
		if info.Result {
			if isExist, _ := user.IsExist(); !isExist {
				user.New() // 新建用户
			} else {
				user.FillFields() // 填充字段
			}
			info.Status = user.Status // 已有用户
		}
		userMap.Delete(user.PhoneNumber)
	} else {
		info.Result = false
		info.Status = utils.NotExist
	}
	return info
}
