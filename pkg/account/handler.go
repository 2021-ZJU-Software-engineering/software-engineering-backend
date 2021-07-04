package account

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/AsterNighT/software-engineering-backend/pkg/database/models"
	"github.com/AsterNighT/software-engineering-backend/pkg/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AccountHandler struct{}

// @Summary create and account based on email(as id), type, name and password
// @Description will check primarykey other, then add to accountList if possible
// @Tags Account
// @Produce json
// @Param email body string true "user e-mail"
// @Param type body string true "user type"
// @Param firstname body string true "user first name"
// @Param lastname body string true "user last name"
// @Param passwd body string true "user password"
// @Param birthday body string true "user birthday string in yyyy-mm-dd"
// @Param gender body string true "user gender"
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/create [POST]
func (h *AccountHandler) CreateAccount(c echo.Context) error {
	type RequestBody struct {
		Email string `json:"email" validate:"required"`

		Type      models.AcountType `json:"type" validate:"required"`
		FirstName string            `json:"firstname" validate:"required"`
		LastName  string            `json:"lastname" validate:"required"`
		Passwd    string            `json:"passwd" validate:"required"`

		BirthString string `json:"birthday" validate:"required"`
		Gender      string `json:"gender" validate:"required"`
	}

	var body RequestBody
	if err := utils.ExtractDataWithValidating(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err.Error()))
	}
	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, body.Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid E-mail Address", nil))
	}
	if body.Type != models.PatientType && body.Type != models.DoctorType && body.Type != models.AdminType {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid Account Type", nil))
	}

	AccountPasswdLen, _ := strconv.Atoi(os.Getenv("PASSWD_MIN_LEN"))
	if len(body.Passwd) < AccountPasswdLen {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid Password Length", nil))
	}

	db, _ := c.Get("db").(*gorm.DB)
	if err := db.Where("email = ?", body.Email).First(&models.Account{}).Error; err == nil {
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail or AccountID occupied", nil))
	}

	birthDay, _ := time.Parse("2006-01-02", body.BirthString)

	account := models.Account{
		Email: body.Email,

		Type:      body.Type,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Passwd:    body.Passwd,

		Gender:   body.Gender,
		Birthday: birthDay,
	}
	account.HashPassword()

	// Add account
	if result := db.Create(&account); result.Error != nil {
		return c.JSON(http.StatusBadRequest, api.Return("DB error", result.Error.Error()))
	}

	// Update account id since it's auto-gen by db
	if err := db.Where("email = ?", account.Email).First(&account).Error; err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("DB error", err.Error()))
	}

	if account.Type == "doctor" {
		doctor := models.Doctor{AccountID: account.ID, Email: account.Email}
		if result := db.Create(&doctor); result.Error != nil {
			return c.JSON(http.StatusBadRequest, api.Return("DB error", result.Error.Error()))
		}
	} else if account.Type == "patient" {
		patient := models.Patient{AccountID: account.ID, Email: account.Email}
		if result := db.Create(&patient); result.Error != nil {
			return c.JSON(http.StatusBadRequest, api.Return("DB error", result.Error.Error()))
		}
	} else {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid account type:", account.Type))
	}

	token, err := account.GenerateToken()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Return("fail to generete token", err.Error()))
	}

	return c.JSON(http.StatusOK, api.Return("Created", echo.Map{
		"account": account,
		"token":   token,
	}))
}

// @Summary set doctor info
// @Description
// @Tags Account
// @Produce json
// @Param departmentid body uint true "doctor department id"
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/setdoctor [POST]
func (h *AccountHandler) SetDoctor(c echo.Context) error {
	accountID := c.Get("id")
	type RequestBody struct {
		DepartmentID uint `json:"departmentid" validate:"required"`
	}

	var body RequestBody
	if err := utils.ExtractDataWithValidating(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err.Error()))
	}

	db, _ := c.Get("db").(*gorm.DB)

	if result := db.Model(&models.Doctor{}).Where("account_id = ?", accountID).Update("department", body.DepartmentID); result.Error != nil {
		return c.JSON(http.StatusBadRequest, api.Return("DB error", result.Error.Error()))
	}
	return c.JSON(http.StatusOK, api.Return("Doctor set", nil))
}

