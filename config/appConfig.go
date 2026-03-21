package config

import ( "errors"
      "github.com/joho/godotenv"
      "os"
)

type AppConfig struct{
	ServerPort string
	Dsn string
	AppSecret string
	TwilioAccountSid string
	TwilioAuthToken string
	TwilioFromPhoneNumber string
}

func SetupEnv() (cfg AppConfig, err error){

//	if os.Getenv("APP_ENV") == "dev"{
		godotenv.Load()
	//}

	httpPort := os.Getenv("HTTP_PORT")
	Dsn:= os.Getenv("DSN")
	appSecret := os.Getenv("APP_SECRET")
	 
	if len(httpPort) < 1{
		return AppConfig{}, errors.New("env variable not found")
	}

	if len(Dsn)<1{
		return AppConfig{}, errors.New("env variable not found")
	}
	if len(appSecret)<1{
		return AppConfig{}, errors.New("app secret not found in the env variable")
	}
	 
	return AppConfig{ServerPort: httpPort, Dsn:Dsn, AppSecret: appSecret,
		TwilioAccountSid: os.Getenv("TWILIO_ACCOUNT_SID"),
		TwilioAuthToken: os.Getenv("TWILIO_AUTH_TOKEN"),
		TwilioFromPhoneNumber: os.Getenv("TWILIO_FROM_PHONE_NUMBER"),

		},nil
}