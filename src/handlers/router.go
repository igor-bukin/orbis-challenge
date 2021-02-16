package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/orbis-challenge/src/config"
	"github.com/orbis-challenge/src/handlers/etf"
	"github.com/orbis-challenge/src/handlers/user"
	middleware "github.com/orbis-challenge/src/middlewares"
)

// NewRouter creates a router for URL-to-service mapping
func NewRouter() *mux.Router { //nolint:funlen
	var (
		router    = mux.NewRouter()
		apiRouter = router.PathPrefix(config.Config.URLPrefix).Subrouter()
		v1Router  = apiRouter.PathPrefix("/v1").Subrouter()

		publicChain = alice.New()
		authChain   = publicChain.Append(middleware.Auth)
	)

	// public routes
	v1Router.Handle("/user/login", publicChain.ThenFunc(user.Login)).Methods(http.MethodPost)

	// public route, sign customer
	v1Router.Handle("/user/signup", publicChain.ThenFunc(user.SignUp)).Methods(http.MethodPost)

	v1Router.Handle("/etf/load", publicChain.ThenFunc(etf.Load)).Methods(http.MethodPost)

	// holding
	v1Router.Handle("/etf", authChain.ThenFunc(etf.GetAll)).Methods(http.MethodGet)

	v1Router.Handle("/etf/ticker", authChain.ThenFunc(etf.GetByTicker)).Methods(http.MethodGet)

	return router
}
