package stores

import (
	"fmt"
	"testing"

	"github.com/tuxounet/k2/types"
)

func TestNewActionPlan(t *testing.T) {
	plan := NewActionPlan(nil)
	if plan == nil {
		t.Fatal("expected non-nil plan")
	}
	if len(plan.Tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(plan.Tasks))
	}
	if len(plan.Entities) != 0 {
		t.Fatalf("expected 0 entities, got %d", len(plan.Entities))
	}
	if len(plan.Templates) != 0 {
		t.Fatalf("expected 0 templates, got %d", len(plan.Templates))
	}
}

func TestActionPlan_AddEntity(t *testing.T) {
	plan := NewActionPlan(nil)
	meta := types.IK2Metadata{ID: "entity-1"}
	plan.AddEntity(meta, "some data")

	if len(plan.Entities) != 1 {
		t.Fatalf("expected 1 entity, got %d", len(plan.Entities))
	}
	if plan.Entities["entity-1"] != "some data" {
		t.Fatalf("unexpected entity value: %v", plan.Entities["entity-1"])
	}
}

func TestActionPlan_AddTask(t *testing.T) {
	plan := NewActionPlan(nil)
	task := ActionTask{
		Type:   ActionTaskTypeApply,
		Params: map[string]interface{}{"id": "test"},
	}
	plan.AddTask(task)

	if len(plan.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(plan.Tasks))
	}
	if plan.Tasks[0].Type != ActionTaskTypeApply {
		t.Fatalf("expected type 'apply', got '%s'", plan.Tasks[0].Type)
	}
}

func TestActionPlan_Dedup(t *testing.T) {
	plan := NewActionPlan(nil)
	task := ActionTask{
		Type:   ActionTaskTypeApply,
		Params: map[string]interface{}{"id": "test"},
	}
	plan.AddTask(task)
	plan.AddTask(task) // duplicate
	plan.AddTask(task) // duplicate

	if len(plan.Tasks) != 3 {
		t.Fatalf("expected 3 tasks before dedup, got %d", len(plan.Tasks))
	}

	plan.Dedup()

	if len(plan.Tasks) != 1 {
		t.Fatalf("expected 1 task after dedup, got %d", len(plan.Tasks))
	}
}

func TestActionPlan_Dedup_DifferentTasks(t *testing.T) {
	plan := NewActionPlan(nil)
	plan.AddTask(ActionTask{
		Type:   ActionTaskTypeApply,
		Params: map[string]interface{}{"id": "a"},
	})
	plan.AddTask(ActionTask{
		Type:   ActionTaskTypeApply,
		Params: map[string]interface{}{"id": "b"},
	})
	plan.AddTask(ActionTask{
		Type:   ActionTaskTypeLocalResolve,
		Params: map[string]interface{}{"hash": "xxx"},
	})

	plan.Dedup()

	if len(plan.Tasks) != 3 {
		t.Fatalf("expected 3 tasks (all different), got %d", len(plan.Tasks))
	}
}

func TestActionPlan_Dedup_MixedDuplicates(t *testing.T) {
	plan := NewActionPlan(nil)
	taskA := ActionTask{Type: ActionTaskTypeApply, Params: map[string]interface{}{"id": "a"}}
	taskB := ActionTask{Type: ActionTaskTypeApply, Params: map[string]interface{}{"id": "b"}}
	plan.AddTask(taskA)
	plan.AddTask(taskB)
	plan.AddTask(taskA)
	plan.AddTask(taskB)

	plan.Dedup()

	if len(plan.Tasks) != 2 {
		t.Fatalf("expected 2 unique tasks, got %d", len(plan.Tasks))
	}
}

func TestActionPlan_GetEntityAsTemplate(t *testing.T) {
	plan := NewActionPlan(nil)
	tpl := &types.IK2Template{
		K2: types.IK2TemplateRoot{
			Metadata: types.IK2Metadata{ID: "tpl-1"},
			Body: types.IK2TemplateBody{
				Name: "kind1",
			},
		},
	}
	plan.AddEntity(types.IK2Metadata{ID: "tpl-1"}, tpl)

	result, err := plan.GetEntityAsTemplate("tpl-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.K2.Body.Name != "kind1" {
		t.Fatalf("expected 'kind1', got '%s'", result.K2.Body.Name)
	}
}

func TestActionPlan_GetEntityAsTemplate_NotFound(t *testing.T) {
	plan := NewActionPlan(nil)
	_, err := plan.GetEntityAsTemplate("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing entity")
	}
}

func TestActionPlan_GetEntityAsTemplateApply(t *testing.T) {
	plan := NewActionPlan(nil)
	apply := &types.IK2TemplateApply{
		K2: types.IK2TemplateApplyRoot{
			Metadata: types.IK2Metadata{ID: "apply-1"},
			Body: types.IK2TemplateApplyBody{
				Vars: map[string]any{"name": "comp"},
			},
		},
	}
	plan.AddEntity(types.IK2Metadata{ID: "apply-1"}, apply)

	result, err := plan.GetEntityAsTemplateApply("apply-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.K2.Body.Vars["name"] != "comp" {
		t.Fatalf("expected name='comp', got '%v'", result.K2.Body.Vars["name"])
	}
}

func TestActionPlan_GetEntityAsTemplateApply_NotFound(t *testing.T) {
	plan := NewActionPlan(nil)
	_, err := plan.GetEntityAsTemplateApply("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing entity")
	}
}

func TestActionTaskType_Constants(t *testing.T) {
	if ActionTaskTypeLocalResolve != "local-resolve" {
		t.Fatalf("expected 'local-resolve', got '%s'", ActionTaskTypeLocalResolve)
	}
	if ActionTaskTypeGitResolve != "git-resolve" {
		t.Fatalf("expected 'git-resolve', got '%s'", ActionTaskTypeGitResolve)
	}
	if ActionTaskTypeApply != "apply" {
		t.Fatalf("expected 'apply', got '%s'", ActionTaskTypeApply)
	}
}

func TestActionTask_ParamsAccess(t *testing.T) {
	task := ActionTask{
		Type: ActionTaskTypeApply,
		Params: map[string]interface{}{
			"id":  "my-apply",
			"ref": "abc123",
		},
	}
	if task.Params["id"] != "my-apply" {
		t.Fatalf("expected 'my-apply', got '%v'", task.Params["id"])
	}
	if task.Params["ref"] != "abc123" {
		t.Fatalf("expected 'abc123', got '%v'", task.Params["ref"])
	}
}

func TestActionPlan_Dedup_EmptyPlan(t *testing.T) {
	plan := NewActionPlan(nil)
	plan.Dedup()
	if len(plan.Tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(plan.Tasks))
	}
}

func TestActionPlan_MultipleEntities(t *testing.T) {
	plan := NewActionPlan(nil)
	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("entity-%d", i)
		plan.AddEntity(types.IK2Metadata{ID: id}, i)
	}
	if len(plan.Entities) != 10 {
		t.Fatalf("expected 10 entities, got %d", len(plan.Entities))
	}
}
