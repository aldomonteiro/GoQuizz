package model

import (
	"time"

	"GoQuizz/app/shared/database"

	"github.com/asdine/storm/q"
	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Question
// *****************************************************************************

// Question table contains the information for each question
type Question struct {
	ObjectID   bson.ObjectId `storm:"id"`
	Content    string
	CorrectMsg string
	WrongMsg   string
	UserID     bson.ObjectId
	UID        uint32
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Deleted    uint8
}

// QuestionID returns the question id
func (u *Question) QuestionID() string {
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

// QuestionByID gets question by ID
func QuestionByID(userID string, questionID string) (Question, error) {
	var err error
	var result Question

	switch database.ReadConfig().Type {
	case database.TypeBolt:
		err = database.BoltDB.Select(q.And(q.Eq("ObjectID", bson.ObjectIdHex(questionID)), q.Eq("UserID", bson.ObjectIdHex(userID)))).Find(&result)
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
		err = database.BoltDB.Find("UserID", userID, &result)
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
		err = database.BoltDB.All(&result)
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

		err = database.BoltDB.Save(&question)
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
				err = database.BoltDB.Save(&question)
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
				err = database.BoltDB.DeleteStruct(&question)
			} else {
				err = ErrUnauthorized
			}
		}
	default:
		err = ErrCode
	}

	return standardizeError(err)
}
