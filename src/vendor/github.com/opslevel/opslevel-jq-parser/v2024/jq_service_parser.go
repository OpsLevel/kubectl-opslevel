package opslevel_jq_parser

type JQServiceParser struct {
	name         *JQFieldParser
	description  *JQFieldParser
	owner        *JQFieldParser
	lifecycle    *JQFieldParser
	tier         *JQFieldParser
	product      *JQFieldParser
	language     *JQFieldParser
	framework    *JQFieldParser
	system       *JQFieldParser
	aliases      *JQArrayParser
	tags         *JQTagsParser
	tools        *JQToolsParser
	repositories *JQRepositoryParser
}

func NewJQServiceParser(cfg ServiceRegistrationConfig) *JQServiceParser {
	return &JQServiceParser{
		name:         NewJQFieldParser(cfg.Name),
		description:  NewJQFieldParser(cfg.Description),
		owner:        NewJQFieldParser(cfg.Owner),
		lifecycle:    NewJQFieldParser(cfg.Lifecycle),
		tier:         NewJQFieldParser(cfg.Tier),
		product:      NewJQFieldParser(cfg.Product),
		language:     NewJQFieldParser(cfg.Language),
		framework:    NewJQFieldParser(cfg.Framework),
		system:       NewJQFieldParser(cfg.System),
		aliases:      NewJQArrayParser(cfg.Aliases),
		tags:         NewJQTagsParser(cfg.Tags),
		tools:        NewJQToolsParser(cfg.Tools),
		repositories: NewJQRepositoryParser(cfg.Repositories),
	}
}

func (p *JQServiceParser) Run(json string) (*ServiceRegistration, error) {
	name, err := p.name.Run(json)
	if err != nil {
		return nil, err
	}
	description, err := p.description.Run(json)
	if err != nil {
		return nil, err
	}
	Owner, err := p.owner.Run(json)
	if err != nil {
		return nil, err
	}
	Lifecycle, err := p.lifecycle.Run(json)
	if err != nil {
		return nil, err
	}
	Tier, err := p.tier.Run(json)
	if err != nil {
		return nil, err
	}
	Product, err := p.product.Run(json)
	if err != nil {
		return nil, err
	}
	Language, err := p.language.Run(json)
	if err != nil {
		return nil, err
	}
	Framework, err := p.framework.Run(json)
	if err != nil {
		return nil, err
	}
	System, err := p.system.Run(json)
	if err != nil {
		return nil, err
	}
	Aliases, err := p.aliases.Run(json)
	if err != nil {
		return nil, err
	}
	TagCreates, TagAssigns, err := p.tags.Run(json)
	if err != nil {
		return nil, err
	}
	Tools, err := p.tools.Run(json)
	if err != nil {
		return nil, err
	}
	Repositories, err := p.repositories.Run(json)
	if err != nil {
		return nil, err
	}
	return &ServiceRegistration{
		Name:         name,
		Description:  description,
		Owner:        Owner,
		Lifecycle:    Lifecycle,
		Tier:         Tier,
		Product:      Product,
		Language:     Language,
		Framework:    Framework,
		System:       System,
		Aliases:      Aliases,
		TagCreates:   TagCreates,
		TagAssigns:   TagAssigns,
		Tools:        Tools,
		Repositories: Repositories,
	}, nil
}
