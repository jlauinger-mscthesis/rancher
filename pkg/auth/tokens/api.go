package tokens

import (
	"context"
	"net/http"

	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	rancherapi "github.com/rancher/rancher/pkg/api"
	managementSchema "github.com/rancher/types/apis/management.cattle.io/v3/schema"
	client "github.com/rancher/types/client/management/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
)

const (
	CookieName      = "R_SESS"
	AuthHeaderName  = "Authorization"
	AuthValuePrefix = "Bearer"
	BasicAuthPrefix = "Basic"
	CSRFCookie      = "CSRF"
)

var crdVersions = []*types.APIVersion{
	&managementSchema.Version,
}

func NewAPIHandler(ctx context.Context, apiContext *config.ScaledContext) (http.Handler, error) {
	api := &tokenAPI{
		mgr: NewManager(ctx, apiContext),
	}

	schemas := types.NewSchemas().AddSchemas(managementSchema.TokenSchemas)
	schema := schemas.Schema(&managementSchema.Version, client.TokenType)
	schema.CollectionActions = map[string]types.Action{
		"logout": {},
	}

	schema.ActionHandler = api.tokenActionHandler
	schema.ListHandler = api.tokenListHandler
	schema.CreateHandler = api.tokenCreateHandler
	schema.DeleteHandler = api.tokenDeleteHandler

	return rancherapi.NewServer(schemas)
}

type tokenAPI struct {
	mgr *Manager
}

func (t *tokenAPI) tokenActionHandler(actionName string, action *types.Action, request *types.APIContext) error {
	logrus.Debugf("TokenActionHandler called for action %v", actionName)
	if actionName == "logout" {
		return t.mgr.logout(actionName, action, request)
	}
	return httperror.NewAPIError(httperror.ActionNotAvailable, "")
}

func (t *tokenAPI) tokenCreateHandler(request *types.APIContext, _ types.RequestHandler) error {
	logrus.Debugf("TokenCreateHandler called")
	return t.mgr.deriveToken(request)
}

func (t *tokenAPI) tokenListHandler(request *types.APIContext, _ types.RequestHandler) error {
	logrus.Debugf("TokenListHandler called")
	if request.ID != "" {
		return t.mgr.getTokenFromRequest(request)
	}
	return t.mgr.listTokens(request)
}

func (t *tokenAPI) tokenDeleteHandler(request *types.APIContext, _ types.RequestHandler) error {
	logrus.Debugf("TokenDeleteHandler called")
	return t.mgr.removeToken(request)
}
