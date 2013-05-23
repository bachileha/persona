package main

import (
	"encoding/json"
	"github.com/sauerbraten/persona"
	"log"
	"net/http"
)

// signs the client in by checking with the persona verification API and setting a secure session cookie.
// passes the persona verifiation API response down to the client so the javascript can act on it.
func signIn(resp http.ResponseWriter, req *http.Request) {
	enc := json.NewEncoder(resp)

	response, err := persona.VerifyAssertion("http://localhost:8080/", req.FormValue("assertion"))
	if err != nil {
		log.Println(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response.OK() {
		setSessionCookie(resp, response.Email, response.Expires)
		resp.WriteHeader(http.StatusOK)
	} else {
		resp.WriteHeader(http.StatusUnauthorized)
	}

	enc.Encode(response)
}

// revokes the cookie → client is signed out
func signOut(resp http.ResponseWriter, req *http.Request) {
	revokeSessionCookie(resp)
	resp.WriteHeader(http.StatusOK)
}

// here you could (and probably should) check wether you already the user.
// for this example, we use a map; in production you'd have a table (see addUser()).
func userExists(email string) bool {
	return knownUsers[email]
}

// here you should add the user to some sort of database so you can save preferences/personalisations/additional data you might want to know about a user.
func addUser(email string) error {
	knownUsers[email] = true
	return nil
}