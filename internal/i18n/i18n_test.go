package i18n

import (
	"encoding/json"
	"testing"
	"testing/fstest"
)

func TestLanguageSelection(t *testing.T) {
	for _, tt := range []struct{ name, preference, all, messages, lang, languages, want string }{
		{"default", "auto", "", "", "", "", "en"},
		{"Russian", "auto", "", "", "ru_RU.UTF-8", "", "ru"},
		{"manual", "en", "ru_RU", "", "ru_RU", "ru", "en"},
		{"messages", "auto", "", "ru_RU", "en_US", "", "ru"},
		{"all", "auto", "en_US", "ru_RU", "ru_RU", "", "en"},
		{"priority list", "auto", "", "", "de_DE", "fr:ru:en", "fr"},
		{"skip unavailable", "auto", "", "", "de_DE", "zz:ru:en", "ru"},
		{"Ukrainian", "auto", "", "", "uk_UA.UTF-8", "", "uk"},
		{"Polish", "auto", "", "", "pl_PL.UTF-8", "", "pl"},
		{"Portuguese", "auto", "", "", "pt_PT.UTF-8", "", "pt"},
		{"Brazilian regional fallback", "auto", "", "", "pt_BR.UTF-8", "", "pt"},
		{"Italian", "auto", "", "", "it_IT.UTF-8", "", "it"},
		{"Belarusian", "auto", "", "", "be_BY.UTF-8", "", "be"},
		{"Kazakh", "auto", "", "", "kk_KZ.UTF-8", "", "kk"},
		{"German", "auto", "", "", "de_DE.UTF-8", "", "de"},
		{"French regional fallback", "auto", "", "", "fr_CA.UTF-8", "", "fr"},
		{"Spanish regional fallback", "auto", "", "", "es_MX.UTF-8", "", "es"},
		{"manual German", "de", "ru_RU", "", "ru_RU", "ru", "de"},
		{"C ignores priority", "auto", "C", "", "ru_RU", "ru", "en"},
		{"POSIX", "auto", "POSIX", "", "", "ru", "en"},
		{"C UTF8", "auto", "C.UTF-8", "", "", "ru", "en"},
		{"unknown", "auto", "", "", "zz_ZZ", "", "en"},
		{"unknown manual", "zz", "", "", "ru_RU", "", "en"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			env := map[string]string{"LC_ALL": tt.all, "LC_MESSAGES": tt.messages, "LANG": tt.lang, "LANGUAGE": tt.languages}
			got := Builtin.Select(tt.preference, func(k string) string { return env[k] })
			if got.Code() != tt.want {
				t.Fatalf("language = %s, want %s", got.Code(), tt.want)
			}
		})
	}
}

func TestCatalogExtensionAndFallback(t *testing.T) {
	fsys := fstest.MapFS{}
	for name, cat := range map[string]Catalog{
		"en": {Code: "en", Name: "English", Direction: "ltr", Messages: map[string]string{"hello": "Hello", "pair": "%s / %d"}},
		"de": {Code: "de", Name: "Deutsch", Direction: "ltr", Messages: map[string]string{"pair": "%[2]d / %[1]s"}},
	} {
		data, err := json.Marshal(cat)
		if err != nil {
			t.Fatal(err)
		}
		fsys[name+".json"] = &fstest.MapFile{Data: data}
	}
	bundle, err := Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	tr := bundle.Select("de-DE", func(string) string { return "" })
	if tr.Code() != "de" || tr.Text("hello") != "Hello" || tr.Format("pair", "A", 3) != "3 / A" {
		t.Fatal("extension, fallback or parameter ordering failed")
	}
	if len(bundle.Languages()) != 2 {
		t.Fatal("catalog not discoverable")
	}
}

func TestInvalidCatalogs(t *testing.T) {
	for _, body := range []string{
		`{"code":"ru","name":"Русский","direction":"ltr","messages":{"x":"a","x":"b"}}`,
		`{"code":"ru","name":"Русский","direction":"ltr","messages":{"x":"%d"}}`,
		`{"code":"ru","name":"Русский","direction":"sideways","messages":{"x":"%s"}}`,
		`{"code":"ru","name":"Русский","direction":"ltr","messages":{"typo":"%s"}}`,
	} {
		fsys := fstest.MapFS{"en.json": {Data: []byte(`{"code":"en","name":"English","direction":"ltr","messages":{"x":"%s"}}`)}, "ru.json": {Data: []byte(body)}}
		if _, err := Load(fsys); err == nil {
			t.Errorf("accepted invalid catalog %s", body)
		}
	}
}

func TestShippedCatalogsComplete(t *testing.T) {
	base := Builtin.catalogs["en"]
	for code, cat := range Builtin.catalogs {
		for key := range base.Messages {
			if cat.Messages[key] == "" {
				t.Errorf("%s missing %s", code, key)
			}
		}
	}
}
