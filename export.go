package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed template.pot
var PotHeader string

func ExportMsgIDs(msgids []MsgID) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(PotHeader, proyectVersion, language))

	for _, msgid := range msgids {
		builder.WriteString(fmt.Sprintf("#: %s:%d\n", msgid.File, msgid.Line))
		builder.WriteString(fmt.Sprintf(`msgid "%s"`+"\n", msgid.ID))
		builder.WriteString(`msgstr ""` + "\n\n")
	}

	return builder.String()
}
