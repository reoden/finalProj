package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"restaurant/common"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"
)

const (
	verificationLimitPerDay = 10
)

// LoginService ...
type LoginService struct {
	db               *domain.LoginRepo
	account          *AccountService
	verificationCode VerificationCodeService
	cacheService     *CacheService
}

// NewLoginService
func NewLoginService() *LoginService {
	return &LoginService{
		db:               domain.NewLoginRepo(),
		account:          NewAccountService(),
		verificationCode: *NewVerificationCodeService(),
		cacheService:     NewCacheService(),
	}
}

// Login 登陆
func (l *LoginService) UserLogin(ctx context.Context, param domain.Login) (interface{}, error) {
	loginInfo, err := l.db.UserLoginInfo(ctx, param.Phone)
	if err != nil {
		log.Printf("get loginInfo error:%v", err)
		return nil, pkgs.LoginInfoGetError
	}

	if loginInfo == nil {
		log.Println("get loginInfo nil")
		return nil, pkgs.LoginInfoGetError
	}

	code, err := l.verificationCode.GetCode(ctx, loginInfo.Phone)
	if err != nil {
		return nil, err
	}
	if code.Code != param.Code {
		return nil, pkgs.VarifyCodeError
	}

	l.verificationCode.UpdateStatus(ctx, loginInfo.Phone, code.Code)

	_, err = l.account.GetAccount(ctx, loginInfo.ID)
	if err != nil {
		log.Printf("get loginInfo GetAccount error:%v", err)
		return nil, pkgs.ServerError
	}

	token := pkgs.CreateJWTToken(loginInfo.ID)
	return map[string]interface{}{
		"access_token": token.AccessToken,
		"expired_at":   token.ExpiredAt,
	}, nil
}

// Login 登陆
func (l *LoginService) AdminLogin(ctx context.Context, param domain.Login) (interface{}, error) {
	loginInfo, err := l.db.AdminLoginInfo(ctx, param.Phone)
	if err != nil {
		log.Printf("get loginInfo error:%v", err)
		return nil, pkgs.LoginInfoGetError
	}

	if loginInfo == nil {
		log.Println("get loginInfo nil")
		return nil, pkgs.LoginInfoGetError
	}

	code, err := l.verificationCode.GetCode(ctx, loginInfo.Phone)
	if err != nil {
		return nil, err
	}
	if code.Code != param.Code {
		return nil, pkgs.VarifyCodeError
	}

	l.verificationCode.UpdateStatus(ctx, loginInfo.Phone, code.Code)

	_, err = l.account.GetAccount(ctx, loginInfo.ID)
	if err != nil {
		log.Printf("get loginInfo GetAccount error:%v", err)
		return nil, pkgs.ServerError
	}

	token := pkgs.CreateJWTToken(loginInfo.ID)
	return map[string]interface{}{
		"access_token": token.AccessToken,
		"expired_at":   token.ExpiredAt,
	}, nil
}

func (l *LoginService) GetFetchingSeconds(ctx context.Context, phone string) (int64, error) {
	leftSeconds, err := l.cacheService.TTLKey(ctx, fmt.Sprintf("fetch_seconds_%d", phone))
	if err != nil {
		return 0, err
	}

	return leftSeconds, err
}

func (l *LoginService) SendCode(ctx context.Context, param domain.Login) (interface{}, error) {
	phone := param.Phone

	secTime, err := l.GetFetchingSeconds(ctx, phone)
	if err != nil {
		log.Println("GetFetchingSeconds error:", err)
		return nil, err
	}

	if secTime > 0 {
		return nil, errors.New("请等待" + fmt.Sprintf("%d", secTime) + "秒后再试")
	}

	data, err := l.verificationCode.GetCode(ctx, phone)
	if err != nil {
		return nil, err
	}

	if data != nil {
		return nil, errors.New("上一个验证码尚未过期")
	}

	count, _ := l.verificationCode.CountTimes(ctx, phone)
	if count > verificationLimitPerDay {
		return nil, errors.New("今日该手机发送验证码次数已达上限")
	}
	code := common.RandDigit(6)

	l.cacheService.SetKey(ctx, fmt.Sprintf("fetch_seconds_%d", phone), "1", 60)
	log.Println("==============", code, "=======================")
	// resp, err := pkgs.SendMsg(phone, code)
	// if err != nil {
	// 	return nil, err
	// }
	// if resp.Message != "OK" {
	// 	return resp.Message, err
	// }

	_, err = l.verificationCode.Create(ctx, entity.VerificationCode{
		Code:  code,
		Phone: phone,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
