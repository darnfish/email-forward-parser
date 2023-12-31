package emailforwardparser

import (
	regexp "github.com/wasilibs/go-re2"
)

var _MailboxesSeparators = []string{
	",",
	";",
}

var (
	_QuoteLineBreak           = regexp.MustCompile(`(?m)^(>+)\s?$`)
	_Quote                    = regexp.MustCompile(`(?m)^(>+)\s?`)
	_FourSpaces               = regexp.MustCompile(`(?m)^(\ {4})\s?`)
	_CarriageReturn           = regexp.MustCompile(`(?m)\r\n`)
	_ByteOrderMark            = regexp.MustCompile(`(?m)\xFEFF`)
	_TrailingNonBreakingSpace = regexp.MustCompile(`(?m)\xA0$`)
	_NonBreakingSpace         = regexp.MustCompile(`(?m)\xA0`)
)

var _Subject = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^Fw:(.*)`),         // Outlook Live / 365 (cs, en, hr, hu, sk), Yahoo Mail (all locales)
	regexp.MustCompile(`(?m)^VS:(.*)`),         // Outlook Live / 365 (da), New Outlook 2019 (da)
	regexp.MustCompile(`(?m)^WG:(.*)`),         // Outlook Live / 365 (de), New Outlook 2019 (de)
	regexp.MustCompile(`(?m)^RV:(.*)`),         // Outlook Live / 365 (es), New Outlook 2019 (es)
	regexp.MustCompile(`(?m)^TR:(.*)`),         // Outlook Live / 365 (fr), New Outlook 2019 (fr)
	regexp.MustCompile(`(?m)^I:(.*)`),          // Outlook Live / 365 (it), New Outlook 2019 (it)
	regexp.MustCompile(`(?m)^FW:(.*)`),         // Outlook Live / 365 (nl, pt), New Outlook 2019 (cs, en, hu, nl, pt, ru, sk), Outlook 2019 (all locales)
	regexp.MustCompile(`(?m)^Vs:(.*)`),         // Outlook Live / 365 (no)
	regexp.MustCompile(`(?m)^PD:(.*)`),         // Outlook Live / 365 (pl), New Outlook 2019 (pl)
	regexp.MustCompile(`(?m)^ENC:(.*)`),        // Outlook Live / 365 (pt-br), New Outlook 2019 (pt-br)
	regexp.MustCompile(`(?m)^Redir.:(.*)`),     // Outlook Live / 365 (ro)
	regexp.MustCompile(`(?m)^VB:(.*)`),         // Outlook Live / 365 (sv), New Outlook 2019 (sv)
	regexp.MustCompile(`(?m)^VL:(.*)`),         // New Outlook 2019 (fi)
	regexp.MustCompile(`(?m)^Videresend:(.*)`), // New Outlook 2019 (no)
	regexp.MustCompile(`(?m)^İLT:(.*)`),        // New Outlook 2019 (tr)
	regexp.MustCompile(`(?m)^Fwd:(.*)`),        // Gmail (all locales), Thunderbird (all locales), Missive (en)
}

var _Separator = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^>?\s*Begin forwarded message\s?:`),                             // Apple Mail (en)
	regexp.MustCompile(`(?m)^>?\s*Začátek přeposílané zprávy\s?:`),                          // Apple Mail (cs)
	regexp.MustCompile(`(?m)^>?\s*Start på videresendt besked\s?:`),                         // Apple Mail (da)
	regexp.MustCompile(`(?m)^>?\s*Anfang der weitergeleiteten Nachricht\s?:`),               // Apple Mail (de)
	regexp.MustCompile(`(?m)^>?\s*Inicio del mensaje reenviado\s?:`),                        // Apple Mail (es)
	regexp.MustCompile(`(?m)^>?\s*Välitetty viesti alkaa\s?:`),                              // Apple Mail (fi)
	regexp.MustCompile(`(?m)^>?\s*Début du message réexpédié\s?:`),                          // Apple Mail (fr)
	regexp.MustCompile(`(?m)^>?\s*Début du message transféré\s?:`),                          // Apple Mail iOS (fr)
	regexp.MustCompile(`(?m)^>?\s*Započni proslijeđenu poruku\s?:`),                         // Apple Mail (hr)
	regexp.MustCompile(`(?m)^>?\s*Továbbított levél kezdete\s?:`),                           // Apple Mail (hu)
	regexp.MustCompile(`(?m)^>?\s*Inizio messaggio inoltrato\s?:`),                          // Apple Mail (it)
	regexp.MustCompile(`(?m)^>?\s*Begin doorgestuurd bericht\s?:`),                          // Apple Mail (nl)
	regexp.MustCompile(`(?m)^>?\s*Videresendt melding\s?:`),                                 // Apple Mail (no)
	regexp.MustCompile(`(?m)^>?\s*Początek przekazywanej wiadomości\s?:`),                   // Apple Mail (pl)
	regexp.MustCompile(`(?m)^>?\s*Início da mensagem reencaminhada\s?:`),                    // Apple Mail (pt)
	regexp.MustCompile(`(?m)^>?\s*Início da mensagem encaminhada\s?:`),                      // Apple Mail (pt-br)
	regexp.MustCompile(`(?m)^>?\s*Începe mesajul redirecționat\s?:`),                        // Apple Mail (ro)
	regexp.MustCompile(`(?m)^>?\s*Начало переадресованного сообщения\s?:`),                  // Apple Mail (ro)
	regexp.MustCompile(`(?m)^>?\s*Začiatok preposlanej správy\s?:`),                         // Apple Mail (sk)
	regexp.MustCompile(`(?m)^>?\s*Vidarebefordrat mejl\s?:`),                                // Apple Mail (sv)
	regexp.MustCompile(`(?m)^>?\s*İleti başlangıcı\s?:`),                                    // Apple Mail (tr)
	regexp.MustCompile(`(?m)^>?\s*Початок листа, що пересилається\s?:`),                     // Apple Mail (uk)
	regexp.MustCompile(`(?m)^\s*-{8,10}\s*Forwarded message\s*-{8,10}\s*`),                  // Gmail (all locales), Missive (en), HubSpot (en)
	regexp.MustCompile(`(?m)^\s*_{32}\s*$`),                                                 // Outlook Live / 365 (all locales)
	regexp.MustCompile(`(?m)^\s?Dne\s?.+\,\s?.+\s*[\[|<].+[\]|>]\s?napsal\(a\)\s?:`),        // Outlook 2019 (cz)
	regexp.MustCompile(`(?m)^\s?D.\s?.+\s?skrev\s?\".+\"\s*[\[|<].+[\]|>]\s?:`),             // Outlook 2019 (da)
	regexp.MustCompile(`(?m)^\s?Am\s?.+\s?schrieb\s?\".+\"\s*[\[|<].+[\]|>]\s?:`),           // Outlook 2019 (de)
	regexp.MustCompile(`(?m)^\s?On\s?.+\,\s?\".+\"\s*[\[|<].+[\]|>]\s?wrote\s?:`),           // Outlook 2019 (en)
	regexp.MustCompile(`(?m)^\s?El\s?.+\,\s?\".+\"\s*[\[|<].+[\]|>]\s?escribió\s?:`),        // Outlook 2019 (es)
	regexp.MustCompile(`(?m)^\s?Le\s?.+\,\s?«.+»\s*[\[|<].+[\]|>]\s?a écrit\s?:`),           // Outlook 2019 (fr)
	regexp.MustCompile(`(?m)^\s?.+\s*[\[|<].+[\]|>]\s?kirjoitti\s?.+\s?:`),                  // Outlook 2019 (fi)
	regexp.MustCompile(`(?m)^\s?.+\s?időpontban\s?.+\s*[\[|<|(].+[\]|>|)]\s?ezt írta\s?:`),  // Outlook 2019 (hu)
	regexp.MustCompile(`(?m)^\s?Il giorno\s?.+\s?\".+\"\s*[\[|<].+[\]|>]\s?ha scritto\s?:`), // Outlook 2019 (it)
	regexp.MustCompile(`(?m)^\s?Op\s?.+\s?heeft\s?.+\s*[\[|<].+[\]|>]\s?geschreven\s?:`),    // Outlook 2019 (nl)
	regexp.MustCompile(`(?m)^\s?.+\s*[\[|<].+[\]|>]\s?skrev følgende den\s?.+\s?:`),         // Outlook 2019 (no)
	regexp.MustCompile(`(?m)^\s?Dnia\s?.+\s?„.+”\s*[\[|<].+[\]|>]\s?napisał\s?:`),           // Outlook 2019 (pl)
	regexp.MustCompile(`(?m)^\s?Em\s?.+\,\s?\".+\"\s*[\[|<].+[\]|>]\s?escreveu\s?:`),        // Outlook 2019 (pt)
	regexp.MustCompile(`(?m)^\s?.+\s?пользователь\s?\".+\"\s*[\[|<].+[\]|>]\s?написал\s?:`), // Outlook 2019 (ru)
	regexp.MustCompile(`(?m)^\s?.+\s?používateľ\s?.+\s*\([\[|<].+[\]|>]\)\s?napísal\s?:`),   // Outlook 2019 (sk)
	regexp.MustCompile(`(?m)^\s?Den\s?.+\s?skrev\s?\".+\"\s*[\[|<].+[\]|>]\s?följande\s?:`), // Outlook 2019 (sv)
	regexp.MustCompile(`(?m)^\s?\".+\"\s*[\[|<].+[\]|>]\,\s?.+\s?tarihinde şunu yazdı\s?:`), // Outlook 2019 (tr)
	regexp.MustCompile(`(?m)^\s*-{5,8} Přeposlaná zpráva -{5,8}\s*`),                        // Yahoo Mail (cs), Thunderbird (cs)
	regexp.MustCompile(`(?m)^\s*-{5,8} Videresendt meddelelse -{5,8}\s*`),                   // Yahoo Mail (da), Thunderbird (da)
	regexp.MustCompile(`(?m)^\s*-{5,10} Weitergeleitete Nachricht -{5,10}\s*`),              // Yahoo Mail (de), Thunderbird (de), HubSpot (de)
	regexp.MustCompile(`(?m)^\s*-{5,8} Forwarded Message -{5,8}\s*`),                        // Yahoo Mail (en), Thunderbird (en)
	regexp.MustCompile(`(?m)^\s*-{5,10} Mensaje reenviado -{5,10}\s*`),                      // Yahoo Mail (es), Thunderbird (es), HubSpot (es)
	regexp.MustCompile(`(?m)^\s*-{5,10} Edelleenlähetetty viesti -{5,10}\s*`),               // Yahoo Mail (fi), HubSpot (fi)
	regexp.MustCompile(`(?m)^\s*-{5} Message transmis -{5}\s*`),                             // Yahoo Mail (fr)
	regexp.MustCompile(`(?m)^\s*-{5,8} Továbbított üzenet -{5,8}\s*`),                       // Yahoo Mail (hu), Thunderbird (hu)
	regexp.MustCompile(`(?m)^\s*-{5,10} Messaggio inoltrato -{5,10}\s*`),                    // Yahoo Mail (it), HubSpot (it)
	regexp.MustCompile(`(?m)^\s*-{5,10} Doorgestuurd bericht -{5,10}\s*`),                   // Yahoo Mail (nl), Thunderbird (nl), HubSpot (nl)
	regexp.MustCompile(`(?m)^\s*-{5,8} Videresendt melding -{5,8}\s*`),                      // Yahoo Mail (no), Thunderbird (no)
	regexp.MustCompile(`(?m)^\s*-{5} Przekazana wiadomość -{5}\s*`),                         // Yahoo Mail (pl)
	regexp.MustCompile(`(?m)^\s*-{5,8} Mensagem reencaminhada -{5,8}\s*`),                   // Yahoo Mail (pt), Thunderbird (pt)
	regexp.MustCompile(`(?m)^\s*-{5,10} Mensagem encaminhada -{5,10}\s*`),                   // Yahoo Mail (pt-br), Thunderbird (pt-br), HubSpot (pt-br)
	regexp.MustCompile(`(?m)^\s*-{5,8} Mesaj redirecționat -{5,8}\s*`),                      // Yahoo Mail (ro)
	regexp.MustCompile(`(?m)^\s*-{5} Пересылаемое сообщение -{5}\s*`),                       // Yahoo Mail (ru)
	regexp.MustCompile(`(?m)^\s*-{5} Preposlaná správa -{5}\s*`),                            // Yahoo Mail (sk)
	regexp.MustCompile(`(?m)^\s*-{5,10} Vidarebefordrat meddelande -{5,10}\s*`),             // Yahoo Mail (sv), Thunderbird (sv), HubSpot (sv)
	regexp.MustCompile(`(?m)^\s*-{5} İletilmiş Mesaj -{5}\s*`),                              // Yahoo Mail (tr)
	regexp.MustCompile(`(?m)^\s*-{5} Перенаправлене повідомлення -{5}\s*`),                  // Yahoo Mail (uk)
	regexp.MustCompile(`(?m)^\s*-{8} Välitetty viesti \/ Fwd.Msg -{8}\s*`),                  // Thunderbird (fi)
	regexp.MustCompile(`(?m)^\s*-{8,10} Message transféré -{8,10}\s*`),                      // Thunderbird (fr), HubSpot (fr)
	regexp.MustCompile(`(?m)^\s*-{8} Proslijeđena poruka -{8}\s*`),                          // Thunderbird (hr)
	regexp.MustCompile(`(?m)^\s*-{8} Messaggio Inoltrato -{8}\s*`),                          // Thunderbird (it)
	regexp.MustCompile(`(?m)^\s*-{3} Treść przekazanej wiadomości -{3}\s*`),                 // Thunderbird (pl)
	regexp.MustCompile(`(?m)^\s*-{8} Перенаправленное сообщение -{8}\s*`),                   // Thunderbird (ru)
	regexp.MustCompile(`(?m)^\s*-{8} Preposlaná správa --- Forwarded Message -{8}\s*`),      // Thunderbird (sk)
	regexp.MustCompile(`(?m)^\s*-{8} İletilen İleti -{8}\s*`),                               // Thunderbird (tr)
	regexp.MustCompile(`(?m)^\s*-{8} Переслане повідомлення -{8}\s*`),                       // Thunderbird (uk)
	regexp.MustCompile(`(?m)^\s*-{9,10} メッセージを転送 -{9,10}\s*`),                               // HubSpot (ja)
	regexp.MustCompile(`(?m)^\s*-{9,10} Wiadomość przesłana dalej -{9,10}\s*`),              // HubSpot (pl)
	regexp.MustCompile(`(?m)^>?\s*-{10} Original Message -{10}\s*`),                         // IONOS by 1 & 1 (en)
}

