package debuglog

import (
	"fmt"

	"github.com/oaswrap/spec/openapi"
)

type Logger struct {
	prefix string
	logger openapi.Logger
}

func NewLogger(prefix string, logger openapi.Logger) *Logger {
	return &Logger{logger: logger, prefix: fmt.Sprintf("[%s]", prefix)}
}

func (l *Logger) Printf(format string, v ...any) {
	l.logger.Printf(l.prefix+" "+format, v...)
}

func (l *Logger) LogOp(method, path, action, value string) {
	l.Printf("%s %s → %s: %s", method, path, action, value)
}

func (l *Logger) LogAction(action, value string) {
	l.Printf("%s: %s", action, value)
}

func (l *Logger) LogContact(contact *openapi.Contact) {
	if contact == nil {
		return
	}
	var contactInfo string
	if contact.Name != "" {
		contactInfo += "name: " + contact.Name + ", "
	}
	if contact.Email != "" {
		contactInfo += "email: " + contact.Email + ", "
	}
	if contact.URL != "" {
		contactInfo += "url: " + contact.URL
	}
	if contactInfo != "" {
		l.Printf("set contact: %s", contactInfo)
	}
}

func (l *Logger) LogLicense(license *openapi.License) {
	var licenseInfo string
	if license.Name != "" {
		licenseInfo += "name: " + license.Name + ", "
	}
	if license.URL != "" {
		licenseInfo += "url: " + license.URL
	}
	if licenseInfo != "" {
		l.Printf("set license: %s", licenseInfo)
	}
}

func (l *Logger) LogExternalDocs(externalDocs *openapi.ExternalDocs) {
	var docsInfo string
	if externalDocs.URL != "" {
		docsInfo += "url: " + externalDocs.URL
	}
	if externalDocs.Description != "" {
		docsInfo += ", description: " + externalDocs.Description
	}
	if docsInfo != "" {
		l.Printf("set external docs: %s", docsInfo)
	}
}

func (l *Logger) LogServer(server openapi.Server) {
	var serverInfo string
	serverInfo += "url: " + server.URL
	if server.Description != nil {
		serverInfo += ", description: " + *server.Description
	}
	if len(server.Variables) > 0 {
		serverInfo += ", variables: "
		for name, variable := range server.Variables {
			serverInfo += name + ": " + variable.Default + ", "
		}
		serverInfo = serverInfo[:len(serverInfo)-2] // Remove trailing comma and space
	}
	l.Printf("set server: %s", serverInfo)
}

func (l *Logger) LogTag(tag openapi.Tag) {
	tagInfo := "name: " + tag.Name
	if tag.Description != "" {
		tagInfo += ", description: " + tag.Description
	}
	if tag.ExternalDocs != nil {
		tagInfo += ", external docs: " + tag.ExternalDocs.URL
		if tag.ExternalDocs.Description != "" {
			tagInfo += " (" + tag.ExternalDocs.Description + ")"
		}
	}
	l.Printf("add tag: %s", tagInfo)
}

func (l *Logger) LogSecurityScheme(name string, scheme *openapi.SecurityScheme) {
	var typeInfo string
	switch {
	case scheme.APIKey != nil:
		typeInfo = "APIKey"
	case scheme.HTTPBearer != nil:
		typeInfo = "HTTPBearer"
	case scheme.OAuth2 != nil:
		typeInfo = "OAuth2"
	default:
		typeInfo = "Unknown"
	}
	schemeInfo := "name: " + name + ", type: " + typeInfo
	if scheme.Description != nil {
		schemeInfo += ", description: " + *scheme.Description
	}
	l.Printf("add security scheme: %s", schemeInfo)
}
