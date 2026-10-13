package handlers

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rhl/businessos-backend/internal/founderos/pages/pgtest"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCreatedTaskSurvivesWorkspaceReload(t *testing.T) {
	pool := pgtest.ThrowawayDB(t)
	ctx := context.Background()
	user := uuid.NewString()
	ws := uuid.New()
	other := uuid.New()
	project := uuid.New()
	_, err := pool.Exec(ctx, `INSERT INTO "user" (id,name,email) VALUES ($1,'Test owner',$2)`, user, user+"@example.test")
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO workspaces (id,name,slug,owner_id) VALUES ($1,'Company',$2,$3),($4,'Other',$5,$3)`, ws, ws.String(), user, other, other.String())
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO workspace_members (workspace_id,user_id,role,status) VALUES ($1,$2,'owner','active'),($3,$2,'owner','active')`, ws, user, other)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO projects (id,user_id,name,workspace_id) VALUES ($1,$2,'Workshop',$3)`, project, user, ws)
	require.NoError(t, err)
	h := &DashboardItemHandler{pool: pool}
	body, _ := json.Marshal(map[string]string{"title": "Choose date", "project_id": project.String(), "due_date": "2026-11-26"})
	w, c := setupGin("POST", "/dashboard/tasks", body)
	injectUser(c, user, "Test owner", user+"@example.test")
	c.Request.Header.Set("X-Workspace-ID", ws.String())
	h.CreateTask(c)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	for _, tc := range []struct {
		id    uuid.UUID
		count int
	}{{ws, 1}, {other, 0}} {
		w, c = setupGin("GET", "/dashboard/tasks", nil)
		injectUser(c, user, "Test owner", user+"@example.test")
		c.Request.Header.Set("X-Workspace-ID", tc.id.String())
		h.ListTasks(c)
		require.Equal(t, 200, w.Code, w.Body.String())
		var tasks []map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &tasks))
		require.Len(t, tasks, tc.count)
		if tc.count > 0 {
			require.Equal(t, "todo", tasks[0]["status"])
			require.Equal(t, project.String(), tasks[0]["project_id"])
			require.NotNil(t, tasks[0]["due_date"])
		}
	}
}
