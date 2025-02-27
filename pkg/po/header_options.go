package po

type HeaderOption func(*HeaderConfig)

func WithHewaderConfig(c HeaderConfig) HeaderOption {
	return func(hc *HeaderConfig) {
		*hc = c
	}
}

func WithNplurals(n uint) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.Nplurals = n
	}
}

func WithProjectIDVersion(v string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.ProjectIDVersion = v
	}
}

func WithReportMsgidBugsTo(r string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.ReportMsgidBugsTo = r
	}
}

func WithLanguage(lang string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.Language = lang
	}
}

func WithLanguageTeam(team string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.LanguageTeam = team
	}
}

func WithLastTranslator(translator string) HeaderOption {
	return func(hc *HeaderConfig) {
		hc.LastTranslator = translator
	}
}
