package services

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"menthalhealthjournal/models"
	"menthalhealthjournal/repositories"
)

const sessionNameMaxRunes = 100

func englishOrdinalDay(day int) string {
	if day%100 >= 11 && day%100 <= 13 {
		return fmt.Sprintf("%dth", day)
	}
	switch day % 10 {
	case 1:
		return fmt.Sprintf("%dst", day)
	case 2:
		return fmt.Sprintf("%dnd", day)
	case 3:
		return fmt.Sprintf("%drd", day)
	default:
		return fmt.Sprintf("%dth", day)
	}
}

func formatEnglishSessionFallbackDate(t time.Time) string {
	wd := t.Weekday().String()
	month := t.Format("January")
	return fmt.Sprintf("%s %s of %s %d", wd, englishOrdinalDay(t.Day()), month, t.Year())
}

func sessionNameFallbackTime(session models.ChatSession) time.Time {
	if session.ContextType == models.ContextTypeSingleEntry && session.EntryID != nil {
		entry, err := repositories.GetEntryByID(*session.EntryID)
		if err == nil {
			return entry.CreatedAt
		}
		log.Println("session name fallback: could not load entry, using session time:", err)
	}
	return session.CreatedAt
}

func truncateSessionName(s string, maxRunes int) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.Join(strings.Fields(s), " ")
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxRunes])
}

func buildSessionName(session models.ChatSession, firstUserMessage string) string {
	title, err := GetChatSessionTitleLLM(firstUserMessage)
	if err != nil {
		log.Println("session name: LLM failed, using date fallback:", err)
		return formatEnglishSessionFallbackDate(sessionNameFallbackTime(session))
	}
	title = truncateSessionName(title, sessionNameMaxRunes)
	if title == "" {
		return formatEnglishSessionFallbackDate(sessionNameFallbackTime(session))
	}
	return title
}

// MaybeSetSessionNameOnFirstUserMessage assigns a session title when this is the first user message in the session.
func MaybeSetSessionNameOnFirstUserMessage(session models.ChatSession, firstUserMessage string) {
	n, err := repositories.CountUserMessagesForSession(session.ID)
	if err != nil {
		log.Println("session name: count user messages:", err)
		return
	}
	if n != 1 {
		return
	}
	name := buildSessionName(session, firstUserMessage)
	if err := repositories.UpdateChatSessionNameIfUnset(session.ID, name); err != nil {
		log.Println("session name: could not persist:", err)
	}
}
