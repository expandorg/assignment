package assignment

import "testing"

func TestNewAssignment_IsAllowed(t *testing.T) {
	limit := 2
	permutations := map[string]*Settings{
		"OneAssignmentRepeat": &Settings{
			ID:     1,
			JobID:  1,
			Limit:  0,
			Repeat: true,
			Singly: true,
		},
		"OneNoAssignmentRepeat": &Settings{
			ID:     1,
			JobID:  1,
			Limit:  0,
			Repeat: false,
			Singly: true,
		},
		"AllAssignmentRepeat": &Settings{
			ID:     1,
			JobID:  1,
			Limit:  0,
			Repeat: true,
			Singly: false,
		},
		"AssignmentWithLimit": &Settings{
			ID:     1,
			JobID:  1,
			Limit:  limit,
			Repeat: true,
			Singly: false,
		},
	}

	type args struct {
		set *Settings
	}
	tests := []struct {
		name    string
		fields  NewAssignment
		args    args
		want    bool
		wantErr bool
	}{
		{
			"singly assignment: returns true if worker is not assigned",
			MakeNewAssignment(false, false, true, true, 0),
			args{
				permutations["OneAssignmentRepeat"],
			},
			true,
			false,
		},
		{
			"singly assignment: returns true if worker is already assigned",
			MakeNewAssignment(true, false, true, true, 0),
			args{
				permutations["OneAssignmentRepeat"],
			},
			false,
			true,
		},
		{
			"singly assignment with repeats: returns true if worker has already responded",
			MakeNewAssignment(false, true, true, true, 0),
			args{
				permutations["OneAssignmentRepeat"],
			},
			true,
			false,
		},
		{
			"singly assignment without repeats: returns false if worker has already responded",
			MakeNewAssignment(false, true, true, true, 0),
			args{
				permutations["OneNoAssignmentRepeat"],
			},
			false,
			true,
		},
		{
			"multi assignment: returns true if worker has not been assigned",
			MakeNewAssignment(false, false, true, true, 0),
			args{
				permutations["AllAssignmentRepeat"],
			},
			true,
			false,
		},
		{
			"multi assignment: returns true if worker has been assigned",
			MakeNewAssignment(true, false, true, true, 0),
			args{
				permutations["AllAssignmentRepeat"],
			},
			true,
			false,
		},
		{
			"multi assignment: returns true if worker has responded but not assigned",
			MakeNewAssignment(false, true, true, true, 0),
			args{
				permutations["AllAssignmentRepeat"],
			},
			true,
			false,
		},
		{
			"multi assignment: returns true if worker has responded and assigned",
			MakeNewAssignment(true, true, true, true, 0),
			args{
				permutations["AllAssignmentRepeat"],
			},
			true,
			false,
		},
		{
			"returns true if job assignment limit has not been reached",
			MakeNewAssignment(false, false, true, true, limit-1),
			args{
				permutations["AssignmentWithLimit"],
			},
			true,
			false,
		},
		{
			"returns false if job assignment limit has been reached",
			MakeNewAssignment(false, false, true, true, limit),
			args{
				permutations["AssignmentWithLimit"],
			},
			false,
			true,
		},
		{
			"returns false if onboarding status is false",
			MakeNewAssignment(false, false, false, true, 0),
			args{
				permutations["AllAssignmentRepeat"],
			},
			false,
			true,
		},
		{
			"returns false if worker doesn't have enough funds",
			MakeNewAssignment(false, false, true, false, 0),
			args{
				permutations["AllAssignmentRepeat"],
			},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAssignment{
				JobID:                  tt.fields.JobID,
				TaskID:                 tt.fields.TaskID,
				WorkerID:               tt.fields.WorkerID,
				JobAssignmentCount:     tt.fields.JobAssignmentCount,
				OnboardingStatus:       tt.fields.OnboardingStatus,
				WorkerAlreadyAssigned:  tt.fields.WorkerAlreadyAssigned,
				WorkerAlreadyResponded: tt.fields.WorkerAlreadyResponded,
				WorkerHasFunds:         tt.fields.WorkerHasFunds,
			}
			got, err := a.IsAllowed(tt.args.set)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAssignment.IsAllowed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewAssignment.IsAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func MakeNewAssignment(assigned, responded, onboarding, funds bool, aCount int) NewAssignment {
	return NewAssignment{
		JobID:                  1,
		TaskID:                 1,
		WorkerID:               1,
		JobAssignmentCount:     aCount,
		OnboardingStatus:       onboarding,
		WorkerAlreadyAssigned:  assigned,
		WorkerAlreadyResponded: responded,
		WorkerHasFunds:         funds,
	}
}
