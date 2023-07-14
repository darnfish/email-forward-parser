package emailforwardparser

import (
	"strings"

	regexp "github.com/wasilibs/go-re2"
)

func ParseSubject(subject string) string {
	match := LoopRegexesMatch(Subject, subject, true)

	if len(match) > 0 {
		return trimString(match[1])
	}

	return ""
}

type ParseBodyResult struct {
	Body    string
	Message string
	Email   string
}

func ParseBody(body string, forwarded bool) ParseBodyResult {
	body = CarriageReturn.ReplaceAllString(body, "\n")
	body = ByteOrderMark.ReplaceAllString(body, "")
	body = TrailingNonBreakingSpace.ReplaceAllString(body, "")
	body = NonBreakingSpace.ReplaceAllString(body, " ")

	match := LoopRegexesSplit(Separator, body, true)

	if len(match) > 2 {
		email := reconciliateSplitMatch(match, 3, []int{2}, nil)

		return ParseBodyResult{
			Body:    body,
			Message: trimString(match[0]),
			Email:   trimString(email),
		}
	}

	if forwarded {
		match = LoopRegexesSplit(OriginalFrom, body, true)

		if len(match) > 3 {
			email := reconciliateSplitMatch(match, 4, []int{1, 3}, func(i int) bool { return i%3 == 2 })

			return ParseBodyResult{
				Body:    body,
				Message: trimString(match[0]),
				Email:   trimString(email),
			}
		}
	}

	return ParseBodyResult{}
}

func ParseOriginalBody(text string) string {
	regexeses := [][]*regexp.Regexp{
		OriginalSubject,
		OriginalCC,
		OriginalTo,
		OriginalReplyTo,
	}

	for _, regexes := range regexeses {
		match := LoopRegexesSplit(regexes, text, true)

		if len(match) > 2 && strings.HasPrefix(match[3], "\n\n") {
			body := reconciliateSplitMatch(match, 4, []int{3}, func(i int) bool { return i%3 == 2 })

			return trimString(body)
		}
	}

	match := LoopRegexesSplit(append(OriginalSubject, OriginalSubjectLax...), text, true)

	if len(match) > 3 {
		body := reconciliateSplitMatch(match, 4, []int{3}, func(i int) bool { return i%3 == 2 })

		return trimString(body)
	}

	return text
}

type ParseOriginalEmailResult struct {
	Body string

	From MailboxResult
	To   []MailboxResult
	CC   []MailboxResult

	Subject string
	Date    string
}

func ParseOriginalEmail(text string, body string) ParseOriginalEmailResult {
	text = ByteOrderMark.ReplaceAllString(text, "")
	text = QuoteLineBreak.ReplaceAllString(text, "")
	text = Quote.ReplaceAllString(text, "")
	text = FourSpaces.ReplaceAllString(text, "")

	return ParseOriginalEmailResult{
		Body: ParseOriginalBody(text),

		From: ParseOriginalFrom(text, body),
		To:   ParseOriginalTo(text),
		CC:   ParseOriginalCC(text),

		Subject: ParseOriginalSubject(text),
		Date:    ParseOriginalDate(text, body),
	}
}

func ParseOriginalFrom(text string, body string) MailboxResult {
	var name string
	var address string

	authors := ParseMailbox(OriginalFrom, text)

	if len(authors) > 0 {
		author := authors[0]

		if len(author.Name) > 0 || len(author.Address) > 0 {
			return author
		}
	}

	match := LoopRegexesMatch(SeparatorWithInformation, body, true)

	if len(match) == 4 {
		// TODO - match.groups?
	}

	match = LoopRegexesMatch(OriginalFromLax, text, true)

	if len(match) > 1 {
		name = match[2]
		address = match[3]

		return PrepareMailbox(name, address)
	}

	return PrepareMailbox("", "")
}

func ParseOriginalTo(text string) []MailboxResult {
	recipients := ParseMailbox(OriginalTo, text)

	if len(recipients) > 0 {
		return recipients
	}

	text = LoopRegexesReplace(OriginalSubjectLax, text)
	text = LoopRegexesReplace(OriginalDateLax, text)
	text = LoopRegexesReplace(OriginalCCLax, text)

	return ParseMailbox(OriginalToLax, text)
}

