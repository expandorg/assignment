package server

import (
	"net/http"

	"github.com/expandorg/assignment/pkg/api/assignmentcreator"
	"github.com/expandorg/assignment/pkg/api/assignmentdestroyer"
	"github.com/expandorg/assignment/pkg/api/assignmentfetcher"
	"github.com/expandorg/assignment/pkg/api/assignmentupdater"
	"github.com/expandorg/assignment/pkg/api/assignmentvalidator"
	"github.com/expandorg/assignment/pkg/api/settingcreator"
	"github.com/expandorg/assignment/pkg/api/settingfetcher"

	"github.com/expandorg/assignment/pkg/authentication"

	"github.com/jmoiron/sqlx"

	"github.com/expandorg/assignment/pkg/api/healthchecker"
	"github.com/expandorg/assignment/pkg/service"
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
	r.Handle("/assignments/{assignment_id}", assignmentdestroyer.MakeHandler(s)).Methods("DELETE")
	r.Handle("/assignments", assignmentupdater.MakeHandler(s)).Methods("PATCH")
	r.Handle("/settings/{job_id}", settingfetcher.MakeHandler(s)).Methods("GET")
	r.Handle("/settings", settingcreator.MakeHandler(s)).Methods("PUT")
	r.Handle("/validate", assignmentvalidator.MakeHandler(s)).Methods("GET")
	r.Use(authentication.AuthMiddleware)
	return withHandlers(r)
}
