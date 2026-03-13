package handlers

import (
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

	code, err := h.svc.GetVerificationCode(user)
	if err !=nil{
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message":"Unable to generate verification code",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"get verification code",
		"data":code,
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