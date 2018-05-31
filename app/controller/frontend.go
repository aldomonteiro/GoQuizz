package controller

import (
	"GoQuizz/app/model"
	"GoQuizz/app/shared/view"
	"encoding/json"
	"log"
	"net/http"
)

type SAnswer struct {
	Option  string `json:"option"`
	Correct bool   `json:"correct"`
}

type SQuestion struct {
	Q         string    `json:"q"`
	A         []SAnswer `json:"a"`
	Correct   string    `json:"correct"`
	Incorrect string    `json:"incorrect"`
}

type SInfo struct {
	Name    string `json:"name"`
	Main    string `json:"main"`
	Results string `json:"results"`
	Level1  string `json:"level1"`
	Level2  string `json:"level2"`
	Level3  string `json:"level3"`
	Level4  string `json:"level4"`
	Level5  string `json:"level5"`
}

type SQuiz struct {
	I SInfo       `json:"info"`
	Q []SQuestion `json:"questions"`
}

// Static maps static files
func FrontendGET(w http.ResponseWriter, r *http.Request) {
	// Funciona para static file
	// http.ServeFile(w, r, r.URL.Path[1:])

	quest_db, err := model.QuestionsList()

	if err != nil {
		log.Println(err)
		quest_db = []model.Question{}
	}

	// Create the questions json struct to pass to the view
	quest_js := make([]SQuestion, len(quest_db))

	// Retrieve all questions from database
	for i := 0; i < len(quest_db); i++ {
		// Retrieve all answers from this question
		answers_db, _ := model.AnswersByQuestionID(quest_db[i].QuestionID())

		// Create the answer json struct to pass to the view
		ans_js := make([]SAnswer, len(answers_db))

		for j := 0; j < len(answers_db); j++ {
			ans_js[j] = SAnswer{Option: answers_db[j].Content, Correct: answers_db[j].IsCorrect}
		}

		quest_js[i] = SQuestion{Q: quest_db[i].Content,
			A:         ans_js,
			Correct:   quest_db[i].CorrectMsg,
			Incorrect: quest_db[i].WrongMsg}
	}

	// ans1 := Answer{Option: "Yes", Correct: true}
	// ans2 := Answer{Option: "No", Correct: false}

	// quest1 := Question{Q: "Is Earth bigger than a basketball?",
	// 	A:         []Answer{ans1, ans2},
	// 	Correct:   "<p><span>Good Job!</span> You must be very observant!</p>",
	// 	Incorrect: "<p><span>ERRRR!</span> What planet Earth are <em>you</em> living on?!?</p>"}

	info := SInfo{Name: "Test Your Knowledge!!",
		Main:    "<p>Think you're smart enough to be on Jeopardy? Find out with this super crazy knowledge quiz!</p>",
		Results: "<h5>Learn More</h5><p>Etiam scelerisque, nunc ac egestas consequat, odio nibh euismod nulla, eget auctor orci nibh vel nisi. Aliquam erat volutpat. Mauris vel neque sit amet nunc gravida congue sed sit amet purus.</p>",
		Level1:  "Jeopardy Ready",
		Level2:  "Jeopardy Contender",
		Level3:  "Jeopardy Amateur",
		Level4:  "Jeopardy Newb",
		Level5:  "Stay in school, kid..."}

	quiz := &SQuiz{I: info, Q: quest_js}

	j, _ := json.Marshal(quiz)

	// Display the view
	v := view.New(r)
	v.Name = "frontend/frontend"
	v.Vars["jsonquiz"] = string(j)
	v.Render(w)
}
