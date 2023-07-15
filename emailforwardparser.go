package emailforwardparser

import (
	"strings"

	regexp "github.com/wasilibs/go-re2"
)

func _ParseSubject(subject string) string {
	match, _ := _LoopRegexesMatch(_Subject, subject, true)

	if len(match) > 0 {
		return trimString(match[1])
	}

	return ""
}

type _ParseBodyResult struct {
	Body    string
	Message string
	Email   string
}

func _ParseBody(body string, forwarded bool) _ParseBodyResult {
	body = _CarriageReturn.ReplaceAllString(body, "\n")
	body = _ByteOrderMark.ReplaceAllString(body, "")
	body = _TrailingNonBreakingSpace.ReplaceAllString(body, "")
	body = _NonBreakingSpace.ReplaceAllString(body, " ")

	match := _LoopRegexesSplit(_Separator, body, true)

	if len(match) > 2 {
		email := reconciliateSplitMatch(match, 3, []int{2}, nil)

		return _ParseBodyResult{
			Body:    body,
			Message: trimString(match[0]),
			Email:   trimString(email),
		}
	}

	if forwarded {
		match = _LoopRegexesSplit(_OriginalFrom, body, true)

		if len(match) > 3 {
			email := reconciliateSplitMatch(match, 4, []int{1, 3}, func(i int) bool { return i%3 == 2 })

			return _ParseBodyResult{
				Body:    body,
				Message: trimString(match[0]),
				Email:   trimString(email),
			}
		}
	}

	return _ParseBodyResult{}
}

func _ParseOriginalBody(text string) string {
	regexeses := [][]*regexp.Regexp{
		_OriginalSubject,
		_OriginalCC,
		_OriginalTo,
		_OriginalReplyTo,
	}

	for _, regexes := range regexeses {
		match := _LoopRegexesSplit(regexes, text, true)

		if len(match) > 2 && strings.HasPrefix(match[3], "\n\n") {
			body := reconciliateSplitMatch(match, 4, []int{3}, func(i int) bool { return i%3 == 2 })

			return trimString(body)
		}
	}

	match := _LoopRegexesSplit(append(_OriginalSubject, _OriginalSubjectLax...), text, true)

	if len(match) > 3 {
		body := reconciliateSplitMatch(match, 4, []int{3}, func(i int) bool { return i%3 == 2 })

		return trimString(body)
	}

	return text
}

type _ParseOriginalEmailResult struct {
	Body string

	From _MailboxResult
	To   []_MailboxResult
	CC   []_MailboxResult

	Subject string
	Date    string
}

func _ParseOriginalEmail(text string, body string) _ParseOriginalEmailResult {
	text = _ByteOrderMark.ReplaceAllString(text, "")
	text = _QuoteLineBreak.ReplaceAllString(text, "")
	text = _Quote.ReplaceAllString(text, "")
	text = _FourSpaces.ReplaceAllString(text, "")

	return _ParseOriginalEmailResult{
		Body: _ParseOriginalBody(text),

		From: _ParseOriginalFrom(text, body),
		To:   _ParseOriginalTo(text),
		CC:   _ParseOriginalCC(text),

		Subject: _ParseOriginalSubject(text),
		Date:    _ParseOriginalDate(text, body),
	}
}

func _ParseOriginalFrom(text string, body string) _MailboxResult {
	var name string
	var address string

	authors := _ParseMailbox(_OriginalFrom, text)

	if len(authors) > 0 {
		author := authors[0]

		if len(author.Name) > 0 || len(author.Address) > 0 {
			return author
		}
	}

	match, pattern := _LoopRegexesMatch(_SeparatorWithInformation, body, true)

	if len(match) == 4 {
		namedMatches := findNamedMatches(pattern, body)

		return _PrepareMailbox(namedMatches["from_name"], namedMatches["from_address"])
	}

	match, _ = _LoopRegexesMatch(_OriginalFromLax, text, true)

	if len(match) > 1 {
		name = match[2]
		address = match[3]

		return _PrepareMailbox(name, address)
	}

	return _PrepareMailbox("", "")
}

func _ParseOriginalTo(text string) []_MailboxResult {
	recipients := _ParseMailbox(_OriginalTo, text)

	if len(recipients) > 0 {
		return recipients
	}

	text = _LoopRegexesReplace(_OriginalSubjectLax, text)
	text = _LoopRegexesReplace(_OriginalDateLax, text)
	text = _LoopRegexesReplace(_OriginalCCLax, text)

	return _ParseMailbox(_OriginalToLax, text)
}

