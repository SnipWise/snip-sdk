package files

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ============================================================================
// Helper functions for tests
// ============================================================================

// createTestDirectory creates a temporary test directory with sample files
func createTestDirectory(t *testing.T) string {
	t.Helper()

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "files_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create test files
	testFiles := map[string]string{
		"file1.txt":         "Content of file1",
		"file2.md":          "# Markdown content",
		"file3.go":          "package main",
		"subdir/file4.txt":  "Content in subdirectory",
		"subdir/file5.html": "<html>HTML content</html>",
		"subdir/nested/file6.txt": "Nested content",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)

		// Create subdirectories if needed
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", fullPath, err)
		}
	}

	return tempDir
}

// ============================================================================
// Tests for FindFiles
// ============================================================================

func TestFindFiles(t *testing.T) {
	tempDir := createTestDirectory(t)
	defer os.RemoveAll(tempDir)

	t.Run("find txt files", func(t *testing.T) {
		files, err := FindFiles(tempDir, ".txt")
		if err != nil {
			t.Fatalf("FindFiles() error = %v", err)
		}

		// Should find 3 .txt files
		if len(files) != 3 {
			t.Errorf("FindFiles() found %d files, want 3", len(files))
		}

		// Check that all files end with .txt
		for _, file := range files {
			if !strings.HasSuffix(file, ".txt") {
				t.Errorf("File %s does not have .txt extension", file)
			}
		}
	})

	t.Run("find all files with .*", func(t *testing.T) {
		files, err := FindFiles(tempDir, ".*")
		if err != nil {
			t.Fatalf("FindFiles() error = %v", err)
		}

		// Should find all 6 files
		if len(files) != 6 {
			t.Errorf("FindFiles() found %d files, want 6", len(files))
		}
	})

	t.Run("find go files", func(t *testing.T) {
		files, err := FindFiles(tempDir, ".go")
		if err != nil {
			t.Fatalf("FindFiles() error = %v", err)
		}

		// Should find 1 .go file
		if len(files) != 1 {
			t.Errorf("FindFiles() found %d files, want 1", len(files))
		}
	})

	t.Run("find non-existent extension", func(t *testing.T) {
		files, err := FindFiles(tempDir, ".pdf")
		if err != nil {
			t.Fatalf("FindFiles() error = %v", err)
		}

		// Should find 0 files
		if len(files) != 0 {
			t.Errorf("FindFiles() found %d files, want 0", len(files))
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		_, err := FindFiles("/non/existent/path", ".txt")
		if err == nil {
			t.Error("FindFiles() expected error for non-existent directory, got nil")
		}
	})
}

// ============================================================================
// Tests for ForEachFile
// ============================================================================

func TestForEachFile(t *testing.T) {
	tempDir := createTestDirectory(t)
	defer os.RemoveAll(tempDir)

	t.Run("iterate over txt files", func(t *testing.T) {
		count := 0
		files, err := ForEachFile(tempDir, ".txt", func(path string) error {
			count++
			// Verify path exists
			if _, err := os.Stat(path); err != nil {
				t.Errorf("File %s does not exist", path)
			}
			return nil
		})

		if err != nil {
			t.Fatalf("ForEachFile() error = %v", err)
		}

		if len(files) != 3 {
			t.Errorf("ForEachFile() found %d files, want 3", len(files))
		}

		if count != 3 {
			t.Errorf("Callback called %d times, want 3", count)
		}
	})

	t.Run("callback error stops iteration", func(t *testing.T) {
		count := 0
		expectedErr := os.ErrPermission

		_, err := ForEachFile(tempDir, ".txt", func(path string) error {
			count++
			if count >= 2 {
				return expectedErr // Stop after second file
			}
			return nil
		})

		if err != expectedErr {
			t.Errorf("ForEachFile() error = %v, want %v", err, expectedErr)
		}

		if count < 2 {
			t.Errorf("Callback called %d times, want at least 2", count)
		}
	})

	t.Run("iterate over all files", func(t *testing.T) {
		var paths []string
		files, err := ForEachFile(tempDir, ".*", func(path string) error {
			paths = append(paths, path)
			return nil
		})

		if err != nil {
			t.Fatalf("ForEachFile() error = %v", err)
		}

		if len(files) != 6 {
			t.Errorf("ForEachFile() found %d files, want 6", len(files))
		}

		if len(paths) != 6 {
			t.Errorf("Callback called %d times, want 6", len(paths))
		}
	})
}

// ============================================================================
// Tests for GetContentFiles
// ============================================================================

func TestGetContentFiles(t *testing.T) {
	tempDir := createTestDirectory(t)
	defer os.RemoveAll(tempDir)

	t.Run("get content of txt files", func(t *testing.T) {
		contents, err := GetContentFiles(tempDir, ".txt")
		if err != nil {
			t.Fatalf("GetContentFiles() error = %v", err)
		}

		if len(contents) != 3 {
			t.Errorf("GetContentFiles() returned %d contents, want 3", len(contents))
		}

		// Verify that contents are not empty
		for i, content := range contents {
			if content == "" {
				t.Errorf("Content %d is empty", i)
			}
		}
	})

	t.Run("get content of md files", func(t *testing.T) {
		contents, err := GetContentFiles(tempDir, ".md")
		if err != nil {
			t.Fatalf("GetContentFiles() error = %v", err)
		}

		if len(contents) != 1 {
			t.Errorf("GetContentFiles() returned %d contents, want 1", len(contents))
		}

		if !strings.Contains(contents[0], "Markdown") {
			t.Errorf("Content does not contain expected text")
		}
	})

	t.Run("no files found", func(t *testing.T) {
		contents, err := GetContentFiles(tempDir, ".pdf")
		if err != nil {
			t.Fatalf("GetContentFiles() error = %v", err)
		}

		if len(contents) != 0 {
			t.Errorf("GetContentFiles() returned %d contents, want 0", len(contents))
		}
	})
}

// ============================================================================
// Tests for ReadTextFile
// ============================================================================

func TestReadTextFile(t *testing.T) {
	tempDir := createTestDirectory(t)
	defer os.RemoveAll(tempDir)

	t.Run("read existing file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "file1.txt")
		content, err := ReadTextFile(filePath)
		if err != nil {
			t.Fatalf("ReadTextFile() error = %v", err)
		}

		expectedContent := "Content of file1"
		if content != expectedContent {
			t.Errorf("ReadTextFile() = %q, want %q", content, expectedContent)
		}
	})

	t.Run("read non-existent file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "non_existent.txt")
		_, err := ReadTextFile(filePath)
		if err == nil {
			t.Error("ReadTextFile() expected error for non-existent file, got nil")
		}
	})

	t.Run("read empty file", func(t *testing.T) {
		// Create an empty file
		emptyFile := filepath.Join(tempDir, "empty.txt")
		if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
			t.Fatalf("Failed to create empty file: %v", err)
		}

		content, err := ReadTextFile(emptyFile)
		if err != nil {
			t.Fatalf("ReadTextFile() error = %v", err)
		}

		if content != "" {
			t.Errorf("ReadTextFile() = %q, want empty string", content)
		}
	})

	t.Run("read file with unicode", func(t *testing.T) {
		unicodeFile := filepath.Join(tempDir, "unicode.txt")
		expectedContent := "Hello 世界 🌍"
		if err := os.WriteFile(unicodeFile, []byte(expectedContent), 0644); err != nil {
			t.Fatalf("Failed to create unicode file: %v", err)
		}

		content, err := ReadTextFile(unicodeFile)
		if err != nil {
			t.Fatalf("ReadTextFile() error = %v", err)
		}

		if content != expectedContent {
			t.Errorf("ReadTextFile() = %q, want %q", content, expectedContent)
		}
	})
}

