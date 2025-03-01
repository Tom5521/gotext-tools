package cmd

import (
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
)

var (
	headerCfg   po.HeaderConfig
	compilerCfg compiler.PoConfig
)

func initConfig() {
	compilerCfg = compiler.PoConfig{}
}
