package yadisk

import (
	"fmt"
	"strconv"
	"strings"
)

type YaDisk interface {
	// Disk

	// Get user disk meta information.
	GetDisk(fields []string) (d *Disk, e error)

	// Trash

	// Empty trash.
	//
	// If the deletion occurs asynchronously, it will return a response with status 202 and a link to the asynchronous operation.
	// Otherwise, it will return a response with status 204 and an empty body.
	//
	// If the path parameter is not specified or points to the root of the Recycle Bin,
	// the recycle bin will be completely cleared, otherwise only the resource pointed to by the path will be deleted from the Recycle Bin.
	ClearTrash(fields []string, forceAsync bool, path string) (l *Link, e error)

	// Get the contents of the Trash.
	GetTrashResource(path string, fields []string, limit int, offset int, previewCrop bool, previewSize string, sort string) (r *TrashResource, e error)

	// Recover Resource from Trash.
	//
	// If recovery is asynchronous, it will return a response with code 202 and a link to the asynchronous operation.
	// Otherwise, it will return a response with code 201 and a link to the created resource.
	RestoreFromTrash(path string, fields []string, forceAsync bool, name string, overwrite bool) (l *Link, e error)

	// Resource

	// Delete file or folder.
	//
	// By default, delete the resource in the trash.
	// To delete a resource without placing it in the trash, you must specify the parameter permanently = true.
	//
	// If the deletion occurs asynchronously, it will return a response with status 202 and a link to the asynchronous operation.
	// Otherwise, it will return a response with status 204 and an empty body.
	DeleteResource(path string, fields []string, forceAsync bool, md5 string, permanently bool) (l *Link, e error)

	// Get meta information about a file or directory.
	GetResource(path string, fields []string, limit int, offset int, previewCrop bool, previewSize string, sort string) (r *Resource, e error)

	// Create directory.
	CreateResource(path string, fields []string) (l *Link, e error)

	// Update User Resource Data.
	UpdateResource(path string, fields []string, body *ResourcePatch) (r *Resource, e error)

	// Create a copy of the file or folder.
	//
	// If copying occurs asynchronously, it will return a response with code 202 and a link to the asynchronous operation.
	// Otherwise, it will return a response with code 201 and a link to the created resource.
	CopyResource(from string, path string, fields []string, forceAsync bool, overwrite bool) (l *Link, e error)

	// Move a file or folder.
	//
	// If the movement occurs asynchronously, it will return a response with code 202 and a link to the asynchronous operation.
	// Otherwise, it will return a response with code 201 and a link to the created resource.
	MoveResource(from string, path string, fields []string, forceAsync bool, overwrite bool) (l *Link, e error)

	// Get link to download file.
	GetResourceDownloadLink(path string, fields []string) (l *Link, e error)

	// Get file list sorted by name.
	GetFlatFilesList(fields []string, limit int, mediaType string, offset int, previewCrop bool, previewSize string, sort string) (l *FilesResourceList, e error)

	// Get a list of files ordered by download date.
	GetLastUploadedFilesList(fields []string, limit int, mediaType string, previewCrop bool, previewSize string) (l *LastUploadedResourceList, e error)

	// Get a list of published resources.
	//
	// resourceType value: "","dir","file".
	ListPublicResources(fields []string, limit int, offset int, previewCrop bool, previewSize string, resourceType string) (l *PublicResourcesList, e error)

	// Publish a resource.
	PublishResource(path string, fields []string) (l *Link, e error)

	// Unpublish a resource.
	UnpublishResource(path string, fields []string) (l *Link, e error)

	// Upload file to Disk by URL.
	//
	// Download asynchronously.
	//
	// Therefore, in response to the request, a reference to the asynchronous operation is returned.
	UploadExternalResource(path string, externalURL string, disableRedirects bool, fields []string) (l *Link, e error)

	// Get file download link.
	GetResourceUploadLink(path string, fields []string, overwrite bool) (l *ResourceUploadLink, e error)

	// Public

	// Get meta-information about a public file or directory.
	GetPublicResource(publicKey string, fields []string, limit int, offset int, path string, previewCrop bool, previewSize string, sort string) (r *PublicResource, e error)

	// Get a link to download a public resource.
	GetPublicResourceDownloadLink(publicKey string, fields []string, path string) (l *Link, e error)

	// Save the public resource to the Downloads folder.
	//
	// If saving occurs asynchronously, it will return a response with code 202 and a link to the asynchronous operation.
	// Otherwise, it will return a response with code 201 and a link to the created resource.
	SaveToDiskPublicResource(publicKey string, fields []string, forceAsync bool, name string, path string, savePath string) (l *Link, e error)

	// Operations

	// Get the status of an asynchronous operation.
	GetOperationStatus(operationID string, fields []string) (s *OperationStatus, e error)
}

type yandexDisk struct {
	Token  *Token // required
	client *client
}

// Token for access to Yandex.Disk Rest-API
type Token struct {
	AccessToken string
}
type responseInfo struct {
	Status     string
	StatusCode int
}