func _ParseOriginalCC(text string) []_MailboxResult {
	recipients := _ParseMailbox(_OriginalCC, text)

	if len(recipients) > 0 {
		return recipients
	}

	text = _LoopRegexesReplace(_OriginalSubjectLax, text)
	text = _LoopRegexesReplace(_OriginalDateLax, text)

	return _ParseMailbox(_OriginalCCLax, text)
}

func _ParseOriginalSubject(text string) string {
	match, _ := _LoopRegexesMatch(_OriginalSubject, text, true)

	if len(match) > 0 {
		return trimString(match[1])
	}

	match, _ = _LoopRegexesMatch(_OriginalSubjectLax, text, true)

	if len(match) > 0 {
		return trimString(match[1])
	}

	return ""
}

func _ParseOriginalDate(text string, body string) string {
	match, _ := _LoopRegexesMatch(_OriginalDate, text, true)

	if len(match) > 0 {
		return trimString(match[1])
	}

	match, pattern := _LoopRegexesMatch(_SeparatorWithInformation, body, true)

	if len(match) == 4 {
		namedMatches := findNamedMatches(pattern, body)

		return trimString(namedMatches["date"])
	}

	text = _LoopRegexesReplace(_OriginalSubjectLax, text)
	match, _ = _LoopRegexesMatch(_OriginalDateLax, text, true)

	if len(match) > 0 {
		return trimString(match[1])
	}

	return ""
}

func _ParseMailbox(regexes []*regexp.Regexp, text string) []_MailboxResult {
	match, _ := _LoopRegexesMatch(regexes, text, true)

	if len(match) > 0 {
		mailboxesLine := trimString(match[len(match)-1])

		if len(mailboxesLine) > 0 {
			mailboxes := []_MailboxResult{}

			for len(mailboxesLine) > 0 {
				mailboxMatch, _ := _LoopRegexesMatch(_Mailbox, mailboxesLine, true)

				if len(mailboxMatch) > 0 {
					var name string
					var address string

					if len(mailboxMatch) == 3 {
						name = mailboxMatch[1]
						address = mailboxMatch[2]
					} else {
						address = mailboxMatch[1]
					}

					mailboxes = append(mailboxes, _PrepareMailbox(name, address))

					mailboxesLine = trimString(strings.Replace(mailboxesLine, mailboxMatch[0], "", 1))

					if len(mailboxesLine) > 0 {
						for _, separator := range _MailboxesSeparators {
							if separator == string(mailboxesLine[0]) {
								mailboxesLine = trimString(mailboxesLine[1:])
								break
							}
						}
					}
				} else {
					mailboxes = append(mailboxes, _PrepareMailbox("", mailboxesLine))

					mailboxesLine = ""
				}
			}

			return mailboxes
		}
	}

	return []_MailboxResult{}
}

type _MailboxResult struct {
	Name    string
	Address string
}

func _PrepareMailbox(name string, address string) _MailboxResult {
	name = trimString(name)
	address = trimString(address)

	match, _ := _LoopRegexesMatch(_MailboxAddress, address, true)

	if len(match) == 0 {
		name = address
		address = ""
	}

	if address == name {
		name = ""
	}

	return _MailboxResult{
		Name:    name,
		Address: address,
	}
}

type ReadResultEmail struct {
	Body    string
	From    _MailboxResult
	To      []_MailboxResult
	CC      []_MailboxResult
	Subject string
	Date    string
}

type ReadResult struct {
	Forwarded bool
	Message   string
	Email     ReadResultEmail
}

func Read(body string, subject string) ReadResult {
	email := _ParseOriginalEmailResult{}
	forwarded := false
	bodyResult := _ParseBodyResult{}
	parsedSubject := ""

	if len(subject) > 0 {
		subject = preprocessString(strings.Clone(subject))
		parsedSubject = _ParseSubject(subject)

		if len(parsedSubject) > 0 {
			forwarded = true
		}
	}

	if len(subject) == 0 || forwarded {
		body = preprocessString(strings.Clone(body))
		bodyResult = _ParseBody(body, forwarded)

		if len(bodyResult.Email) > 0 {
			forwarded = true

			email = _ParseOriginalEmail(bodyResult.Email, bodyResult.Body)
		}
	}

	subjectResult := ""

	if len(parsedSubject) > 0 {
		subjectResult = parsedSubject
	} else {
		subjectResult = email.Subject
	}

	return ReadResult{
		Forwarded: forwarded,

		Message: bodyResult.Message,

		Email: ReadResultEmail{
			Body:    email.Body,
			From:    email.From,
			To:      email.To,
			CC:      email.CC,
			Subject: subjectResult,
			Date:    email.Date,
		},
	}
}
