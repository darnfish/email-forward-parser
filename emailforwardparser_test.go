package emailforwardparser

import (
	"log"
	"os"
	"strings"
	"testing"
)

var _Subject = "Integer consequat non purus"
var _Body = "Aenean quis diam urna. Maecenas eleifend vulputate ligula ac consequat. Pellentesque cursus tincidunt mauris non venenatis.\nSed nec facilisis tellus. Nunc eget eros quis ex congue iaculis nec quis massa. Morbi in nisi tincidunt, euismod ante eget, eleifend nisi.\n\nPraesent ac ligula orci. Pellentesque convallis suscipit mi, at congue massa sagittis eget."
var _Message = "Praesent suscipit egestas hendrerit.\n\nAliquam eget dui dui."

var _FromAddress = "john.doe@acme.com"
var _FromName = "John Doe"

var _ToAddress1 = "bessie.berry@acme.com"
var _ToName1 = "Bessie Berry"
var _ToAddress2 = "suzanne@globex.corp"
var _ToName2 = "Suzanne"

var _CcAddress1 = "walter.sheltan@acme.com"
var _CcName1 = "Walter Sheltan"
var _CcAddress2 = "nicholas@globex.corp"
var _CcName2 = "Nicholas"

func _ReadAndParse(emailFile string, subjectFile string) ReadResult {
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

	if result.Email.Subject != _Subject {
		t.Error(entryName, "result.Email.Subject != _Subject, result.Email.Subject=", result.Email.Subject, "_Subject=", _Subject)
	}

	if !skipBody {
		if result.Email.Body != _Body {
			t.Error(entryName, "result.Email.Body != _Body")
		}
	}

	// test.strictEqual(typeof email.date, "string");
	// test.strictEqual((email.date || "").length > 1, true);

	// Opposite of > 1
	if !(len(result.Email.Date) > 1) {
		t.Error(entryName, "!(len(result.Email.Date) > 1)", result.Email.Date)
	}

	if !skipFrom {
		if result.Email.From.Name != _FromName {
			t.Error(entryName, "result.Email.From.Name != _FromName", result.Email.From.Name, _FromName)
		}

		if result.Email.From.Address != _FromAddress {
			t.Error(entryName, "result.Email.From.Address != _FromAddress", result.Email.From.Address, _FromAddress)
		}
	}

	if !skipTo {
		if len(result.Email.To) > 0 {
			if len(result.Email.To[0].Name) > 0 {
				t.Error(entryName, "len(result.Email.To[0].Name) > 0")
			}

			if result.Email.To[0].Address != _ToAddress1 {
				t.Error(entryName, "result.Email.To[0].Address != _ToAddress1")
			}
		} else {
			t.Error(entryName, "len(result.Email.To) == 0, expected", _ToName1, _ToAddress1)
		}
	}

	if !skipCc {
		if len(result.Email.CC) > 0 {
			if result.Email.CC[0].Name != _CcName1 {
				t.Error(entryName, "result.Email.CC[0].Name != _CcName1")
			}

			if result.Email.CC[0].Address != _CcAddress1 {
				t.Error(entryName, "result.Email.CC[0].Address != _CcAddress1")
			}

			if result.Email.CC[1].Name != _CcName2 {
				t.Error(entryName, "result.Email.CC[1].Name != _CcName2")
			}

			if result.Email.CC[1].Address != _CcAddress2 {
				t.Error(entryName, "result.Email.CC[1].Address != _CcAddress2")
			}
		} else {
			t.Error(entryName, "len(result.Email.CC) == 0, expected", _CcName1, _CcAddress1, _CcName2, _CcAddress2)
		}
	}

	if !skipMessage {
		if result.Message != _Message {
			t.Error(entryName, "result.Message != _Message")
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
