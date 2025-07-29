package main

import "mime/multipart"

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Pet struct {
	ID        int64     `json:"id"`
	Category  *Category `json:"category,omitempty"`
	Name      string    `json:"name" validate:"required"`
	PhotoURLs []string  `json:"photoUrls" validate:"required"`
	Tags      []Tag     `json:"tags,omitempty"`
	Status    string    `json:"status,omitempty" enum:"available,pending,sold"`
}

type ApiResponse struct {
	Code    int32  `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type FindPetByIdRequest struct {
	ID int64 `params:"petId" path:"petId"`
}

type FindPetsByStatusRequest struct {
	Status string `query:"status" validate:"required" enum:"available,pending,sold"`
}

type FindPetsByTagsRequest struct {
	Tags []string `query:"tags" required:"false"`
}

type DeletePetRequest struct {
	ApiKey string `header:"api_key"`
	ID     int64  `params:"petId" path:"petId"`
}

type UpdatePetFormDataRequest struct {
	ID     int64  `params:"petId" path:"petId"`
	Name   string `formData:"name" validate:"required"`
	Status string `formData:"status" enum:"available,pending,sold"`
}

type UploadImageRequest struct {
	ID                 int64           `params:"petId" path:"petId"`
	AdditionalMetaData string          `query:"additionalMetadata"`
	_                  *multipart.File `contentType:"application/octet-stream"`
}
