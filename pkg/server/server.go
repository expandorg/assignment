package server

import (
	"net/http"

	"github.com/gemsorg/assignment/pkg/api/assignmentcreator"
	"github.com/gemsorg/assignment/pkg/api/assignmentdestroyer"
	"github.com/gemsorg/assignment/pkg/api/assignmentfetcher"
	"github.com/gemsorg/assignment/pkg/api/assignmentupdater"

	"github.com/gemsorg/assignment/pkg/authentication"

	"github.com/jmoiron/sqlx"

	"github.com/gemsorg/assignment/pkg/api/healthchecker"
	"github.com/gemsorg/assignment/pkg/service"
	"github.com/gorilla/mux"
)

func New(
	db *sqlx.DB,
	s service.AssignmentService,
) http.Handler {
	r := mux.NewRouter()

	r.Handle("/_health", healthchecker.MakeHandler(s)).Methods("GET")
	r.Handle("/assignments", assignmentfetcher.MakeAssignmentsFetcherHandler(s)).Methods("GET")
	r.Handle("/assignments/{assignment_id}", assignmentfetcher.MakeAssignmentFetcherHandler(s)).Methods("GET")
	r.Handle("/assignments", assignmentcreator.MakeHandler(s)).Methods("POST")
	r.Handle("/assignments", assignmentdestroyer.MakeHandler(s)).Methods("DELETE")
	r.Handle("/assignments", assignmentupdater.MakeHandler(s)).Methods("PATCH")
	r.Use(authentication.AuthMiddleware)
	return withHandlers(r)
}
