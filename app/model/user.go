package model

import (
	"time"

	"GoQuizz/app/shared/database"

	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
	ObjectID  bson.ObjectId `storm:"id"`
	FirstName string
	LastName  string
	Email     string
	Password  string
	StatusID  uint8
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   uint8
}

// UserStatus table contains every possible user status (active/inactive)
type UserStatus struct {
	ID        uint8
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   uint8
}

// UserID returns the user id
func (u *User) UserID() string {
	r := ""

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
	case database.TypeMongoDB:
		r = u.ObjectID.Hex()
	case database.TypeBolt:
		r = u.ObjectID.Hex()
	}

	return r
}

// UserByEmail gets user information from email
func UserByEmail(email string) (User, error) {
	var err error

	result := User{}

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		err = database.SQL.Get(&result, "SELECT id, password, status_id, first_name FROM user WHERE email = ? LIMIT 1", email)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
			err = c.Find(bson.M{"email": email}).One(&result)
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		err = database.BoltDB.One("Email", email, &result)
		if err != nil {
			err = ErrNoResult
		}
	default:
		err = ErrCode
	}

	return result, standardizeError(err)
}

// UserCreate creates user
func UserCreate(firstName, lastName, email, password string) error {
	var err error

	now := time.Now()

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.SQL.Exec("INSERT INTO user (first_name, last_name, email, password) VALUES (?,?,?,?)", firstName,
			lastName, email, password)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("user")

			user := &User{
				ObjectID:  bson.NewObjectId(),
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Password:  password,
				StatusID:  1,
				CreatedAt: now,
				UpdatedAt: now,
				Deleted:   0,
			}
			err = c.Insert(user)
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		user := &User{
			ObjectID:  bson.NewObjectId(),
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
			StatusID:  1,
			CreatedAt: now,
			UpdatedAt: now,
			Deleted:   0,
		}

		err = database.BoltDB.Save(&user)
	default:
		err = ErrCode
	}

	return standardizeError(err)
}
