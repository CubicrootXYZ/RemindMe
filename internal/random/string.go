package random

var urlSaveRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-")

// URLSaveString generates a random URL-save string with n characters.
func URLSaveString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = urlSaveRunes[Intn(len(urlSaveRunes))]
	}

	return string(b)
}
