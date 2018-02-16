package main

import "go.uber.org/zap"

var logger, _ = zap.NewProduction()
var sLogger = logger.Sugar()
