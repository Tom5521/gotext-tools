package po

type HeaderOption func(*HeaderConfig)

func HeaderWithConfig(c HeaderConfig) HeaderOption {
	return func(hc *HeaderConfig) {
		*hc = c
	}
}

func HeaderWithTemplate(t bool) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.Template = t
	}
}

func HeaderWithPlural(p string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.Plural = p
	}
}

func HeaderWithPOTCreationDate(c string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.POTCreationDate = c
	}
}

func HeaderWithPORevisionDate(r string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.PORevisionDate = r
	}
}

func HeaderWithXGenerator(g string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.XGenerator = g
	}
}

func HeaderWtihContentTransferEncoding(c string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.ContentTransferEncoding = c
	}
}

func HeaderWithMediaType(t string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.MediaType = t
	}
}

func HeaderWithCharset(charset string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.Charset = charset
	}
}

func HeaderWithMimeVersion(version string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.MimeVersion = version
	}
}

func HeaderWithNplurals(n uint) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.Nplurals = n
	}
}

func HeaderWithProjectIDVersion(v string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.ProjectIDVersion = v
	}
}

func HeaderWithReportMsgidBugsTo(r string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.ReportMsgidBugsTo = r
	}
}

func HeaderWithLanguage(lang string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.Language = lang
	}
}

func HeaderWithLanguageTeam(team string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.LanguageTeam = team
	}
}

func HeaderWithLastTranslator(translator string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.LastTranslator = translator
	}
}
