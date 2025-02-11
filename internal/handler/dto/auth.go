package dto

import (
	"github.com/PosokhovVadim/stawberry/internal/domain/entity"
)

type RegisterReq struct {
	Username string `json:"username" form:"username" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func (r RegisterReq) ConvertToSvc() entity.Register {
	return entity.Register{
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}
}

type LoginReq struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func (r LoginReq) ConvertToSvc() entity.Credentials {
	return entity.Credentials{
		Email:    r.Email,
		Password: r.Password,
	}
}

type RefreshTokensReq struct {
	RefreshToken string `header:"Refresh-Token" json:"refresh_token" binding:"required"`
}

func (r *RefreshTokensReq) ConvertToEntity() entity.RefreshSession {
	return entity.RefreshSession{
		Token: r.RefreshToken,
	}
}
