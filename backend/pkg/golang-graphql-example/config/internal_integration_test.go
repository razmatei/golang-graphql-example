//go:build integration

package config

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/log"
	"github.com/stretchr/testify/assert"
)

const reloadEventuallyWait = 10 * time.Second
const reloadEventuallyTick = 500 * time.Millisecond

func Test_managerimpl_Load(t *testing.T) {
	tests := []struct {
		name           string
		configs        map[string]string
		envVariables   map[string]string
		secretFiles    map[string]string
		expectedResult *Config
		wantErr        bool
	}{
		{
			name: "Configuration not found",
			configs: map[string]string{
				"config": "",
			},
			wantErr: true,
		},
		{
			name: "Not a yaml",
			configs: map[string]string{
				"config.yaml": "notayaml",
			},
			wantErr: true,
		},
		{
			name: "Empty",
			configs: map[string]string{
				"config.yaml": "",
			},
			wantErr: true,
		},
		{
			name: "default config",
			configs: map[string]string{
				"log.yaml": `
log:
  level: error
  format: text
`,
				"database.yaml": `
database:
  connectionUrl:
    value: host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable

`,
			},
			expectedResult: &Config{
				Log: &LogConfig{
					Format: "text",
					Level:  "error",
				},
				Tracing: &TracingConfig{Enabled: false, Type: TracingOtelHTTPType},
				Database: &DatabaseConfig{
					Driver:        "POSTGRES",
					ConnectionURL: &CredentialConfig{Value: "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"},
				},
				Server:         &ServerConfig{Port: 8080},
				InternalServer: &ServerConfig{Port: 9090},
				LockDistributor: &LockDistributorConfig{
					HeartbeatFrequency: "1s",
					LeaseDuration:      "3s",
					TableName:          "locks",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "config")
			if err != nil {
				t.Error(err)
				return
			}

			defer os.RemoveAll(dir) // clean up
			for k, v := range tt.configs {
				tmpfn := filepath.Join(dir, k)
				err = ioutil.WriteFile(tmpfn, []byte(v), 0666)
				if err != nil {
					t.Error(err)
					return
				}
			}

			// Set environment variables
			for k, v := range tt.envVariables {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// Create secret files
			for k, v := range tt.secretFiles {
				dirToCr := filepath.Dir(k)
				err = os.MkdirAll(dirToCr, 0666)
				if err != nil {
					t.Error(err)
					return
				}
				err = ioutil.WriteFile(k, []byte(v), 0666)
				if err != nil {
					t.Error(err)
					return
				}
				defer os.Remove(k)
			}

			ctx := &managerimpl{
				logger: log.NewLogger(),
			}

			// Load config
			err = ctx.Load(dir)

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Get configuration
			res := ctx.GetConfig()

			assert.Equal(t, tt.expectedResult, res)
		})
	}
}

func Test_Load_reload_config(t *testing.T) {
	dir, err := ioutil.TempDir("", "config-reload")
	assert.NoError(t, err)

	configs := map[string]string{
		"log.yaml": `
log:
  level: error
`,
		"database.yaml": `
database:
  connectionUrl:
    value: host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable

`,
		"tracing.yaml": `
tracing:
  enabled: true
  otelHttp:
    serverUrl: http://localhost:4318/v1/traces
`,
	}

	defer os.RemoveAll(dir) // clean up
	for k, v := range configs {
		tmpfn := filepath.Join(dir, k)
		err = ioutil.WriteFile(tmpfn, []byte(v), 0666)
		assert.NoError(t, err)
	}

	secretFiles := map[string]string{
		os.TempDir() + "/secret1": "VALUE1",
	}
	// Create secret files
	for k, v := range secretFiles {
		dirToCr := filepath.Dir(k)
		err = os.MkdirAll(dirToCr, 0666)
		assert.NoError(t, err)
		err = ioutil.WriteFile(k, []byte(v), 0666)
		assert.NoError(t, err)
		defer os.Remove(k)
	}

	ctx := &managerimpl{
		logger: log.NewLogger(),
	}

	reloadHookCalled := false
	ctx.AddOnChangeHook(&HookDefinition{
		Hook: func() error {
			reloadHookCalled = true
			return nil
		},
	})

	// Load config
	err = ctx.Load(dir)
	assert.NoError(t, err)
	// Get configuration
	res := ctx.GetConfig()

	assert.Equal(t, &Config{
		Log: &LogConfig{
			Level:  "error",
			Format: "json",
		},
		Server: &ServerConfig{
			Port: 8080,
		},
		InternalServer: &ServerConfig{
			Port: 9090,
		},
		Tracing: &TracingConfig{Enabled: true, Type: TracingOtelHTTPType, OtelHTTP: &TracingOtelHTTPConfig{ServerURL: "http://localhost:4318/v1/traces"}},
		Database: &DatabaseConfig{
			Driver:        "POSTGRES",
			ConnectionURL: &CredentialConfig{Value: "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"},
		},
		LockDistributor: &LockDistributorConfig{
			HeartbeatFrequency: "1s",
			LeaseDuration:      "3s",
			TableName:          "locks",
		},
	}, res)

	configs = map[string]string{
		"log.yaml": `
log:
  level: debug
  format: text
`,
	}

	defer os.RemoveAll(dir) // clean up
	for k, v := range configs {
		tmpfn := filepath.Join(dir, k)
		err = ioutil.WriteFile(tmpfn, []byte(v), 0666)
		assert.NoError(t, err)
	}

	assert.Eventually(
		t,
		func() bool {
			return reloadHookCalled
		},
		reloadEventuallyWait,
		reloadEventuallyTick,
	)

	// Get configuration
	res = ctx.GetConfig()

	assert.Equal(t, &Config{
		Log: &LogConfig{
			Level:  "debug",
			Format: "text",
		},
		Server: &ServerConfig{
			Port: 8080,
		},
		InternalServer: &ServerConfig{
			Port: 9090,
		},
		Tracing: &TracingConfig{Enabled: true, Type: TracingOtelHTTPType, OtelHTTP: &TracingOtelHTTPConfig{ServerURL: "http://localhost:4318/v1/traces"}},
		Database: &DatabaseConfig{
			Driver:        "POSTGRES",
			ConnectionURL: &CredentialConfig{Value: "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"},
		},
		LockDistributor: &LockDistributorConfig{
			HeartbeatFrequency: "1s",
			LeaseDuration:      "3s",
			TableName:          "locks",
		},
	}, res)
	assert.True(t, reloadHookCalled)
}

func Test_Load_reload_secret(t *testing.T) {
	dir, err := ioutil.TempDir("", "config-reload-secret")
	assert.NoError(t, err)

	configs := map[string]string{
		"log.yaml": `
log:
  level: error
  format: text
`,
		"database.yaml": `
database:
  connectionUrl:
    value: host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable

`,
		"tracing.yaml": `
tracing:
  enabled: true
  otelHttp:
    serverUrl: http://localhost:4318/v1/traces
`,
		"auth.yaml": `
oidcAuthentication:
  clientID: client-with-secret
  state: my-secret-state-key
  issuerUrl: http://localhost:8088/auth/realms/integration
  redirectUrl: http://localhost:8080/ # /auth/oidc/callback will be added
  emailVerified: true
  clientSecret:
    path: ` + os.TempDir() + `/secret1
`,
	}

	defer os.RemoveAll(dir) // clean up
	for k, v := range configs {
		tmpfn := filepath.Join(dir, k)
		err = ioutil.WriteFile(tmpfn, []byte(v), 0666)
		assert.NoError(t, err)
	}

	secretFiles := map[string]string{
		os.TempDir() + "/secret1": "VALUE1",
	}
	// Create secret files
	for k, v := range secretFiles {
		dirToCr := filepath.Dir(k)
		err = os.MkdirAll(dirToCr, 0666)
		assert.NoError(t, err)
		err = ioutil.WriteFile(k, []byte(v), 0666)
		assert.NoError(t, err)
		defer os.Remove(k)
	}

	ctx := &managerimpl{
		logger: log.NewLogger(),
	}

	// Initialize
	err = ctx.InitializeOnce()
	assert.NoError(t, err)

	reloadHookCalled := false
	ctx.AddOnChangeHook(&HookDefinition{
		Hook: func() error {
			reloadHookCalled = true

			return nil
		},
	})

	// Load config
	err = ctx.Load(dir)
	assert.NoError(t, err)
	// Get configuration
	res := ctx.GetConfig()

	assert.Equal(t, &Config{
		Log: &LogConfig{
			Level:  "error",
			Format: "text",
		},
		Server: &ServerConfig{
			Port: 8080,
		},
		InternalServer: &ServerConfig{
			Port: 9090,
		},
		Tracing: &TracingConfig{Enabled: true, Type: TracingOtelHTTPType, OtelHTTP: &TracingOtelHTTPConfig{ServerURL: "http://localhost:4318/v1/traces"}},
		Database: &DatabaseConfig{
			Driver:        "POSTGRES",
			ConnectionURL: &CredentialConfig{Value: "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"},
		},
		LockDistributor: &LockDistributorConfig{
			HeartbeatFrequency: "1s",
			LeaseDuration:      "3s",
			TableName:          "locks",
		},
		OIDCAuthentication: &OIDCAuthConfig{
			ClientID: "client-with-secret",
			ClientSecret: &CredentialConfig{
				Path:  os.TempDir() + "/secret1",
				Value: "VALUE1",
			},
			CookieName:    "oidc",
			State:         "my-secret-state-key",
			IssuerURL:     "http://localhost:8088/auth/realms/integration",
			RedirectURL:   "http://localhost:8080/",
			EmailVerified: true,
			Scopes:        []string{"openid", "email", "profile"},
		},
	}, res)

	secretFiles = map[string]string{
		os.TempDir() + "/secret1": "SECRET1",
	}
	// Create secret files
	for k, v := range secretFiles {
		dirToCr := filepath.Dir(k)
		err = os.MkdirAll(dirToCr, 0666)
		assert.NoError(t, err)
		err = ioutil.WriteFile(k, []byte(v), 0666)
		assert.NoError(t, err)
		defer os.Remove(k)
	}

	assert.Eventually(
		t,
		func() bool {
			return reloadHookCalled
		},
		reloadEventuallyWait,
		reloadEventuallyTick,
	)

	// Get configuration
	res = ctx.GetConfig()

	assert.Equal(t, &Config{
		Log: &LogConfig{
			Level:  "error",
			Format: "text",
		},
		Server: &ServerConfig{
			Port: 8080,
		},
		InternalServer: &ServerConfig{
			Port: 9090,
		},
		Tracing: &TracingConfig{Enabled: true, Type: TracingOtelHTTPType, OtelHTTP: &TracingOtelHTTPConfig{ServerURL: "http://localhost:4318/v1/traces"}},
		Database: &DatabaseConfig{
			Driver:        "POSTGRES",
			ConnectionURL: &CredentialConfig{Value: "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"},
		},
		LockDistributor: &LockDistributorConfig{
			HeartbeatFrequency: "1s",
			LeaseDuration:      "3s",
			TableName:          "locks",
		},
		OIDCAuthentication: &OIDCAuthConfig{
			ClientID: "client-with-secret",
			ClientSecret: &CredentialConfig{
				Path:  os.TempDir() + "/secret1",
				Value: "SECRET1",
			},
			CookieName:    "oidc",
			State:         "my-secret-state-key",
			IssuerURL:     "http://localhost:8088/auth/realms/integration",
			RedirectURL:   "http://localhost:8080/",
			EmailVerified: true,
			Scopes:        []string{"openid", "email", "profile"},
		},
	}, res)
	assert.True(t, reloadHookCalled)
}

func Test_Load_reload_config_with_wrong_config(t *testing.T) {
	dir, err := ioutil.TempDir("", "config-reload-wrong-config")
	assert.NoError(t, err)

	configs := map[string]string{
		"log.yaml": `
log:
  level: error
  format: text
`,
		"database.yaml": `
database:
  connectionUrl:
    value: host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable

`,
		"tracing.yaml": `
tracing:
  enabled: true
  otelHttp:
    serverUrl: http://localhost:4318/v1/traces
`,
	}

	defer os.RemoveAll(dir) // clean up
	for k, v := range configs {
		tmpfn := filepath.Join(dir, k)
		err = ioutil.WriteFile(tmpfn, []byte(v), 0666)
		assert.NoError(t, err)
	}

	secretFiles := map[string]string{
		os.TempDir() + "/secret1": "VALUE1",
	}
	// Create secret files
	for k, v := range secretFiles {
		dirToCr := filepath.Dir(k)
		err = os.MkdirAll(dirToCr, 0666)
		assert.NoError(t, err)
		err = ioutil.WriteFile(k, []byte(v), 0666)
		assert.NoError(t, err)
		defer os.Remove(k)
	}

	ctx := &managerimpl{
		logger: log.NewLogger(),
	}

	// Initialize
	err = ctx.InitializeOnce()
	assert.NoError(t, err)

	reloadHookCalled := false
	ctx.AddOnChangeHook(&HookDefinition{
		Hook: func() error {
			reloadHookCalled = true

			return nil
		},
	})

	// Load config
	err = ctx.Load(dir)
	assert.NoError(t, err)
	// Get configuration
	res := ctx.GetConfig()

	assert.Equal(t, &Config{
		Log: &LogConfig{
			Level:  "error",
			Format: "text",
		},
		Server: &ServerConfig{
			Port: 8080,
		},
		InternalServer: &ServerConfig{
			Port: 9090,
		},
		Tracing: &TracingConfig{Enabled: true, Type: TracingOtelHTTPType, OtelHTTP: &TracingOtelHTTPConfig{ServerURL: "http://localhost:4318/v1/traces"}},
		Database: &DatabaseConfig{
			Driver:        "POSTGRES",
			ConnectionURL: &CredentialConfig{Value: "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"},
		},
		LockDistributor: &LockDistributorConfig{
			HeartbeatFrequency: "1s",
			LeaseDuration:      "3s",
			TableName:          "locks",
		},
	}, res)

	configs = map[string]string{
		"log.yaml": `
configuration with error
`,
	}

	defer os.RemoveAll(dir) // clean up
	for k, v := range configs {
		tmpfn := filepath.Join(dir, k)
		err = ioutil.WriteFile(tmpfn, []byte(v), 0666)
		assert.NoError(t, err)
	}

	time.Sleep(5 * time.Second)

	// Get configuration
	res = ctx.GetConfig()

	assert.Equal(t, &Config{
		Log: &LogConfig{
			Level:  "error",
			Format: "text",
		},
		Server: &ServerConfig{
			Port: 8080,
		},
		InternalServer: &ServerConfig{
			Port: 9090,
		},
		Tracing: &TracingConfig{Enabled: true, Type: TracingOtelHTTPType, OtelHTTP: &TracingOtelHTTPConfig{ServerURL: "http://localhost:4318/v1/traces"}},
		Database: &DatabaseConfig{
			Driver:        "POSTGRES",
			ConnectionURL: &CredentialConfig{Value: "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"},
		},
		LockDistributor: &LockDistributorConfig{
			HeartbeatFrequency: "1s",
			LeaseDuration:      "3s",
			TableName:          "locks",
		},
	}, res)
	assert.False(t, reloadHookCalled)
}

func Test_Load_reload_config_map_structure(t *testing.T) {
	dir, err := ioutil.TempDir("", "config-reload-map-structure")
	assert.NoError(t, err)

	configs := map[string]string{
		"log.yaml": `
log:
  level: error
`,
		"database.yaml": `
database:
  connectionUrl:
    value: host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable

`,
		"tracing.yaml": `
tracing:
  enabled: true
  otelHttp:
    serverUrl: http://localhost:4318/v1/traces
`,
		"opa1.yaml": `
opaServerAuthorization:
  url: http://fake.com
  tags:
    t1: v1
`,
		"opa2.yaml": `
opaServerAuthorization:
  tags:
    t2: v2
`,
	}

	defer os.RemoveAll(dir) // clean up
	for k, v := range configs {
		tmpfn := filepath.Join(dir, k)
		err = ioutil.WriteFile(tmpfn, []byte(v), 0666)
		assert.NoError(t, err)
	}

	secretFiles := map[string]string{
		os.TempDir() + "/secret1": "VALUE1",
	}
	// Create secret files
	for k, v := range secretFiles {
		dirToCr := filepath.Dir(k)
		err = os.MkdirAll(dirToCr, 0666)
		assert.NoError(t, err)
		err = ioutil.WriteFile(k, []byte(v), 0666)
		assert.NoError(t, err)
		defer os.Remove(k)
	}

	ctx := &managerimpl{
		logger: log.NewLogger(),
	}

	// Initialize
	err = ctx.InitializeOnce()
	assert.NoError(t, err)

	reloadHookCalled := false
	ctx.AddOnChangeHook(&HookDefinition{
		Hook: func() error {
			reloadHookCalled = true

			return nil
		},
	})

	// Load config
	err = ctx.Load(dir)
	assert.NoError(t, err)
	// Get configuration
	res := ctx.GetConfig()

	assert.Equal(t, &Config{
		Log: &LogConfig{
			Level:  "error",
			Format: "json",
		},
		Server: &ServerConfig{
			Port: 8080,
		},
		InternalServer: &ServerConfig{
			Port: 9090,
		},
		Tracing: &TracingConfig{Enabled: true, Type: TracingOtelHTTPType, OtelHTTP: &TracingOtelHTTPConfig{ServerURL: "http://localhost:4318/v1/traces"}},
		Database: &DatabaseConfig{
			Driver:        "POSTGRES",
			ConnectionURL: &CredentialConfig{Value: "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"},
		},
		LockDistributor: &LockDistributorConfig{
			HeartbeatFrequency: "1s",
			LeaseDuration:      "3s",
			TableName:          "locks",
		},
		OPAServerAuthorization: &OPAServerAuthorization{
			URL: "http://fake.com",
			Tags: map[string]string{
				"t1": "v1",
				"t2": "v2",
			},
		},
	}, res)

	configs = map[string]string{
		"opa2.yaml": `
opaServerAuthorization:
  tags:
    t3: v3
`,
	}

	defer os.RemoveAll(dir) // clean up
	for k, v := range configs {
		tmpfn := filepath.Join(dir, k)
		err = ioutil.WriteFile(tmpfn, []byte(v), 0666)
		assert.NoError(t, err)
	}

	assert.Eventually(
		t,
		func() bool {
			return reloadHookCalled
		},
		reloadEventuallyWait,
		reloadEventuallyTick,
	)

	// Get configuration
	res = ctx.GetConfig()

	assert.Equal(t, &Config{
		Log: &LogConfig{
			Level:  "error",
			Format: "json",
		},
		Server: &ServerConfig{
			Port: 8080,
		},
		InternalServer: &ServerConfig{
			Port: 9090,
		},
		Tracing: &TracingConfig{Enabled: true, Type: TracingOtelHTTPType, OtelHTTP: &TracingOtelHTTPConfig{ServerURL: "http://localhost:4318/v1/traces"}},
		Database: &DatabaseConfig{
			Driver:        "POSTGRES",
			ConnectionURL: &CredentialConfig{Value: "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"},
		},
		LockDistributor: &LockDistributorConfig{
			HeartbeatFrequency: "1s",
			LeaseDuration:      "3s",
			TableName:          "locks",
		},
		OPAServerAuthorization: &OPAServerAuthorization{
			URL: "http://fake.com",
			Tags: map[string]string{
				"t1": "v1",
				"t3": "v3",
			},
		},
	}, res)
	assert.True(t, reloadHookCalled)
}

func Test_Load_reload_config_ignore_hidden_file_and_directory(t *testing.T) {
	dir, err := ioutil.TempDir("", "config-reload-ignore")
	assert.NoError(t, err)
	err = os.MkdirAll(path.Join(dir, "dir1"), os.ModePerm)
	assert.NoError(t, err)

	configs := map[string]string{
		"..log.yaml": `
log:
  level: error
`,
		".log2.yaml": `
log:
  format: fake
`,
		"dir1/log2.yaml": `
server:
  port: 8181
`,
		"log.yaml": `
log:
  format: humanfriendly
  level: debug
`,
		"database.yaml": `
database:
  connectionUrl:
    value: host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable

`,
	}

	defer os.RemoveAll(dir) // clean up
	for k, v := range configs {
		tmpfn := filepath.Join(dir, k)
		err = ioutil.WriteFile(tmpfn, []byte(v), 0666)
		assert.NoError(t, err)
	}

	secretFiles := map[string]string{
		os.TempDir() + "/secret1": "VALUE1",
	}
	// Create secret files
	for k, v := range secretFiles {
		dirToCr := filepath.Dir(k)
		err = os.MkdirAll(dirToCr, 0666)
		assert.NoError(t, err)
		err = ioutil.WriteFile(k, []byte(v), 0666)
		assert.NoError(t, err)
		defer os.Remove(k)
	}

	ctx := &managerimpl{
		logger: log.NewLogger(),
	}

	// Initialize
	err = ctx.InitializeOnce()
	assert.NoError(t, err)

	// Load config
	err = ctx.Load(dir)
	assert.NoError(t, err)
	// Get configuration
	res := ctx.GetConfig()

	assert.Equal(t, &Config{
		Log: &LogConfig{
			Level:  "debug",
			Format: "humanfriendly",
		},
		Server: &ServerConfig{
			Port: 8080,
		},
		InternalServer: &ServerConfig{
			Port: 9090,
		},
		Tracing: &TracingConfig{Enabled: false, Type: TracingOtelHTTPType},
		Database: &DatabaseConfig{
			Driver:        "POSTGRES",
			ConnectionURL: &CredentialConfig{Value: "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"},
		},
		LockDistributor: &LockDistributorConfig{
			HeartbeatFrequency: "1s",
			LeaseDuration:      "3s",
			TableName:          "locks",
		},
	}, res)
}

func Test_Load_reload_config_map_structure_with_one_error(t *testing.T) {
	dir, err := ioutil.TempDir("", "config-reload-map-structure-with-one-error")
	assert.NoError(t, err)

	configs := map[string]string{
		"log.yaml": `
log:
  level: error
`,
		"database.yaml": `
database:
  connectionUrl:
    value: host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable

`,
		"tracing.yaml": `
tracing:
  enabled: true
  otelHttp:
    serverUrl: http://localhost:4318/v1/traces
`,
		"opa1.yaml": `
opaServerAuthorization:
  url: http://fake.com
  tags:
    t1: v1
`,
		"opa2.yaml": `
opaServerAuthorization:
  tags:
    t2: v2
`,
	}

	defer os.RemoveAll(dir) // clean up
	for k, v := range configs {
		tmpfn := filepath.Join(dir, k)
		err = ioutil.WriteFile(tmpfn, []byte(v), 0666)
		assert.NoError(t, err)
	}

	secretFiles := map[string]string{
		os.TempDir() + "/secret1": "VALUE1",
	}
	// Create secret files
	for k, v := range secretFiles {
		dirToCr := filepath.Dir(k)
		err = os.MkdirAll(dirToCr, 0666)
		assert.NoError(t, err)
		err = ioutil.WriteFile(k, []byte(v), 0666)
		assert.NoError(t, err)
		defer os.Remove(k)
	}

	ctx := &managerimpl{
		logger: log.NewLogger(),
	}

	// Initialize
	err = ctx.InitializeOnce()
	assert.NoError(t, err)

	reloadHookCalledCount := 0
	ctx.AddOnChangeHook(&HookDefinition{
		Hook: func() error {
			if reloadHookCalledCount == 0 {
				reloadHookCalledCount += 1
				return errors.New("fake error")
			}

			return nil
		},
	})

	// Load config
	err = ctx.Load(dir)
	assert.NoError(t, err)
	// Get configuration
	res := ctx.GetConfig()

	assert.Equal(t, &Config{
		Log: &LogConfig{
			Level:  "error",
			Format: "json",
		},
		Server: &ServerConfig{
			Port: 8080,
		},
		InternalServer: &ServerConfig{
			Port: 9090,
		},
		Tracing: &TracingConfig{Enabled: true, Type: TracingOtelHTTPType, OtelHTTP: &TracingOtelHTTPConfig{ServerURL: "http://localhost:4318/v1/traces"}},
		Database: &DatabaseConfig{
			Driver:        "POSTGRES",
			ConnectionURL: &CredentialConfig{Value: "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"},
		},
		LockDistributor: &LockDistributorConfig{
			HeartbeatFrequency: "1s",
			LeaseDuration:      "3s",
			TableName:          "locks",
		},
		OPAServerAuthorization: &OPAServerAuthorization{
			URL: "http://fake.com",
			Tags: map[string]string{
				"t1": "v1",
				"t2": "v2",
			},
		},
	}, res)

	configs = map[string]string{
		"opa2.yaml": `
opaServerAuthorization:
  tags:
    t3: v3
`,
	}

	defer os.RemoveAll(dir) // clean up
	for k, v := range configs {
		tmpfn := filepath.Join(dir, k)
		err = ioutil.WriteFile(tmpfn, []byte(v), 0666)
		assert.NoError(t, err)
	}

	assert.Eventually(
		t,
		func() bool {
			return reloadHookCalledCount > 0
		},
		reloadEventuallyWait,
		reloadEventuallyTick,
	)

	// Get configuration
	res = ctx.GetConfig()

	assert.Equal(t, &Config{
		Log: &LogConfig{
			Level:  "error",
			Format: "json",
		},
		Server: &ServerConfig{
			Port: 8080,
		},
		InternalServer: &ServerConfig{
			Port: 9090,
		},
		Tracing: &TracingConfig{Enabled: true, Type: TracingOtelHTTPType, OtelHTTP: &TracingOtelHTTPConfig{ServerURL: "http://localhost:4318/v1/traces"}},
		Database: &DatabaseConfig{
			Driver:        "POSTGRES",
			ConnectionURL: &CredentialConfig{Value: "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"},
		},
		LockDistributor: &LockDistributorConfig{
			HeartbeatFrequency: "1s",
			LeaseDuration:      "3s",
			TableName:          "locks",
		},
		OPAServerAuthorization: &OPAServerAuthorization{
			URL: "http://fake.com",
			Tags: map[string]string{
				"t1": "v1",
				"t3": "v3",
			},
		},
	}, res)
	assert.Equal(t, 1, reloadHookCalledCount)
}
