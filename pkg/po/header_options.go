package po

type HeaderOption func(*HeaderConfig)

func HeaderWithConfig(c HeaderConfig) HeaderOption {
	return func(hc *HeaderConfig) {
		*hc = c
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
