// Package service provides the business logic layer exposed to the Wails frontend.
// It wraps repository operations with validation and exposes methods bound to
// the JavaScript runtime.
package service

var appVersion = "dev"

// AppService exposes application-level information such as the current version.
type AppService struct{}

// NewAppService creates a new AppService instance.
func NewAppService() *AppService { return &AppService{} }

// SetVersion sets the global application version string. It is called once at startup.
func SetVersion(v string) { appVersion = v }

// Version returns the current application version string.
func (s *AppService) Version() string { return appVersion }