// @Summary set patient info
// @Description
// @Tags Account
// @Produce json
// @Param allergy body string true "patient allergy history"
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/setpatient [POST]
func (h *AccountHandler) SetPatient(c echo.Context) error {
	accountID := c.Get("id")
	type RequestBody struct {
		Allergy string `json:"allergy" validate:"required"`
	}

	var body RequestBody
	if err := utils.ExtractDataWithValidating(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err.Error()))
	}

	db, _ := c.Get("db").(*gorm.DB)

	if result := db.Model(&models.Patient{}).Where("account_id = ?", accountID).Updates(map[string]interface{}{"allergy": body.Allergy}); result.Error != nil {
		return c.JSON(http.StatusBadRequest, api.Return("DB error", result.Error.Error()))
	}
	return c.JSON(http.StatusOK, api.Return("Patient set", nil))
}

// @Summary check email's existense
// @Description
// @Tags Account
// @Produce json
// @Param email body string true "user e-mail"
// @Success 200 {string} api.ReturnedData{data=echo.Map{"emailok": true}}
// @Failure 400 {string} api.ReturnedData{data=echo.Map{"emailok": false}}
// @Router /account/checkemail [POST]
func (h *AccountHandler) CheckEmail(c echo.Context) error {
	type RequestBody struct {
		Email string `json:"email" validate:"required"`
	}
	var body RequestBody

	if err := utils.ExtractDataWithValidating(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err.Error()))
	}

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, body.Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid E-mail Address", nil))
	}

	db, _ := c.Get("db").(*gorm.DB)
	var account models.Account
	if err := db.Where("email = ?", body.Email).First(&account).Error; err != nil { // not found
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail", echo.Map{"emailok": false}))
	}
	return c.JSON(http.StatusOK, api.Return("E-Mail", echo.Map{"emailok": true}))

}

// @Summary login using email and passwd
// @Description
// @Tags Account
// @Produce json
// @Param email body string true "user e-mail"
// @Param passwd body string true "user password"
// @Success 200 {string} api.ReturnedData{data=echo.Map{"account": account, "token": token,}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/login [POST]
func (h *AccountHandler) LoginAccount(c echo.Context) error {
	type RequestBody struct {
		Email  string `json:"email" validate:"required"`
		Passwd string `json:"passwd" validate:"required"`
	}
	var body RequestBody

	if err := utils.ExtractDataWithValidating(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err.Error()))
	}

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, body.Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid E-mail Address", nil))
	}

	AccountPasswdLen, _ := strconv.Atoi(os.Getenv("PASSWD_MIN_LEN"))
	if len(body.Passwd) < AccountPasswdLen {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid Password Length", nil))
	}

	db, _ := c.Get("db").(*gorm.DB)
	var account models.Account
	if err := db.Where("email = ?", body.Email).First(&account).Error; err != nil { // not found
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail", echo.Map{"emailok": false}))
	}
	if bcrypt.CompareHashAndPassword([]byte(account.Passwd), []byte(body.Passwd)) != nil {
		return c.JSON(http.StatusBadRequest, api.Return("Wrong Password", nil))
	}

	token, err := account.GenerateToken()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Return("fail to generete token", err.Error()))
	}

	return c.JSON(http.StatusOK, api.Return("Logged in", echo.Map{
		"account": account,
		"token":   token,
	}))
}

