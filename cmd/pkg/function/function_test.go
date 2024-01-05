package function

import (
	"fmt"
	"testing"

	"approval-demo/cmd/internal/activity"
)

func Test_GetFunctionName(t *testing.T) {
	name := GetFunctionName(activity.PostApproveActivity)
	fmt.Println(name)
}
