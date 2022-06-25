package random

import "math/rand"

var motivationalSentences = []string{
	"Have a nice day! 👋",
	"Enjoy your time. ✋",
	"Have fun today, I will be here if you need me. 🙌",
	"You are doing great! 🤗",
	"Keep up your hard work. 💪",
}

// MotivationalSentence returns a random motivational sentence.
func MotivationalSentence() string {
	return motivationalSentences[rand.Int31n(int32(len(motivationalSentences)-1))]
}
