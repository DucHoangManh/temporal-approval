package workflow

type ApprovalStatus string

const (
	ApprovalSignal string = "APPROVAL_SIGNAL_CHANNEL"
	TaskQueueName  string = "approval_workflow"

	ApprovalStatusPending  ApprovalStatus = "pending"
	ApprovalStatusApproved ApprovalStatus = "approved"
	ApprovalStatusRejected ApprovalStatus = "rejected"
)

var DefaultApprovalDefinition = ApprovalDefinition{
	Status: ApprovalStatusPending,
	Approves: []ApproveGroup{
		{
			ApprovalUsers: []ApprovalUser{
				{
					Email: "duc@encapital.io",
				},
			},
		},
		{
			ApprovalUsers: []ApprovalUser{
				{
					Email: "nmh@encapital.io",
				},
			},
		},
	},
}

type ApprovalDefinition struct {
	Status   ApprovalStatus `json:"status"`
	Approves []ApproveGroup `json:"approves"`
}

func (a *ApprovalDefinition) ContainsEmail(email string) bool {
	for _, approveGroup := range a.Approves {
		for _, approvalUser := range approveGroup.ApprovalUsers {
			if approvalUser.Email == email {
				return true
			}
		}
	}
	return false
}

func (a *ApprovalDefinition) HandleApprove(email string) {
	for groupIndex, approveGroup := range a.Approves {
		for approvalIndex, approvalUser := range approveGroup.ApprovalUsers {
			if approvalUser.Email == email {
				approveGroup.ApprovalUsers[approvalIndex].Approved = true
			}
			if approveGroup.IsApproved() {
				a.Approves[groupIndex].Approved = true
				if groupIndex == len(a.Approves)-1 {
					a.Status = ApprovalStatusApproved
				}
			}
		}
	}
}

func (a *ApprovalDefinition) HandleReject() {
	a.Status = ApprovalStatusRejected
}

type ApproveGroup struct {
	Approved      bool           `json:"approved"`
	ApprovalUsers []ApprovalUser `json:"approvalUsers"`
}

func (a ApproveGroup) IsApproved() bool {
	for _, user := range a.ApprovalUsers {
		if !user.Approved {
			return false
		}
	}
	return true
}

type ApprovalUser struct {
	Email    string `json:"email"`
	Approved bool   `json:"approved"`
}
