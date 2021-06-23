package admin

import (
	"github.com/huminghe/infini-framework/core/api"
	"github.com/huminghe/infini-framework/core/api/router"
	"github.com/huminghe/infini-framework/modules/ui/admin/console"
	"github.com/huminghe/infini-framework/modules/ui/admin/dashboard"
	"github.com/huminghe/infini-framework/modules/ui/common"
	"net/http"
)

type AdminUI struct {
	api.Handler
}

func (h AdminUI) DashboardAction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	dashboard.Index(w, r)
}

func (h AdminUI) ConsolePageAction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	console.Index(w, r)
}

func (h AdminUI) ExplorePageAction(w http.ResponseWriter, r *http.Request) {
	common.Message(w, r, "hello", "world")
	//explore.Index(w, r)
}
