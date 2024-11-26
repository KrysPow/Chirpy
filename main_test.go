package main

import "testing"

func TestCleanWords(t *testing.T) {
	cases := map[string]string{
		"I hear Mastodon is better than Chirpy. sharbert I need to migrate": "I hear Mastodon is better than Chirpy. **** I need to migrate",
		"I really need a kerfuffle to go to bed sooner, Fornax !":           "I really need a **** to go to bed sooner, **** !",
	}

	for k, v := range cases {
		if censorProfaneWords(k) != v {
			t.Errorf("- %s - does not match \n - %s -", censorProfaneWords(k), v)
		}
	}
}
