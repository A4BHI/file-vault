package email

import (
	"fmt"
	"math/rand"
	"strconv"

	"gopkg.in/mail.v2"
)

func Otp() string {
	otp := rand.Intn(900000) + 100000
	return strconv.Itoa(otp)
}

func SendMail(email string, username string) {

	dm := mail.NewDialer("smtp.gmail.com", 587, "vaultx000@gmail.com", "rvcm inec dgxh eypu")
	otp := Otp()

	mess := mail.NewMessage()

	mess.SetHeader("From", "vaultx000@gmail.com")
	mess.SetHeader("To", email)
	mess.SetHeader("Subject", "VaultX Email Verification")
	mess.SetBody("text/plain", "Your email verification code for VaultX is \n \n OTP CODE: "+otp)
	// mess.SetBody("text/plain",otp)
	err := dm.DialAndSend(mess)

	if err != nil {
		fmt.Println(err)
	}

}
