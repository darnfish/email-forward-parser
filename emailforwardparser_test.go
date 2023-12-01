package emailforwardparser

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	regexp "github.com/wasilibs/go-re2"
)

var _TestSubject = "Integer consequat non purus"
var _TestBody = "Aenean quis diam urna. Maecenas eleifend vulputate ligula ac consequat. Pellentesque cursus tincidunt mauris non venenatis.\nSed nec facilisis tellus. Nunc eget eros quis ex congue iaculis nec quis massa. Morbi in nisi tincidunt, euismod ante eget, eleifend nisi.\n\nPraesent ac ligula orci. Pellentesque convallis suscipit mi, at congue massa sagittis eget."
var _TestMessage = "Praesent suscipit egestas hendrerit.\n\nAliquam eget dui dui."

var _TestFromAddress = "john.doe@acme.com"
var _TestFromName = "John Doe"

var _TestToAddress1 = "bessie.berry@acme.com"
var _TestToName1 = "Bessie Berry"
var _TestToAddress2 = "suzanne@globex.corp"
var _TestToName2 = "Suzanne"

var _TestCcAddress1 = "walter.sheltan@acme.com"
var _TestCcName1 = "Walter Sheltan"
var _TestCcAddress2 = "nicholas@globex.corp"
var _TestCcName2 = "Nicholas"

