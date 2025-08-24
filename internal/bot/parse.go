package bot

import "strings"

func ParseSlashCommand(text string, self string) (cmd string, args string, ok bool) {
	text = strings.TrimSpace(text)

	if !strings.HasPrefix(text, "/") {
		return "", "", false
	}

	after := text[1:] // remove leading '/'

	// find first whitespace to split token and args
	ws := strings.IndexFunc(after, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})

	token := after

	if ws >= 0 {
		token = after[:ws]
		args = strings.TrimLeft(after[ws:], " \t\r\n")
	} else {
		args = ""
	}

	if token == "" {
		return "", "", false
	}

	// optional @bot suffix
	cmdPart := token

	if at := strings.IndexByte(token, '@'); at >= 0 {
		cmdPart = token[:at]
		username := token[at+1:]
		// normalize self (allow passing with '@' by mistake)
		if self != "" {
			selfName := strings.TrimPrefix(self, "@")
			if username == "" || !strings.EqualFold(username, selfName) {
				return "", "", false
			}
		}
	}

	if cmdPart == "" {
		return "", "", false
	}

	// validate [A-Za-z0-9_]+ and lowercase
	for i := range len(cmdPart) {
		c := cmdPart[i]
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_') {
			return "", "", false
		}
	}
	cmd = strings.ToLower(cmdPart)

	return cmd, args, true
}
