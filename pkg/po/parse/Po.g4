grammar Po;

start: entry*;

entry: 
    comment*
    msgctxt?
    msgid
    (
        msgstr
        | (plural_msgid plural_msgstr+)
    )
    NL?
    ;



msgctxt: MSGCTXT string;
msgid: MSGID string;
msgstr: MSGSTR string;
plural_msgid: PLURAL_MSGID string;
plural_msgstr: PLURAL_MSGSTR string;
string: (STRING '\n'?)+;
comment: COMMENT
    | FLAG_COMMENT
    | EXTRACTED_COMMENT
    | REFERENCE_COMMENT
    | PREVIOUS_COMMENT
    ;


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
FLAG_COMMENT: '#,' ~[\n]*;
EXTRACTED_COMMENT: '#.' ~[\n]*;
REFERENCE_COMMENT: '#:' ~[\n]*;
PREVIOUS_COMMENT: '#|' ~[\n]*;