// ============================================================================
// Tests for WriteTextFile
// ============================================================================

func TestWriteTextFile(t *testing.T) {
	tempDir := createTestDirectory(t)
	defer os.RemoveAll(tempDir)

	t.Run("write new file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "new_file.txt")
		content := "This is new content"

		err := WriteTextFile(filePath, content)
		if err != nil {
			t.Fatalf("WriteTextFile() error = %v", err)
		}

		// Verify file was created
		readContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read written file: %v", err)
		}

		if string(readContent) != content {
			t.Errorf("File content = %q, want %q", string(readContent), content)
		}
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "file1.txt")
		newContent := "Overwritten content"

		err := WriteTextFile(filePath, newContent)
		if err != nil {
			t.Fatalf("WriteTextFile() error = %v", err)
		}

		// Verify content was overwritten
		readContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if string(readContent) != newContent {
			t.Errorf("File content = %q, want %q", string(readContent), newContent)
		}
	})

	t.Run("write empty content", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "empty_write.txt")

		err := WriteTextFile(filePath, "")
		if err != nil {
			t.Fatalf("WriteTextFile() error = %v", err)
		}

		// Verify empty file was created
		readContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if len(readContent) != 0 {
			t.Errorf("File should be empty, got %d bytes", len(readContent))
		}
	})

	t.Run("write unicode content", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "unicode_write.txt")
		content := "Hello 世界 🌍"

		err := WriteTextFile(filePath, content)
		if err != nil {
			t.Fatalf("WriteTextFile() error = %v", err)
		}

		readContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if string(readContent) != content {
			t.Errorf("File content = %q, want %q", string(readContent), content)
		}
	})

	t.Run("write to invalid path", func(t *testing.T) {
		filePath := "/invalid/path/that/does/not/exist/file.txt"

		err := WriteTextFile(filePath, "content")
		if err == nil {
			t.Error("WriteTextFile() expected error for invalid path, got nil")
		}
	})
}

