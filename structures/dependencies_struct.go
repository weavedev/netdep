// Package structures defines structs used in the project
// Copyright © 2022 TW Group 13C, Weave BV, TU Delft/*
package structures

import (
	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/callanalyzer"
	"lab.weave.nl/internships/tud-2022/netDep/stages/discovery/natsanalyzer"
)

type Dependencies struct {
	// stores dependencies for generic method
	// in callanalyzer package
	Calls     []*callanalyzer.CallTarget
	Endpoints []*callanalyzer.CallTarget

	// stores dependencies for nats analyzer
	Consumers []*natsanalyzer.NatsCall
	Producers []*natsanalyzer.NatsCall
}
