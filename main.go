package main

import (
	"fmt"
	"learning-go/db"
	"learning-go/env"
	"learning-go/wealthsimple"
)

func main() {
	err := env.InitEnv()
	if err != nil {
		panic(err)
	}

	database := db.Database{}
	err = database.Connect(
		env.GetEnv("DB_USERNAME"),
		env.GetEnv("DB_PASSWORD"),
		env.GetEnv("DB_HOST"),
		env.GetEnv("DB_NAME"),
	)
	if err != nil {
		panic(err)
	}

	authTokenDb := db.NewAuthTokenDB(&database)

	email := env.GetEnv("EMAIL")
	password := env.GetEnv("PASSWORD")

	_, err = authTokenDb.FetchAuthToken(email)
	if err != nil {
		fmt.Println("No token found, need to get a new one")

		wsLogin := wealthsimple.WealthSimpleLogin{Email: email, Password: password}
		err := wsLogin.InitOtpClaim()
		if err != nil {
			panic(err)
		}

	}

}