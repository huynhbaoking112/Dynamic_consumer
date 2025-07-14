package dto

type WorkspaceCreatedPayload struct {
	WorkspaceID   string `json:"workspaceId"`
	WorkspaceName string `json:"workspaceName"`
	CreatedByID   string `json:"createdById"`
}