func _Read(emailFile string, subjectFile string) (string, string) {
	var email string
	var subject string

	emailBytes, err := os.ReadFile("./fixtures/" + emailFile + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	email = string(emailBytes)

	if len(subjectFile) > 0 {
		subjectBytes, err := os.ReadFile("./fixtures/" + subjectFile + ".txt")
		if err != nil {
			log.Fatal(err)
		}

		subject = string(subjectBytes)
	}

	return email, subject
}

func _ReadAndParse(emailFile string, subjectFile string) ReadResult {
	email, subject := _Read(emailFile, subjectFile)

	return Read(email, subject)
}

func _LoopTests(entries []string, testFn func(ReadResult, string)) {
	for _, entry := range entries {
		var result ReadResult
		var entryName string

		if strings.Contains(entry, ",") {
			entryParts := strings.Split(entry, ",")

			result = _ReadAndParse(entryParts[0], entryParts[1])
			entryName = entryParts[0]
		} else {
			result = _ReadAndParse(entry, "")
			entryName = entry
		}

		testFn(result, entryName)
	}
}

func _TestEmail(t *testing.T, result ReadResult, entryName string, skipFrom bool, skipTo bool, skipCc bool, skipMessage bool, skipBody bool) {
	if result.Forwarded != true {
		t.Error(entryName, "result.Forwarded != true", result.Forwarded)
	}

	if result.Email.Subject != _TestSubject {
		t.Error(entryName, "result.Email.Subject != _TestSubject, result.Email.Subject=", result.Email.Subject, "_TestSubject=", _TestSubject)
	}

	if !skipBody {
		if result.Email.Body != _TestBody {
			t.Error(entryName, "result.Email.Body != _Body", result.Email.Body)
		}
	}

	// test.strictEqual(typeof email.date, "string");
	// test.strictEqual((email.date || "").length > 1, true);

	// Opposite of > 1
	if !(len(result.Email.Date) > 1) {
		t.Error(entryName, "!(len(result.Email.Date) > 1)", result.Email.Date)
	}

	if !skipFrom {
		if result.Email.From.Name != _TestFromName {
			t.Error(entryName, "result.Email.From.Name != _TestFromName", result.Email.From.Name, _TestFromName)
		}

		if result.Email.From.Address != _TestFromAddress {
			t.Error(entryName, "result.Email.From.Address != _TestFromAddress", result.Email.From.Address, _TestFromAddress)
		}
	}

	if !skipTo {
		if len(result.Email.To) > 0 {
			if len(result.Email.To[0].Name) > 0 {
				t.Error(entryName, "len(result.Email.To[0].Name) > 0")
			}

			if result.Email.To[0].Address != _TestToAddress1 {
				t.Error(entryName, "result.Email.To[0].Address != _TestToAddress1")
			}
		} else {
			t.Error(entryName, "len(result.Email.To) == 0, expected", _TestToName1, _TestToAddress1)
		}
	}

	if !skipCc {
		if len(result.Email.CC) > 0 {
			if result.Email.CC[0].Name != _TestCcName1 {
				t.Error(entryName, "result.Email.CC[0].Name != _TestCcName1")
			}

			if result.Email.CC[0].Address != _TestCcAddress1 {
				t.Error(entryName, "result.Email.CC[0].Address != _TestCcAddress1")
			}

			if result.Email.CC[1].Name != _TestCcName2 {
				t.Error(entryName, "result.Email.CC[1].Name != _TestCcName2")
			}

			if result.Email.CC[1].Address != _TestCcAddress2 {
				t.Error(entryName, "result.Email.CC[1].Address != _TestCcAddress2")
			}
		} else {
			t.Error(entryName, "len(result.Email.CC) == 0, expected", _TestCcName1, _TestCcAddress1, _TestCcName2, _TestCcAddress2)
		}
	}

	if !skipMessage {
		if result.Message != _TestMessage {
			t.Error(entryName, "result.Message != _TestMessage")
		}
	}
}

func TestCommon(t *testing.T) {
	_LoopTests([]string{
		"apple_mail_cs_body",
		"apple_mail_da_body",
		"apple_mail_de_body",
		"apple_mail_en_body",
		"apple_mail_es_body",
		"apple_mail_fi_body",
		"apple_mail_fr_body",
		"apple_mail_hr_body",
		"apple_mail_hu_body",
		"apple_mail_it_body",
		"apple_mail_nl_body",
		"apple_mail_no_body",
		"apple_mail_pl_body",
		"apple_mail_pt_br_body",
		"apple_mail_pt_body",
		"apple_mail_ro_body",
		"apple_mail_ru_body",
		"apple_mail_sk_body",
		"apple_mail_sv_body",
		"apple_mail_tr_body",
		"apple_mail_uk_body",

		"gmail_cs_body",
		"gmail_da_body",
		"gmail_de_body",
		"gmail_en_body",
		"gmail_es_body",
		"gmail_et_body",
		"gmail_fi_body",
		"gmail_fr_body",
		"gmail_hr_body",
		"gmail_hu_body",
		"gmail_it_body",
		"gmail_nl_body",
		"gmail_no_body",
		"gmail_pl_body",
		"gmail_pt_br_body",
		"gmail_pt_body",
		"gmail_ro_body",
		"gmail_ru_body",
		"gmail_sk_body",
		"gmail_sv_body",
		"gmail_tr_body",
		"gmail_uk_body",

		"hubspot_de_body",
		"hubspot_en_body",
		"hubspot_es_body",
		"hubspot_fi_body",
		"hubspot_fr_body",
		"hubspot_it_body",
		"hubspot_ja_body",
		"hubspot_nl_body",
		"hubspot_pl_body",
		"hubspot_pt_br_body",
		"hubspot_sv_body",

		"ionos_one_and_one_en_body",

		"missive_en_body",

		"outlook_live_body,outlook_live_cs_subject",
		"outlook_live_body,outlook_live_da_subject",
		"outlook_live_body,outlook_live_de_subject",
		"outlook_live_body,outlook_live_en_subject",
		"outlook_live_body,outlook_live_es_subject",
		"outlook_live_body,outlook_live_fr_subject",
		"outlook_live_body,outlook_live_hr_subject",
		"outlook_live_body,outlook_live_hu_subject",
		"outlook_live_body,outlook_live_it_subject",
		"outlook_live_body,outlook_live_nl_subject",
		"outlook_live_body,outlook_live_no_subject",
		"outlook_live_body,outlook_live_pl_subject",
		"outlook_live_body,outlook_live_pt_br_subject",
		"outlook_live_body,outlook_live_pt_subject",
		"outlook_live_body,outlook_live_ro_subject",
		"outlook_live_body,outlook_live_sk_subject",
		"outlook_live_body,outlook_live_sv_subject",

		"outlook_2013_en_body,outlook_2013_en_subject",

		"new_outlook_2019_cs_body,new_outlook_2019_cs_subject",
		"new_outlook_2019_da_body,new_outlook_2019_da_subject",
		"new_outlook_2019_de_body,new_outlook_2019_de_subject",
		"new_outlook_2019_en_body,new_outlook_2019_en_subject",
		"new_outlook_2019_es_body,new_outlook_2019_es_subject",
		"new_outlook_2019_fi_body,new_outlook_2019_fi_subject",
		"new_outlook_2019_fr_body,new_outlook_2019_fr_subject",
		"new_outlook_2019_hu_body,new_outlook_2019_hu_subject",
		"new_outlook_2019_it_body,new_outlook_2019_it_subject",
		"new_outlook_2019_nl_body,new_outlook_2019_nl_subject",
		"new_outlook_2019_no_body,new_outlook_2019_no_subject",
		"new_outlook_2019_pl_body,new_outlook_2019_pl_subject",
		"new_outlook_2019_pt_br_body,new_outlook_2019_pt_br_subject",
		"new_outlook_2019_pt_body,new_outlook_2019_pt_subject",
		"new_outlook_2019_ru_body,new_outlook_2019_ru_subject",
		"new_outlook_2019_sk_body,new_outlook_2019_sk_subject",
		"new_outlook_2019_sv_body,new_outlook_2019_sv_subject",
		"new_outlook_2019_tr_body,new_outlook_2019_tr_subject",

		"outlook_2019_cz_body,outlook_2019_subject",
		"outlook_2019_da_body,outlook_2019_subject",
		"outlook_2019_de_body,outlook_2019_subject",
		"outlook_2019_en_body,outlook_2019_subject",
		"outlook_2019_es_body,outlook_2019_subject",
		"outlook_2019_fi_body,outlook_2019_subject",
		"outlook_2019_fr_body,outlook_2019_subject",
		"outlook_2019_hu_body,outlook_2019_subject",
		"outlook_2019_it_body,outlook_2019_subject",
		"outlook_2019_nl_body,outlook_2019_subject",
		"outlook_2019_no_body,outlook_2019_subject",
		"outlook_2019_pl_body,outlook_2019_subject",
		"outlook_2019_pt_body,outlook_2019_subject",
		"outlook_2019_ru_body,outlook_2019_subject",
		"outlook_2019_sk_body,outlook_2019_subject",
		"outlook_2019_sv_body,outlook_2019_subject",
		"outlook_2019_tr_body,outlook_2019_subject",

		"thunderbird_cs_body",
		"thunderbird_da_body",
		"thunderbird_de_body",
		"thunderbird_en_body",
		"thunderbird_es_body",
		"thunderbird_fi_body",
		"thunderbird_fr_body",
		"thunderbird_hr_body",
		"thunderbird_hu_body",
		"thunderbird_it_body",
		"thunderbird_nl_body",
		"thunderbird_no_body",
		"thunderbird_pl_body",
		"thunderbird_pt_br_body",
		"thunderbird_pt_body",
		"thunderbird_ro_body",
		"thunderbird_ru_body",
		"thunderbird_sk_body",
		"thunderbird_sv_body",
		"thunderbird_tr_body",
		"thunderbird_uk_body",

		"yahoo_cs_body",
		"yahoo_da_body",
		"yahoo_de_body",
		"yahoo_en_body",
		"yahoo_es_body",
		"yahoo_fi_body",
		"yahoo_fr_body",
		"yahoo_hu_body",
		"yahoo_it_body",
		"yahoo_nl_body",
		"yahoo_no_body",
		"yahoo_pl_body",
		"yahoo_pt_body",
		"yahoo_pt_br_body",
		"yahoo_ro_body",
		"yahoo_ru_body",
		"yahoo_sk_body",
		"yahoo_sv_body",
		"yahoo_tr_body",
		"yahoo_uk_body",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, false, strings.HasPrefix(entryName, "outlook_2019_"), strings.HasPrefix(entryName, "outlook_2019_") || strings.HasPrefix(entryName, "ionos_one_and_one_"), true, false)

		if len(result.Message) > 0 {
			t.Error("result.Message != null")
		}
	})
}

