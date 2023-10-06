package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-semantic-release/plugin-registry/pkg/client"
	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
)

var version = "dev"

var defaultPluginRegistryURLs = []string{
	client.DefaultStagingEndpoint,
	client.DefaultProductionEndpoint,
}

type PluginRegistryUpdate struct {
	pluginName                     string
	pluginRegistryAdminAccessToken string
	log                            *log.Logger
}

func (p *PluginRegistryUpdate) Init(m map[string]string) error {
	p.pluginName = os.Getenv("PLUGIN_NAME")
	if m["plugin_name"] != "" {
		p.pluginName = m["plugin_name"]
	}
	p.pluginRegistryAdminAccessToken = os.Getenv("PLUGIN_REGISTRY_ADMIN_ACCESS_TOKEN")
	if m["plugin_registry_admin_access_token"] != "" {
		p.pluginRegistryAdminAccessToken = m["plugin_registry_admin_access_token"]
	}
	if p.pluginRegistryAdminAccessToken == "" {
		return fmt.Errorf("plugin registry admin access token is not configured")
	}
	return nil
}

func (p *PluginRegistryUpdate) Name() string {
	return "plugin-registry-update"
}

func (p *PluginRegistryUpdate) Version() string {
	return version
}

func (p *PluginRegistryUpdate) Success(config *hooks.SuccessHookConfig) error {
	pluginName := p.pluginName
	if pluginName == "" {
		pluginName = config.RepoInfo.Repo
	}
	if pluginName == "" {
		return fmt.Errorf("plugin name is not configured")
	}
	pluginVersion := config.NewRelease.Version
	p.log.Printf("triggering plugin registry update for %s@%s", pluginName, pluginVersion)

	for _, url := range defaultPluginRegistryURLs {
		p.log.Printf("updating plugin registry: %s", url)
		c := client.New(url)
		err := c.UpdatePluginRelease(context.Background(), p.pluginRegistryAdminAccessToken, pluginName, pluginVersion)
		if err != nil {
			p.log.Printf("error: failed to update plugin registry %s: %v", url, err)
		}
	}

	return nil
}

func (p *PluginRegistryUpdate) NoRelease(_ *hooks.NoReleaseConfig) error {
	return nil
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Hooks: func() hooks.Hooks {
			return &PluginRegistryUpdate{
				log: log.New(os.Stderr, "", 0),
			}
		},
	})
}
