package api

import (
	"encoding/json"
	"io/ioutil"

	"cloud.google.com/go/firestore"
	pfb "github.com/jomy10/poll/firebase"

	"context"
	"fmt"
	"log"
	"net/http"
)

func VoteHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/api/vote" {
		http.Error(res, "Wrong endpoint", http.StatusTeapot)
		return
	}

	if req.Method != "POST" {
		http.Error(res, "Method is not supported", http.StatusMethodNotAllowed)
		fmt.Println("Method not supported")
		return
	}
	res.Header().Add("Access-Control-Allow-Origin", "*")
	res.Header().Add("Access-Control-Allow-Methods", "POST")
	res.Header().Add("Access-Control-Allow-Headers", "*")

	// CORS
	if req.Method == "OPTIONS" {
		// CORS
		method := req.Header.Get("Access-Control-Request-Method")
		if method == "POST" {
			log.Println("Awaiting POST request...")
			res.WriteHeader(http.StatusOK)
			return
		} else {
			http.Error(res, "Method is not supported", http.StatusMethodNotAllowed)
			log.Println("Unsupported CORS method:", method)
			return
		}
	}

	// Firebase
	_, firestoreDB, err := pfb.InitFirebase("jomy-database")
	if err != nil {
		http.Error(res, "Couldn't connect to firebase", http.StatusBadGateway)
		return
	}

	// Get parameters
	params, err := getParams(req)
	if err != nil {
		http.Error(res, "No or invalid parameters", http.StatusBadRequest)
		return
	}

	// Poll data
	doc := firestoreDB.Collection("polls").Doc(params.PollId)
	dsnap, err := doc.Get(context.Background())
	if err != nil {
		http.Error(res, "Couln't read firestore document", http.StatusBadGateway)
	}

	var pollData Poll
	dsnap.DataTo(&pollData)

	pollData.Votes[params.Vote] += 1

	// update votes
	_, err = doc.Update(context.Background(), []firestore.Update{
		{Path: "votes", Value: pollData.Votes},
	})

	if err != nil {
		log.Println(err)
		http.Error(res, "Couldn't update document", http.StatusBadGateway)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Succesfully voted"))
}

type Poll struct {
	Title string
	// Amount of votes mapped to the options
	Votes map[string]int
}

type Params struct {
	// The id of the poll
	PollId string `json:"pollId"`
	// The name of the option the user chose
	Vote string `json:"vote"`
}

func getParams(req *http.Request) (Params, error) {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return Params{}, err
	}

	var data Params
	err = json.Unmarshal(b, &data)
	if err != nil {
		return Params{}, err
	}

	return data, nil
}
