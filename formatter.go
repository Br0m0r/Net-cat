package main

import "fmt"

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"

	BrightBlue  = "\033[94m"
	BrightGreen = "\033[92m"
	BrightCyan  = "\033[96m"
)

// FormatWelcomeMessage returns a formatted welcome message with an ASCII art logo. (Defined in formatter.go)
func FormatWelcomeMessage() string {
	logo := "         _nnnn_\n" +
		"        dGGGGMMb\n" +
		"       @p~qp~~qMb\n" +
		"       M|@||@) M|\n" +
		"       @,----.JM|\n" +
		"      JS^\\__/  qKL\n" +
		"     dZP        qKRb\n" +
		"    dZP          qKKb\n" +
		"   fZP            SMMb\n" +
		"   HZM            MMMM\n" +
		"   FqM            MMMM\n" +
		" __| \".        |\\dS\"qML\n" +
		" |    `.       | `' \\Zq\n" +
		"_)      \\.___.,|     .'\n" +
		"\\____   )MMMMMP|   .'\n" +
		"     `-'       `--'\n"

	return fmt.Sprintf("%sWelcome to TCP-Chat!%s\n%s%s%s\n",
		BrightCyan, ColorReset, BrightBlue, logo, ColorReset)
}

// FormatChatMessage formats a chat message with a timestamp, username, and message. (Defined in formatter.go)
func FormatChatMessage(timestamp, name, message string) string {
	return fmt.Sprintf("%s[%s]%s %s[%s]:%s %s",
		ColorYellow, timestamp, ColorReset,
		BrightGreen, name, ColorReset,
		message)
}

// FormatSystemMessage formats system messages (e.g., join/leave notifications) in magenta. (Defined in formatter.go)
func FormatSystemMessage(message string) string {
	return fmt.Sprintf("%s%s%s", ColorMagenta, message, ColorReset)
}
