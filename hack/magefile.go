//go:build mage

package main

import (
	//mage:import
	_ "github.com/eunomie/dague/mage/golang"
	//mage:import
	_ "github.com/eunomie/dague/mage/lint"
	//mage:import
	_ "github.com/eunomie/dague/mage/gofumpt"
)
