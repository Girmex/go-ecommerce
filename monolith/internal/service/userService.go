package service

import (
	"errors"
	"fmt"
	"log"
	"time"
	"github.com/Girmex/go-ecommerce/monolith/config"
	"github.com/Girmex/go-ecommerce/monolith/internal/domain"
	"github.com/Girmex/go-ecommerce/monolith/internal/dto"
	"github.com/Girmex/go-ecommerce/monolith/internal/helper"
	"github.com/Girmex/go-ecommerce/monolith/internal/repository"
	"github.com/Girmex/go-ecommerce/monolith/pkg/notification"
)

type UserService struct{
	Repo   repository.UserRepository
	CRepo  repository.CatalogRepository
	Auth   helper.Auth
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
func (s UserService) VerifyCode(id uint, code string) error{

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

func (s UserService) CreateProfile(id uint, input dto.ProfileInput) error {

	// update user
	user, err := s.Repo.FindUserById(id)

	if err != nil {
		return err
	}
	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}

	_, err = s.Repo.UpdateUser(id, user)

	if err != nil {
		return err
	}

	// create address
	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostCode:     input.AddressInput.PostCode,
		UserId:       id,
	}

	err = s.Repo.CreateProfile(address)
	if err != nil {
		return err
	}

	return nil
}

func (s UserService) GetProfile(id uint) (*domain.User, error) {

	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s UserService) UpdateProfile(id uint, input dto.ProfileInput) error {

	// find the user
	user, err := s.Repo.FindUserById(id)

	if err != nil {
		return err
	}
	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}

	_, err = s.Repo.UpdateUser(id, user)
	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostCode:     input.AddressInput.PostCode,
		UserId:       id,
	}

	err = s.Repo.UpdateProfile(address)
	if err != nil {
		return err
	}
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

func (s UserService) FindCart(id uint) ([]domain.Cart, float64, error) {

	cartItems, err := s.Repo.FindCartItems(id)

	if err != nil {
		return nil, 0, errors.New("error on finding cart items")
	}

	var totalAmount float64

	for _, item := range cartItems {
		totalAmount += item.Price * float64(item.Qty)
	}

	return cartItems, totalAmount, err
}

func (s UserService) CreateCart(input dto.CreateCartRequest, u domain.User) ([]domain.Cart, error) {
	// check if the cart is Exist
	cart, _ := s.Repo.FindCartItem(u.ID, input.ProductId)

	if cart.ID > 0 {
		if input.ProductId == 0 {
			return nil, errors.New("please provide a valid product id")
		}
		//  => delete the cart item
		if input.Qty < 1 {
			err := s.Repo.DeleteCartById(cart.ID)
			if err != nil {
				log.Printf("Error on deleting cart item %v", err)
				return nil, errors.New("error on deleting cart item")
			}
		} else {
			//  => update the cart item
			cart.Qty = input.Qty
			err := s.Repo.UpdateCart(cart)
			if err != nil {
				// log error
				return nil, errors.New("error on updating cart item")
			}
		}

	} else {
		// check if product exist
		product, _ := s.CRepo.FindProductById(int(input.ProductId))
		if product.ID < 1 {
			return nil, errors.New("product not found to create cart item")
		}
		// create cart

		err := s.Repo.CreateCart(domain.Cart{
			UserId:    u.ID,
			ProductId: input.ProductId,
			Name:      product.Name,
			ImageUrl:  product.ImageUrl,
			Qty:       input.Qty,
			Price:     product.Price,
			SellerId:  uint(product.UserId),
		})

		if err != nil {
			return nil, errors.New("error on creating cart item")
		}
	}

	return s.Repo.FindCartItems(u.ID)

}

func (s UserService) CreateOrder(uId uint, orderRef string, pId string, amount float64) error {

	// find cart items for the user
	cartItems, _, err := s.FindCart(uId)
	if err != nil {
		return errors.New("error on finding cart items")
	}

	if len(cartItems) == 0 {
		return errors.New("cart is empty cannot create the order")
	}

	// create order with generated OrderNumber
	var orderItems []domain.OrderItem

	for _, item := range cartItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductId: item.ProductId,
			Qty:       item.Qty,
			Price:     item.Price,
			Name:      item.Name,
			ImageUrl:  item.ImageUrl,
			SellerId:  item.SellerId,
		})
	}

	order := domain.Order{
		UserId:         uId,
		PaymentId:      pId,
		OrderRefNumber: orderRef,
		Amount:         amount,
		Items:          orderItems,
	}

	err = s.Repo.CreateOrder(order)
	if err != nil {
		return err
	}
	// send email to user with order details

	// remove cart items from the cart
	err = s.Repo.DeleteCartItems(uId)
	log.Printf("Deleting cart items Error %v", err)

	// return order number
	return err
}

func (s UserService) GetOrders(u domain.User) ([]domain.Order, error) {
	orders, err := s.Repo.FindOrders(u.ID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s UserService) GetOrderById(id uint, uId uint) (domain.Order, error) {
	order, err := s.Repo.FindOrderById(id, uId)
	if err != nil {
		return order, err
	}
	return order, nil
}


