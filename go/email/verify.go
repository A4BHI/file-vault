package email

import (
	"math/rand"
	"strconv"

	"gopkg.in/mail.v2"
)

func Otp() string {
	otp := rand.Intn(900000) + 100000
	return strconv.Itoa(otp)
}

func SendMail(email string, username string) (string, error) {

	dm := mail.NewDialer("smtp.gmail.com", 587, "vaultx000@gmail.com", "")
	otp := Otp()

	mess := mail.NewMessage()

	mess.SetHeader("From", "vaultx000@gmail.com")
	mess.SetHeader("To", email)
	mess.SetHeader("Subject", "VaultX Email Verification")
	mess.SetBody("text/plain", "Hi "+username+"\nYour Email Verification code for VaultX is \n \n OTP CODE: "+otp+"\nPlease enter this code in VaultX Email Verification form.\nThis code is confidential - DO NOT SHARE!!!")
	err := dm.DialAndSend(mess)

	return otp, err

}
