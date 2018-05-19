package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"app/shared/database"

	"github.com/boltdb/bolt"
	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Answer
// *****************************************************************************

// Answer table contains the information for each answer
type Answer struct {
	ObjectID   bson.ObjectId `bson:"_id"`
	ID         uint32        `db:"id" bson:"id,omitempty"` // Don't use Id, use AnswerID() instead for consistency with MongoDB
	Content    string        `db:"content" bson:"content"`
	QuestionID bson.ObjectId `bson:"question_id"`
	QID        uint32        `db:"question_id" bson:"questionid,omitempty"`
	IsCorrect  bool          `db:"iscorrect" bson:"iscorrect"`
	CreatedAt  time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted    uint8         `db:"deleted" bson:"deleted"`
}

// AnswerID returns the answer id
func (q *Answer) AnswerID() string {
	r := ""

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		r = fmt.Sprintf("%v", q.ID)
	case database.TypeMongoDB:
		r = q.ObjectID.Hex()
	case database.TypeBolt:
		r = q.ObjectID.Hex()
	}

	return r
}

// QuestionID returns the question id
func (q *Answer) MyQuestionID() string {
	r := ""

	switch database.ReadConfig().Type {
	case database.TypeBolt:
		r = q.QuestionID.Hex()
	}

	return r
}

// QuestionID returns the question id
func (a *Answer) IsCorrectAnswer() string {
	r := ""
	switch database.ReadConfig().Type {
	case database.TypeBolt:
		if a.IsCorrect {
			r = "Yes"
		} else {
			r = "No"
		}
	}
	return r
}

// AnswerByID gets answer by ID
func AnswerByID(questionID string, answerID string) (Answer, error) {
	var err error

	result := Answer{}

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		// err = database.SQL.Get(&result, "SELECT id, content, user_id, created_at, updated_at, deleted FROM note WHERE id = ? AND user_id = ? LIMIT 1", answerID, userID)
	case database.TypeMongoDB:
		// if database.CheckConnection() {
		// 	// Create a copy of mongo
		// 	session := database.Mongo.Copy()
		// 	defer session.Close()
		// 	c := session.DB(database.ReadConfig().MongoDB.Database).C("answer")

		// 	// Validate the object id
		// 	if bson.IsObjectIdHex(answerID) {
		// 		err = c.FindId(bson.ObjectIdHex(answerID)).One(&result)
		// 		if result.UserID != bson.ObjectIdHex(userID) {
		// 			result = Answer{}
		// 			err = ErrUnauthorized
		// 		}
		// 	} else {
		// 		err = ErrNoResult
		// 	}
		// } else {
		// 	err = ErrUnavailable
		// }
	case database.TypeBolt:
		err = database.View("answer", questionID+answerID, &result)
		if err != nil {
			err = ErrNoResult
		}
		if result.QuestionID != bson.ObjectIdHex(questionID) {
			result = Answer{}
			err = ErrUnauthorized
		}
	default:
		err = ErrCode
	}

	return result, standardizeError(err)
}