// ============================================================================
// Tests for GetAllFilesInDirectory
// ============================================================================

func TestGetAllFilesInDirectory(t *testing.T) {
	tempDir := createTestDirectory(t)
	defer os.RemoveAll(tempDir)

	t.Run("get files in root directory only", func(t *testing.T) {
		files, err := GetAllFilesInDirectory(tempDir)
		if err != nil {
			t.Fatalf("GetAllFilesInDirectory() error = %v", err)
		}

		// Should find only 3 files in root (not in subdirectories)
		if len(files) != 3 {
			t.Errorf("GetAllFilesInDirectory() found %d files, want 3", len(files))
		}

		// Verify all returned paths are files
		for _, file := range files {
			info, err := os.Stat(file)
			if err != nil {
				t.Errorf("File %s does not exist: %v", file, err)
			}
			if info.IsDir() {
				t.Errorf("Path %s is a directory, expected file", file)
			}
		}
	})

	t.Run("get files in subdirectory", func(t *testing.T) {
		subdirPath := filepath.Join(tempDir, "subdir")
		files, err := GetAllFilesInDirectory(subdirPath)
		if err != nil {
			t.Fatalf("GetAllFilesInDirectory() error = %v", err)
		}

		// Should find 2 files in subdir (not nested)
		if len(files) != 2 {
			t.Errorf("GetAllFilesInDirectory() found %d files, want 2", len(files))
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		emptyDir := filepath.Join(tempDir, "empty_dir")
		if err := os.Mkdir(emptyDir, 0755); err != nil {
			t.Fatalf("Failed to create empty directory: %v", err)
		}

		files, err := GetAllFilesInDirectory(emptyDir)
		if err != nil {
			t.Fatalf("GetAllFilesInDirectory() error = %v", err)
		}

		if len(files) != 0 {
			t.Errorf("GetAllFilesInDirectory() found %d files in empty directory, want 0", len(files))
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		_, err := GetAllFilesInDirectory("/non/existent/directory")
		if err == nil {
			t.Error("GetAllFilesInDirectory() expected error for non-existent directory, got nil")
		}
	})
}

// ============================================================================
// Integration test
// ============================================================================

func TestIntegration_WriteReadFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "integration_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "integration.txt")
	originalContent := "Integration test content"

	// Write file
	if err := WriteTextFile(filePath, originalContent); err != nil {
		t.Fatalf("WriteTextFile() error = %v", err)
	}

	// Read file
	readContent, err := ReadTextFile(filePath)
	if err != nil {
		t.Fatalf("ReadTextFile() error = %v", err)
	}

	if readContent != originalContent {
		t.Errorf("Content mismatch: got %q, want %q", readContent, originalContent)
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkFindFiles(b *testing.B) {
	tempDir, _ := os.MkdirTemp("", "benchmark_*")
	defer os.RemoveAll(tempDir)

	// Create some test files
	for i := 0; i < 10; i++ {
		os.WriteFile(filepath.Join(tempDir, "file"+string(rune(i))+".txt"), []byte("test"), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindFiles(tempDir, ".txt")
	}
}

func BenchmarkReadTextFile(b *testing.B) {
	tempDir, _ := os.MkdirTemp("", "benchmark_*")
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "benchmark.txt")
	os.WriteFile(filePath, []byte("benchmark content"), 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReadTextFile(filePath)
	}
}