func TestAlternative1(t *testing.T) {
	_LoopTests([]string{
		"apple_mail_en_body_variant_1",
		"gmail_en_body_variant_1",
		"hubspot_en_body_variant_1",
		"missive_en_body_variant_1",
		"outlook_live_en_body_variant_1,outlook_live_en_subject",
		"new_outlook_2019_en_body_variant_1,new_outlook_2019_en_subject",
		"yahoo_en_body_variant_1",
		"thunderbird_en_body_variant_1",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, false, true, true, true, false)

		if len(result.Email.To[0].Name) > 0 {
			t.Error(entryName)
		}

		if result.Email.To[0].Address != _TestToAddress1 {
			t.Error(entryName)
		}

		if len(result.Email.To[1].Name) > 0 {
			t.Error(entryName)
		}

		if result.Email.To[1].Address != _TestToAddress2 {
			t.Error(entryName)
		}

		if len(result.Email.CC) > 0 {
			t.Error(entryName)
		}
	})
}

func TestAlternative2(t *testing.T) {
	_LoopTests([]string{
		"apple_mail_en_body_variant_2",
		"gmail_en_body_variant_2",
		"hubspot_en_body_variant_2",
		"ionos_one_and_one_en_body_variant_2",
		"missive_en_body_variant_2",
		"outlook_live_en_body_variant_2,outlook_live_en_subject",
		"new_outlook_2019_en_body_variant_2,new_outlook_2019_en_subject",
		"outlook_2019_en_body_variant_2,outlook_2019_subject",
		"yahoo_en_body_variant_2",
		"thunderbird_en_body_variant_2",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, false, true, entryName == "outlook_2019_en_body_variant_2" || entryName == "ionos_one_and_one_en_body_variant_2", false, false)

		switch entryName {
		case "outlook_2019_en_body_variant_2":
		default:
			if result.Email.To[0].Address != _TestToAddress1 {
				t.Error(entryName)
			}

			if result.Email.To[1].Address != _TestToAddress2 {
				t.Error(entryName)
			}
		}
	})
}

