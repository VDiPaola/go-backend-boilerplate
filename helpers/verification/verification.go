package verification

import (
	"boilerplate/backend/database"
	"boilerplate/backend/helpers/mail"
	"boilerplate/backend/models"
	"crypto/rand"
	"errors"
	"io"
	"time"

	"github.com/VDiPaola/go-backend-module/module_helpers"
	"github.com/VDiPaola/go-backend-module/module_models"
)

func GenerateCode(length int, expirationLengthMins int64) module_models.Code {
	var nums = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = nums[int(b[i])%len(nums)]
	}

	return module_models.Code{
		Value:     string(b),
		ExpiresAt: module_helpers.UnixNanoToJS(time.Now().Add(time.Minute * time.Duration(expirationLengthMins)).UnixNano()),
	}
}

func VerifyCode(code module_models.Code, sentCode string) error {

	//check that codes match
	if code.Value != sentCode {
		return errors.New("code doesn't match")
	}

	//check code hasnt expired
	if module_helpers.NowJS() > code.ExpiresAt {
		return errors.New("code has expired")
	}

	return nil
}

func GenerateCodeAndSend(user models.User, expirationLengthMins int64) error {
	//generate code
	code := GenerateCode(6, expirationLengthMins)

	//set code in user
	user.Code = code

	//update user
	if dbc := database.Connection.Updates(&user); dbc.Error != nil {
		return dbc.Error
	}

	//send email
	go mail.Send([]string{user.Email}, "Code", `<div>
	<div style="display:flex; flex-direction:column; font-family: sans-serif; padding:10px; padding-left: 50px;">
		<div style="display:flex; align-items:center;">
			<img src="https://iili.io/6JIJvR.png" width="50" height="50" style="max-width:50px"/>
			<span style="align-self: center; font-weight: 600; font-size: large;">Story Craft</span>
		</div>
		<div style="display:flex; margin-top: 200px; flex-direction: column; margin-bottom: 200px;">
			<h1 style="font-size:x-large">This is your code</h1>
			<div style="display:flex; background-color: lightgray; flex-direction: column; border-radius: 10px; padding-left: 20px; width:30%; align-items: center;">
				<p style="font-size:large">`+user.Code.Value+`</p>
			</div>
		</div>
	</div>
	<div style=" padding:4px; background-color: lightgray; padding-top:20px; overflow-x: hidden; font-family: sans-serif;">
		<div style="display:flex; justify-content: space-between;">
			<div style="display:flex; justify-content: space-evenly; width:30%;">
				<div>
					<h2 style=" font-size:small; font-weight:700;">Resources</h2>
					<ul style="color:gray; list-style:none">
						<li style="padding-bottom:10px">
							<a href="http://www.skillingsociety.com/contactus" style="text-decoration: none; color:black;">Contact us</a>
						</li>
						<li style="padding-bottom:10px">
							<a href="http://www.skillingsociety.com/" style="text-decoration: none; color:black;">Careers</a>
						</li>
					</ul>
				</div>
				<div>
					<h2 style=" font-size:small; font-weight:700;">Follow Us</h2>
					<ul style="color:gray; list-style:none">
						<li style="padding-bottom:10px">
							<a href="https://twitter.com/boilerplate" style="text-decoration: none; color:black;">Twitter</a>
						</li>
						<li style="padding-bottom:10px">
							<a href="https://www.linkedin.com/company/skillingsociety/" style="text-decoration: none; color:black;">Linked In</a>
						</li>
					</ul>
				</div>

			</div>
		</div>
		<hr style=" margin-top:10px; margin-bottom:10px; border-color: gray;"/>
		<div style="display:flex; align-items: center; justify-content: space-between;">
			<span style="font-size:small; color:gray; text-align:center">© 2023 boilerplate™. All Rights Reserved.
			</span>
		</div>
	</div>
</div>`)

	return nil
}
