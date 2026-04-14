package agent_memory

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"gorm.io/gorm"

	"github.com/Notailab/Notailab/internal/model"
	agent_core "github.com/Notailab/go-agent/agent/core"
)

const projectLongMemoryFileName = ".MEMORY"

type longMemoryFileStore interface {
	GetFileByNameAndProjectID(projectID uint, name string) (*model.File, error)
	CreateFile(file *model.File) error
	UpdateFileContent(fileID uint, content string) error
}

type ProjectLongStore struct {
	fileStore longMemoryFileStore
	projectID uint
	fileID    uint
	mu        sync.RWMutex
	memories  []string
}

func NewProjectLongStore(fileStore longMemoryFileStore, projectID uint) (*ProjectLongStore, error) {
	if fileStore == nil {
		return nil, fmt.Errorf("file store is nil")
	}

	store := &ProjectLongStore{
		fileStore: fileStore,
		projectID: projectID,
		memories:  []string{},
	}
	if err := store.loadOrCreate(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *ProjectLongStore) loadOrCreate() error {
	file, err := s.fileStore.GetFileByNameAndProjectID(s.projectID, projectLongMemoryFileName)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := s.fileStore.CreateFile(&model.File{
			ProjectID: s.projectID,
			Name:      projectLongMemoryFileName,
			Content:   "",
			IsHidden:  true,
		}); err != nil {
			file, err = s.fileStore.GetFileByNameAndProjectID(s.projectID, projectLongMemoryFileName)
			if err != nil {
				return err
			}
		} else {
			file, err = s.fileStore.GetFileByNameAndProjectID(s.projectID, projectLongMemoryFileName)
			if err != nil {
				return err
			}
		}
	}

	memories, err := decodeMemories(file.Content)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.fileID = file.FileID
	s.memories = memories
	s.mu.Unlock()
	return nil
}

func decodeMemories(content string) ([]string, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return []string{}, nil
	}

	var memories []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		var entry string
		if err := json.Unmarshal(line, &entry); err != nil {
			memories = append(memories, string(line))
			continue
		}
		memories = append(memories, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return memories, nil
}

func encodeMemories(memories []string) (string, error) {
	if len(memories) == 0 {
		return "", nil
	}

	var builder strings.Builder
	for i, memory := range memories {
		if i > 0 {
			builder.WriteByte('\n')
		}
		line, err := json.Marshal(memory)
		if err != nil {
			return "", err
		}
		builder.Write(line)
	}
	return builder.String(), nil
}

func (s *ProjectLongStore) persistLocked(memories []string) error {
	content, err := encodeMemories(memories)
	if err != nil {
		return err
	}
	if err := s.fileStore.UpdateFileContent(s.fileID, content); err != nil {
		return err
	}
	s.memories = append([]string(nil), memories...)
	return nil
}

func (s *ProjectLongStore) Get(index int) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if index < 0 || index >= len(s.memories) {
		return "", errors.New("index out of bounds")
	}
	return s.memories[index], nil
}

func (s *ProjectLongStore) Append(memory string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	updated := append(append([]string(nil), s.memories...), memory)
	return s.persistLocked(updated)
}

func (s *ProjectLongStore) Update(index int, memory string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.memories) {
		return errors.New("index out of bounds")
	}
	updated := append([]string(nil), s.memories...)
	updated[index] = memory
	return s.persistLocked(updated)
}

func (s *ProjectLongStore) Replace(start, end int, memories []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if start < 0 || start > len(s.memories) {
		return errors.New("start index out of bounds")
	}
	if end < 0 || end > len(s.memories) {
		return errors.New("end index out of bounds")
	}
	if start > end {
		return errors.New("start index cannot be greater than end index")
	}
	updated := append([]string(nil), s.memories[:start]...)
	updated = append(updated, memories...)
	updated = append(updated, s.memories[end:]...)
	return s.persistLocked(updated)
}

func (s *ProjectLongStore) Delete(index int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.memories) {
		return errors.New("index out of bounds")
	}
	updated := append([]string(nil), s.memories[:index]...)
	updated = append(updated, s.memories[index+1:]...)
	return s.persistLocked(updated)
}

func (s *ProjectLongStore) List() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]string(nil), s.memories...), nil
}

func (s *ProjectLongStore) Count() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.memories), nil
}

func (s *ProjectLongStore) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.persistLocked([]string{})
}

func (s *ProjectLongStore) Clone() agent_core.LongMemoryStore {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &ProjectLongStore{
		fileStore: s.fileStore,
		projectID: s.projectID,
		fileID:    s.fileID,
		memories:  append([]string(nil), s.memories...),
	}
}
