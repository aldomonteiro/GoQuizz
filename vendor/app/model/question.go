package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"app/shared/database"

	bolt "github.com/coreos/bbolt"
	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Question
// *****************************************************************************

// Question table contains the information for each question
type Question struct {
	ObjectID   bson.ObjectId `bson:"_id"`
	ID         uint32        `db:"id" bson:"id,omitempty"`
	Content    string        `db:"content" bson:"content"`
	CorrectMsg string        `db:"correctmsg" bson:"correctmsg"`
	WrongMsg   string        `db:"wrongmsg" bson:"wrongmsg"`
	UserID     bson.ObjectId `bson:"user_id"`
	UID        uint32        `db:"user_id" bson:"userid,omitempty"`
	CreatedAt  time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted    uint8         `db:"deleted" bson:"deleted"`
}

// QuestionID returns the question id
func (u *Question) QuestionID() string {
	r := ""

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		r = fmt.Sprintf("%v", u.ID)
	case database.TypeMongoDB:
		r = u.ObjectID.Hex()
	case database.TypeBolt:
		r = u.ObjectID.Hex()
	}

	return r
}

// QuestionByID gets question by ID
func QuestionByID(userID string, questionID string) (Question, error) {
	var err error
	var result Question

	switch database.ReadConfig().Type {
	case database.TypeBolt:
		err = database.View("question", userID+questionID, &result)
		if err != nil {
			result = Question{}
			err = ErrNoResult
		} else if result.UserID != bson.ObjectIdHex(userID) {
			result = Question{}
			err = ErrUnauthorized
		}
	default:
		err = ErrCode
	}

	return result, standardizeError(err)
}

// QuestionsByUserID gets all questions for a user
func QuestionsByUserID(userID string) ([]Question, error) {
	var err error

	var result []Question

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		err = database.SQL.Select(&result, "SELECT id, content, user_id, created_at, updated_at, deleted FROM note WHERE user_id = ?", userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			// Create a copy of mongo
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("question")

			// Validate the object id
			if bson.IsObjectIdHex(userID) {
				err = c.Find(bson.M{"user_id": bson.ObjectIdHex(userID)}).All(&result)
			} else {
				err = ErrNoResult
			}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		// View retrieves a record set in Bolt
		err = database.BoltDB.View(func(tx *bolt.Tx) error {
			// Get the bucket
			b := tx.Bucket([]byte("question"))
			if b == nil {
				return bolt.ErrBucketNotFound
			}

			// Get the iterator
			c := b.Cursor()

			prefix := []byte(userID)
			for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
				var single Question

				// Decode the record
				err := json.Unmarshal(v, &single)
				if err != nil {
					log.Println(err)
					continue
				}

				result = append(result, single)
			}

			return nil
		})
	default:
		err = ErrCode
	}

	return result, standardizeError(err)
}

// QuestionsList gets all questions
func QuestionsList() ([]Question, error) {
	var err error

	var result []Question

	switch database.ReadConfig().Type {
	case database.TypeBolt:
		// View retrieves a record set in Bolt
		err = database.BoltDB.View(func(tx *bolt.Tx) error {
			// Get the bucket
			b := tx.Bucket([]byte("question"))
			if b == nil {
				return bolt.ErrBucketNotFound
			}

			b.ForEach(func(k, v []byte) error {
				var single Question

				// Decode the record
				err := json.Unmarshal(v, &single)
				if err != nil {
					log.Println(err)
				} else {
					result = append(result, single)
				}
				return nil
			})
			return nil
		})
	default:
		err = ErrCode
	}

	return result, standardizeError(err)
}

// QuestionCreate creates a question
func QuestionCreate(content string, correctmsg string, wrongmsg string, userID string) error {
	var err error

	now := time.Now()

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.SQL.Exec("INSERT INTO note (content, user_id) VALUES (?,?)", content, userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			// Create a copy of mongo
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("question")

			question := &Question{
				ObjectID:   bson.NewObjectId(),
				Content:    content,
				CorrectMsg: correctmsg,
				WrongMsg:   wrongmsg,
				UserID:     bson.ObjectIdHex(userID),
				CreatedAt:  now,
				UpdatedAt:  now,
				Deleted:    0,
			}
			err = c.Insert(question)
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		question := &Question{
			ObjectID:   bson.NewObjectId(),
			Content:    content,
			CorrectMsg: correctmsg,
			WrongMsg:   wrongmsg,
			UserID:     bson.ObjectIdHex(userID),
			CreatedAt:  now,
			UpdatedAt:  now,
			Deleted:    0,
		}

		err = database.Update("question", userID+question.ObjectID.Hex(), &question)
	default:
		err = ErrCode
	}

	return standardizeError(err)
}

// QuestionUpdate updates a question
func QuestionUpdate(content string, correctmsg string, wrongmsg string, userID string, questionID string) error {
	var err error

	now := time.Now()

	switch database.ReadConfig().Type {
	case database.TypeBolt:
		var question Question
		question, err = QuestionByID(userID, questionID)
		if err == nil {
			// Confirm the owner is attempting to modify the question
			if question.UserID.Hex() == userID {
				question.UpdatedAt = now
				question.Content = content
				question.CorrectMsg = correctmsg
				question.WrongMsg = wrongmsg
				err = database.Update("question", userID+question.ObjectID.Hex(), &question)
			} else {
				err = ErrUnauthorized
			}
		}
	default:
		err = ErrCode
	}

	return standardizeError(err)
}

// QuestionDelete deletes a question
func QuestionDelete(userID string, questionID string) error {
	var err error

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.SQL.Exec("DELETE FROM note WHERE id = ? AND user_id = ?", questionID, userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			// Create a copy of mongo
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("question")

			var question Question
			question, err = QuestionByID(userID, questionID)
			if err == nil {
				// Confirm the owner is attempting to modify the question
				if question.UserID.Hex() == userID {
					err = c.RemoveId(bson.ObjectIdHex(questionID))
				} else {
					err = ErrUnauthorized
				}
			}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		var question Question
		question, err = QuestionByID(userID, questionID)
		if err == nil {
			// Confirm the owner is attempting to modify the question
			if question.UserID.Hex() == userID {
				err = database.Delete("question", userID+question.ObjectID.Hex())
			} else {
				err = ErrUnauthorized
			}
		}
	default:
		err = ErrCode
	}

	return standardizeError(err)
}
