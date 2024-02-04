package services

import (
	"fmt"

	"github.com/letitloose/nsdtr-club-us/internal/models"
	"github.com/letitloose/nsdtr-club-us/internal/validator"
)

type UserService struct {
	*models.UserModel
	*Email
}

type UserForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (service *UserService) InsertUser(uf *UserForm) error {

	// Validate the form contents using our helper functions.
	uf.CheckField(validator.NotBlank(uf.Email), "email", "This field cannot be blank")
	uf.CheckField(validator.ValidEmail(uf.Email), "email", "You must enter a valid email: name@domain.ext")
	uf.CheckField(validator.NotBlank(uf.Password), "password", "This field cannot be blank")
	uf.CheckField(validator.MinChars(uf.Password, 8), "password", "This field must be at least 8 characters long")

	if !uf.Valid() {
		return models.ErrBadData
	}

	_, err := service.Insert(uf.Email, uf.Password)
	if err != nil {
		return err
	}

	//user created successfully,  send an email with the validation link
	verificationHash, err := service.GetVerificationHashByEmail(uf.Email)
	body := fmt.Sprintf(
		`<html>
			<body>
				<h1>Hello!</h1>
				<p>Please <a href="https://localhost:8080/user/activate?hash=%s">click here</a> to validate your email and activate your account.<p>
			</body>
		</html>`, verificationHash)
	err = service.SendEmail("Welcome to NSDTRC-USA Membership", "", body)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) AuthenticateUser(uf *UserForm) (int, error) {

	uf.CheckField(validator.NotBlank(uf.Email), "email", "Please enter your email to login")

	if !uf.Valid() {
		return 0, models.ErrBadData
	}

	id, err := service.Authenticate(uf.Email, uf.Password)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (service *UserService) ActivateUser(hash string) error {
	userID, err := service.GetByVerificationHash(hash)
	if err != nil {
		return err
	}

	return service.Activate(userID)
}
