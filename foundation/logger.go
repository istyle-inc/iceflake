package foundation

import "go.uber.org/zap"

var Logger, _ = zap.NewProduction()
var SLogger = Logger.Sugar()
