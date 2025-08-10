package officialaccount

import (
	originalContext "context"
	"fmt"

	"github.com/silenceper/wechat/v2/credential"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	opContext "github.com/silenceper/wechat/v2/openplatform/context"
	"github.com/silenceper/wechat/v2/openplatform/officialaccount/js"
	"github.com/silenceper/wechat/v2/openplatform/officialaccount/oauth"
)

// OfficialAccount 代公众号实现业务
type OfficialAccount struct {
	// 授权的公众号的appID
	AppID       string
	openContext *opContext.Context
	*officialaccount.OfficialAccount
	authorizerRefreshToken string
}

// GetAccessToken 获取ak
func (officialAccount *OfficialAccount) GetAccessToken() (string, error) {
	ak, akErr := officialAccount.openContext.GetAuthrAccessToken(officialAccount.AppID)
	if akErr == nil {
		return ak, nil
	}
	if officialAccount.authorizerRefreshToken == "" {
		return "", fmt.Errorf("please set the authorizer_refresh_token first")
	}
	akRes, akResErr := officialAccount.openContext.RefreshAuthrToken(officialAccount.AppID, officialAccount.authorizerRefreshToken)
	if akResErr != nil {
		return "", akResErr
	}
	return akRes.AccessToken, nil
}

// GetAccessTokenContext 利用ctx获取ak
func (officialAccount *OfficialAccount) GetAccessTokenContext(ctx originalContext.Context) (string, error) {
	ak, akErr := officialAccount.openContext.GetAuthrAccessTokenContext(ctx, officialAccount.AppID)
	if akErr == nil {
		return ak, nil
	}
	if officialAccount.authorizerRefreshToken == "" {
		return "", fmt.Errorf("please set the authorizer_refresh_token first")
	}
	akRes, akResErr := officialAccount.openContext.RefreshAuthrTokenContext(ctx, officialAccount.AppID, officialAccount.authorizerRefreshToken)
	if akResErr != nil {
		return "", akResErr
	}
	return akRes.AccessToken, nil
}

// SetAuthorizerRefreshToken 设置代执操作业务授权账号authorizer_refresh_token
func (officialAccount *OfficialAccount) SetAuthorizerRefreshToken(authorizerRefreshToken string) *OfficialAccount {
	officialAccount.authorizerRefreshToken = authorizerRefreshToken
	return officialAccount
}

// NewOfficialAccount 实例化
// appID :为授权方公众号 APPID，非开放平台第三方平台 APPID
func NewOfficialAccount(opCtx *opContext.Context, appID string) *OfficialAccount {
	officialAccount := officialaccount.NewOfficialAccount(&offConfig.Config{
		AppID:          opCtx.AppID,
		EncodingAESKey: opCtx.EncodingAESKey,
		Token:          opCtx.Token,
		Cache:          opCtx.Cache,
	})
	// 设置获取access_token的函数
	ret := &OfficialAccount{AppID: appID, OfficialAccount: officialAccount, openContext: opCtx}
	officialAccount.SetAccessTokenHandle(ret)
	return ret
}

// PlatformOauth 平台代发起oauth2网页授权
func (officialAccount *OfficialAccount) PlatformOauth() *oauth.Oauth {
	return oauth.NewOauth(officialAccount.GetContext())
}

// PlatformJs 平台代获取js-sdk配置
func (officialAccount *OfficialAccount) PlatformJs() *js.Js {
	return js.NewJs(officialAccount.GetContext(), officialAccount.AppID)
}

// DefaultAuthrAccessToken 默认获取授权ak的方法
type DefaultAuthrAccessToken struct {
	opCtx *opContext.Context
	appID string
}

// NewDefaultAuthrAccessToken New
func NewDefaultAuthrAccessToken(opCtx *opContext.Context, appID string) credential.AccessTokenHandle {
	return &DefaultAuthrAccessToken{
		opCtx: opCtx,
		appID: appID,
	}
}

// GetAccessToken 获取ak
func (ak *DefaultAuthrAccessToken) GetAccessToken() (string, error) {
	return ak.opCtx.GetAuthrAccessToken(ak.appID)
}