var _SeparatorWithInformation = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^\s?Dne\s?(?P<date>.+)\,\s?(?P<from_name>.+)\s*[\[|<](?P<from_address>.+)[\]|>]\s?napsal\(a\)\s?:`),        // Outlook 2019 (cz)
	regexp.MustCompile(`(?m)^\s?D.\s?(?P<date>.+)\s?skrev\s?\"(?P<from_name>.+)\"\s*[\[|<](?P<from_address>.+)[\]|>]\s?:`),             // Outlook 2019 (da)
	regexp.MustCompile(`(?m)^\s?Am\s?(?P<date>.+)\s?schrieb\s?\"(?P<from_name>.+)\"\s*[\[|<](?P<from_address>.+)[\]|>]\s?:`),           // Outlook 2019 (de)
	regexp.MustCompile(`(?m)^\s?On\s?(?P<date>.+)\,\s?\"(?P<from_name>.+)\"\s*[\[|<](?P<from_address>.+)[\]|>]\s?wrote\s?:`),           // Outlook 2019 (en)
	regexp.MustCompile(`(?m)^\s?El\s?(?P<date>.+)\,\s?\"(?P<from_name>.+)\"\s*[\[|<](?P<from_address>.+)[\]|>]\s?escribió\s?:`),        // Outlook 2019 (es)
	regexp.MustCompile(`(?m)^\s?Le\s?(?P<date>.+)\,\s?«(?P<from_name>.+)»\s*[\[|<](?P<from_address>.+)[\]|>]\s?a écrit\s?:`),           // Outlook 2019 (fr)
	regexp.MustCompile(`(?m)^\s?(?P<from_name>.+)\s*[\[|<](?P<from_address>.+)[\]|>]\s?kirjoitti\s?(?P<date>.+)\s?:`),                  // Outlook 2019 (fi)
	regexp.MustCompile(`(?m)^\s?(?P<date>.+)\s?időpontban\s?(?P<from_name>.+)\s*[\[|<|(](?P<from_address>.+)[\]|>|)]\s?ezt írta\s?:`),  // Outlook 2019 (hu)
	regexp.MustCompile(`(?m)^\s?Il giorno\s?(?P<date>.+)\s?\"(?P<from_name>.+)\"\s*[\[|<](?P<from_address>.+)[\]|>]\s?ha scritto\s?:`), // Outlook 2019 (it)
	regexp.MustCompile(`(?m)^\s?Op\s?(?P<date>.+)\s?heeft\s?(?P<from_name>.+)\s*[\[|<](?P<from_address>.+)[\]|>]\s?geschreven\s?:`),    // Outlook 2019 (nl)
	regexp.MustCompile(`(?m)^\s?(?P<from_name>.+)\s*[\[|<](?P<from_address>.+)[\]|>]\s?skrev følgende den\s?(?P<date>.+)\s?:`),         // Outlook 2019 (no)
	regexp.MustCompile(`(?m)^\s?Dnia\s?(?P<date>.+)\s?„(?P<from_name>.+)”\s*[\[|<](?P<from_address>.+)[\]|>]\s?napisał\s?:`),           // Outlook 2019 (pl)
	regexp.MustCompile(`(?m)^\s?Em\s?(?P<date>.+)\,\s?\"(?P<from_name>.+)\"\s*[\[|<](?P<from_address>.+)[\]|>]\s?escreveu\s?:`),        // Outlook 2019 (pt)
	regexp.MustCompile(`(?m)^\s?(?P<date>.+)\s?пользователь\s?\"(?P<from_name>.+)\"\s*[\[|<](?P<from_address>.+)[\]|>]\s?написал\s?:`), // Outlook 2019 (ru)
	regexp.MustCompile(`(?m)^\s?(?P<date>.+)\s?používateľ\s?(?P<from_name>.+)\s*\([\[|<](?P<from_address>.+)[\]|>]\)\s?napísal\s?:`),   // Outlook 2019 (sk)
	regexp.MustCompile(`(?m)^\s?Den\s?(?P<date>.+)\s?skrev\s?\"(?P<from_name>.+)\"\s*[\[|<](?P<from_address>.+)[\]|>]\s?följande\s?:`), // Outlook 2019 (sv)
	regexp.MustCompile(`(?m)^\s?\"(?P<from_name>.+)\"\s*[\[|<](?P<from_address>.+)[\]|>]\,\s?(?P<date>.+)\s?tarihinde şunu yazdı\s?:`), // Outlook 2019 (tr)
}

