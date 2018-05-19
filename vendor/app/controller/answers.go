package controller

import (
	"fmt"
	"log"
	"net/http"

	"app/model"
	"app/shared/session"
	"app/shared/view"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
)

// AnswersReadGET displays the answers from a Question
func AnswersReadGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	userID := fmt.Sprintf("%s", sess.Values["id"])

	// Get the question id
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	questionID := params.ByName("id")

	question, err1 := model.QuestionByID(userID, questionID)
	if err1 != nil {
		log.Println(err1)
		question = model.Question{}
	}

	answers, err2 := model.AnswersByQuestionID(questionID)
	if err2 != nil {
		log.Println(err2)
		answers = []model.Answer{}
	}

	// Display the view
	v := view.New(r)
	v.Name = "answers/read"
	v.Vars["question_header"] = question.Content
	v.Vars["question_id"] = question.QuestionID()
	v.Vars["answers"] = answers
	v.Render(w)
}

// AnswerCreateGET displays the question creation page
func AnswersCreateGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	userID := fmt.Sprintf("%s", sess.Values["id"])

	// Get the question id
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	questionID := params.ByName("id")

	// Get the question
	question, err := model.QuestionByID(userID, questionID)
	if err != nil { // If the question doesn't exist
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/answers/list", http.StatusFound)
		return
	}

	// Display the view, passing to the html the question ID and its text (content)
	v := view.New(r)
	v.Name = "answers/create"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["id"] = question.QuestionID
	v.Vars["question_content"] = question.Content
	v.Render(w)
}

// AnswersCreatePOST handles the note creation form submission
func AnswersCreatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"answer_content"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		AnswersCreateGET(w, r)
		return
	}

	// Get the question id
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	questionID := params.ByName("id")

	// Get form values
	content := r.FormValue("answer_content")

	var iscorrect bool

	if r.FormValue("answer_correct") == "on" {
		iscorrect = true
	} else {
		iscorrect = false
	}

	//userID := fmt.Sprintf("%s", sess.Values["id"])

	// Get database result
	err := model.AnswerCreate(questionID, content, iscorrect)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Answer added!", view.FlashSuccess})
		sess.Save(r, w)
		AnswersReadGET(w, r)
		return
	}

	// Display the same page
	AnswersCreateGET(w, r)
}

// AnswersUpdateGET displays the note update page
func AnswersUpdateGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get the note id
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	questionID := params.ByName("q_id")
	answerID := params.ByName("a_id")

	// userID := fmt.Sprintf("%s", sess.Values["id"])

	// Get the answer
	answer, err := model.AnswerByID(questionID, answerID)

	if err != nil { // If the answer doesn't exist
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/questionslist", http.StatusFound)
		return
	}

	// Display the view
	v := view.New(r)
	v.Name = "answers/update"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["answer_content"] = answer.Content
	if answer.IsCorrect {
		v.Vars["answer_correct"] = "checked"
	} else {
		v.Vars["answer_correct"] = ""
	}

	v.Render(w)
}

// AnswersUpdatePOST handles the answer update form submission
func AnswersUpdatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"answer_content"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		AnswersUpdateGET(w, r)
		return
	}

	// Get form values
	content := r.FormValue("answer_content")

	var iscorrect bool
	if r.FormValue("answer_correct") == "on" {
		iscorrect = true
	} else {
		iscorrect = false
	}

	// userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	questionID := params.ByName("q_id")
	answerID := params.ByName("a_id")

	// Get database result
	err := model.AnswerUpdate(content, iscorrect, questionID, answerID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Answer updated!", view.FlashSuccess})
		sess.Save(r, w)
	}
	// Display the answers list
	http.Redirect(w, r, "/answers/list/"+questionID, http.StatusFound)
}

// AnswersDeleteGET handles the question deletion
func AnswersDeleteGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	questionID := params.ByName("q_id")
	answerID := params.ByName("a_id")

	// Get database result
	err := model.AnswerDelete(questionID, answerID)

	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Answer deleted!", view.FlashSuccess})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/answers/list/"+questionID, http.StatusFound)
	return
}