// @Summary the interface of modifying password
// @Description
// @Tags Account
// @Produce json
// @Param email body string true "user e-mail"
// @Param passwd body string true "user password (the old one)"
// @Param newpasswd body string true "user password (the new one)"
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/modifypasswd [POST]
func (h *AccountHandler) ModifyPasswd(c echo.Context) error {
	type RequestBody struct {
		Email     string `json:"email" validate:"required"`
		Passwd    string `json:"passwd" validate:"required"`
		NewPasswd string `json:"newpasswd" validate:"required"`
	}
	var body RequestBody

	if err := utils.ExtractDataWithValidating(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err.Error()))
	}

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, body.Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid E-mail Address", nil))
	}

	// Check old passwd
	db, _ := c.Get("db").(*gorm.DB)
	var account models.Account
	if err := db.Where("email = ?", body.Email).First(&account).Error; err != nil { // not found
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail", echo.Map{"emailok": false}))
	}
	if bcrypt.CompareHashAndPassword([]byte(account.Passwd), []byte(body.Passwd)) != nil {
		return c.JSON(http.StatusBadRequest, api.Return("Wrong Password", nil))
	}

	AccountPasswdLen, _ := strconv.Atoi(os.Getenv("PASSWD_MIN_LEN"))
	if len(body.NewPasswd) < AccountPasswdLen {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid Password Length", nil))
	}

	account.Passwd = body.NewPasswd
	account.HashPassword()

	if result := db.Model(&models.Account{}).Where("id = ?", account.ID).Update("passwd", account.Passwd); result.Error != nil {
		return c.JSON(http.StatusBadRequest, api.Return("DB error", result.Error.Error()))
	}

	return c.JSON(http.StatusOK, api.Return("Successfully modified", nil))
}

// @Summary the interface of sending email to reset password
// @Description can only be called during logged-in status since there is no password check
// @Tags Account
// @Produce json
// @Param email body string true "user e-mail"
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=echo.Map{"authCode": authCode}}
// @Router /account/sendemail [POST]
func (h *AccountHandler) SendEmail(c echo.Context) error {
	type RequestBody struct {
		Email string `json:"email" validate:"required"`
	}
	var body RequestBody

	if err := utils.ExtractDataWithValidating(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err.Error()))
	}

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, body.Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid E-mail Address", nil))
	}

	db, _ := c.Get("db").(*gorm.DB)

	authCode := ""
	for i := 0; i < 6; i++ {
		nBig, _ := rand.Int(rand.Reader, big.NewInt(10))
		authCode += string("0123456789"[nBig.Int64()])
	}
	c.Logger().Debug(authCode)

	if tmp := db.Model(&models.Auth{}).Where("email = ?", body.Email).Update("auth_code", authCode); tmp.Error != nil {
		return c.JSON(http.StatusBadRequest, api.Return("DB error", tmp.Error))
	}

	emailServerHost := os.Getenv("EMAIL_SERVER_HOST")
	emailServerPort := os.Getenv("EMAIL_SERVER_PORT")
	emailUser := os.Getenv("EMAIL_USER")
	emailPasswd := os.Getenv("EMAIL_PASSWD")
	expireMin, _ := strconv.Atoi(os.Getenv("EMAIL_VALID_MIN"))

	auth := models.Auth{
		Email:           body.Email,
		AuthCode:        authCode,
		AuthCodeExpires: time.Now().Add(time.Duration(expireMin) * time.Minute),
	}

	db.Delete(&auth)

	if tmp := db.Model(&models.Auth{}).Where("email = ?", body.Email).Update("auth_code", authCode); tmp.Error != nil {
		return c.JSON(http.StatusBadRequest, api.Return("DB error", tmp.Error))
	}

	if result := db.Create(&auth); result.Error != nil {
		return c.JSON(http.StatusBadRequest, api.Return("DB error", result.Error.Error()))
	}

	emailAuth := smtp.PlainAuth("", emailUser, emailPasswd, emailServerHost)
	to := body.Email
	msg := []byte("From: \"MediConnect\" <noreply@mediconnect.com>\n" +
		"To: " + to + "\n" +
		"Subject: MediConnect Account Reset\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\n" +
		"\n" +
		"Your verification code is " + authCode + " (Only valid in " + strconv.Itoa(expireMin) + " minutes)\n")
	if err := smtp.SendMail(emailServerHost+":"+emailServerPort, emailAuth, emailUser, []string{to}, msg); err != nil {
		return c.JSON(http.StatusOK, api.Return("Email server error", echo.Map{"err": err, "msg": msg}))
	}

	return c.JSON(http.StatusOK, api.Return("Successfully send reset email", nil))
}