var _OriginalSubject = []*regexp.Regexp{
	regexp.MustCompile(`(?im)^\*?Subject\s?:\*?(.+)`), // Apple Mail (en), Gmail (all locales), Outlook Live / 365 (all locales), New Outlook 2019 (en), Thunderbird (da, en), Missive (en), HubSpot (en)
	regexp.MustCompile(`(?im)^Předmět\s?:(.+)`),       // Apple Mail (cs), New Outlook 2019 (cs), Thunderbird (cs)
	regexp.MustCompile(`(?im)^Emne\s?:(.+)`),          // Apple Mail (da, no), New Outlook 2019 (da), Thunderbird (no)
	regexp.MustCompile(`(?im)^Betreff\s?:(.+)`),       // Apple Mail (de), New Outlook 2019 (de), Thunderbird (de), HubSpot (de)
	regexp.MustCompile(`(?im)^Asunto\s?:(.+)`),        // Apple Mail (es), New Outlook 2019 (es), Thunderbird (es), HubSpot (es)
	regexp.MustCompile(`(?im)^Aihe\s?:(.+)`),          // Apple Mail (fi), New Outlook 2019 (fi), Thunderbird (fi), HubSpot (fi)
	regexp.MustCompile(`(?im)^Objet\s?:(.+)`),         // Apple Mail (fr), New Outlook 2019 (fr), HubSpot (fr)
	regexp.MustCompile(`(?im)^Predmet\s?:(.+)`),       // Apple Mail (hr, sk), New Outlook 2019 (sk), Thunderbird (sk)
	regexp.MustCompile(`(?im)^Tárgy\s?:(.+)`),         // Apple Mail (hu), New Outlook 2019 (hu), Thunderbird (hu)
	regexp.MustCompile(`(?im)^Oggetto\s?:(.+)`),       // Apple Mail (it), New Outlook 2019 (it), Thunderbird (it), HubSpot (it)
	regexp.MustCompile(`(?im)^Onderwerp\s?:(.+)`),     // Apple Mail (nl), New Outlook 2019 (nl), Thunderbird (nl), HubSpot (nl)
	regexp.MustCompile(`(?im)^Temat\s?:(.+)`),         // Apple Mail (pl), New Outlook 2019 (pl), Thunderbird (pl), HubSpot (pl)
	regexp.MustCompile(`(?im)^Assunto\s?:(.+)`),       // Apple Mail (pt, pt-br), New Outlook 2019 (pt, pt-br), Thunderbird (pt, pt-br), HubSpot (pt-br)
	regexp.MustCompile(`(?im)^Subiectul\s?:(.+)`),     // Apple Mail (ro), Thunderbird (ro)
	regexp.MustCompile(`(?im)^Тема\s?:(.+)`),          // Apple Mail (ru, uk), New Outlook 2019 (ru), Thunderbird (ru, uk)
	regexp.MustCompile(`(?im)^Ämne\s?:(.+)`),          // Apple Mail (sv), New Outlook 2019 (sv), Thunderbird (sv), HubSpot (sv)
	regexp.MustCompile(`(?im)^Konu\s?:(.+)`),          // Apple Mail (tr), Thunderbird (tr)
	regexp.MustCompile(`(?im)^Sujet\s?:(.+)`),         // Thunderbird (fr)
	regexp.MustCompile(`(?im)^Naslov\s?:(.+)`),        // Thunderbird (hr)
	regexp.MustCompile(`(?im)^件名：(.+)`),               // HubSpot (ja)
}

