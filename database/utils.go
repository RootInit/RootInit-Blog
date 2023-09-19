package database

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

/* Thrown when an a Cache getter is called with an invalid key */
type InvalidCacheKey struct {
	CacheType string
	Key       interface{}
	Message   string
}

func (r *InvalidCacheKey) Error() string {
	return fmt.Sprintf(
		"%s key does not exist: %v Error: %s",
		r.CacheType, r.Key, r.Message,
	)
}

/* Thrown when an unexpected result occurs from a DB update */
type InvalidRowUpdate struct {
	Table      string
	PrimaryKey int
	Message    string
}

func (r *InvalidRowUpdate) Error() string {
	return fmt.Sprintf(
		"Unexpected update result updating Table: %s Key: %d Error: %s",
		r.Table, r.PrimaryKey, r.Message,
	)
}

func parseIntList(data string) []int {
	splitStr := strings.Split(data, ",")
	var values []int
	for _, strVal := range splitStr {
		intVal, err := strconv.Atoi(strVal)
		if err != nil {
			// This should never happen
			log.Println("Error parsing int list. Atoi error: ", err)
			return nil
		}
		values = append(values, intVal)
	}
	return values
}

func placeholders(num int) string {
	phArray := make([]string, num)
	for i := range phArray {
		phArray[i] = "?"
	}
	placeholders := strings.Join(phArray, ", ")
	return placeholders
}

func getUrlTitle(title string) string {
	if len(title) > 100 {
		// Remove the subtitle E.G. "Farming Potatos: How to grow and harvest" > "Farming Potatos"
		removeSubTitle := regexp.MustCompile(`^(.{30,})(?:\:\s.+)`)
		title = removeSubTitle.ReplaceAllString(title, `$1`)
	}
	// Strip common stop words
	stopWords := regexp.MustCompile(`(?i)(\b(?:a|about|above|actually|after|again|against|all|almost|also|although|always|am|an|and|any|are|as|at|be|became|become|because|been|before|being|below|between|both|but|by|can|could|did|do|does|doing|down|during|each|either|else|few|for|from|further|had|has|have|having|he|he'd|he'll|hence|he's|her|here|here's|hers|herself|him|himself|his|how|how's|I'd|I'll|I'm|I've|if|in|into|is|it|it's|its|itself|just|let's|may|maybe|me|might|mine|more|most|must|my|myself|neither|nor|not|of|oh|on|once|only|ok|or|other|ought|our|ours|ourselves|out|over|own|same|she|she'd|she'll|she's|should|so|some|such|than|that|that's|the|their|theirs|them|themselves|then|there|there's|these|they|they'd|they'll|they're|they've|this|those|through|to|too|under|until|up|very|was|we|we'd|we'll|we're|we've|were|what|what's|when|whenever|when's|where|whereas|wherever|where's|whether|which|while|who|whoever|who's|whose|whom|why|why's|will|with|within|would|yes|yet|you|you'd|you'll|you're|you've|your|yours|yourself|yourselves)\b)`)
	title = stopWords.ReplaceAllString(title, `_`)
	// Strip non alpha-numeric
	specialChars := regexp.MustCompile(`([^A-z0-9])`)
	title = specialChars.ReplaceAllString(title, `_`)
	// Strip duplicate spaces
	dupeSpace := regexp.MustCompile(`(_{2,})`)
	title = dupeSpace.ReplaceAllString(title, `_`)
	// Remove leading and trailing underscores
	title = strings.Trim(title, `_`)
	return title
}
