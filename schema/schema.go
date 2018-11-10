// Code generated by go-jsonschema-compiler. DO NOT EDIT.

package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sourcegraph/go-jsonschema/jsonschema"
)

type AWSCodeCommitConnection struct {
	AccessKeyID                 string `json:"accessKeyID"`
	InitialRepositoryEnablement bool   `json:"initialRepositoryEnablement,omitempty"`
	Region                      string `json:"region"`
	RepositoryPathPattern       string `json:"repositoryPathPattern,omitempty"`
	SecretAccessKey             string `json:"secretAccessKey"`
}
type Action struct {
	ActionItem       *ActionItem   `json:"actionItem,omitempty"`
	Category         string        `json:"category,omitempty"`
	Command          string        `json:"command,omitempty"`
	CommandArguments []interface{} `json:"commandArguments,omitempty"`
	IconURL          string        `json:"iconURL,omitempty"`
	Id               string        `json:"id,omitempty"`
	Title            string        `json:"title,omitempty"`
}

// ActionItem description: The action item.
type ActionItem struct {
	Description string `json:"description,omitempty"`
	IconURL     string `json:"iconURL,omitempty"`
	Label       string `json:"label,omitempty"`
}

// AuthAccessTokens description: Settings for access tokens, which enable external tools to access the Sourcegraph API with the privileges of the user.
type AuthAccessTokens struct {
	Allow string `json:"allow,omitempty"`
}

// AuthProviderCommon description: Common properties for authentication providers.
type AuthProviderCommon struct {
	DisplayName string `json:"displayName,omitempty"`
}
type AuthProviders struct {
	Builtin       *BuiltinAuthProvider
	Saml          *SAMLAuthProvider
	Openidconnect *OpenIDConnectAuthProvider
	HttpHeader    *HTTPHeaderAuthProvider
}

