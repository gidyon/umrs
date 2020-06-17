package token

import (
	"context"
	"encoding/json"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/google/uuid"
	"net/http"
)

// Write writes a token with given IDs
func Write(wr http.ResponseWriter, groupID int32, actorID string) {
	if actorID == "" {
		actorID = uuid.New().String()
	}
	payloadData := &auth.Payload{
		ID:           actorID,
		FirstName:    randomdata.FirstName(randomdata.Male),
		LastName:     randomdata.LastName(),
		PhoneNumber:  randomdata.PhoneNumber(),
		EmailAddress: randomdata.Email(),
		Group:        groupID,
	}

	token, err := auth.GenToken(context.Background(), payloadData, groupID, 0)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}
	wr.Header().Set("content-type", "application/json")
	err = json.NewEncoder(wr).Encode(map[string]string{"token": token, "tokenId": actorID})
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}
}
