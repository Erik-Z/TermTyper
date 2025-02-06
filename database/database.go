package database

import "database/sql"

type Context struct {
	userRepository *sql.DB
}

func InitDB() Context {
	context := Context{}
	var userRepositoryError error

	context.userRepository, userRepositoryError = initUserDB()
	if userRepositoryError != nil {
		panic(userRepositoryError)
	}

	return context
}
