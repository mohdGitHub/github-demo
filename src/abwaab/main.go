package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/golang-jwt/jwt"
	_ "github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"io/ioutil"
	_ "io/ioutil"
	"log"
	"net/http"
	"net/mail"
	_ "sync"
	"time"
)

// User Defined Types
type User struct {
	gorm.Model
	Email string `json:"email"`
	Password string `json:"password"`
}

type AbwaabTweet struct {
	gorm.Model
	Description string `json:"description"`
}

type Credentials struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type DBParam struct {
	Limit string `json:"limit"`
	Offset string `json:"offset"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Variables to be used across all over REST
var (AbwaabTweets [] AbwaabTweet
	signedKey = "abwaab"
	privateKey []byte
	publicKey []byte
	db *gorm.DB
	client *twitter.Client
	err error
	rows *sql.Rows
	cookie *http.Cookie
	tkn *jwt.Token)

// Util
func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Init
func init(){
	initTwitter()
	initJWT()
	initDatabase()
}

func initTwitter(){
	config := oauth1.NewConfig("0I68KfhHKxQdA7wKinkESsFu3", "qegXI3TrBNqNsn0zgpWkjNs8tLYFbmrlEFxpgenYLhXYNJBCZ7")
	token := oauth1.NewToken("753752674810654721-IVYPe3GJQYHgYu43QcldIbp4npMGyAE", "1hgngP1Y5TnzJGVFjc9nlIXOUi49cV8W7h4Abrfofhmq0")
	httpClient := config.Client(oauth1.NoContext, token)
	client = twitter.NewClient(httpClient)
}

func initJWT(){
	privateKey, _ = ioutil.ReadFile("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiU2VuaW9yIFNvZnR3YXJlIEVuZ2luZWVyIiwibmFtZSI6IkpvaG4gRG9lIn0.KPq0QAfOAXF0W8qCceQUBz1JOAP4arMv3yHWl0CXk5Y")
}

func initDatabase(){
	db, err = gorm.Open( "postgres", "host=localhost port=5432 dbname=postgres sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
}

// Endpoints
func getTweets(w http.ResponseWriter, r *http.Request){
	if !checkTokenValidity(w,r){
		return
	}
	tweets, _, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: 5,
	})

	if err != nil {
		fmt.Println("Failed to execute getTweets: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, s := range tweets {
		var tweet AbwaabTweet = AbwaabTweet{
			Description: s.Text,
		}
		AbwaabTweets = append(AbwaabTweets,tweet)
	}
}

func createTweet(w http.ResponseWriter, r *http.Request) {
	if !checkTokenValidity(w,r){
		return
	}
	var abwaabTweet AbwaabTweet
	err := json.NewDecoder(r.Body).Decode(&abwaabTweet)
	if err != nil {
		fmt.Println("Failed to execute setTweet: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db.AutoMigrate(&AbwaabTweet{})
	db.Create(&abwaabTweet)
}

func saveCurrentFetchedTweets(w http.ResponseWriter, r *http.Request){
	if !checkTokenValidity(w,r){
		return
	}
	if len(AbwaabTweets) == 0 {
		fmt.Println("No AbwaabTweets To Save")
		w.WriteHeader(http.StatusAccepted)
		return
	}
	db.AutoMigrate(&AbwaabTweet{})
	for index := range AbwaabTweets {
		db.Create(&AbwaabTweets[index])
	}
	// Empty
	AbwaabTweets = nil
}

func navigateTweets(w http.ResponseWriter, r *http.Request){
	if !checkTokenValidity(w,r){
		return
	}
	// Reset list of results
	AbwaabTweets = nil
	var dbparam DBParam
	err := json.NewDecoder(r.Body).Decode(&dbparam)
	if err != nil {
		fmt.Println("Failed to execute navigation: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rows, err = db.DB().Query("SELECT Description FROM ABWAAB_TWEETS ABT LIMIT $1 OFFSET $2", dbparam.Limit, dbparam.Offset)
	if err != nil {
		fmt.Println("Failed to execute query: ", err)
	}

	for rows.Next() {
		var abwaabTweet AbwaabTweet
		rows.Scan(&abwaabTweet.Description)
		AbwaabTweets = append(AbwaabTweets, abwaabTweet)
	}

	json.NewEncoder(w).Encode(&AbwaabTweets)
}

func getUser(email string) User {
	var user User
	rows, err = db.DB().Query("SELECT email, password FROM USERS WHERE email=$1", email)
	if err != nil {
		fmt.Println("Failed to execute query: ", err)
	}
	for rows.Next() {
		rows.Scan(&user.Email, &user.Password)
		break
	}
	return user
}

func register(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		fmt.Println("Failed to execute register: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !valid(credentials.Email) {
		fmt.Println("Wrong email")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u := getUser(credentials.Email)

	if u.Email != "" {
		fmt.Println("User already exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.AutoMigrate(&User{})

	user := User {
		Email: credentials.Email,
		Password: credentials.Password,
	}

	db.Create(&user)
}
/**
Name: login
Purpose: Provide users of REST services the needed authentication
**/
func login(w http.ResponseWriter, r *http.Request){
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := getUser(credentials.Email)
	if &user == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Password != credentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expirationTime := time.Now().Add(time.Second * 15)
	claims := &Claims{
		Email: credentials.Email,
		StandardClaims : jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w,&http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires: expirationTime,
	})
}

func checkTokenValidity(w http.ResponseWriter, r *http.Request) bool{
	cookie, err = r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	tkn, err = jwt.ParseWithClaims(tokenStr,claims, func(token *jwt.Token) (interface{}, error) {
		return privateKey,nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}

func handleRequests(){
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/getTweets", getTweets)
	http.HandleFunc("/saveTweets", saveCurrentFetchedTweets)
	http.HandleFunc("/navigateTweets", navigateTweets)
	http.HandleFunc("/createTweet", createTweet)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main(){
	handleRequests()
}