// return info of doctor to front end
// func (h *AccountHandler) RetInfo(c echo.Context)

// @Summary check authcode's correctness
// @Description
// @Tags Account
// @Produce json
// @Param email body string true "user e-mail"
// @Param authcode body string true "given auth code"
// @Success 200 {string} api.ReturnedData{data=echo.Map{"authcodeok": false}}
// @Failure 400 {string} api.ReturnedData{data=echo.Map{"authcodeok": true}}
// @Router /account/checkauthcode [POST]
func (h *AccountHandler) CheckAuthCode(c echo.Context) error {
	type RequestBody struct {
		Email    string `json:"email" validate:"required"`
		AuthCode string `json:"authcode" validate:"required"`
	}
	var body RequestBody

	if err := utils.ExtractDataWithValidating(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err.Error()))
	}

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, body.Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid E-mail Address", nil))
	}

	// Check authcode
	db, _ := c.Get("db").(*gorm.DB)
	var auth models.Auth
	if err := db.Where("email = ?", body.Email).First(&auth).Error; err != nil { // not found
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail", echo.Map{"emailok": false}))
	}
	if auth.AuthCode == body.AuthCode && time.Now().Before(auth.AuthCodeExpires) {
		return c.JSON(http.StatusOK, api.Return("AuthCode", echo.Map{"authcodeok": true}))
	}
	return c.JSON(http.StatusBadRequest, api.Return("AuthCode", echo.Map{"authcodeok": false}))

}

// @Summary the interface of reset password
// @Description
// @Tags Account
// @Produce json
// @Param email body string true "user e-mail"
// @Param authCode body string true "given auth code"
// @Param passwd body string true "user password (the new one)"
// @Success 200 {string} api.ReturnedData{data=nil}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/resetpasswd [POST]
func (h *AccountHandler) ResetPasswd(c echo.Context) error {
	type RequestBody struct {
		Email     string `json:"email" validate:"required"`
		AuthCode  string `json:"authcode" validate:"required"`
		NewPasswd string `json:"newpasswd" validate:"required"`
	}
	var body RequestBody

	if err := utils.ExtractDataWithValidating(c, &body); err != nil {
		return c.JSON(http.StatusBadRequest, api.Return("error", err.Error()))
	}

	if ok, _ := regexp.MatchString(`^\w+@\w+[.\w+]+$`, body.Email); !ok {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid E-mail Address", nil))
	}

	// Check authcode
	db, _ := c.Get("db").(*gorm.DB)
	var account models.Account
	var auth models.Auth
	if err := db.Where("email = ?", body.Email).First(&account).Error; err != nil { // not found
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail", echo.Map{"emailok": false}))
	}
	if err := db.Where("email = ?", body.Email).First(&auth).Error; err != nil { // not found
		return c.JSON(http.StatusBadRequest, api.Return("E-Mail", echo.Map{"emailok": false}))
	}

	if auth.AuthCode != body.AuthCode || time.Now().After(auth.AuthCodeExpires) {
		return c.JSON(http.StatusBadRequest, api.Return("AuthCode", echo.Map{"authcodeok": false}))
	}

	AccountPasswdLen, _ := strconv.Atoi(os.Getenv("PASSWD_MIN_LEN"))
	if len(body.NewPasswd) < AccountPasswdLen {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid Password Length", nil))
	}

	account.Passwd = body.NewPasswd
	account.HashPassword()

	if result := db.Model(&models.Account{}).Where("id = ?", account.ID).Update("passwd", account.Passwd); result.Error != nil {
		return c.JSON(http.StatusBadRequest, api.Return("DB error", result.Error.Error()))
	}

	return c.JSON(http.StatusOK, api.Return("Successfully modified", nil))
}

