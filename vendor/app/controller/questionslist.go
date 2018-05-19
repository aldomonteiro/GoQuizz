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

// QuestionsReadGET displays the questions in the QuestionsList
func QuestionsListReadGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	userID := fmt.Sprintf("%s", sess.Values["id"])

	questions, err := model.QuestionsByUserID(userID)

	if err != nil {
		log.Println(err)
		questions = []model.Question{}
	}

	// Display the view
	v := view.New(r)
	v.Name = "questionslist/read"
	v.Vars["first_name"] = sess.Values["first_name"]
	v.Vars["questions"] = questions
	v.Render(w)
}

// NotepadCreateGET displays the question creation page
func QuestionsListCreateGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "questionslist/createHeader"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
}

// QuestionsListCreatePOST handles the note creation form submission
func QuestionsListCreatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"question_header"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		QuestionsListCreateGET(w, r)
		return
	}

	// Get form values
	content := r.FormValue("question_header")

	userID := fmt.Sprintf("%s", sess.Values["id"])

	// Get database result
	err := model.QuestionCreate(content, userID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Question added!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/questionslist", http.StatusFound)
		return
	}

	// Display the same page
	QuestionsListCreateGET(w, r)
}

// QuestionsListUpdateGET displays the note update page
func QuestionsListUpdateGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get the note id
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	questionID := params.ByName("id")

	userID := fmt.Sprintf("%s", sess.Values["id"])

	// Get the note
	question, err := model.QuestionByID(userID, questionID)
	if err != nil { // If the note doesn't exist
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/questionslist", http.StatusFound)
		return
	}

	// Display the view
	v := view.New(r)
	v.Name = "questionslist/update"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["question_content"] = question.Content
	v.Render(w)
}

// QuestionsListUpdatePOST handles the question update form submission
func QuestionsListUpdatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"question_content"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		QuestionsListUpdateGET(w, r)
		return
	}

	// Get form values
	content := r.FormValue("question")

	userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	questionID := params.ByName("id")

	// Get database result
	err := model.QuestionUpdate(content, userID, questionID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Question updated!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/questionslist", http.StatusFound)
		return
	}

	// Display the same page
	QuestionsListUpdateGET(w, r)
}

// QuestionsListDeleteGET handles the question deletion
func QuestionsListDeleteGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	questionID := params.ByName("id")

	// Get database result
	err := model.QuestionDelete(userID, questionID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Question deleted!", view.FlashSuccess})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/questionslist", http.StatusFound)
	return
}
