package responses

type SignInResponse struct {
	AccessToken                      string `json:"token"`
	AccessTokenExpiredTimeInSeconds  uint64 `json:"tokenExpiredTimeInSeconds"`
	RefreshToken                     string `json:"refreshToken"`
	RefreshTokenExpiredTimeInSeconds uint64 `json:"refreshTokenExpiredTimeInSeconds"`
}