func TestAlternative3(t *testing.T) {
	_LoopTests([]string{
		"apple_mail_en_body_variant_3",
		"gmail_en_body_variant_3",
		"missive_en_body_variant_3",
		"outlook_live_en_body_variant_3,outlook_live_en_subject",
		"new_outlook_2019_en_body_variant_3,new_outlook_2019_en_subject",
		"yahoo_en_body_variant_3",
		"thunderbird_en_body_variant_3",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, false, true, true, true, false)

		if result.Email.To[0].Name != _TestToName1 {
			t.Error(entryName)
		}

		if result.Email.To[0].Address != _TestToAddress1 {
			t.Error(entryName)
		}

		if len(result.Email.To[1].Name) > 0 {
			t.Error(entryName)
		}

		if result.Email.To[1].Address != _TestToAddress2 {
			t.Error(entryName)
		}

		if len(result.Email.CC[0].Name) > 0 {
			t.Error(entryName)
		}

		if result.Email.CC[0].Address != _TestCcAddress1 {
			t.Error(entryName)
		}

		if result.Email.CC[1].Name != _TestCcName2 {
			t.Error(entryName)
		}

		if result.Email.CC[1].Address != _TestCcAddress2 {
			t.Error(entryName)
		}
	})
}

func TestAlternative4(t *testing.T) {
	_LoopTests([]string{
		"apple_mail_en_body_variant_4",
		"gmail_en_body_variant_4",
		"hubspot_en_body_variant_4",
		"missive_en_body_variant_4",
		"outlook_live_en_body_variant_4,outlook_live_en_subject_variant_4",
		"new_outlook_2019_en_body_variant_4,new_outlook_2019_en_subject_variant_4",
		"outlook_2019_en_body_variant_4,outlook_2019_en_subject_variant_4",
		"yahoo_en_body_variant_4",
		"thunderbird_en_body_variant_4",
	}, func(result ReadResult, entryName string) {
		if result.Forwarded {
			t.Error(entryName)
		}
	})
}

func TestAlternative5(t *testing.T) {
	_LoopTests([]string{
		"apple_mail_en_body_variant_5",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, true, false, false, true, false)

		if len(result.Email.From.Name) > 0 {
			t.Error(entryName)
		}

		if result.Email.From.Address != _TestFromAddress {
			t.Error(entryName)
		}
	})
}

func TestAlternative6(t *testing.T) {
	_LoopTests([]string{
		"apple_mail_en_body_variant_6",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, false, false, true, true, false)
	})
}

func TestAlternative7(t *testing.T) {
	_LoopTests([]string{
		"apple_mail_en_body_variant_7",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, false, true, true, true, false)

		if result.Email.To[0].Name != "Bessie, Berry" {
			t.Error(entryName)
		}

		if result.Email.To[0].Address != _TestToAddress1 {
			t.Error(entryName)
		}

		if result.Email.To[1].Name != _TestToName2 {
			t.Error(entryName)
		}

		if result.Email.To[1].Address != _TestToAddress2 {
			t.Error(entryName)
		}

		if len(result.Email.CC[0].Name) > 0 {
			t.Error(entryName)
		}

		if result.Email.CC[0].Address != _TestCcAddress1 {
			t.Error(entryName)
		}

		if len(result.Email.CC[1].Name) > 0 {
			t.Error(entryName)
		}

		if result.Email.CC[1].Address != _TestCcAddress2 {
			t.Error(entryName)
		}
	})
}

func TestAlternative8(t *testing.T) {
	_LoopTests([]string{
		"outlook_live_en_body_variant_8",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, false, false, false, false, true)

		email, _ := _Read(entryName, "")
		separator := fmt.Sprintf("Subject: %s\n", _TestSubject)

		body := trimString(strings.Split(email, separator)[1])

		if result.Email.Body != body {
			t.Error(entryName)
		}
	})
}

func TestAlternative9(t *testing.T) {
	_LoopTests([]string{
		"outlook_live_en_body_variant_9",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, false, false, false, false, false)
	})
}

