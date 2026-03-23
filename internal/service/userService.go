package service

import (
	"errors"
	"fmt"
	"time"
	"github.com/Girmex/go-ecommerce/config"
	"github.com/Girmex/go-ecommerce/internal/domain"
	"github.com/Girmex/go-ecommerce/internal/dto"
	"github.com/Girmex/go-ecommerce/internal/helper"
	"github.com/Girmex/go-ecommerce/internal/repository"
	"github.com/Girmex/go-ecommerce/pkg/notification"
)

type UserService struct{
	Repo repository.UserRepository
	Auth helper.Auth
	Config config.AppConfig

}

func (s UserService) Signup(input dto.UserSignup) (string, error){

   hPassword, err:= s.Auth.CreateHashedPassword(input.Password)

   if err != nil {
	return "", err
   }
	user, err := s.Repo.CreateUser(domain.User{
		Email: input.Email,
		Phone: input.Phone,
		Password: hPassword,
	})
	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s UserService) Login(email string, password string) (string, error){

	user, err := s.findUserByEmail(email)
	// compare password and generate token
	if err != nil{
		return "", err
	}
	err = s.Auth.VerifyPassword(password, user.Password)
	if err != nil{
		return "", err
	}
	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s UserService) findUserByEmail(email string) (*domain.User, error){
	user, err :=s.Repo.FindUser(email)
	return &user, err
}

func(s UserService) IsVerifiedUser(id uint) bool{
	currentUser, err:= s.Repo.FindUserById(id)

	return err !=nil && currentUser.Verified
}

func (s UserService) GetVerificationCode(e domain.User) error{
	//check if the user is already verified
 if s.IsVerifiedUser(e.ID){
	return nil
 }


	//generate verification

	code, err :=  s.Auth.GenerateCode()
	if err != nil{
		return nil
	}
	user := domain.User{
        Expiry: time.Now().Add(30 * time.Minute),
		Code: code,
	}



	//update the user to verified

  _, err = s.Repo.UpdateUser(e.ID, user)
  if err !=nil{
	return errors.New("unable to update verification code")
  }

  user, _ = s.Repo.FindUserById(e.ID)
// send SMS
msg :=fmt.Sprintf( "Your verification code is: %v " ,code)

notificationClient := notification.NewNotificationClient(s.Config)
 err = notificationClient.SendSMS(user.Phone, msg)
 if err !=nil{
	return errors.New("unable to update verification code")
  }
	//return the verification code

	return nil
}
func (s UserService) VerifyCode(id uint, code int) error{

	if s.IsVerifiedUser(id){
		return errors.New("user is already verified!")
	}

	user, err := s.Repo.FindUserById(id)

	if err!= nil{
		return err
	}

	if user.Code != code {
		return errors.New("Verification code does not match!")
	}

	if !time.Now().Before(user.Expiry){
		return errors.New("Verification code expired")

	}

	updateUser := domain.User{
		Verified:true,
	}
	 _, err = s.Repo.UpdateUser(id, updateUser)
	 if err !=nil{
		return errors.New("unable to verify user!")
	 }

	return nil
}

func (s UserService) CreateProfile(input any) error{

	return nil
}

func (s UserService) GetProfile(id uint) (*domain.User, error){

	return nil,nil
}

func (s UserService) UpdateProfile(id uint, input any) error{

	return nil
}

func (s UserService) BecomeSeller(id uint, input dto.SellerInput) (string, error) {
	// find existing user
	user, _ := s.Repo.FindUserById(id)

	if user.UserType == domain.SELLER {
		return "", errors.New("you have already joined seller program")
	}

	// update user
	seller, err := s.Repo.UpdateUser(id, domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.PhoneNumber,
		UserType:  domain.SELLER,
	})

	if err != nil {
		return "", err
	}

	// generating token
	token, err := s.Auth.GenerateToken(user.ID, user.Email, seller.UserType)

	// create bank account information

	err = s.Repo.CreateBankAccount(domain.BankAccount{
		BankAccount: input.BankAccountNumber,
		SwiftCode:   input.SwiftCode,
		PaymentType: input.PaymentType,
		UserId:      id,
	})

	return token, err
}

func (s UserService) FindCart(id uint) ([] interface{},error){

	return nil, nil
}

func (s UserService) Createcart(input any, u domain.User) ([] interface{},error){

	return nil, nil
}

func (s UserService) CreateOrder(u  domain.User) (int,error){

	return 0, nil
}

func (s UserService) GetOrders(u domain.User) ([] interface{},error){

	return nil, nil
}

func (s UserService) GetOrderById(id uint, uId uint) ([] interface{},error){

	return nil, nil
}


