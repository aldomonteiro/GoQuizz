# GoQuizz

[![Go Report Card](https://goreportcard.com/badge/github.com/aldomonteiro/GoQuizz)](https://goreportcard.com/report/github.com/aldomonteiro/GoQuizz)

I am proud to bring you GoQuizz: an online quiz built in Go. Forked from the GoWebApp [https://github.com/josephspurrier/gowebapp] and integrated with Slickquiz[https://github.com/jewlofthelotus/SlickQuiz], it is a complete system that allows the user to create questions for a quiz and answer them.

## Why GoQuizz? 

I would like to learn Go, so, I found a really simple MVC application written in Go, the GoWebApp [https://github.com/josephspurrier/gowebapp], forked and extended it to do a bit more. The idea was build a system that could be used in real life applications and an online quiz fit this intent.

## GoWebApp

GoWebApp demonstrates how to structure and build a website using the Go language without a framework. There is a blog article you can read at [http://www.josephspurrier.com/go-web-app-example/](http://www.josephspurrier.com/go-web-app-example/).

## Current state

You can create a new user and login in the system. There, it is possible to create questions and answers for them. 

All questions created can be anwswered using the Frontend option in the menu. The Frontend is an almost untouched version of the SlickQuiz project [https://github.com/jewlofthelotus/SlickQuiz], written in JavaScript. This beautiful project came with an option to integrate with CMS system, where I just needed to format the created questions in a specfic JSon format.

## Get started

Getting started with GoQuizz is quite straight. Download it and run the following command:

~~~
go get github.com/aldomonteiro/goquizz
~~~

## Quick Start with Bolt

The forked project GoWebApp was written to work with Bolt, MongoDB or MySQL. When I started to write the GoQuizz, I kept these 3 database options alive. However, maintain this compatibility started to become hard. My intent was learn Go, so, I chose the simpler database option: Bolt.

The goquizz.db file will be created once you start the application.

Build and run from the root directory. Open your web browser to: http://localhost. You should see the welcome page.

Navigate to the login page, and then to the register page. Create a new user and you should be able to login. That's it.

## Overview

The web app has a public home page, authenticated home page, login page, register page,
about page, and a simple notepad to demonstrate the CRUD operations.

The entrypoint for the web app is goquizz.go. The file loads the application settings, 
starts the session, connects to the database, sets up the templates, loads 
the routes, attaches the middleware, and starts the web server.

The front end is built using Bootstrap with a few small changes to fonts and spacing. The flash 
messages are customized so they show up at the bottom right of the screen.

All of the error and warning messages should be either displayed either to the 
user or in the console. Informational messages are displayed to the user via 
flash messages that disappear after 4 seconds. The flash messages are controlled 
by JavaScript in the static folder.

## Structure

The project structure is basically the same as the GoWebApp project [https://github.com/josephspurrier/gowebapp]. Visit the project page and read carefully the README to understand the structure, components, templates and model. 

## Quizz model

Question and Answer structures were added to the project. To do that, the files question.go and answer.go were created under the model folder.

Question.go defines the structure where the questions will be saved:

~~~ go

type Question struct {
	ObjectID  	bson.ObjectId 	`bson:"_id"`
	ID        	uint32        	`db:"id" bson:"id,omitempty"`
	Content   	string        	`db:"content" bson:"content"`
	CorrectMsg 	string			`db:"correctmsg" bson:"correctmsg"`
	WrongMsg 	string			`db:"wrongmsg" bson:"wrongmsg"`
	UserID    	bson.ObjectId 	`bson:"user_id"`
	UID       	uint32        	`db:"user_id" bson:"userid,omitempty"`
	CreatedAt 	time.Time     	`db:"created_at" bson:"created_at"`
	UpdatedAt 	time.Time     	`db:"updated_at" bson:"updated_at"`
	Deleted   	uint8         	`db:"deleted" bson:"deleted"`
}

~~~

Answer.go defines the structure where the answers will be saved (note that it has a link to the question structure)

~~~ go

type Answer struct {
	ObjectID   bson.ObjectId `bson:"_id"`
	ID         uint32        `db:"id" bson:"id,omitempty"` 
	Content    string        `db:"content" bson:"content"`
	QuestionID bson.ObjectId `bson:"question_id"`
	QID        uint32        `db:"question_id" bson:"questionid,omitempty"`
	IsCorrect  bool          `db:"iscorrect" bson:"iscorrect"`
	CreatedAt  time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted    uint8         `db:"deleted" bson:"deleted"`
}

~~~

## Database

This project is focused in the Go language, so, to keep it simple, the support for MongoDB and MySQL are no longer available. The external package used to handle the database is the BBolt [https://github.com/coreos/bbolt]. Bolt is a pure Go key/value store.

## Screenshots

Public Home:

![Image of Public Home](https://cloud.githubusercontent.com/assets/2394539/11319464/e2cd0eac-9045-11e5-9b24-5e480240cd69.jpg)

About:

![Image of About](https://cloud.githubusercontent.com/assets/2394539/11319462/e2c4d2d2-9045-11e5-805f-8b40598c92c3.jpg)

Register:

![Image of Register](https://cloud.githubusercontent.com/assets/2394539/11319466/e2d03500-9045-11e5-9c8e-c28fe663ed0f.jpg)

Login:

![Image of Login](https://cloud.githubusercontent.com/assets/2394539/11319463/e2cd1a00-9045-11e5-8b8e-68030d870cbe.jpg)

Authenticated Home:

![Image of Auth Home](https://cloud.githubusercontent.com/assets/2394539/14809208/75f340d2-0b59-11e6-8d2a-cd26ee872281.PNG)

Create Questions:

![Image of Question Creation]
(https://user-images.githubusercontent.com/26613925/40624812-1d1a49f8-6285-11e8-9bad-59f63e3fa3f8.png)

Show Questions:

![Image of Questions List](https://user-images.githubusercontent.com/26613925/40624816-26462fd8-6285-11e8-9762-5a29263e01f7.png)

Show Answers:

![Image of Answers List]https://user-images.githubusercontent.com/26613925/40624825-2c528138-6285-11e8-94c8-c61bf4ecb42c.png)

Quiz in action:

![Image of Quiz in action]
(https://user-images.githubusercontent.com/26613925/40624831-31731fc4-6285-11e8-955c-e67ed275296b.png)


## Feedback

All feedback is welcome. Let me know if you have any suggestions, questions, or criticisms. 
If something is not idiomatic to Go, please let me know know so we can make it better.
