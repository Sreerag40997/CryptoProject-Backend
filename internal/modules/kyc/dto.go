package kyc

import "mime/multipart"


// for the request
type SubmitKYCRequest struct {
	// Personal Info
	FullName string `form:"full_name" validate:"required,min=3"`
	DOB      string `form:"dob" validate:"required"`

	// Identity
	AadhaarNumber string `form:"aadhaar" validate:"required,len=12,numeric"`
	PANNumber     string `form:"pan" validate:"required,len=10"`

	// Banking
	AccountNumber string `form:"account_number" validate:"required,min=8,max=18,numeric"`
	IFSCCode      string `form:"ifsc" validate:"required,len=11"`

	// Files
	AadhaarFront *multipart.FileHeader `form:"aadhaar_front" validate:"required"`
	AadhaarBack  *multipart.FileHeader `form:"aadhaar_back" validate:"required"`
	PANFile      *multipart.FileHeader `form:"pan_file" validate:"required"`
	Selfie       *multipart.FileHeader `form:"selfie" validate:"required"`
}

// for update the kyc
type UpdateKYCRequest struct {
	FullName string `form:"full_name"`
	DOB      string `form:"dob"`

	AadhaarNumber string `form:"aadhaar"`
	PANNumber     string `form:"pan"`

	AccountNumber string `form:"account_number"`
	IFSCCode      string `form:"ifsc"`

	AadhaarFront *multipart.FileHeader `form:"aadhaar_front"`
	AadhaarBack  *multipart.FileHeader `form:"aadhaar_back"`
	PANFile      *multipart.FileHeader `form:"pan_file"`
	Selfie       *multipart.FileHeader `form:"selfie"`
}

// for admin review
type ReviewKYCRequest struct {
	Status string `json:"status" validate:"required,oneof=approved rejected"`
	Reason string `json:"reason"`
}

// for response
type KYCResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Status   string `json:"status"`

	AadhaarMasked string `json:"aadhaar"`
	PANMasked     string `json:"pan"`

	AccountMasked string `json:"account_number"`
	IFSCCode      string `json:"ifsc"`

	AadhaarFrontURL string `json:"aadhaar_front_url"`
	AadhaarBackURL  string `json:"aadhaar_back_url"`
	PANURL          string `json:"pan_url"`
	SelfieURL       string `json:"selfie_url"`
}