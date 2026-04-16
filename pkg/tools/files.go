package tools

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/Notailab/Notailab/internal/model"
	agent_core "github.com/Notailab/go-agent/agent/core"
)

type projectFileStore interface {
	GetFilesByProjectID(projectID uint) ([]model.File, error)
	GetFileByNameAndProjectID(projectID uint, name string) (*model.File, error)
	CreateFile(file *model.File) error
	UpdateFileContent(fileID uint, content string) error
}

func NewFileTools(fileStore projectFileStore, projectID uint) ([]agent_core.Tool, error) {
	if fileStore == nil {
		return nil, fmt.Errorf("file store is nil")
	}

	return []agent_core.Tool{
		NewListFilesTool(fileStore, projectID),
		NewReadFileTool(fileStore, projectID),
		NewWriteFileTool(fileStore, projectID),
		NewEditFileTool(fileStore, projectID),
	}, nil
}

type ListFilesTool struct {
	store     projectFileStore
	projectID uint
}

func NewListFilesTool(store projectFileStore, projectID uint) *ListFilesTool {
	return &ListFilesTool{store: store, projectID: projectID}
}

func (t *ListFilesTool) Name() string {
	return "List_files"
}

func (t *ListFilesTool) Description() string {
	return "List files stored in the current Notailab project database."
}

func (t *ListFilesTool) Parameters() agent_core.Parameters {
	return agent_core.Parameters{
		Type:       "object",
		Properties: map[string]agent_core.Param{},
	}
}

func (t *ListFilesTool) Execute(paramsJson string) (string, error) {
	files, err := t.store.GetFilesByProjectID(t.projectID)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "(empty)", nil
	}

	lines := make([]string, 0, len(files))
	for _, file := range files {
		lines = append(lines, fmt.Sprintf("name=%s size=%d", file.Name, len(file.Content)))
	}
	return strings.Join(lines, "\n"), nil
}

type ReadFileTool struct {
	store     projectFileStore
	projectID uint
}

func NewReadFileTool(store projectFileStore, projectID uint) *ReadFileTool {
	return &ReadFileTool{store: store, projectID: projectID}
}

func (t *ReadFileTool) Name() string {
	return "Read_file"
}

func (t *ReadFileTool) Description() string {
	return "Read a file record from the current Notailab project database."
}

func (t *ReadFileTool) Parameters() agent_core.Parameters {
	return agent_core.Parameters{
		Type: "object",
		Properties: map[string]agent_core.Param{
			"file_name": {
				Type:        "string",
				Description: "The file name inside the project.",
			},
		},
		Required: []string{"file_name"},
	}
}

func (t *ReadFileTool) Execute(paramsJson string) (string, error) {
	params, err := agent_core.ParseToolParams(paramsJson, t.Parameters())
	if err != nil {
		return "", err
	}
	
	fileName := params["file_name"].(string)

	file, err := t.store.GetFileByNameAndProjectID(t.projectID, fileName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("file_id=%d name=%s\ncontent:\n%s", file.FileID, file.Name, file.Content), nil
}

type WriteFileTool struct {
	store     projectFileStore
	projectID uint
}

func NewWriteFileTool(store projectFileStore, projectID uint) *WriteFileTool {
	return &WriteFileTool{store: store, projectID: projectID}
}

func (t *WriteFileTool) Name() string {
	return "Write_file"
}

func (t *WriteFileTool) Description() string {
	return "Create or overwrite a file record in the current Notailab project database."
}

func (t *WriteFileTool) Parameters() agent_core.Parameters {
	return agent_core.Parameters{
		Type: "object",
		Properties: map[string]agent_core.Param{
			"file_name": {
				Type:        "string",
				Description: "The file name inside the project.",
			},
			"content": {
				Type:        "string",
				Description: "The full file content.",
			},
		},
		Required: []string{"file_name", "content"},
	}
}

func (t *WriteFileTool) Execute(paramsJson string) (string, error) {
	params, err := agent_core.ParseToolParams(paramsJson, t.Parameters())
	if err != nil {
		return "", err
	}
	
	fileName := params["file_name"].(string)
	content := params["content"].(string)

	if existing, err := t.store.GetFileByNameAndProjectID(t.projectID, fileName); err == nil {
		if err := t.store.UpdateFileContent(existing.FileID, content); err != nil {
			return "", err
		}
		return fmt.Sprintf("updated file_id=%d name=%s", existing.FileID, existing.Name), nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	file := &model.File{
		ProjectID: t.projectID,
		Name:      fileName,
		Content:   content,
	}
	if err := t.store.CreateFile(file); err != nil {
		return "", err
	}

	return fmt.Sprintf("created file_id=%d name=%s", file.FileID, file.Name), nil
}

type EditFileTool struct {
	store     projectFileStore
	projectID uint
}

func NewEditFileTool(store projectFileStore, projectID uint) *EditFileTool {
	return &EditFileTool{store: store, projectID: projectID}
}

func (t *EditFileTool) Name() string {
	return "Edit_file"
}

func (t *EditFileTool) Description() string {
	return "Edit an existing file record in the current Notailab project database by replacing text."
}

func (t *EditFileTool) Parameters() agent_core.Parameters {
	return agent_core.Parameters{
		Type: "object",
		Properties: map[string]agent_core.Param{
			"file_name": {
				Type:        "string",
				Description: "The file name inside the project.",
			},
			"old_text": {
				Type:        "string",
				Description: "The text to replace.",
			},
			"new_text": {
				Type:        "string",
				Description: "The replacement text.",
			},
		},
		Required: []string{"file_name", "old_text", "new_text"},
	}
}

func (t *EditFileTool) Execute(paramsJson string) (string, error) {
	params, err := agent_core.ParseToolParams(paramsJson, t.Parameters())
	if err != nil {
		return "", err
	}
	
	fileName := params["file_name"].(string)
	oldText := params["old_text"].(string)
	newText := params["new_text"].(string)

	if strings.TrimSpace(oldText) == "" {
		return "", fmt.Errorf("old_text is required")
	}

	file, err := t.store.GetFileByNameAndProjectID(t.projectID, fileName)
	if err != nil {
		return "", err
	}

	if !strings.Contains(file.Content, oldText) {
		return "", fmt.Errorf("old_text not found in file")
	}

	updated := strings.ReplaceAll(file.Content, oldText, newText)
	if err := t.store.UpdateFileContent(file.FileID, updated); err != nil {
		return "", err
	}

	return fmt.Sprintf("edited file_id=%d name=%s", file.FileID, file.Name), nil
}