func TestAlternative10(t *testing.T) {
	_LoopTests([]string{
		"outlook_live_en_body_variant_10,outlook_live_en_subject_variant_10",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, false, false, false, false, false)
	})
}

func TestAlternative11(t *testing.T) {
	_LoopTests([]string{
		"outlook_live_en_body_variant_11",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, true, true, true, true, false)

		if result.Email.From.Name != "John, Doe" {
			t.Error(entryName)
		}

		if result.Email.From.Address != _TestFromAddress {
			t.Error(entryName)
		}

		if result.Email.To[0].Name != "Bessie, Berry" {
			t.Error(entryName)
		}

		if result.Email.To[0].Address != _TestToAddress1 {
			t.Error(entryName)
		}

		if len(result.Email.To[1].Name) > 0 {
			t.Error(entryName)
		}

		if result.Email.To[1].Address != _TestToAddress2 {
			t.Error(entryName)
		}

		if result.Email.CC[0].Name != "Walter, Sheltan" {
			t.Error(entryName)
		}

		if result.Email.CC[0].Address != _TestCcAddress1 {
			t.Error(entryName)
		}

		if result.Email.CC[1].Name != "Nicholas, Landers" {
			t.Error(entryName)
		}

		if result.Email.CC[1].Address != _TestCcAddress2 {
			t.Error(entryName)
		}
	})
}

func TestAlternative12(t *testing.T) {
	_LoopTests([]string{
		"unknown_en_body_variant_12,unknown_en_subject",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, true, true, true, true, false)

		if result.Email.From.Name != _TestFromName {
			t.Error(entryName)
		}

		if len(result.Email.From.Address) > 0 {
			t.Error(entryName)
		}

		if result.Email.To[0].Name != _TestToName1 {
			t.Error(entryName)
		}

		if len(result.Email.To[0].Address) > 0 {
			t.Error(entryName)
		}
	})
}

func TestAlternative13(t *testing.T) {
	_LoopTests([]string{
		"apple_mail_en_body_variant_13",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, false, false, false, false, true)

		email, _ := _Read(entryName, "")
		separator := fmt.Sprintf("Subject: %s\n", _TestSubject)

		body := strings.Split(email, separator)[1]

		body = regexp.MustCompile(`(?m)^(>+)\s?$`).ReplaceAllString(body, "")
		body = regexp.MustCompile(`(?m)^(>+)\s?`).ReplaceAllString(body, "")
		body = preprocessString(body) // original test is .trim(), but weird js stuff

		if result.Email.Body != body {
			t.Error(entryName)
		}

		if !strings.HasPrefix(result.Email.Body, _TestBody) {
			t.Error(entryName)
		}
	})
}

func TestAlternative14(t *testing.T) {
	_LoopTests([]string{
		"gmail_en_body_variant_14",
		"outlook_live_en_body_variant_14,outlook_live_en_subject",
		"new_outlook_2019_en_body_variant_14,new_outlook_2019_fr_subject",
		"new_outlook_2019_en_body_variant_14_1,new_outlook_2019_fr_subject",
		"thunderbird_en_body_variant_14",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, false, false, false, false, true)

		email, _ := _Read(entryName, "")
		separator := fmt.Sprintf("Subject: %s\n", _TestSubject)

		switch entryName {
		case "thunderbird_en_body_variant_14":
			separator = fmt.Sprintf("CC: 	%s <%s>, %s <%s>\n", _TestCcName1, _TestCcAddress1, _TestCcName2, _TestCcAddress2)
		}

		body := strings.Split(email, separator)[1]
		body = preprocessString(body) // original test is .trim(), but weird js stuff

		if result.Email.Body != body {
			t.Error(entryName)
		}
	})
}

func TestAlternative15(t *testing.T) {
	_LoopTests([]string{
		"gmail_en_body_variant_15",
		"outlook_live_en_body_variant_15,outlook_live_en_subject",
		"new_outlook_2019_en_body_variant_15,new_outlook_2019_fr_subject",
		"thunderbird_en_body_variant_15",
	}, func(result ReadResult, entryName string) {
		_TestEmail(t, result, entryName, false, false, false, false, true)

		email, _ := _Read(entryName, "")
		separator := fmt.Sprintf("Subject: %s\n", _TestSubject)

		switch entryName {
		case "thunderbird_en_body_variant_15":
			separator = fmt.Sprintf("CC: 	%s <%s>, %s <%s>\n", _TestCcName1, _TestCcAddress1, _TestCcName2, _TestCcAddress2)
		}

		body := strings.Split(email, separator)[1]
		body = preprocessString(body) // original test is .trim(), but weird js stuff

		if result.Email.Body != body {
			t.Error(entryName)
		}
	})
}
