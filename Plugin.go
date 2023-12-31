package tags

import (
	"bytes"

	"github.com/go-bolo/bolo"
	"github.com/gookit/event"
	"github.com/sirupsen/logrus"
)

type Plugin struct {
	Name string

	VocabularyController *VocabularyController
	TermController       *TermController

	RenderRelatedRecord func(mt *ModelstermsModel, ctx *bolo.RequestContext) (bytes.Buffer, error)
}

func (r *Plugin) GetName() string {
	return r.Name
}

func (r *Plugin) Init(app bolo.App) error {
	logrus.Debug(r.GetName() + " Init")

	r.VocabularyController = NewVocabularyController(&VocabularyControllerCfg{App: app})
	r.TermController = NewTermController(&TermControllerCfg{App: app})

	app.GetEvents().On("bindRoutes", event.ListenerFunc(func(e event.Event) error {
		return r.BindRoutes(app)
	}), event.Normal)

	return nil
}

func (r *Plugin) BindRoutes(app bolo.App) error {
	logrus.Debug(r.GetName() + " On BindRoutes")

	vocabularyCTL := r.VocabularyController
	termCTL := r.TermController

	mainRouter := app.GetRouterGroup("main")
	mainRouter.GET("api/v1/term-texts", termCTL.TermTexts)

	routerApi := app.SetRouterGroup("vocabulary-api", "/api/vocabulary")

	app.SetResource("vocabulary", vocabularyCTL, routerApi)

	routerVocTermApi := app.SetRouterGroup("vocabulary-term-api", "/api/vocabulary/:vocabulary/term")
	app.SetResource("vocabulary-term", termCTL, routerVocTermApi)

	mainRouter.GET("vocabulary/:vocabulary/term/:id", termCTL.FindOnePageHandler)

	return nil
}

func (r *Plugin) SetTemplateFuncMap(app bolo.App) error {
	return nil
}

func (r *Plugin) GetMigrations() []*bolo.Migration {
	return []*bolo.Migration{}
}

type PluginCfgs struct {
	RenderRelatedRecord func(mt *ModelstermsModel, ctx *bolo.RequestContext) (bytes.Buffer, error)
}

func NewPlugin(cfg *PluginCfgs) *Plugin {
	p := Plugin{Name: "taxonomy", RenderRelatedRecord: cfg.RenderRelatedRecord}

	if p.RenderRelatedRecord == nil {
		p.RenderRelatedRecord = func(mt *ModelstermsModel, ctx *bolo.RequestContext) (bytes.Buffer, error) {
			return bytes.Buffer{}, nil
		}
	}
	return &p
}
