package models

import (
	"regexp"
	"sort"
)

var (
	mentionRegexp, _ = regexp.Compile(`@([\w\-\_]{3,20})`)
)

func searchMentionLogins(body string) []string {
	logins := []string{}
	matches := mentionRegexp.FindAllStringSubmatch(body, 10)
	for _, match := range matches {
		if sort.SearchStrings(logins, match[1]) < len(logins) {
			continue
		}
		logins = append(logins, match[1])
	}

	return logins
}

func searchMentionUserIds(body string) (userIds []int32) {
	logins := searchMentionLogins(body)
	if len(logins) > 0 {
		//DB.Model(&User{}).Where("login in (?)", logins).Pluck("id", &userIds)
	}
	return
}

func (r *Reply) CheckMention() {
	if r.NewRecord() {
		return
	}
	mentionUserIds := searchMentionUserIds(r.Body)
	for _, userId := range mentionUserIds {
		if userId == r.User.Id {
			continue
		}
		NotifyMention(userId, r.User.Id, "Reply", r.Id)
	}
}

func (t *Topic) CheckMention() {
	if t.NewRecord() {
		return
	}
	mentionUserIds := searchMentionUserIds(t.Body)
	for _, userId := range mentionUserIds {
		if userId == t.User.Id {
			continue
		}

		NotifyMention(userId, t.User.Id, "Topic", t.Id)
	}
}