var _OriginalSubjectLax = []*regexp.Regexp{
	regexp.MustCompile(`(?i)Subject\s?:(.+)`),   // Yahoo Mail (en)
	regexp.MustCompile(`(?i)Emne\s?:(.+)`),      // Yahoo Mail (da, no)
	regexp.MustCompile(`(?i)Předmět\s?:(.+)`),   // Yahoo Mail (cs)
	regexp.MustCompile(`(?i)Betreff\s?:(.+)`),   // Yahoo Mail (de)
	regexp.MustCompile(`(?i)Asunto\s?:(.+)`),    // Yahoo Mail (es)
	regexp.MustCompile(`(?i)Aihe\s?:(.+)`),      // Yahoo Mail (fi)
	regexp.MustCompile(`(?i)Objet\s?:(.+)`),     // Yahoo Mail (fr)
	regexp.MustCompile(`(?i)Tárgy\s?:(.+)`),     // Yahoo Mail (hu)
	regexp.MustCompile(`(?i)Oggetto\s?:(.+)`),   // Yahoo Mail (it)
	regexp.MustCompile(`(?i)Onderwerp\s?:(.+)`), // Yahoo Mail (nl)
	regexp.MustCompile(`(?i)Assunto\s?:?(.+)`),  // Yahoo Mail (pt, pt-br)
	regexp.MustCompile(`(?i)Temat\s?:(.+)`),     // Yahoo Mail (pl)
	regexp.MustCompile(`(?i)Subiect\s?:(.+)`),   // Yahoo Mail (ro)
	regexp.MustCompile(`(?i)Тема\s?:(.+)`),      // Yahoo Mail (ru, uk)
	regexp.MustCompile(`(?i)Predmet\s?:(.+)`),   // Yahoo Mail (sk)
	regexp.MustCompile(`(?i)Ämne\s?:(.+)`),      // Yahoo Mail (sv)
	regexp.MustCompile(`(?i)Konu\s?:(.+)`),      // Yahoo Mail (tr)
}

