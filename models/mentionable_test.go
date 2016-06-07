package models

import "testing"

func TestSearchMentionLogins(t *testing.T) {
	body := `@admin 你好啊 @science09 @admin`
	logins := searchMentionLogins(body)
	if logins[0] != "admin" && logins[1] != "science09" {
		t.Error("not match result:", logins)
	}
}
