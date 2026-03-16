package main

func (c *CTags) LspSymbolKind(kind string) int {
	k, ok := symbolKindByTagKind[kind]
	if !ok {
		return 0
	}
	return k
}

// Numeric values match LSP 3.17 `SymbolKind`.
const (
	SymbolKindFile          = 1
	SymbolKindModule        = 2
	SymbolKindNamespace     = 3
	SymbolKindPackage       = 4
	SymbolKindClass         = 5
	SymbolKindMethod        = 6
	SymbolKindProperty      = 7
	SymbolKindField         = 8
	SymbolKindConstructor   = 9
	SymbolKindEnum          = 10
	SymbolKindInterface     = 11
	SymbolKindFunction      = 12
	SymbolKindVariable      = 13
	SymbolKindConstant      = 14
	SymbolKindString        = 15
	SymbolKindNumber        = 16
	SymbolKindBoolean       = 17
	SymbolKindArray         = 18
	SymbolKindObject        = 19
	SymbolKindKey           = 20
	SymbolKindNull          = 21
	SymbolKindEnumMember    = 22
	SymbolKindStruct        = 23
	SymbolKindEvent         = 24
	SymbolKindOperator      = 25
	SymbolKindTypeParameter = 26
)

var symbolKindByTagKind = map[string]int{
	"alias":            SymbolKindVariable,
	"arg":              SymbolKindVariable,
	"attribute":        SymbolKindProperty,
	"boolean":          SymbolKindConstant,
	"callback":         SymbolKindFunction,
	"category":         SymbolKindEnum,
	"ccflag":           SymbolKindConstant,
	"cell":             SymbolKindVariable,
	"class":            SymbolKindClass,
	"collection":       SymbolKindClass,
	"command":          SymbolKindFunction,
	"component":        SymbolKindStruct,
	"config":           SymbolKindConstant,
	"const":            SymbolKindConstant,
	"constant":         SymbolKindConstant,
	"constructor":      SymbolKindConstructor,
	"context":          SymbolKindVariable,
	"counter":          SymbolKindVariable,
	"data":             SymbolKindVariable,
	"dataset":          SymbolKindVariable,
	"def":              SymbolKindFunction,
	"define":           SymbolKindConstant,
	"delegate":         SymbolKindClass,
	"enum":             SymbolKindEnum,
	"enumConstant":     SymbolKindEnumMember,
	"enumerator":       SymbolKindEnum,
	"environment":      SymbolKindVariable,
	"error":            SymbolKindEnum,
	"event":            SymbolKindEvent,
	"exception":        SymbolKindClass,
	"externvar":        SymbolKindVariable,
	"face":             SymbolKindInterface,
	"feature":          SymbolKindProperty,
	"field":            SymbolKindField,
	"fn":               SymbolKindFunction,
	"fun":              SymbolKindFunction,
	"func":             SymbolKindFunction,
	"function":         SymbolKindFunction,
	"functionVar":      SymbolKindVariable,
	"functor":          SymbolKindClass,
	"generic":          SymbolKindTypeParameter,
	"getter":           SymbolKindMethod,
	"global":           SymbolKindVariable,
	"globalVar":        SymbolKindVariable,
	"group":            SymbolKindEnum,
	"guard":            SymbolKindVariable,
	"handler":          SymbolKindFunction,
	"icon":             SymbolKindEnum,
	"id":               SymbolKindVariable,
	"implementation":   SymbolKindClass,
	"index":            SymbolKindVariable,
	"infoitem":         SymbolKindVariable,
	"instance":         SymbolKindVariable,
	"interface":        SymbolKindInterface,
	"it":               SymbolKindVariable,
	"jurisdiction":     SymbolKindVariable,
	"library":          SymbolKindModule,
	"list":             SymbolKindVariable,
	"local":            SymbolKindVariable,
	"localVariable":    SymbolKindVariable,
	"locale":           SymbolKindVariable,
	"localvar":         SymbolKindVariable,
	"macro":            SymbolKindVariable,
	"macroParameter":   SymbolKindVariable,
	"macrofile":        SymbolKindFile,
	"macroparam":       SymbolKindVariable,
	"makefile":         SymbolKindFile,
	"map":              SymbolKindVariable,
	"method":           SymbolKindMethod,
	"methodSpec":       SymbolKindMethod,
	"misc":             SymbolKindVariable,
	"module":           SymbolKindModule,
	"name":             SymbolKindVariable,
	"namespace":        SymbolKindModule,
	"nettype":          SymbolKindTypeParameter,
	"newFile":          SymbolKindFile,
	"node":             SymbolKindVariable,
	"object":           SymbolKindClass,
	"oneof":            SymbolKindEnum,
	"operator":         SymbolKindOperator,
	"output":           SymbolKindVariable,
	"package":          SymbolKindModule,
	"param":            SymbolKindVariable,
	"parameter":        SymbolKindVariable,
	"paramEntity":      SymbolKindVariable,
	"part":             SymbolKindVariable,
	"placeholder":      SymbolKindVariable,
	"port":             SymbolKindVariable,
	"process":          SymbolKindFunction,
	"property":         SymbolKindProperty,
	"prototype":        SymbolKindVariable,
	"protocol":         SymbolKindClass,
	"provider":         SymbolKindClass,
	"publication":      SymbolKindVariable,
	"qkey":             SymbolKindVariable,
	"receiver":         SymbolKindVariable,
	"record":           SymbolKindStruct,
	"region":           SymbolKindVariable,
	"register":         SymbolKindVariable,
	"repoid":           SymbolKindVariable,
	"report":           SymbolKindVariable,
	"repositoryId":     SymbolKindVariable,
	"repr":             SymbolKindVariable,
	"resource":         SymbolKindVariable,
	"response":         SymbolKindFunction,
	"role":             SymbolKindClass,
	"rpc":              SymbolKindVariable,
	"schema":           SymbolKindVariable,
	"script":           SymbolKindFile,
	"sequence":         SymbolKindVariable,
	"server":           SymbolKindClass,
	"service":          SymbolKindClass,
	"setter":           SymbolKindMethod,
	"signal":           SymbolKindFunction,
	"singletonMethod":  SymbolKindMethod,
	"slot":             SymbolKindVariable,
	"software":         SymbolKindClass,
	"sourcefile":       SymbolKindFile,
	"standard":         SymbolKindVariable,
	"string":           SymbolKindString,
	"structure":        SymbolKindStruct,
	"stylesheet":       SymbolKindVariable,
	"submethod":        SymbolKindMethod,
	"submodule":        SymbolKindModule,
	"subprogram":       SymbolKindFunction,
	"subprogspec":      SymbolKindVariable,
	"subroutine":       SymbolKindFunction,
	"subsection":       SymbolKindVariable,
	"subst":            SymbolKindVariable,
	"substdef":         SymbolKindVariable,
	"tag":              SymbolKindVariable,
	"template":         SymbolKindVariable,
	"test":             SymbolKindVariable,
	"theme":            SymbolKindVariable,
	"theorem":          SymbolKindVariable,
	"thriftFile":       SymbolKindFile,
	"throwsparam":      SymbolKindVariable,
	"title":            SymbolKindVariable,
	"token":            SymbolKindVariable,
	"toplevelVariable": SymbolKindVariable,
	"trait":            SymbolKindVariable,
	"type":             SymbolKindStruct,
	"typealias":        SymbolKindVariable,
	"typedef":          SymbolKindTypeParameter,
	"typespec":         SymbolKindTypeParameter,
	"union":            SymbolKindStruct,
	"username":         SymbolKindVariable,
	"val":              SymbolKindVariable,
	"value":            SymbolKindVariable,
	"var":              SymbolKindVariable,
	"variable":         SymbolKindVariable,
	"vector":           SymbolKindVariable,
	"version":          SymbolKindVariable,
	"video":            SymbolKindFile,
	"view":             SymbolKindVariable,
	"wrapper":          SymbolKindVariable,
	"xdata":            SymbolKindVariable,
	"xinput":           SymbolKindVariable,
	"xtask":            SymbolKindVariable,
}
