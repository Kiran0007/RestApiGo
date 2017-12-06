package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
  "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/julienschmidt/httprouter"
)

func (a SignUpInput) Validate() error {
	return validation.ValidateStruct(&a,
		// Name cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Name, validation.Required, validation.Length(5, 50)),
		// Email cannot be empty
		validation.Field(&a.Email, validation.Required, is.Email),
		// Password cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Password, validation.Required, validation.Length(5, 50)),
	)
}


type SignUpInput struct {
	Name     *string `json:"userName"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}



type User struct {
	UserId      *int64  `json:"userId"`
	UserName    *string `json:"userName"`
	Email       *string `json:"email"`
	CreatedDate *int64  `json:"createdDate"`
}

type UserSession struct {
	UserId int64  `json:"userId"`
	Name   string `json:"userName"`
	Email  string `json:"email"`
	Token string `json:"token"`
}

func checkIfEmailExist(emailId string) (bool, error) {

	stmt, err := db.Prepare("SELECT count(*) FROM users WHERE email = ?")

	defer stmt.Close()

	if err != nil {
		return false, err
	}

	var count uint64
	err = stmt.QueryRow(emailId).Scan(&count)
	if count != 0 {
		return true, nil
	}

	return false, nil
}

func createNewUser(u *SignUpInput) (*UserSession, error) {

	//insert user record
	transaction, _ := db.Begin()

	stringTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
	password := GetSHA1(*u.Password + PASSWORD_SALT)

	res, err := db.Exec("INSERT INTO users SET user_name = '" + *u.Name + "', email = '" + *u.Email + "', password = '" + password + "', created_date = " + stringTimestamp)

	if err != nil {
		Log.Println("transaction failed")
		transaction.Rollback()
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		Log.Println("failed while getting last insertid")
		transaction.Rollback()
		return nil, err
	}

	stmt, err := db.Prepare("INSERT INTO user_session SET user_id = ?, session_token = ?, login_time = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	token := RandToken()

	_, err = stmt.Exec(lastId, token, time.Now().Unix())
	if err != nil {
		return nil,err
	}

	transaction.Commit()

	userSession := new(UserSession)
	userSession.Email = *u.Email
	userSession.UserId = lastId
	userSession.Name = *u.Name
	userSession.Token = token

	return userSession, nil
}

func signUpHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) *HttpError {

	newUser := new(SignUpInput)

	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		return NewErrorBadRequestError(fmt.Sprintf("Error decoding %s", err))
	}

	err := newUser.Validate()

	if err != nil {
		return NewErrorBadRequestError(err.Error())
	}

	if booll, _ := checkIfEmailExist(*newUser.Email); booll {
		return NewErrorBadRequestError("User Already Exist with That email.")
	}

	session, err := createNewUser(newUser)

	if err != nil {
		Log.Println("Error while adding user to database " + err.Error())
		return NewErrorInternalServerError(err)
	}

	type Output struct {
		Data *UserSession `json:"data"`
	}

	output := new(Output)

	output.Data = session

	b, err := json.Marshal(output)
	if err != nil {
		Log.Println("Error while marshalling")
		return NewErrorInternalServerError(err)
	}

	PrintSuccessJson(w, b)
	return nil
}

func getUserByToken(token string) (*int64,error){
	rows, err := db.Query("SELECT user_id from user_session where session_token = '" +token+"'" )
	defer rows.Close()

	if err != nil {
		print("Error ",err.Error())
		return nil, err
	}

	var userId *int64

	if rows.Next() != true {
		return nil, nil
	} else {
		err := rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		return userId,nil
	}
}