// @Summary the interface of getting current token's info
// @Description
// @Tags Account
// @Produce json
// @Success 200 {string} api.ReturnedData{data=echo.Map{"id": account.ID, "email": account.Email, "type": account.Type, "firstname": account.FirstName, "lastname": account.LastName, "gender": account.Gender, "birthday": account.Birthday, "department": doctor.Department, "allergy": patient.Allergy}}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/getinfo [GET]
func (h *AccountHandler) GetInfo(c echo.Context) error {
	id := c.Get("id")

	db, _ := c.Get("db").(*gorm.DB)
	var account models.Account
	if err := db.Where("id = ?", id).First(&account).Error; err != nil { // not found
		return c.JSON(http.StatusBadRequest, api.Return("Not logged in", nil))
	}

	if account.Type == "doctor" {
		var doctor models.Doctor
		if err := db.Where("account_id = ?", id).First(&doctor).Error; err != nil { // not found
			return c.JSON(http.StatusBadRequest, api.Return("Wrong account id for doctor", nil))
		}
		return c.JSON(http.StatusOK, api.Return("Successfully Get", echo.Map{"id": account.ID, "email": account.Email, "type": account.Type, "firstname": account.FirstName, "lastname": account.LastName, "gender": account.Gender, "birthday": account.Birthday, "department": doctor.Department}))
	} else if account.Type == "patient" {
		var patient models.Patient
		if err := db.Where("account_id = ?", id).First(&patient).Error; err != nil { // not found
			return c.JSON(http.StatusBadRequest, api.Return("Wrong account id for patient", nil))
		}
		return c.JSON(http.StatusOK, api.Return("Successfully Get", echo.Map{"id": account.ID, "email": account.Email, "type": account.Type, "firstname": account.FirstName, "lastname": account.LastName, "gender": account.Gender, "birthday": account.Birthday, "allergy": patient.Allergy}))
	} else {
		return c.JSON(http.StatusBadRequest, api.Return("Invalid account type: ", account.Type))
	}
}

// @Summary the interface of getting given pat's info
// @Description
// @Tags Account
// @Produce json
// @Param patientID path uint true "patient ID"
// @Success 200 {string} api.ReturnedData{data=echo.Map{"id": account.ID, "email": account.Email, "type": account.Type, "firstname": account.FirstName, "lastname": account.LastName, "gender": account.Gender, "birthday": account.Birthday, "allergy": patient.Allergy}}
// @Failure 400 {string} api.ReturnedData{data=nil}
// @Router /account/getinfobypatid/{patientID} [GET]
func (h *AccountHandler) GetInfoByPatID(c echo.Context) error {
	var account models.Account
	var patient models.Patient

	patID, _ := strconv.Atoi(c.Param("patientID"))
	db, _ := c.Get("db").(*gorm.DB)

	if err := db.Where("id = ?", patID).First(&patient).Error; err != nil { // not found
		return c.JSON(http.StatusBadRequest, api.Return("Wrong Patient ID", nil))
	}

	if err := db.Where("id = ?", patient.AccountID).First(&account).Error; err != nil { // not found
		return c.JSON(http.StatusBadRequest, api.Return("Wrong Account ID", nil))
	}
	return c.JSON(http.StatusOK, api.Return("Successfully Get", echo.Map{"id": account.ID, "email": account.Email, "type": account.Type, "firstname": account.FirstName, "lastname": account.LastName, "gender": account.Gender, "birthday": account.Birthday, "allergy": patient.Allergy}))
}

/**
 * @brief public method for getting current logged-in account's ID.
 */
func getAccountID(c echo.Context) (uint, error) {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		auth = c.QueryParam("token")
	} else {
		words := strings.Split(auth, "Bearer ")
		if len(words) != 2 {
			return 0, errors.New("invalid authorization header")
		}
		auth = words[1]
	}
	c.Logger().Debug("get token: ", auth)
	id, err := ParseToken(auth)
	if err != nil {
		return 0, err
	}
	return id, nil
}

/**
 * @brief middleware for getting current logged-in account's ID.
 */
func CheckAccountID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := getAccountID(c)
		if err != nil {
			return c.JSON(403, api.Return("fail to get id from token", err.Error()))
		}
		c.Set("id", id)
		return next(c)
	}
}

func ParseToken(tokenString string) (uint, error) {
	if tokenString == "" {
		return 0, fmt.Errorf("cannot find auth token")
	}
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		return uint((*claims)["id"].(float64)), nil
	}
	return 0, err
}
