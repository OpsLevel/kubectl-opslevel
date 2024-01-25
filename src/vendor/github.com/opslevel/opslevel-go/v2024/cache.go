package opslevel

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type Cacher struct {
	mutex        sync.Mutex
	Tiers        map[string]Tier
	Lifecycles   map[string]Lifecycle
	Teams        map[string]Team
	Categories   map[string]Category
	Levels       map[string]Level
	Filters      map[string]Filter
	Integrations map[string]Integration
	Repositories map[string]Repository
	InfraSchemas map[string]InfrastructureResourceSchema
}

func (c *Cacher) TryGetTier(alias string) (*Tier, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Tiers[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *Cacher) TryGetLifecycle(alias string) (*Lifecycle, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Lifecycles[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *Cacher) TryGetTeam(alias string) (*Team, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Teams[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *Cacher) TryGetCategory(alias string) (*Category, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Categories[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *Cacher) TryGetLevel(alias string) (*Level, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Levels[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *Cacher) TryGetFilter(alias string) (*Filter, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Filters[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *Cacher) TryGetIntegration(alias string) (*Integration, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Integrations[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *Cacher) TryGetRepository(alias string) (*Repository, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Repositories[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *Cacher) TryGetInfrastructureSchema(alias string) (*InfrastructureResourceSchema, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.InfraSchemas[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *Cacher) doCacheTiers(client *Client) {
	log.Debug().Msg("Caching 'Tier' lookup table from API ...")

	data, dataErr := client.ListTiers()
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Tier' from API - REASON: %s", dataErr.Error())
	}
	for _, item := range data {
		c.Tiers[string(item.Alias)] = item
	}
}

func (c *Cacher) doCacheLifecycles(client *Client) {
	log.Debug().Msg("Caching 'Lifecycle' lookup table from API ...")

	data, dataErr := client.ListLifecycles()
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Lifecycle' from API - REASON: %s", dataErr.Error())
	}
	for _, item := range data {
		c.Lifecycles[string(item.Alias)] = item
	}
}

func (c *Cacher) doCacheTeams(client *Client) {
	log.Debug().Msg("Caching 'Team' lookup table from API ...")

	data, dataErr := client.ListTeams(nil)
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Team' from API - REASON: %s", dataErr.Error())
	}

	for _, item := range data.Nodes {
		for _, alias := range item.Aliases {
			c.Teams[string(alias)] = item
		}
	}
}

func (c *Cacher) doCacheCategories(client *Client) {
	log.Debug().Msg("Caching 'Category' lookup table from API ...")

	data, dataErr := client.ListCategories(nil)
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Category' from API - REASON: %s", dataErr.Error())
	}

	for _, item := range data.Nodes {
		c.Categories[item.Alias()] = item
	}
}

func (c *Cacher) doCacheLevels(client *Client) {
	log.Debug().Msg("Caching 'Level' lookup table from API ...")

	data, dataErr := client.ListLevels()
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Level' from API - REASON: %s", dataErr.Error())
	}

	for _, item := range data {
		c.Levels[string(item.Alias)] = item
	}
}

func (c *Cacher) doCacheFilters(client *Client) {
	log.Debug().Msg("Caching 'Filter' lookup table from API ...")

	data, dataErr := client.ListFilters(nil)
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Filter' from API - REASON: %s", dataErr.Error())
	}

	for _, item := range data.Nodes {
		c.Filters[item.Alias()] = item
	}
}

func (c *Cacher) doCacheIntegrations(client *Client) {
	log.Debug().Msg("Caching 'Integration' lookup table from API ...")

	data, dataErr := client.ListIntegrations(nil)
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Integration' from API - REASON: %s", dataErr.Error())
	}

	for _, item := range data.Nodes {
		c.Integrations[item.Alias()] = item
	}
}

func (c *Cacher) doCacheRepositories(client *Client) {
	log.Debug().Msg("Caching 'Repository' lookup table from API ...")

	data, dataErr := client.ListRepositories(nil)
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'Repository' from API - REASON: %s", dataErr.Error())
	}

	for _, item := range data.Nodes {
		c.Repositories[item.DefaultAlias] = item
	}
}

func (c *Cacher) doCacheInfraSchemas(client *Client) {
	log.Debug().Msg("Caching 'InfrastructureSchema' lookup table from API ...")

	data, dataErr := client.ListInfrastructureSchemas(nil)
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to list all 'InfrastructureSchema' from API - REASON: %s", dataErr.Error())
	}
	for _, item := range data.Nodes {
		// log.Info().Msgf("Caching 'InfrastructureSchema' '%s' ...", item.Type)
		c.InfraSchemas[item.Type] = item
	}
}

func (c *Cacher) CacheTiers(client *Client) {
	c.mutex.Lock()
	c.doCacheTiers(client)
	c.mutex.Unlock()
}

func (c *Cacher) CacheLifecycles(client *Client) {
	c.mutex.Lock()
	c.doCacheLifecycles(client)
	c.mutex.Unlock()
}

func (c *Cacher) CacheTeams(client *Client) {
	c.mutex.Lock()
	c.doCacheTeams(client)
	c.mutex.Unlock()
}

func (c *Cacher) CacheCategories(client *Client) {
	c.mutex.Lock()
	c.doCacheCategories(client)
	c.mutex.Unlock()
}

func (c *Cacher) CacheLevels(client *Client) {
	c.mutex.Lock()
	c.doCacheLevels(client)
	c.mutex.Unlock()
}

func (c *Cacher) CacheFilters(client *Client) {
	c.mutex.Lock()
	c.doCacheFilters(client)
	c.mutex.Unlock()
}

func (c *Cacher) CacheIntegrations(client *Client) {
	c.mutex.Lock()
	c.doCacheIntegrations(client)
	c.mutex.Unlock()
}

func (c *Cacher) CacheRepositories(client *Client) {
	c.mutex.Lock()
	c.doCacheRepositories(client)
	c.mutex.Unlock()
}

func (c *Cacher) CacheInfraSchemas(client *Client) {
	c.mutex.Lock()
	c.doCacheInfraSchemas(client)
	c.mutex.Unlock()
}

func (c *Cacher) CacheAll(client *Client) {
	c.mutex.Lock()
	c.doCacheTiers(client)
	c.doCacheLifecycles(client)
	c.doCacheTeams(client)
	c.doCacheCategories(client)
	c.doCacheLevels(client)
	c.doCacheFilters(client)
	c.doCacheIntegrations(client)
	c.doCacheRepositories(client)
	c.doCacheInfraSchemas(client)
	c.mutex.Unlock()
}

var Cache = &Cacher{
	mutex:        sync.Mutex{},
	Tiers:        make(map[string]Tier),
	Lifecycles:   make(map[string]Lifecycle),
	Teams:        make(map[string]Team),
	Categories:   make(map[string]Category),
	Levels:       make(map[string]Level),
	Filters:      make(map[string]Filter),
	Integrations: make(map[string]Integration),
	Repositories: make(map[string]Repository),
	InfraSchemas: make(map[string]InfrastructureResourceSchema),
}
