package database

import "database/sql"

type Context struct {
	UserRepository *sql.DB
}

func InitDB() Context {
	context := Context{}
	var userRepositoryError error

	context.UserRepository, userRepositoryError = initUserDB()
	if userRepositoryError != nil {
		panic(userRepositoryError)
	}

	return context
}