func (v AuthProviders) MarshalJSON() ([]byte, error) {
	if v.Builtin != nil {
		return json.Marshal(v.Builtin)
	}
	if v.Saml != nil {
		return json.Marshal(v.Saml)
	}
	if v.Openidconnect != nil {
		return json.Marshal(v.Openidconnect)
	}
	if v.HttpHeader != nil {
		return json.Marshal(v.HttpHeader)
	}
	return nil, errors.New("tagged union type must have exactly 1 non-nil field value")
}
func (v *AuthProviders) UnmarshalJSON(data []byte) error {
	var d struct {
		DiscriminantProperty string `json:"type"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	switch d.DiscriminantProperty {
	case "builtin":
		return json.Unmarshal(data, &v.Builtin)
	case "http-header":
		return json.Unmarshal(data, &v.HttpHeader)
	case "openidconnect":
		return json.Unmarshal(data, &v.Openidconnect)
	case "saml":
		return json.Unmarshal(data, &v.Saml)
	}
	return fmt.Errorf("tagged union type must have a %q property whose value is one of %s", "type", []string{"builtin", "saml", "openidconnect", "http-header"})
}

// AuthnProvider description: Identifies the authentication provider to use to identify users to GitLab.
type AuthnProvider struct {
	ConfigID       string `json:"configID"`
	GitlabProvider string `json:"gitlabProvider"`
	Type           string `json:"type"`
}

// Authorization description: If non-null, enables GitLab permission checks. This requires that the value of `token` be an access token with "sudo" and "api" scopes.
type Authorization struct {
	AuthnProvider AuthnProvider `json:"authnProvider"`
	Ttl           string        `json:"ttl,omitempty"`
}
type BitbucketServerConnection struct {
	Certificate                 string `json:"certificate,omitempty"`
	ExcludePersonalRepositories bool   `json:"excludePersonalRepositories,omitempty"`
	GitURLType                  string `json:"gitURLType,omitempty"`
	InitialRepositoryEnablement bool   `json:"initialRepositoryEnablement,omitempty"`
	Password                    string `json:"password,omitempty"`
	RepositoryPathPattern       string `json:"repositoryPathPattern,omitempty"`
	Token                       string `json:"token,omitempty"`
	Url                         string `json:"url"`
	Username                    string `json:"username,omitempty"`
}

// BuiltinAuthProvider description: Configures the builtin username-password authentication provider.
type BuiltinAuthProvider struct {
	AllowSignup bool   `json:"allowSignup,omitempty"`
	Type        string `json:"type"`
}

// CloneURLToRepositoryName description: Describes a mapping from clone URL to repository name. The `from` field contains a regular expression with named capturing groups. The `to` field contains a template string that references capturing group names. For instance, if `from` is "^../(?P<name>\w+)$" and `to` is "github.com/user/{name}", the clone URL "../myRepository" would be mapped to the repository name "github.com/user/myRepository".
type CloneURLToRepositoryName struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Contributions description: Features contributed by this extension. Extensions may also register certain types of contributions dynamically.
type Contributions struct {
	Actions       []*Action          `json:"actions,omitempty"`
	Configuration *jsonschema.Schema `json:"configuration,omitempty"`
	Menus         *Menus             `json:"menus,omitempty"`
}

// Discussions description: Configures Sourcegraph code discussions.
type Discussions struct {
	AbuseEmails     []string `json:"abuseEmails,omitempty"`
	AbuseProtection bool     `json:"abuseProtection,omitempty"`
}

// ExperimentalFeatures description: Experimental features to enable or disable. Features that are now enabled by default are marked as deprecated.
type ExperimentalFeatures struct {
	CanonicalURLRedirect string `json:"canonicalURLRedirect,omitempty"`
	Discussions          string `json:"discussions,omitempty"`
	JumpToDefOSSIndex    string `json:"jumpToDefOSSIndex,omitempty"`
	UpdateScheduler2     string `json:"updateScheduler2,omitempty"`
}

// ExtensionRepository description: The location of the version control repository for this extension.
type ExtensionRepository struct {
	Type string `json:"type,omitempty"`
	Url  string `json:"url"`
}

// Extensions description: Configures Sourcegraph extensions.
type Extensions struct {
	AllowRemoteExtensions []string    `json:"allowRemoteExtensions,omitempty"`
	Disabled              *bool       `json:"disabled,omitempty"`
	RemoteRegistry        interface{} `json:"remoteRegistry,omitempty"`
}
type GitHubConnection struct {
	Certificate                 string   `json:"certificate,omitempty"`
	GitURLType                  string   `json:"gitURLType,omitempty"`
	InitialRepositoryEnablement bool     `json:"initialRepositoryEnablement,omitempty"`
	Repos                       []string `json:"repos,omitempty"`
	RepositoryPathPattern       string   `json:"repositoryPathPattern,omitempty"`
	RepositoryQuery             []string `json:"repositoryQuery,omitempty"`
	Token                       string   `json:"token"`
	Url                         string   `json:"url"`
}
type GitLabConnection struct {
	Authorization               *Authorization `json:"authorization,omitempty"`
	Certificate                 string         `json:"certificate,omitempty"`
	GitURLType                  string         `json:"gitURLType,omitempty"`
	InitialRepositoryEnablement bool           `json:"initialRepositoryEnablement,omitempty"`
	ProjectQuery                []string       `json:"projectQuery,omitempty"`
	RepositoryPathPattern       string         `json:"repositoryPathPattern,omitempty"`
	Token                       string         `json:"token"`
	Url                         string         `json:"url"`
}
type GitoliteConnection struct {
	Blacklist                  string `json:"blacklist,omitempty"`
	Host                       string `json:"host"`
	PhabricatorMetadataCommand string `json:"phabricatorMetadataCommand,omitempty"`
	Prefix                     string `json:"prefix"`
}

// HTTPHeaderAuthProvider description: Configures the HTTP header authentication provider (which authenticates users by consulting an HTTP request header set by an authentication proxy such as https://github.com/bitly/oauth2_proxy).
type HTTPHeaderAuthProvider struct {
	Type           string `json:"type"`
	UsernameHeader string `json:"usernameHeader"`
}

// IMAPServerConfig description: Optional. The IMAP server used to retrieve emails (such as code discussion reply emails).
type IMAPServerConfig struct {
	Host     string `json:"host"`
	Password string `json:"password,omitempty"`
	Port     int    `json:"port"`
	Username string `json:"username,omitempty"`
}
type Langservers struct {
	Address               string                 `json:"address,omitempty"`
	Disabled              bool                   `json:"disabled,omitempty"`
	InitializationOptions map[string]interface{} `json:"initializationOptions,omitempty"`
	Language              string                 `json:"language"`
	Metadata              *Metadata              `json:"metadata,omitempty"`
}
type Links struct {
	Blob       string `json:"blob,omitempty"`
	Commit     string `json:"commit,omitempty"`
	Repository string `json:"repository,omitempty"`
	Tree       string `json:"tree,omitempty"`
}

// Log description: Configuration for logging and alerting, including to external services.
type Log struct {
	Sentry *Sentry `json:"sentry,omitempty"`
}
type MenuItem struct {
	Action string `json:"action,omitempty"`
	Alt    string `json:"alt,omitempty"`
	When   string `json:"when,omitempty"`
}

// Menus description: Describes where to place actions in menus.
type Menus struct {
	CommandPalette []*MenuItem `json:"commandPalette,omitempty"`
	EditorTitle    []*MenuItem `json:"editor/title,omitempty"`
	Help           []*MenuItem `json:"help,omitempty"`
}

// Metadata description: Language server metadata. Used to populate various UI elements.
type Metadata struct {
	DocsURL      string `json:"docsURL,omitempty"`
	Experimental bool   `json:"experimental,omitempty"`
	HomepageURL  string `json:"homepageURL,omitempty"`
	IssuesURL    string `json:"issuesURL,omitempty"`
}

// OpenIDConnectAuthProvider description: Configures the OpenID Connect authentication provider for SSO.
type OpenIDConnectAuthProvider struct {
	ClientID           string `json:"clientID"`
	ClientSecret       string `json:"clientSecret"`
	ConfigID           string `json:"configID,omitempty"`
	DisplayName        string `json:"displayName,omitempty"`
	Issuer             string `json:"issuer"`
	RequireEmailDomain string `json:"requireEmailDomain,omitempty"`
	Type               string `json:"type"`
}

// ParentSourcegraph description: URL to fetch unreachable repository details from. Defaults to "https://sourcegraph.com"
type ParentSourcegraph struct {
	Url string `json:"url,omitempty"`
}
type Phabricator struct {
	Repos []*Repos `json:"repos,omitempty"`
	Token string   `json:"token,omitempty"`
	Url   string   `json:"url,omitempty"`
}
type Repos struct {
	Callsign string `json:"callsign"`
	Path     string `json:"path"`
}
type Repository struct {
	Links *Links `json:"links,omitempty"`
	Path  string `json:"path"`
	Type  string `json:"type,omitempty"`
	Url   string `json:"url"`
}
type ReviewBoard struct {
	Url string `json:"url,omitempty"`
}

// SAMLAuthProvider description: Configures the SAML authentication provider for SSO.
//
// Note: if you are using IdP-initiated login, you must have *at most one* SAMLAuthProvider in the `auth.providers` array.
type SAMLAuthProvider struct {
	ConfigID                                 string `json:"configID,omitempty"`
	DisplayName                              string `json:"displayName,omitempty"`
	IdentityProviderMetadata                 string `json:"identityProviderMetadata,omitempty"`
	IdentityProviderMetadataURL              string `json:"identityProviderMetadataURL,omitempty"`
	InsecureSkipAssertionSignatureValidation bool   `json:"insecureSkipAssertionSignatureValidation,omitempty"`
	NameIDFormat                             string `json:"nameIDFormat,omitempty"`
	ServiceProviderCertificate               string `json:"serviceProviderCertificate,omitempty"`
	ServiceProviderIssuer                    string `json:"serviceProviderIssuer,omitempty"`
	ServiceProviderPrivateKey                string `json:"serviceProviderPrivateKey,omitempty"`
	SignRequests                             *bool  `json:"signRequests,omitempty"`
	Type                                     string `json:"type"`
}

// SMTPServerConfig description: The SMTP server used to send transactional emails (such as email verifications, reset-password emails, and notifications).
type SMTPServerConfig struct {
	Authentication string `json:"authentication"`
	Domain         string `json:"domain,omitempty"`
	Host           string `json:"host"`
	Password       string `json:"password,omitempty"`
	Port           int    `json:"port"`
	Username       string `json:"username,omitempty"`
}
type SearchSavedQueries struct {
	Description    string `json:"description"`
	Key            string `json:"key"`
	Notify         bool   `json:"notify,omitempty"`
	NotifySlack    bool   `json:"notifySlack,omitempty"`
	Query          string `json:"query"`
	ShowOnHomepage bool   `json:"showOnHomepage,omitempty"`
}
type SearchScope struct {
	Description string `json:"description,omitempty"`
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Value       string `json:"value"`
}

// Sentry description: Configuration for Sentry
type Sentry struct {
	Dsn string `json:"dsn,omitempty"`
}

// Settings description: Configuration settings for users and organizations on Sourcegraph.
type Settings struct {
	Extensions             map[string]bool           `json:"extensions,omitempty"`
	Motd                   []string                  `json:"motd,omitempty"`
	NotificationsSlack     *SlackNotificationsConfig `json:"notifications.slack,omitempty"`
	SearchRepositoryGroups map[string][]string       `json:"search.repositoryGroups,omitempty"`
	SearchSavedQueries     []*SearchSavedQueries     `json:"search.savedQueries,omitempty"`
	SearchScopes           []*SearchScope            `json:"search.scopes,omitempty"`
}

// SiteConfiguration description: Configuration for a Sourcegraph site.
type SiteConfiguration struct {
	AuthAccessTokens                  *AuthAccessTokens            `json:"auth.accessTokens,omitempty"`
	AuthDisableAccessTokens           bool                         `json:"auth.disableAccessTokens,omitempty"`
	AuthProviders                     []AuthProviders              `json:"auth.providers,omitempty"`
	AuthPublic                        bool                         `json:"auth.public,omitempty"`
	AuthSessionExpiry                 string                       `json:"auth.sessionExpiry,omitempty"`
	AuthUserOrgMap                    map[string][]string          `json:"auth.userOrgMap,omitempty"`
	AwsCodeCommit                     []*AWSCodeCommitConnection   `json:"awsCodeCommit,omitempty"`
	BitbucketServer                   []*BitbucketServerConnection `json:"bitbucketServer,omitempty"`
	BlacklistGoGet                    []string                     `json:"blacklistGoGet,omitempty"`
	CorsOrigin                        string                       `json:"corsOrigin,omitempty"`
	DisableAutoGitUpdates             bool                         `json:"disableAutoGitUpdates,omitempty"`
	DisableBrowserExtension           bool                         `json:"disableBrowserExtension,omitempty"`
	DisableBuiltInSearches            bool                         `json:"disableBuiltInSearches,omitempty"`
	DisablePublicRepoRedirects        bool                         `json:"disablePublicRepoRedirects,omitempty"`
	Discussions                       *Discussions                 `json:"discussions,omitempty"`
	DontIncludeSymbolResultsByDefault bool                         `json:"dontIncludeSymbolResultsByDefault,omitempty"`
	EmailAddress                      string                       `json:"email.address,omitempty"`
	EmailImap                         *IMAPServerConfig            `json:"email.imap,omitempty"`
	EmailSmtp                         *SMTPServerConfig            `json:"email.smtp,omitempty"`
	ExecuteGradleOriginalRootPaths    string                       `json:"executeGradleOriginalRootPaths,omitempty"`
	ExperimentalFeatures              *ExperimentalFeatures        `json:"experimentalFeatures,omitempty"`
	Extensions                        *Extensions                  `json:"extensions,omitempty"`
	ExternalURL                       string                       `json:"externalURL,omitempty"`
	GitCloneURLToRepositoryName       []*CloneURLToRepositoryName  `json:"git.cloneURLToRepositoryName,omitempty"`
	GitMaxConcurrentClones            int                          `json:"gitMaxConcurrentClones,omitempty"`
	Github                            []*GitHubConnection          `json:"github,omitempty"`
	GithubClientID                    string                       `json:"githubClientID,omitempty"`
	GithubClientSecret                string                       `json:"githubClientSecret,omitempty"`
	Gitlab                            []*GitLabConnection          `json:"gitlab,omitempty"`
	Gitolite                          []*GitoliteConnection        `json:"gitolite,omitempty"`
	HtmlBodyBottom                    string                       `json:"htmlBodyBottom,omitempty"`
	HtmlBodyTop                       string                       `json:"htmlBodyTop,omitempty"`
	HtmlHeadBottom                    string                       `json:"htmlHeadBottom,omitempty"`
	HtmlHeadTop                       string                       `json:"htmlHeadTop,omitempty"`
	HttpToHttpsRedirect               interface{}                  `json:"httpToHttpsRedirect,omitempty"`
	Langservers                       []*Langservers               `json:"langservers,omitempty"`
	LicenseKey                        string                       `json:"licenseKey,omitempty"`
	LightstepAccessToken              string                       `json:"lightstepAccessToken,omitempty"`
	LightstepProject                  string                       `json:"lightstepProject,omitempty"`
	Log                               *Log                         `json:"log,omitempty"`
	MaxReposToSearch                  int                          `json:"maxReposToSearch,omitempty"`
	NoGoGetDomains                    string                       `json:"noGoGetDomains,omitempty"`
	ParentSourcegraph                 *ParentSourcegraph           `json:"parentSourcegraph,omitempty"`
	Phabricator                       []*Phabricator               `json:"phabricator,omitempty"`
	PrivateArtifactRepoID             string                       `json:"privateArtifactRepoID,omitempty"`
	PrivateArtifactRepoPassword       string                       `json:"privateArtifactRepoPassword,omitempty"`
	PrivateArtifactRepoURL            string                       `json:"privateArtifactRepoURL,omitempty"`
	PrivateArtifactRepoUsername       string                       `json:"privateArtifactRepoUsername,omitempty"`
	RepoListUpdateInterval            int                          `json:"repoListUpdateInterval,omitempty"`
	ReposList                         []*Repository                `json:"repos.list,omitempty"`
	ReviewBoard                       []*ReviewBoard               `json:"reviewBoard,omitempty"`
	SearchIndexEnabled                *bool                        `json:"search.index.enabled,omitempty"`
	SiteID                            string                       `json:"siteID,omitempty"`
	TlsCert                           string                       `json:"tlsCert,omitempty"`
	TlsKey                            string                       `json:"tlsKey,omitempty"`
	UpdateChannel                     string                       `json:"update.channel,omitempty"`
	UseJaeger                         bool                         `json:"useJaeger,omitempty"`
}

// SlackNotificationsConfig description: Configuration for sending notifications to Slack.
type SlackNotificationsConfig struct {
	WebhookURL string `json:"webhookURL"`
}

// SourcegraphExtensionManifest description: The Sourcegraph extension manifest describes the extension and the features it provides.
type SourcegraphExtensionManifest struct {
	ActivationEvents []string                `json:"activationEvents"`
	Args             *map[string]interface{} `json:"args,omitempty"`
	Contributes      *Contributions          `json:"contributes,omitempty"`
	Description      string                  `json:"description,omitempty"`
	Icon             string                  `json:"icon,omitempty"`
	Readme           string                  `json:"readme,omitempty"`
	Repository       *ExtensionRepository    `json:"repository,omitempty"`
	Title            string                  `json:"title,omitempty"`
	Url              string                  `json:"url"`
}
