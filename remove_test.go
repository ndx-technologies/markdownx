package markdownx_test

import (
	"testing"

	"github.com/ndx-technologies/markdownx"
)

func TestRemove(t *testing.T) {
	tests := map[string]string{
		"leite e ovos":                          "leite e ovos",
		"   ":                                   "",
		"one\ntwo   three":                      "one\ntwo three",
		"leite\n\novos":                         "leite\n\novos",
		"**Bold** and `code`":                   "Bold and code",
		"## Summary\nthe rest":                  "Summary\nthe rest",
		"- leite\n- ovos":                       "leite\novos",
		"1. leite\n2. ovos":                     "leite\novos",
		"- [ ] leite\n- [x] ovos":               "leite\novos",
		"- [X] ovos":                            "ovos",
		"- [-] leite\n- [+] ovos\n- [?] pão":    "leite\novos\npão",
		"+ leite\n• ovos":                       "leite\novos",
		"> leite\n> ovos":                       "leite\novos",
		"*leite* e ~~ovos~~":                    "leite e ovos",
		"_leite_ e _ovos_":                      "leite e ovos",
		"snake_case and a_b":                    "snake_case and a_b",
		"```\nscanner.go\n```":                  "scanner.go",
		"~~~\nscanner.go\n~~~":                  "scanner.go",
		"```\nleite * ovos\n```":                "leite * ovos",
		"leite\n---\novos":                      "leite\n\novos",
		"| A | B |\n| - | - |\n| leite | 5 |":   "A B\nleite 5",
		"| preço | 10-20 |":                     "preço 10-20",
		"Hello[^1]":                             "Hello",
		"leite [^1":                             "leite [^1",
		"Título\n====\nresto":                   "Título\n\nresto",
		"## Summary ##":                         "Summary",
		"C# top":                                "C# top",
		"## C# top":                             "C# top",
		"leite\n[1]: https://ndx.one/x\novos":   "leite\novos",
		"[^1]: nota":                            "",
		"[leite":                                "[leite",
		"veja [isso]: é caro":                   "veja [isso]: é caro",
		"veja [recibo](mm://receipt/abc) agora": "veja recibo agora",
		"veja [recibo][1]":                      "veja recibo",
		"![foto](mm://image/abc)":               "foto",
		"![](mm://image/abc)":                   "",
		"[1](https://ndx.one/budgets) veja":     "1 veja",
		"-5% hoje\n1.5 kg":                      "-5% hoje\n1.5 kg",
		"2026. ano novo":                        "2026. ano novo",
		"42":                                    "42",
		"#Hora de escanear":                     "#Hora de escanear",
	}
	for md, exp := range tests {
		if s := markdownx.Remove(md); s != exp {
			t.Error(md, exp, s)
		}
	}
}
