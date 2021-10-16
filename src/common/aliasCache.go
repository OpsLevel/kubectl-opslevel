package common

import (
	"sync"

	"github.com/opslevel/opslevel-go"
	"github.com/rs/zerolog/log"
)

type AliasCacher struct {
	mutex      sync.Mutex
	Tiers      map[string]opslevel.Tier
	Lifecycles map[string]opslevel.Lifecycle
	Teams      map[string]opslevel.Team
}

func (c *AliasCacher) TryGetTier(alias string) (*opslevel.Tier, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Tiers[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *AliasCacher) TryGetLifecycle(alias string) (*opslevel.Lifecycle, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Lifecycles[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *AliasCacher) TryGetTeam(alias string) (*opslevel.Team, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.Teams[alias]; ok {
		return &v, ok
	}
	return nil, false
}

func (c *AliasCacher) doCacheTiers(client *opslevel.Client) {
	log.Info().Msg("Caching 'Tiers' lookup table from OpsLevel API ...")

	data, dataErr := client.ListTiers()
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to retrive tiers from OpsLevel API - Unable to assign field 'Tier' to services. REASON: %s", dataErr.Error())
	}
	for _, item := range data {
		c.Tiers[string(item.Alias)] = item
	}
}

func (c *AliasCacher) doCacheLifecycles(client *opslevel.Client) {
	log.Info().Msg("Caching 'Lifecycles' lookup table from OpsLevel API ...")

	data, dataErr := client.ListLifecycles()
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to retrive lifecycles from OpsLevel API - Unable to assign field 'Lifecycle' to services. REASON: %s", dataErr.Error())
	}
	for _, item := range data {
		c.Lifecycles[string(item.Alias)] = item
	}
}

func (c *AliasCacher) doCacheTeams(client *opslevel.Client) {
	log.Info().Msg("Caching 'Teams' lookup table from OpsLevel API ...")

	data, dataErr := client.ListTeams()
	if dataErr != nil {
		log.Warn().Msgf("===> Failed to retrive teams from OpsLevel API - Unable to assign field 'Owner' to services. REASON: %s", dataErr.Error())
	}

	for _, item := range data {
		c.Teams[string(item.Alias)] = item
	}
}

func (c *AliasCacher) CacheTiers(client *opslevel.Client) {
	c.mutex.Lock()
	c.doCacheTiers(client)
	c.mutex.Unlock()
}

func (c *AliasCacher) CacheLifecycles(client *opslevel.Client) {
	c.mutex.Lock()
	c.doCacheLifecycles(client)
	c.mutex.Unlock()
}

func (c *AliasCacher) CacheTeams(client *opslevel.Client) {
	c.mutex.Lock()
	c.doCacheTeams(client)
	c.mutex.Unlock()
}

func (c *AliasCacher) CacheAll(client *opslevel.Client) {
	c.mutex.Lock()
	c.doCacheTiers(client)
	c.doCacheLifecycles(client)
	c.doCacheTeams(client)
	c.mutex.Unlock()
}

var aliasCacher *AliasCacher

func GetOrCreateAliasCache() *AliasCacher {
	if aliasCacher != nil {
		return aliasCacher
	}
	aliasCacher := &AliasCacher{
		mutex:      sync.Mutex{},
		Tiers:      make(map[string]opslevel.Tier),
		Lifecycles: make(map[string]opslevel.Lifecycle),
		Teams:      make(map[string]opslevel.Team),
	}
	return aliasCacher
}
