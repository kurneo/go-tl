//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/kurneo/go-template/internal/auth"
	"github.com/kurneo/go-template/internal/category"
	"github.com/kurneo/go-template/pkg"
)

// InitializeApp resolve and create global application that will start to serve request
func InitializeApp() App {
	wire.Build(
		pkg.WireSet,
		auth.WireSet,
		category.WireSet,
		NewApplication,
	)
	return &application{}
}