func ParseOriginalCC(text string) []MailboxResult {
	recipients := ParseMailbox(OriginalCC, text)

	if len(recipients) > 0 {
		return recipients
	}

	text = LoopRegexesReplace(OriginalSubjectLax, text)
	text = LoopRegexesReplace(OriginalDateLax, text)

	return ParseMailbox(OriginalCCLax, text)
}

func ParseOriginalSubject(text string) string {
	match := LoopRegexesMatch(OriginalSubject, text, true)

	if len(match) > 0 {
		return trimString(match[1])
	}

	match = LoopRegexesMatch(OriginalSubjectLax, text, true)

	if len(match) > 0 {
		return trimString(match[1])
	}

	return ""
}

func ParseOriginalDate(text string, body string) string {
	match := LoopRegexesMatch(OriginalDate, text, true)

	if len(match) > 0 {
		return trimString(match[1])
	}

	match = LoopRegexesMatch(SeparatorWithInformation, body, true)

	if len(match) == 4 {
		// TODO - match.groups?
	}

	text = LoopRegexesReplace(OriginalSubjectLax, text)
	match = LoopRegexesMatch(OriginalDateLax, text, true)

	if len(match) > 0 {
		return trimString(match[1])
	}

	return ""
}

func ParseMailbox(regexes []*regexp.Regexp, text string) []MailboxResult {
	match := LoopRegexesMatch(regexes, text, true)

	if len(match) > 0 {
		mailboxesLine := trimString(match[len(match)-1])

		if len(mailboxesLine) > 0 {
			mailboxes := []MailboxResult{}

			for len(mailboxesLine) > 0 {
				mailboxMatch := LoopRegexesMatch(Mailbox, mailboxesLine, true)

				if len(mailboxMatch) > 0 {
					var name string
					var address string

					if len(mailboxMatch) == 3 {
						name = mailboxMatch[1]
						address = mailboxMatch[2]
					} else {
						address = mailboxMatch[1]
					}

					mailboxes = append(mailboxes, PrepareMailbox(name, address))

					mailboxesLine = trimString(strings.Replace(mailboxesLine, mailboxMatch[0], "", 1))

					if len(mailboxesLine) > 0 {
						for _, separator := range MailboxesSeparators {
							if separator == string(mailboxesLine[0]) {
								mailboxesLine = trimString(mailboxesLine[1:])
								break
							}
						}
					}
				} else {
					mailboxes = append(mailboxes, PrepareMailbox("", mailboxesLine))

					mailboxesLine = ""
				}
			}

			return mailboxes
		}
	}

	return []MailboxResult{}
}

type MailboxResult struct {
	Name    string
	Address string
}

func PrepareMailbox(name string, address string) MailboxResult {
	name = trimString(name)
	address = trimString(address)

	match := LoopRegexesMatch(MailboxAddress, address, true)

	if len(match) == 0 {
		name = address
		address = ""
	}

	if address == name {
		name = ""
	}

	return MailboxResult{
		Name:    name,
		Address: address,
	}
}

type ReadEmailResult struct {
	Body    string
	From    MailboxResult
	To      []MailboxResult
	CC      []MailboxResult
	Subject string
	Date    string
}

type ReadResult struct {
	Forwarded bool
	Message   string
	Email     ReadEmailResult
}

func Read(body string, subject string) ReadResult {
	email := ParseOriginalEmailResult{}
	forwarded := false
	bodyResult := ParseBodyResult{}
	parsedSubject := ""

	if len(subject) > 0 {
		parsedSubject = ParseSubject(subject)

		if len(parsedSubject) > 0 {
			forwarded = true
		}
	}

	if len(subject) == 0 || forwarded {
		bodyResult = ParseBody(body, forwarded)

		if len(bodyResult.Email) > 0 {
			forwarded = true

			email = ParseOriginalEmail(bodyResult.Email, bodyResult.Body)
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

		Email: ReadEmailResult{
			Body:    email.Body,
			From:    email.From,
			To:      email.To,
			CC:      email.CC,
			Subject: subjectResult,
			Date:    email.Date,
		},
	}
}
