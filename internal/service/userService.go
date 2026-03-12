package service

import (
	"github.com/Girmex/go-ecommerce/internal/domain"
	"github.com/Girmex/go-ecommerce/internal/dto"
	"github.com/Girmex/go-ecommerce/internal/helper"
	"github.com/Girmex/go-ecommerce/internal/repository"
)

type UserService struct{
	Repo repository.UserRepository
	Auth helper.Auth

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

func (s UserService) GetVerificationCode(e domain.User) (int, error){
	//check if the user is already verified
 if s.IsVerifiedUser(e.ID){
	return 0, nil
 }


	//generate verification

	//update the user to verified

	//return the verification code

	return 0, nil
}
func (s UserService) VerifyCode(id uint, code int) error{

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

func (s UserService) BecomeSeller(id uint, input any) (string,error){

	return "", nil
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