// AnswersByQuestionID gets all answers for a question
func AnswersByQuestionID(questionID string) ([]Answer, error) {
	var err error

	var result []Answer

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		// err = database.SQL.Select(&result, "SELECT id, content, user_id, created_at, updated_at, deleted FROM note WHERE user_id = ?", questionID)
	case database.TypeMongoDB:
		// if database.CheckConnection() {
		// 	// Create a copy of mongo
		// 	session := database.Mongo.Copy()
		// 	defer session.Close()
		// 	c := session.DB(database.ReadConfig().MongoDB.Database).C("answer")

		// 	// Validate the object id
		// 	if bson.IsObjectIdHex(questionID) {
		// 		err = c.Find(bson.M{"user_id": bson.ObjectIdHex(questionID)}).All(&result)
		// 	} else {
		// 		err = ErrNoResult
		// 	}
		// } else {
		// 	err = ErrUnavailable
		// }
	case database.TypeBolt:
		// View retrieves a record set in Bolt
		err = database.BoltDB.View(func(tx *bolt.Tx) error {
			// Get the bucket
			b := tx.Bucket([]byte("answer"))
			if b == nil {
				return bolt.ErrBucketNotFound
			}

			// Get the iterator
			c := b.Cursor()

			prefix := []byte(questionID)
			for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
				var single Answer

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

// AnswerCreate creates a answer
func AnswerCreate(questionID string, content string, iscorrect bool) error {
	var err error

	now := time.Now()

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		// _, err = database.SQL.Exec("INSERT INTO note (content, user_id) VALUES (?,?)", content, userID)
	case database.TypeMongoDB:
		// if database.CheckConnection() {
		// 	// Create a copy of mongo
		// 	session := database.Mongo.Copy()
		// 	defer session.Close()
		// 	c := session.DB(database.ReadConfig().MongoDB.Database).C("answer")

		// 	answer := &Answer{
		// 		ObjectID:  bson.NewObjectId(),
		// 		Content:   content,
		// 		UserID:    bson.ObjectIdHex(userID),
		// 		CreatedAt: now,
		// 		UpdatedAt: now,
		// 		Deleted:   0,
		// 	}
		// 	err = c.Insert(answer)
		// } else {
		// 	err = ErrUnavailable
		// }
	case database.TypeBolt:
		answer := &Answer{
			ObjectID:   bson.NewObjectId(),
			Content:    content,
			QuestionID: bson.ObjectIdHex(questionID),
			IsCorrect:  iscorrect,
			CreatedAt:  now,
			UpdatedAt:  now,
			Deleted:    0,
		}

		err = database.Update("answer", questionID+answer.ObjectID.Hex(), &answer)
	default:
		err = ErrCode
	}

	return standardizeError(err)
}

// AnswerUpdate updates a answer
func AnswerUpdate(content string, iscorrect bool, questionID string, answerID string) error {
	var err error

	now := time.Now()

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		// _, err = database.SQL.Exec("UPDATE note SET content=? WHERE id = ? AND user_id = ? LIMIT 1", content, answerID, userID)
	case database.TypeMongoDB:
		// if database.CheckConnection() {
		// 	// Create a copy of mongo
		// 	session := database.Mongo.Copy()
		// 	defer session.Close()
		// 	c := session.DB(database.ReadConfig().MongoDB.Database).C("answer")
		// 	var answer Answer
		// 	answer, err = AnswerByID(userID, answerID)
		// 	if err == nil {
		// 		// Confirm the owner is attempting to modify the answer
		// 		if answer.UserID.Hex() == userID {
		// 			answer.UpdatedAt = now
		// 			answer.Content = content
		// 			err = c.UpdateId(bson.ObjectIdHex(answerID), &answer)
		// 		} else {
		// 			err = ErrUnauthorized
		// 		}
		// 	}
		// } else {
		// 	err = ErrUnavailable
		// }
	case database.TypeBolt:
		var answer Answer
		answer, err = AnswerByID(questionID, answerID)
		if err == nil {
			// Confirm the question is correct for the actual the answer
			if answer.QuestionID.Hex() == questionID {
				answer.UpdatedAt = now
				answer.Content = content
				answer.IsCorrect = iscorrect
				err = database.Update("answer", questionID+answer.ObjectID.Hex(), &answer)
			} else {
				err = ErrUnauthorized
			}
		}
	default:
		err = ErrCode
	}

	return standardizeError(err)
}

// AnswerDelete deletes a answer
func AnswerDelete(questionID string, answerID string) error {
	var err error

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		// _, err = database.SQL.Exec("DELETE FROM note WHERE id = ? AND user_id = ?", answerID, userID)
	case database.TypeMongoDB:
		// if database.CheckConnection() {
		// 	// Create a copy of mongo
		// 	session := database.Mongo.Copy()
		// 	defer session.Close()
		// 	c := session.DB(database.ReadConfig().MongoDB.Database).C("answer")

		// 	var answer Answer
		// 	answer, err = AnswerByID(userID, answerID)
		// 	if err == nil {
		// 		// Confirm the owner is attempting to modify the answer
		// 		if answer.UserID.Hex() == userID {
		// 			err = c.RemoveId(bson.ObjectIdHex(answerID))
		// 		} else {
		// 			err = ErrUnauthorized
		// 		}
		// 	}
		// } else {
		// 	err = ErrUnavailable
		// }
	case database.TypeBolt:
		var answer Answer
		answer, err = AnswerByID(questionID, answerID)
		if err == nil {
			// Confirm the question is correct for the actual answer
			if answer.QuestionID.Hex() == questionID {
				err = database.Delete("answer", questionID+answer.ObjectID.Hex())
			} else {
				err = ErrUnauthorized
			}
		}
	default:
		err = ErrCode
	}

	return standardizeError(err)
}
