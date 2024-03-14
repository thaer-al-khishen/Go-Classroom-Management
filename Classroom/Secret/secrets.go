package Secret

import "time"

var JwtKey = []byte("eyJhbGciOiJIUzI1NiJ9.ew0KICAic3ViIjogIjEyMzQ1Njc4OTAiLA0KICAibmFtZSI6ICJUaGFlciIsDQogICJpYXQiOiAxNTE2MjM5MDIyDQp9.cnRtVUVWpVBp1-2jE9yQTKfxKx6hah5nfO5lr0ceg8Q")

const AccessTokenExpiry = 15 * time.Minute
const RefreshTokenExpiry = 24 * time.Hour
