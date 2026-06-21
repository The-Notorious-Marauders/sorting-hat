package constants

const (
	ErrorCode_BadRequest_CouldNotBindMultipartForm = 400001
	ErrorCode_BadRequest_InvalidRequest            = 400002
	ErrorCode_BadRequest_PasswordTooLong           = 400003

	ErrorCode_Unauthorized_InvalidSignature     = 401001
	ErrorCode_Unauthorized_InvalidSignInMessage = 401002
	ErrorCode_Unauthorized_InvalidJWT           = 401003
	ErrorCode_Unauthorized_MalformedJWT         = 401004
	ErrorCode_Unauthorized_InvalidCredentials   = 401005

	ErrorCode_Conflict_UsernameAlreadyExists = 409001

	ErrorCode_InternalError_CouldNotMakeJWTSignedString = 500004
	ErrorCode_InternalError_CouldNotFindUser            = 500005
	ErrorCode_InternalError_CouldNotUpdateLastLogin     = 500006
	ErrorCode_InternalError_CouldNotHashPassword        = 500007
	ErrorCode_InternalError_CouldNotCreateUser          = 500008
)
