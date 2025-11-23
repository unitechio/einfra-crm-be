package domain

import "time"

// FileInfo represents file/directory information
type FileInfo struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	IsDir       bool      `json:"is_dir"`
	ModTime     time.Time `json:"mod_time"`
	Permissions string    `json:"permissions"`
	Owner       string    `json:"owner,omitempty"`
	Group       string    `json:"group,omitempty"`
}

// FileBrowserRequest represents a request to browse files
type FileBrowserRequest struct {
	VolumeName string `json:"volume_name" binding:"required"`
	Path       string `json:"path"`
}

// FileUploadRequest represents a file upload request
type FileUploadRequest struct {
	VolumeName string `json:"volume_name" binding:"required"`
	Path       string `json:"path" binding:"required"`
	Filename   string `json:"filename" binding:"required"`
}

// FileDownloadRequest represents a file download request
type FileDownloadRequest struct {
	VolumeName string `json:"volume_name" binding:"required"`
	Path       string `json:"path" binding:"required"`
}

// FileDeleteRequest represents a file delete request
type FileDeleteRequest struct {
	VolumeName string `json:"volume_name" binding:"required"`
	Path       string `json:"path" binding:"required"`
}

// CreateFolderRequest represents a create folder request
type CreateFolderRequest struct {
	VolumeName string `json:"volume_name" binding:"required"`
	Path       string `json:"path" binding:"required"`
	FolderName string `json:"folder_name" binding:"required"`
}
