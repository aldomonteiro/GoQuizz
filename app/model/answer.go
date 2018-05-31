package model

import (
	"time"

	"github.com/asdine/storm/q"

	"GoQuizz/app/shared/database"

	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Answer
// *****************************************************************************

// Answer table contains the information for each answer
type Answer struct {
	ObjectID   bson.ObjectId `storm:"id"`
	Content    string
	QuestionID bson.ObjectId
	QID        uint32
	IsCorrect  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Deleted    uint8
}

// AnswerID returns the answer id
func (q *Answer) AnswerID() string {
	r := ""

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
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
	case database.TypeMongoDB:
	case database.TypeBolt:
		err = database.BoltDB.Select(q.And(q.Eq("ObjectID", bson.ObjectIdHex(answerID)), q.Eq("QuestionID", bson.ObjectIdHex(questionID)))).Find(&result)

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
	case database.TypeMongoDB:
	case database.TypeBolt:
		// View retrieves a record set in Bolt
		err = database.BoltDB.Find("QuestionID", questionID, &result)
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
	case database.TypeMongoDB:
	case database.TypeBolt:
		answer := &Answer{
			ObjectID:   bson.NewObjectId(),
			Content:    content,
			QuestionID: bson.ObjectIdHex(questionID),
			IsCorrect:  iscorrect,
			CreatedAt:  now,
			Deleted:    0,
		}
		err = database.BoltDB.Save(&answer)
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
	case database.TypeMongoDB:
	case database.TypeBolt:
		answer := &Answer{
			ObjectID:   bson.ObjectIdHex(answerID),
			Content:    content,
			QuestionID: bson.ObjectIdHex(questionID),
			IsCorrect:  iscorrect,
			UpdatedAt:  now,
			Deleted:    0,
		}
		err = database.BoltDB.Save(&answer)
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
	case database.TypeMongoDB:
	case database.TypeBolt:
		var answer Answer
		answer, err = AnswerByID(questionID, answerID)
		if err == nil {
			// Confirm the question is correct for the actual answer
			if answer.QuestionID.Hex() == questionID {
				err = database.BoltDB.DeleteStruct(&answer)
			} else {
				err = ErrUnauthorized
			}
		}
	default:
		err = ErrCode
	}

	return standardizeError(err)
}
