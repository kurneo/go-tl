package auth

import (
	"github.com/google/wire"
	"github.com/kurneo/go-template/internal/auth/data/datasource"
	"github.com/kurneo/go-template/internal/auth/data/repository"
	"github.com/kurneo/go-template/internal/auth/domain"
	domainRepository "github.com/kurneo/go-template/internal/auth/domain/repository"
	"github.com/kurneo/go-template/internal/auth/domain/usecase"
	v1 "github.com/kurneo/go-template/internal/auth/transport/http/v1"
	"github.com/kurneo/go-template/pkg/database"
	"github.com/kurneo/go-template/pkg/jwt"
	"github.com/kurneo/go-template/pkg/log"
)

// WireSet set of DI export from auth
var WireSet = wire.NewSet(
	ResolveUserDatasource,
	ResolveUserRepo,
	ResolveTokenManager,
	ResolveUserUseCase,
	ResolvePasswordChecker,
	ResolveAuthHttpV1Controller,
)

func ResolveUserDatasource(db database.Contract) *datasource.UserDatasource {
	return datasource.NewUserDataSource(db)
}

func ResolveUserRepo(u *datasource.UserDatasource) domainRepository.UserRepositoryContact {
	return repository.NewUserRepo(u)
}

func ResolveTokenManager(tm *jwt.TokenManager[int64]) *domain.TokenManager {
	return domain.NewTokenManager(tm)
}

func ResolvePasswordChecker() *domain.PasswordChecker {
	return domain.NewPasswordChecker()
}

func ResolveUserUseCase(
	r domainRepository.UserRepositoryContact,
	t *domain.TokenManager,
	p *domain.PasswordChecker,
) usecase.UserUseCaseContract {
	return usecase.NewUserUseCase(r, t, p)
}

func ResolveAuthHttpV1Controller(
	u usecase.UserUseCaseContract,
	l log.Contract,
	db database.Contract,
) *v1.Controller {
	return v1.NewHtpV1Controller(u, l, db)
}
