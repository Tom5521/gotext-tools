package parse_test

import (
	"testing"

	"github.com/Tom5521/xgotext/pkg/po/parse"
)

func TestPoParser(t *testing.T) {
	input := `# 
# Copyright (C) 
# This file is distributed under the same license as the PACKAGE NAME package.
#
msgid ""
msgstr ""

msgctxt "My context :3"
msgid "id1"
msgstr "id1"

msgid "id2"
msgid_plural "helooows"
msgstr[0] "Holanda"
msgstr[1] "Holandas"

msgid "id3"
msgstr "id3"`

	parser := parse.NewPoFromString(input, "test.po")
	parser.Parse()

	if parser.Error() != nil {
		t.Error(parser.Error())
		return
	}
}
