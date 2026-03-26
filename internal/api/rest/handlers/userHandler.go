package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Girmex/go-ecommerce/internal/api/rest"
	"github.com/Girmex/go-ecommerce/internal/dto"
	"github.com/Girmex/go-ecommerce/internal/repository"
	"github.com/Girmex/go-ecommerce/internal/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	// user service
	svc service.UserService
}

func SetupUserRoutes(rh *rest.RestHandler) {
	app:= rh.App

	svc:=service.UserService{
		Repo: repository.NewUserRepository(rh.DB),
		Auth: rh.Auth,
		Config: rh.Config,
	}
	handler:=UserHandler{
		svc:svc,
	}

	pubRoutes:= app.Group("/users")
	pubRoutes.Post("/register", handler.Register)
	pubRoutes.Post("/login", handler.Login)

    pvtRoutes := app.Group("/",rh.Auth.Authorize)
	pvtRoutes.Post("/verify", handler.Verify)
	pvtRoutes.Get("/verify", handler.GetVerificationCode)

	pvtRoutes.Post("/cart", handler.AddToCart)
	pvtRoutes.Get("/cart", handler.GetCart)


}

func (h * UserHandler) Register(ctx *fiber.Ctx) error{

	user:=dto.UserSignup{}
	err:=ctx.BodyParser(&user)

	if err !=nil{
		return ctx.Status(http.StatusBadRequest).JSON(*&fiber.Map{
        "message":"Please provide valid input",
		})
	}
	token , err:= h.svc.Signup(user)

	if err !=nil{
		return ctx.Status(http.StatusInternalServerError).JSON(*&fiber.Map{
			"message":"Error on signup",
			})
	}
 return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"register",
		"token":token,
	})
}


func (h * UserHandler) Login(ctx *fiber.Ctx) error{
	loginInput:=dto.UserSignup{}
	err:=ctx.BodyParser(&loginInput)

	if err !=nil{
		return ctx.Status(http.StatusBadRequest).JSON(*&fiber.Map{
        "message":"Please provide valid user id password",
		})
	}
	token, err := h.svc.Login(loginInput.Email,loginInput.Password)
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"login",
	})
	if err!= nil{
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "please provide correct user password",
		})
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"login",
		"token":token,
	})
}

//



func (h * UserHandler) GetProfile(ctx *fiber.Ctx) error{
	user:= h.svc.Auth.GetCurrentUser(ctx)
	log.Println(user)
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"get profile",
		"user":user,
	})
}

func (h * UserHandler) GetVerificationCode(ctx *fiber.Ctx) error{

	user:= h.svc.Auth.GetCurrentUser(ctx)

	err := h.svc.GetVerificationCode(user)
	if err !=nil{
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message":"Unable to generate verification code",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"get verification code",
	})
}

func (h * UserHandler) Verify(ctx *fiber.Ctx) error{

	user := h.svc.Auth.GetCurrentUser(ctx)
     var req dto.VerificationCodeInput
	 if err:= ctx.BodyParser(&req); err !=nil{
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message":"Please provide valid input",
		})
	 }
	 err := h.svc.VerifyCode(user.ID,req.Code)

	 if err != nil{
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err,
		})
	 }

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"verified successfuly",
	})
}

func (h *UserHandler) BecomeSeller(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	req := dto.SellerInput{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(400).JSON(&fiber.Map{
			"message": "request parameters are not valid",
		})
	}

	token, err := h.svc.BecomeSeller(user.ID, req)

	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "fail to become seller",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "become seller",
		"token":   token,
	})
}

func (h *UserHandler) AddToCart(ctx *fiber.Ctx) error {

	req := dto.CreateCartRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide a valid product and qty",
		})
	}

	user := h.svc.Auth.GetCurrentUser(ctx)

	// call user service and perform create cart
	cartItems, err := h.svc.CreateCart(req, user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "cart created successfully", cartItems)

}
func (h *UserHandler) GetCart(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	cart, _, err := h.svc.FindCart(user.ID)
	if err != nil {
		return rest.InternalError(ctx, errors.New("cart does not exist"))
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "get cart",
		"cart":    cart,
	})
}