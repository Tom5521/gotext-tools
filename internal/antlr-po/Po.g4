grammar Po;

start: entry*;

entry:
    comment*
    msgctxt?
    comment*
    msgid
    comment*
    (
        msgstr
        | (
            plural_msgid
            comment*
            plural_msgstr+
        )
    )
    comment*?
    NL?
    ;



msgctxt: MSGCTXT string;
msgid: MSGID string;
msgstr: MSGSTR string;
plural_msgid: PLURAL_MSGID string;
plural_msgstr: PLURAL_MSGSTR string;
string: (STRING '\n'?)+;
comment: COMMENT;


/*
* Lexer tokens.
*/

WS: [\n\r\t ] -> skip;
INT: [0-9]+;
STRING: 
    '"' (~'"' | '\\"')* '"' |
    '\'' (~'\'' | '\\\'')*;
NL: '\n';

MSGCTXT: 'msgctxt';
MSGID: 'msgid';
MSGSTR: 'msgstr';
PLURAL_MSGID: 'msgid_plural';
PLURAL_MSGSTR: 'msgstr['INT']';
COMMENT: '#' ~[\n]*;
