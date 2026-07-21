package handler

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestMatchQueueFlow(t *testing.T) {
	app := fiber.New()
	api := app.Group("/v1/api")
	h := NewRoomHandler()
	h.RegisterRoutes(api)

	// Helper to generate base64 valid P-256 public key
	genPubKey := func() string {
		priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		pubBytes := elliptic.Marshal(elliptic.P256(), priv.X, priv.Y)
		return base64.StdEncoding.EncodeToString(pubBytes)
	}

	pk1Base64 := genPubKey()
	pk2Base64 := genPubKey()

	// 1. User 1 Enqueue -> Status: waiting (17)
	reqBody1, _ := json.Marshal(MatchQueueRequest{
		PublicKey: pk1Base64,
		Nickname:  "User1",
	})
	req1 := httptest.NewRequest("POST", "/v1/api/match/queue", bytes.NewReader(reqBody1))
	req1.Header.Set("Content-Type", "application/json")
	resp1, err := app.Test(req1)
	if err != nil {
		t.Fatalf("Failed request 1: %v", err)
	}

	var resData1 MatchWaitingStandardResponse
	json.NewDecoder(resp1.Body).Decode(&resData1)

	if resData1.ResponseCode != "17" || resData1.Data.Status != "waiting" {
		t.Fatalf("Expected waiting response (17), got code=%s status=%s", resData1.ResponseCode, resData1.Data.Status)
	}
	if resData1.Data.QueueID == "" {
		t.Fatalf("Expected non-empty queueId")
	}

	// 2. User 2 Enqueue -> Status: matched (00)
	reqBody2, _ := json.Marshal(MatchQueueRequest{
		PublicKey: pk2Base64,
		Nickname:  "User2",
	})
	req2 := httptest.NewRequest("POST", "/v1/api/match/queue", bytes.NewReader(reqBody2))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := app.Test(req2)
	if err != nil {
		t.Fatalf("Failed request 2: %v", err)
	}

	var resData2 MatchMatchedStandardResponse
	json.NewDecoder(resp2.Body).Decode(&resData2)

	if resData2.ResponseCode != "00" || resData2.Data.Status != "matched" {
		t.Fatalf("Expected matched response (00), got code=%s status=%s", resData2.ResponseCode, resData2.Data.Status)
	}
	if resData2.Data.RoomID != resData1.Data.QueueID {
		t.Fatalf("Expected roomId=%s to match queueId=%s", resData2.Data.RoomID, resData1.Data.QueueID)
	}
	if resData2.Data.PeerPublicKey != pk1Base64 {
		t.Fatalf("Expected peerPublicKey=%s, got %s", pk1Base64, resData2.Data.PeerPublicKey)
	}
}
