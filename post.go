package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
  "github.com/go-ozzo/ozzo-validation"
	"github.com/julienschmidt/httprouter"
)

type PostInput struct{
	Token *string `json:"token"`
	Description *string `json:"description"`
}

type Post struct {
  PostId *int64 `json:"postId"`
  PostDescription *string `json:"postDescription"`
  NumberOfLikes *int64 `json:"numberOfLikes"`
  NumberOfComments *int64 `json:"numberOfComments"`
  CreatedTime *int64 `json:"createdTime"`
}

func (a PostInput) Validate() error {
	return validation.ValidateStruct(&a,
		// Name cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Description, validation.Required),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Token, validation.Required),
	)
}

func createPost(userId int64,description string)(*int64,error) {

	//insert post record
	transaction, _ := db.Begin()

	stringTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
	stringUserId := strconv.FormatInt(userId, 10)

  print(stringUserId,"\n",stringTimestamp,"\n",description)

	res, err := db.Exec("INSERT INTO posts SET post_description = '" + description + "', user_id = " +stringUserId +", number_of_likes = 0, number_of_comments=0, created_date = " + stringTimestamp)

	if err != nil {
    print(err.Error())
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
	transaction.Commit()
	return &lastId,nil
}

func getPostsListForUser(userId int64) ([]Post, *HttpError) {


  stringUserId := strconv.FormatInt(userId,10)
	var posts []Post
	rows, err := db.Query("SELECT id,post_description,number_of_likes,number_of_comments,created_date from posts where user_id="+stringUserId)
	if err != nil {
		Log.Println("Error while fetching posts", err.Error())
		return nil, NewErrorInternalServerError(err)
	}
	defer rows.Close()

	var (
		postId      *int64
		description    *string
		nooflikes *int64
		noofcomments   *int64
    createdDate *int64
	)
	for rows.Next() {
		err := rows.Scan(&postId, &description, &nooflikes, &noofcomments,&createdDate)
		if err != nil {
			Log.Println("Failed to scan prize list", err.Error())
			return nil, NewErrorInternalServerError(err)
		}
		posts = append(posts, Post{postId, description,nooflikes,noofcomments, createdDate})
	}

	err = rows.Err()

	if err != nil {
		Log.Println("Error while itterating over rows", err.Error())
		return nil, NewErrorInternalServerError(err)
	}

	return posts, nil

}


func createPostHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) *HttpError {

	post := new(PostInput)

	if err := json.NewDecoder(r.Body).Decode(post); err != nil {
		return NewErrorBadRequestError(fmt.Sprintf("Error decoding %s", err))
	}

	err := post.Validate()

	if err != nil {
		return NewErrorBadRequestError(err.Error())
	}

	print(*post.Token,*post.Description)

	userId,err := getUserByToken(*post.Token)

	if err!= nil {
		Log.Println(err.Error())
		return NewErrorInternalServerError(err)
	}

	if userId == nil{
		return NewErrorBadRequestError("Invalid Credentails")
	}

  print("hoisfjvvvvffffffffffffffffffffffffffffffffffffuserId",*userId)

	postId,err := createPost(*userId,*post.Description)
	if err!= nil {
		Log.Println(err.Error())
		return NewErrorInternalServerError(err)
	}

	type Post struct{
		PostId *int64 `json:"postId"`
	}

	type Output struct{
		Data Post `json:"data"`
	}

	var output Output
	output.Data.PostId = postId

	b, err := json.Marshal(output)
	if err != nil {
		Log.Println("Error while marshalling")
		return NewErrorInternalServerError(err)
	}

	PrintSuccessJson(w, b)
	return nil

}

func getPostsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) *HttpError {
  queryValues := r.URL.Query()

  stringuserId := queryValues.Get("userId")
  if stringuserId == "" {
    return NewErrorBadRequestError("No userId in query parameters")
  }

  userId, err := strconv.ParseInt(stringuserId, 10, 64)
  if err != nil {
      return NewErrorBadRequestError("userId is not of type integer")
  }

  posts, httpErr := getPostsListForUser(userId)

	if httpErr != nil {
		return httpErr
	}

	type Output struct {
		Data []Post `json:"data"`
	}

	var output Output
	output.Data = posts

	j, err := json.Marshal(output)

	if err != nil {
		Log.Println("Error while marshalling", err.Error())
		return NewErrorInternalServerError(err)
	}

	PrintSuccessJson(w, j)
	return nil

}