var _OriginalFrom = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^(\*?\s*From\s?:\*?(.+))$`),  // Apple Mail (en), Outlook Live / 365 (all locales), New Outlook 2019 (en), Thunderbird (da, en), Missive (en), HubSpot (en)
	regexp.MustCompile(`(?m)^(\s*Od\s?:(.+))$`),          // Apple Mail (cs, pl, sk), Gmail (cs, pl, sk), New Outlook 2019 (cs, pl, sk), Thunderbird (cs, sk), HubSpot (pl)
	regexp.MustCompile(`(?m)^(\s*Fra\s?:(.+))$`),         // Apple Mail (da, no), Gmail (da, no), New Outlook 2019 (da), Thunderbird (no)
	regexp.MustCompile(`(?m)^(\s*Von\s?:(.+))$`),         // Apple Mail (de), Gmail (de), New Outlook 2019 (de), Thunderbird (de), HubSpot (de)
	regexp.MustCompile(`(?m)^(\s*De\s?:(.+))$`),          // Apple Mail (es, fr, pt, pt-br), Gmail (es, fr, pt, pt-br), New Outlook 2019 (es, fr, pt, pt-br), Thunderbird (fr, pt, pt-br), HubSpot (es, fr, pt-br)
	regexp.MustCompile(`(?m)^(\s*Lähettäjä\s?:(.+))$`),   // Apple Mail (fi), Gmail (fi), New Outlook 2019 (fi), Thunderbird (fi), HubSpot (fi)
	regexp.MustCompile(`(?m)^(\s*Šalje\s?:(.+))$`),       // Apple Mail (hr), Gmail (hr), Thunderbird (hr)
	regexp.MustCompile(`(?m)^(\s*Feladó\s?:(.+))$`),      // Apple Mail (hu), Gmail (hu), New Outlook 2019 (fr), Thunderbird (hu)
	regexp.MustCompile(`(?m)^(\s*Da\s?:(.+))$`),          // Apple Mail (it), Gmail (it), New Outlook 2019 (it), HubSpot (it)
	regexp.MustCompile(`(?m)^(\s*Van\s?:(.+))$`),         // Apple Mail (nl), Gmail (nl), New Outlook 2019 (nl), Thunderbird (nl), HubSpot (nl)
	regexp.MustCompile(`(?m)^(\s*Expeditorul\s?:(.+))$`), // Apple Mail (ro)
	regexp.MustCompile(`(?m)^(\s*Отправитель\s?:(.+))$`), // Apple Mail (ru)
	regexp.MustCompile(`(?m)^(\s*Från\s?:(.+))$`),        // Apple Mail (sv), Gmail (sv), New Outlook 2019 (sv), Thunderbird (sv), HubSpot (sv)
	regexp.MustCompile(`(?m)^(\s*Kimden\s?:(.+))$`),      // Apple Mail (tr), Thunderbird (tr)
	regexp.MustCompile(`(?m)^(\s*Від кого\s?:(.+))$`),    // Apple Mail (uk)
	regexp.MustCompile(`(?m)^(\s*Saatja\s?:(.+))$`),      // Gmail (et)
	regexp.MustCompile(`(?m)^(\s*De la\s?:(.+))$`),       // Gmail (ro)
	regexp.MustCompile(`(?m)^(\s*Gönderen\s?:(.+))$`),    // Gmail (tr)
	regexp.MustCompile(`(?m)^(\s*От\s?:(.+))$`),          // Gmail (ru), New Outlook 2019 (ru), Thunderbird (ru)
	regexp.MustCompile(`(?m)^(\s*Від\s?:(.+))$`),         // Gmail (uk), Thunderbird (uk)
	regexp.MustCompile(`(?m)^(\s*Mittente\s?:(.+))$`),    // Thunderbird (it)
	regexp.MustCompile(`(?m)^(\s*Nadawca\s?:(.+))$`),     // Thunderbird (pl)
	regexp.MustCompile(`(?m)^(\s*de la\s?:(.+))$`),       // Thunderbird (ro)
	regexp.MustCompile(`(?m)^(\s*送信元：(.+))$`),            // HubSpot (ja)
}

var _OriginalFromLax = []*regexp.Regexp{
	regexp.MustCompile(`(\s*From\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`),      // Yahoo Mail (en)
	regexp.MustCompile(`(\s*Od\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`),        // Yahoo Mail (cs, pl, sk)
	regexp.MustCompile(`(\s*Fra\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`),       // Yahoo Mail (da, no)
	regexp.MustCompile(`(\s*Von\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`),       // Yahoo Mail (de)
	regexp.MustCompile(`(\s*De\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`),        // Yahoo Mail (es, fr, pt, pt-br)
	regexp.MustCompile(`(\s*Lähettäjä\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`), // Yahoo Mail (fi)
	regexp.MustCompile(`(\s*Feladó\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`),    // Yahoo Mail (hu)
	regexp.MustCompile(`(\s*Da\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`),        // Yahoo Mail (it)
	regexp.MustCompile(`(\s*Van\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`),       // Yahoo Mail (nl)
	regexp.MustCompile(`(\s*De la\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`),     // Yahoo Mail (ro)
	regexp.MustCompile(`(\s*От\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`),        // Yahoo Mail (ru)
	regexp.MustCompile(`(\s*Från\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`),      // Yahoo Mail (sv)
	regexp.MustCompile(`(\s*Kimden\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`),    // Yahoo Mail (tr)
	regexp.MustCompile(`(\s*Від\s?:(.+?)\s?\n?\s*[\[|<](.+?)[\]|>])`),       // Yahoo Mail (uk)
}

var _OriginalTo = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^\*?\s*To\s?:\*?(.+)$`),      //Apple Mail (en), Gmail (all locales), Outlook Live / 365 (all locales), Thunderbird (da, en), Missive (en), HubSpot (en)
	regexp.MustCompile(`(?m)^\s*Komu\s?:(.+)$`),          //Apple Mail (cs), New Outlook 2019 (cs, sk), Thunderbird (cs)
	regexp.MustCompile(`(?m)^\s*Til\s?:(.+)$`),           //Apple Mail (da, no), New Outlook 2019 (da), Thunderbird (no)
	regexp.MustCompile(`(?m)^\s*An\s?:(.+)$`),            //Apple Mail (de), New Outlook 2019 (de), Thunderbird (de), HubSpot (de)
	regexp.MustCompile(`(?m)^\s*Para\s?:(.+)$`),          //Apple Mail (es, pt, pt-br), New Outlook 2019 (es, pt, pt-br), Thunderbird (es, pt, pt-br), HubSpot (pt-br)
	regexp.MustCompile(`(?m)^\s*Vastaanottaja\s?:(.+)$`), //Apple Mail (fi), New Outlook 2019 (fi), Thunderbird (fi), HubSpot (fi)
	regexp.MustCompile(`(?m)^\s*À\s?:(.+)$`),             //Apple Mail (fr), New Outlook 2019 (fr), HubSpot (fr)
	regexp.MustCompile(`(?m)^\s*Prima\s?:(.+)$`),         //Apple Mail (hr), Thunderbird (hr)
	regexp.MustCompile(`(?m)^\s*Címzett\s?:(.+)$`),       //Apple Mail (hu), New Outlook 2019 (hu), Thunderbird (hu)
	regexp.MustCompile(`(?m)^\s*A\s?:(.+)$`),             //Apple Mail (it), New Outlook 2019 (it), Thunderbird (it), HubSpot (es, it)
	regexp.MustCompile(`(?m)^\s*Aan\s?:(.+)$`),           //Apple Mail (nl), New Outlook 2019 (nl), Thunderbird (nl), HubSpot (nl)
	regexp.MustCompile(`(?m)^\s*Do\s?:(.+)$`),            //Apple Mail (pl), New Outlook 2019 (pl), HubSpot (pl)
	regexp.MustCompile(`(?m)^\s*Destinatarul\s?:(.+)$`),  //Apple Mail (ro)
	regexp.MustCompile(`(?m)^\s*Кому\s?:(.+)$`),          //Apple Mail (ru, uk), New Outlook 2019 (ru), Thunderbird (ru, uk)
	regexp.MustCompile(`(?m)^\s*Pre\s?:(.+)$`),           //Apple Mail (sk), Thunderbird (sk)
	regexp.MustCompile(`(?m)^\s*Till\s?:(.+)$`),          //Apple Mail (sv), New Outlook 2019 (sv), Thunderbird (sv)
	regexp.MustCompile(`(?m)^\s*Kime\s?:(.+)$`),          //Apple Mail (tr), Thunderbird (tr)
	regexp.MustCompile(`(?m)^\s*Pour\s?:(.+)$`),          //Thunderbird (fr)
	regexp.MustCompile(`(?m)^\s*Adresat\s?:(.+)$`),       //Thunderbird (pl)
	regexp.MustCompile(`(?m)^\s*送信先：(.+)$`),              // HubSpot (ja)
}

var _OriginalToLax = []*regexp.Regexp{
	regexp.MustCompile(`(?m)\s*To\s?:(.+)$`),            // Yahook Mail (en)
	regexp.MustCompile(`(?m)\s*Komu\s?:(.+)$`),          // Yahook Mail (cs, sk)
	regexp.MustCompile(`(?m)\s*Til\s?:(.+)$`),           // Yahook Mail (da, no, sv)
	regexp.MustCompile(`(?m)\s*An\s?:(.+)$`),            // Yahook Mail (de)
	regexp.MustCompile(`(?m)\s*Para\s?:(.+)$`),          // Yahook Mail (es, pt, pt-br)
	regexp.MustCompile(`(?m)\s*Vastaanottaja\s?:(.+)$`), // Yahook Mail (fi)
	regexp.MustCompile(`(?m)\s*À\s?:(.+)$`),             // Yahook Mail (fr)
	regexp.MustCompile(`(?m)\s*Címzett\s?:(.+)$`),       // Yahook Mail (hu)
	regexp.MustCompile(`(?m)\s*A\s?:(.+)$`),             // Yahook Mail (it)
	regexp.MustCompile(`(?m)\s*Aan\s?:(.+)$`),           // Yahook Mail (nl)
	regexp.MustCompile(`(?m)\s*Do\s?:(.+)$`),            // Yahook Mail (pl)
	regexp.MustCompile(`(?m)\s*Către\s?:(.+)$`),         // Yahook Mail (ro), Thunderbird (ro)
	regexp.MustCompile(`(?m)\s*Кому\s?:(.+)$`),          // Yahook Mail (ru, uk)
	regexp.MustCompile(`(?m)\s*Till\s?:(.+)$`),          // Yahook Mail (sv)
	regexp.MustCompile(`(?m)\s*Kime\s?:(.+)$`),          // Yahook Mail (tr)
}

var _OriginalReplyTo = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^\s*Reply-To\s?:(.+)$`),        // Apple Mail (en)
	regexp.MustCompile(`(?m)^\s*Odgovori na\s?:(.+)$`),     // Apple Mail (hr)
	regexp.MustCompile(`(?m)^\s*Odpověď na\s?:(.+)$`),      // Apple Mail (cs)
	regexp.MustCompile(`(?m)^\s*Svar til\s?:(.+)$`),        // Apple Mail (da)
	regexp.MustCompile(`(?m)^\s*Antwoord aan\s?:(.+)$`),    // Apple Mail (nl)
	regexp.MustCompile(`(?m)^\s*Vastaus\s?:(.+)$`),         // Apple Mail (fi)
	regexp.MustCompile(`(?m)^\s*Répondre à\s?:(.+)$`),      // Apple Mail (fr)
	regexp.MustCompile(`(?m)^\s*Antwort an\s?:(.+)$`),      // Apple Mail (de)
	regexp.MustCompile(`(?m)^\s*Válaszcím\s?:(.+)$`),       // Apple Mail (hu)
	regexp.MustCompile(`(?m)^\s*Rispondi a\s?:(.+)$`),      // Apple Mail (it)
	regexp.MustCompile(`(?m)^\s*Svar til\s?:(.+)$`),        // Apple Mail (no)
	regexp.MustCompile(`(?m)^\s*Odpowiedź-do\s?:(.+)$`),    // Apple Mail (pl)
	regexp.MustCompile(`(?m)^\s*Responder A\s?:(.+)$`),     // Apple Mail (pt)
	regexp.MustCompile(`(?m)^\s*Responder a\s?:(.+)$`),     // Apple Mail (pt-br, es)
	regexp.MustCompile(`(?m)^\s*Răspuns către\s?:(.+)$`),   // Apple Mail (ro)
	regexp.MustCompile(`(?m)^\s*Ответ-Кому\s?:(.+)$`),      // Apple Mail (ru)
	regexp.MustCompile(`(?m)^\s*Odpovedať-Pre\s?:(.+)$`),   // Apple Mail (sk)
	regexp.MustCompile(`(?m)^\s*Svara till\s?:(.+)$`),      // Apple Mail (sv)
	regexp.MustCompile(`(?m)^\s*Yanıt Adresi\s?:(.+)$`),    // Apple Mail (tr)
	regexp.MustCompile(`(?m)^\s*Кому відповісти\s?:(.+)$`), // Apple Mail (uk)
}

var _OriginalCC = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^\*?\s*Cc\s?:\*?(.+)$`),      // Apple Mail (en, da, es, fr, hr, it, pt, pt-br, ro, sk), Gmail (all locales), Outlook Live / 365 (all locales), New Outlook 2019 (da, de, en, fr, it, pt-br), Missive (en), HubSpot (de, en, es, it, nl, pt-br)
	regexp.MustCompile(`(?m)^\s*CC\s?:(.+)$`),            // New Outlook 2019 (es, nl, pt), Thunderbird (da, en, es, fi, hr, hu, it, nl, no, pt, pt-br, ro, tr, uk)
	regexp.MustCompile(`(?m)^\s*Kopie\s?:(.+)$`),         // Apple Mail (cs, de, nl), New Outlook 2019 (cs), Thunderbird (cs)
	regexp.MustCompile(`(?m)^\s*Kopio\s?:(.+)$`),         // Apple Mail (fi), New Outlook 2019 (fi), HubSpot (fi)
	regexp.MustCompile(`(?m)^\s*Másolat\s?:(.+)$`),       // Apple Mail (hu)
	regexp.MustCompile(`(?m)^\s*Kopi\s?:(.+)$`),          // Apple Mail (no)
	regexp.MustCompile(`(?m)^\s*Dw\s?:(.+)$`),            // Apple Mail (pl)
	regexp.MustCompile(`(?m)^\s*Копия\s?:(.+)$`),         // Apple Mail (ru), New Outlook 2019 (ru), Thunderbird (ru)
	regexp.MustCompile(`(?m)^\s*Kopia\s?:(.+)$`),         // Apple Mail (sv), New Outlook 2019 (sv), Thunderbird (pl, sv), HubSpot (sv)
	regexp.MustCompile(`(?m)^\s*Bilgi\s?:(.+)$`),         // Apple Mail (tr)
	regexp.MustCompile(`(?m)^\s*Копія\s?:(.+)$`),         // Apple Mail (uk),
	regexp.MustCompile(`(?m)^\s*Másolatot kap\s?:(.+)$`), // New Outlook 2019 (hu)
	regexp.MustCompile(`(?m)^\s*Kópia\s?:(.+)$`),         // New Outlook 2019 (sk), Thunderbird (sk)
	regexp.MustCompile(`(?m)^\s*DW\s?:(.+)$`),            // New Outlook 2019 (pl), HubSpot (pl)
	regexp.MustCompile(`(?m)^\s*Kopie \(CC\)\s?:(.+)$`),  // Thunderbird (de)
	regexp.MustCompile(`(?m)^\s*Copie à\s?:(.+)$`),       // Thunderbird (fr)
	regexp.MustCompile(`(?m)^\s*CC：(.+)$`),               // HubSpot (ja)}
}

var _OriginalCCLax = []*regexp.Regexp{
	regexp.MustCompile(`(?m)\s*Cc\s?:(.+)$`),      // Yahoo Mail (da, en, it, nl, pt, pt-br, ro, tr)
	regexp.MustCompile(`(?m)\s*CC\s?:(.+)$`),      // Yahoo Mail (de, es)
	regexp.MustCompile(`(?m)\s*Kopie\s?:(.+)$`),   // Yahoo Mail (cs)
	regexp.MustCompile(`(?m)\s*Kopio\s?:(.+)$`),   // Yahoo Mail (fi)
	regexp.MustCompile(`(?m)\s*Másolat\s?:(.+)$`), // Yahoo Mail (hu)
	regexp.MustCompile(`(?m)\s*Kopi\s?:(.+)$`),    // Yahoo Mail (no)
	regexp.MustCompile(`(?m)\s*Dw\s?(.+)$`),       // Yahoo Mail (pl)
	regexp.MustCompile(`(?m)\s*Копия\s?:(.+)$`),   // Yahoo Mail (ru)
	regexp.MustCompile(`(?m)\s*Kópia\s?:(.+)$`),   // Yahoo Mail (sk)
	regexp.MustCompile(`(?m)\s*Kopia\s?:(.+)$`),   // Yahoo Mail (sv)
	regexp.MustCompile(`(?m)\s*Копія\s?:(.+)$`),   // Yahoo Mail (uk)
}

var _OriginalDate = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^\s*Date\s?:(.+)$`),       // Apple Mail (en, fr), Gmail (all locales), New Outlook 2019 (en, fr), Thunderbird (da, en, fr), Missive (en), HubSpot (en, fr)
	regexp.MustCompile(`(?m)^\s*Datum\s?:(.+)$`),      // Apple Mail (cs, de, hr, nl, sv), New Outlook 2019 (cs, de, nl, sv), Thunderbird (cs, de, hr, nl, sv), HubSpot (de, nl, sv)
	regexp.MustCompile(`(?m)^\s*Dato\s?:(.+)$`),       // Apple Mail (da, no), New Outlook 2019 (da), Thunderbird (no)
	regexp.MustCompile(`(?m)^\s*Envoyé\s?:(.+)$`),     // New Outlook 2019 (fr)
	regexp.MustCompile(`(?m)^\s*Fecha\s?:(.+)$`),      // Apple Mail (es), New Outlook 2019 (es), Thunderbird (es), HubSpot (es)
	regexp.MustCompile(`(?m)^\s*Päivämäärä\s?:(.+)$`), // Apple Mail (fi), New Outlook 2019 (fi), HubSpot (fi)
	regexp.MustCompile(`(?m)^\s*Dátum\s?:(.+)$`),      // Apple Mail (hu, sk), New Outlook 2019 (sk), Thunderbird (hu, sk)
	regexp.MustCompile(`(?m)^\s*Data\s?:(.+)$`),       // Apple Mail (it, pl, pt, pt-br), New Outlook 2019 (it, pl, pt, pt-br), Thunderbird (it, pl, pt, pt-br), HubSpot (it, pl, pt-br)
	regexp.MustCompile(`(?m)^\s*Dată\s?:(.+)$`),       // Apple Mail (ro), Thunderbird (ro)
	regexp.MustCompile(`(?m)^\s*Дата\s?:(.+)$`),       // Apple Mail (ru, uk), New Outlook 2019 (ru), Thunderbird (ru, uk)
	regexp.MustCompile(`(?m)^\s*Tarih\s?:(.+)$`),      // Apple Mail (tr), Thunderbird (tr)
	regexp.MustCompile(`(?m)^\*?\s*Sent\s?:\*?(.+)$`), // Outlook Live / 365 (all locales)
	regexp.MustCompile(`(?m)^\s*Päiväys\s?:(.+)$`),    // Thunderbird (fi)
	regexp.MustCompile(`(?m)^\s*日付：(.+)$`),            // HubSpot (ja)
}

var _OriginalDateLax = []*regexp.Regexp{
	regexp.MustCompile(`(?m)\s*Datum\s?:(.+)$`),       // Yahoo Mail (cs)
	regexp.MustCompile(`(?m)\s*Sendt\s?:(.+)$`),       // Yahoo Mail (da, no)
	regexp.MustCompile(`(?m)\s*Gesendet\s?:(.+)$`),    // Yahoo Mail (de)
	regexp.MustCompile(`(?m)\s*Sent\s?:(.+)$`),        // Yahoo Mail (en)
	regexp.MustCompile(`(?m)\s*Enviado\s?:(.+)$`),     // Yahoo Mail (es, pt, pt-br)
	regexp.MustCompile(`(?m)\s*Envoyé\s?:(.+)$`),      // Yahoo Mail (fr)
	regexp.MustCompile(`(?m)\s*Lähetetty\s?:(.+)$`),   // Yahoo Mail (fi)
	regexp.MustCompile(`(?m)\s*Elküldve\s?:(.+)$`),    // Yahoo Mail (hu)
	regexp.MustCompile(`(?m)\s*Inviato\s?:(.+)$`),     // Yahoo Mail (it)
	regexp.MustCompile(`(?m)\s*Verzonden\s?:(.+)$`),   // Yahoo Mail (it)
	regexp.MustCompile(`(?m)\s*Wysłano\s?:(.+)$`),     // Yahoo Mail (pl)
	regexp.MustCompile(`(?m)\s*Trimis\s?:(.+)$`),      // Yahoo Mail (ro)
	regexp.MustCompile(`(?m)\s*Отправлено\s?:(.+)$`),  // Yahoo Mail (ru)
	regexp.MustCompile(`(?m)\s*Odoslané\s?:(.+)$`),    // Yahoo Mail (sk)
	regexp.MustCompile(`(?m)\s*Skickat\s?:(.+)$`),     // Yahoo Mail (sv)
	regexp.MustCompile(`(?m)\s*Gönderilen\s?:(.+)$`),  // Yahoo Mail (tr)
	regexp.MustCompile(`(?m)\s*Відправлено\s?:(.+)$`), // Yahoo Mail (uk)
}

var _Mailbox = []*regexp.Regexp{
	regexp.MustCompile(`^\s?\n?\s*<.+?<mailto\:(.+?)>>`),           // "<walter.sheltan@acme.com<mailto:walter.sheltan@acme.com>>"
	regexp.MustCompile(`^(.+?)\s?\n?\s*<.+?<mailto\:(.+?)>>`),      // "Walter Sheltan <walter.sheltan@acme.com<mailto:walter.sheltan@acme.com>>"
	regexp.MustCompile(`^(.+?)\s?\n?\s*[\[|<]mailto\:(.+?)[\]|>]`), // "Walter Sheltan <mailto:walter.sheltan@acme.com>" or "Walter Sheltan [mailto:walter.sheltan@acme.com]" or "walter.sheltan@acme.com <mailto:walter.sheltan@acme.com>"
	regexp.MustCompile(`^\'(.+?)\'\s?\n?\s*[\[|<](.+?)[\]|>]`),     // "'Walter Sheltan' <walter.sheltan@acme.com>" or "'Walter Sheltan' [walter.sheltan@acme.com]" or "'walter.sheltan@acme.com' <walter.sheltan@acme.com>"
	regexp.MustCompile(`^\"\'(.+?)\'\"\s?\n?\s*[\[|<](.+?)[\]|>]`), // ""'Walter Sheltan'" <walter.sheltan@acme.com>" or ""'Walter Sheltan'" [walter.sheltan@acme.com]" or ""'walter.sheltan@acme.com'" <walter.sheltan@acme.com>"
	regexp.MustCompile(`^\"(.+?)\"\s?\n?\s*[\[|<](.+?)[\]|>]`),     // ""Walter Sheltan" <walter.sheltan@acme.com>" or ""Walter Sheltan" [walter.sheltan@acme.com]" or ""walter.sheltan@acme.com" <walter.sheltan@acme.com>"
	regexp.MustCompile(`^([^,;]+?)\s?\n?\s*[\[|<](.+?)[\]|>]`),     // "Walter Sheltan <walter.sheltan@acme.com>" or "Walter Sheltan [walter.sheltan@acme.com]" or "walter.sheltan@acme.com <walter.sheltan@acme.com>"
	regexp.MustCompile(`^(.?)\s?\n?\s*[\[|<](.+?)[\]|>]`),          // "<walter.sheltan@acme.com>"
	regexp.MustCompile(`^([^\s@]+@[^\s@]+\.[^\s@,]+)`),             // "walter.sheltan@acme.com"
	regexp.MustCompile(`^([^;].+?)\s?\n?\s*[\[|<](.+?)[\]|>]`),     // "Walter, Sheltan <walter.sheltan@acme.com>" or "Walter, Sheltan [walter.sheltan@acme.com]"
}

var _MailboxAddress = []*regexp.Regexp{
	regexp.MustCompile(`^(([^\s@]+)@([^\s@]+)\.([^\s@]+))$`),
}
