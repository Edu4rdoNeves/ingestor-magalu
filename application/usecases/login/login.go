package login

import (
	"errors"

	"github.com/Edu4rdoNeves/ingestor-magalu/cmd/api/auth"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
)

type ILoginUseCase interface {
	Login(input *dto.LoginInput) (*dto.LoginAuth, error)
}

type LoginUseCase struct {
}

func NewLoginUseCase() ILoginUseCase {
	return &LoginUseCase{}
}

func (uc *LoginUseCase) Login(input *dto.LoginInput) (*dto.LoginAuth, error) {
	if input.Username != env.AppUser || input.Password != env.AppPassword {
		return nil, errors.New("invalid credentials")
	}

	jwtSvc := auth.NewJWTService()
	loginAuth, err := jwtSvc.GenerateToken(1)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	return loginAuth, nil
}
