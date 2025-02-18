package main

import "fmt"

// ANSI escape codes for colors.
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"

	// Bright colors for enhanced readability.
	BrightBlue  = "\033[94m"
	BrightGreen = "\033[92m"
	BrightCyan  = "\033[96m"
)

// FormatWelcomeMessage returns a formatted welcome message including an ASCII art logo.
// The logo is displayed in bright blue, and the welcome text in bright cyan.
func FormatWelcomeMessage() string {
	// Define the ASCII art logo using a double-quoted string with newline escapes.
	// This avoids issues with backticks appearing in the logo.
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

	// Format the welcome message with colored text.
	welcome := fmt.Sprintf("%sWelcome to TCP-Chat!%s\n%s%s%s\n",
		BrightCyan, ColorReset, BrightBlue, logo, ColorReset)
	return welcome
}

// FormatChatMessage formats a chat message with a timestamp, username, and message content.
// The timestamp is shown in yellow, the username in bright green, and the message in the default color.
func FormatChatMessage(timestamp, name, message string) string {
	return fmt.Sprintf("%s[%s]%s %s[%s]:%s %s",
		ColorYellow, timestamp, ColorReset,
		BrightGreen, name, ColorReset,
		message)
}

// FormatSystemMessage formats system messages (such as join or leave notifications)
// in magenta to differentiate them from regular chat messages.
func FormatSystemMessage(message string) string {
	return fmt.Sprintf("%s%s%s", ColorMagenta, message, ColorReset)
}
