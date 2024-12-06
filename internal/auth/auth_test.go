package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckHashPassword(t *testing.T) {
	pwd := "04234"
	hashed_pwd, _ := HashPassword(pwd)
	err := CheckPasswordHash(pwd, hashed_pwd)
	if err != nil {
		t.Errorf("Not matching password and hash")
	}
}

func TestValidateJWT(t *testing.T) {
	id := uuid.New()
	key := "test"
	token, err := MakeJWT(id, key, time.Second)
	if err != nil {
		t.Error("Problem making token ", err)
	}

	v_id, _ := ValidateJWT(token, key)
	if v_id != id {
		t.Errorf("%s does not match %v", id, v_id)
	}

	_, err = ValidateJWT(token, key+"1")
	if err == nil {
		t.Error("Wrong secret, not supposed to be able to validate")
	}

	time.Sleep(time.Millisecond * 1100)
	_, err = ValidateJWT(token, key)
	if err == nil {
		t.Errorf("Token is supposed to be expired")
	}
}

func TestGetBearerToken(t *testing.T) {
	header := http.Header{}
	header.Set("Authorization", "Bearer TOKEN")
	tok, err := GetBearerToken(header)
	if err != nil {
		t.Error(tok, err)
	}
}
