package service

import (
	"testing"

	"github.com/expandorg/assignment/pkg/authorization"
	"github.com/expandorg/assignment/pkg/datastore"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	authorizer := authorization.NewAuthorizer()
	ds := &datastore.AssignmentStore{}
	type args struct {
		s *datastore.AssignmentStore
	}
	tests := []struct {
		name string
		args args
		want *service
	}{
		{
			"it creates a new service",
			args{s: ds},
			&service{ds, authorizer},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.s, authorizer)
			assert.Equal(t, got, tt.want, tt.name)
		})
	}
}

func TestHealthy(t *testing.T) {
	ds := &datastore.AssignmentStore{}
	type fields struct {
		store *datastore.AssignmentStore
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"it returns true if healthy",
			fields{store: ds},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				store: tt.fields.store,
			}
			got := s.Healthy()
			assert.Equal(t, got, tt.want, tt.name)
		})
	}
}