func (ri *responseInfo) setResponseInfo(status string, statusCode int) {
	ri.Status = status
	ri.StatusCode = statusCode
}

type ResourcePatch struct {
	// User Attribute Resource
	// Structure fields must have a `json` tag
	CustomProperties interface{} `json:"custom_properties"`
}

type Error struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	ErrorID     string `json:"error"`
}

func (e *Error) Error() string {
	return e.ErrorID
}

type PerformUpload struct {
}

func (pu *PerformUpload) handleError(ri responseInfo) (e error) {
	err := new(Error)
	err.Description = ri.Status
	err.ErrorID = strconv.Itoa(ri.StatusCode)
	e = err
	switch strings.ToLower(ri.Status) {
	case "201 created",
		"202 accepted":
		return nil
	case "412 precondition failed",
		"413 payload too large",
		"500 internal server error",
		"503 service unavailable",
		"507 insufficient storage":
		return
	default:
		panic(fmt.Sprintf("undefined status: %s", ri.Status))
	}
}

type DeleteResource struct {
}

type Disk struct {
	MaxFileSize                int  `json:"max_file_size"`
	UnlimitedAutouploadEnabled bool `json:"unlimited_autoupload_enabled"`
	TotalSpace                 int  `json:"total_space"`
	TrashSize                  int  `json:"trash_size"`
	IsPaid                     bool `json:"is_paid"`
	UsedSpace                  int  `json:"used_space"`
	SystemFolders              struct {
		Odnoklassniki string `json:"odnoklassniki"`
		Google        string `json:"google"`
		Instagram     string `json:"instagram"`
		Vkontakte     string `json:"vkontakte"`
		Mailru        string `json:"mailru"`
		Downloads     string `json:"downloads"`
		Applications  string `json:"applications"`
		Facebook      string `json:"facebook"`
		Social        string `json:"social"`
		Screenshots   string `json:"screenshots"`
		Photostream   string `json:"photostream"`
	} `json:"system_folders"`
	User     User `json:"user"`
	Revision int  `json:"revision"`
}

type baseResource struct {
	ResourceID     string     `json:"resource_id"`
	Share          Share      `json:"share"`
	File           string     `json:"file"`
	Size           int        `json:"size"`
	PhotosliceTime string     `json:"photoslice_time"`
	Exif           Exif       `json:"exif"`
	MediaType      string     `json:"media_type"`
	Sha256         string     `json:"sha256"`
	Type           string     `json:"type"`
	MimeType       string     `json:"mime_type"`
	Revision       int        `json:"revision"`
	PublicURL      string     `json:"public_url"`
	Path           string     `json:"path"`
	Md5            string     `json:"md5"`
	PublicKey      string     `json:"public_key"`
	Preview        string     `json:"preview"`
	Name           string     `json:"name"`
	Created        string     `json:"created"`
	Modified       string     `json:"modified"`
	CommentIds     CommentIds `json:"comment_ids"`
}

type Link struct {
	Href      string `json:"href"`
	Method    string `json:"method"`
	Templated bool   `json:"templated"`
}

type ResourceUploadLink struct {
	OperationID string `json:"operation_id"`
	Href        string `json:"href"`
	Method      string `json:"method"`
	Templated   bool   `json:"templated"`
}

type Resource struct {
	baseResource
	CustomProperties string   `json:"custom_properties"`
	Embedded         Embedded `json:"_embedded"`
}

type FilesResourceList struct {
	Items  []Resource `json:"items"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

type LastUploadedResourceList struct {
	Items []Resource `json:"items"`
	Limit int        `json:"limit"`
}

type PublicResourcesList struct {
	Items  []Resource `json:"items"`
	Type   string     `json:"type"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

type PublicResource struct {
	baseResource
	ViewsCount int            `json:"views_count"`
	Owner      Owner          `json:"owner"`
	Embedded   PublicEmbedded `json:"_embedded"`
}

type TrashResource struct {
	baseResource
	Embedded         TrashEmbedded `json:"_embedded"`
	CustomProperties string        `json:"custom_properties"`
	OriginPath       string        `json:"origin_path"`
	Deleted          string        `json:"deleted"`
}

type OperationStatus struct {
	Status string `json:"status"`
}

type User struct {
	Country     string `json:"country"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
	UID         string `json:"uid"`
}

type Owner struct {
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
	UID         string `json:"uid"`
}

type baseEmbedded struct {
	Sort   string `json:"sort"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Path   string `json:"path"`
	Total  int    `json:"total"`
}

type TrashEmbedded struct {
	baseEmbedded
	Items []TrashResource `json:"items"`
}

type PublicEmbedded struct {
	baseEmbedded
	Items []PublicResource `json:"items"`
}

type Embedded struct {
	baseEmbedded
	Items []Resource `json:"items"`
}

type Share struct {
	IsRoot  bool   `json:"is_root"`
	IsOwned bool   `json:"is_owned"`
	Rights  string `json:"rights"`
}

type Exif struct {
	DateTime string `json:"date_time"`
}

type CommentIds struct {
	PrivateResource string `json:"private_resource"`
	PublicResource  string `json:"public_resource"`
}
