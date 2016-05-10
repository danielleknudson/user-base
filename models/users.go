package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/gorp.v1"
)

type User struct {
	Id         int64      `json:"id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	CreatedAt  *time.Time `json:"created_at"`
	ModifiedAt *time.Time `json:"modified_at"`
}

// A PublicUser is a user without a password
// Used prevent clients from accessing a user's password hash
type PublicUser struct {
	*User
	Password bool `json:"password,omitempty"`
}

type Users []*User

type PublicUsers []*PublicUser

func NewUser() *User {
	user := new(User)
	return user
}

func NewPublicUser(user *User) *PublicUser {
	publicUser := PublicUser{
		User: user,
	}

	return &publicUser
}

func NewUsers() Users {
	return make([]*User, 0)
}

// PreUpdate is called before inserting rows into the database
// Assigns values to the created_at and modified_at columns
func (user *User) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now()
	user.CreatedAt = &now
	user.ModifiedAt = user.CreatedAt
	return nil
}

// PreUpdate is called before updating rows in the database
// Updates the modified_at timestamp for a user
func (user *User) PreUpdate(s gorp.SqlExecutor) error {
	now := time.Now()
	user.ModifiedAt = &now
	return nil
}

func (user *User) Save() error {
	err := dbmap.Insert(user)
	return err
}

func (user *User) Fetch() error {
	err := dbmap.SelectOne(&user, "select * from users where id=$1", user.Id)
	return err
}

// Fetches all users
// TODO: Implement pagination
func (users Users) FetchAll() (PublicUsers, error) {
	var publicUsers PublicUsers
	_, err := dbmap.Select(&users, "SELECT * FROM users ORDER BY id")

	for _, user := range users {
		publicUsers = append(publicUsers, NewPublicUser(user))
	}

	return publicUsers, err
}

// Updates a user's attributes
// The PreUpdate hook will update the modified_at timestamp for the user
func (user *User) Update() error {
	_, err := dbmap.Update(user)
	return err
}

func (user *User) Destroy() error {
	_, err := dbmap.Exec("DELETE FROM users WHERE id=$1", user.Id)
	return err
}

// Looks up a user via email address
// Used during login
func (user *User) FindByEmail() error {
	err := dbmap.SelectOne(&user, "select * from users where email=$1", user.Email)
	return err
}

// Compares the user's existing password with the user-provided password
func ComparePasswords(inputPassword string, savedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(inputPassword))

	if err != nil {
		panic(err.Error())
	}

	return err == nil
}

// Hashes the user's password before the user is saved to the database
func HashPassword(password string) (string, error) {
	pw := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	return string(hashedPassword), err
}
