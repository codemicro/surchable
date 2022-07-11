package endpoints

import (
	"crypto/sha1"
	"net/url"
	"time"

	"github.com/codemicro/surchable/internal/search"

	db "github.com/codemicro/surchable/internal/libdb"
	"github.com/codemicro/surchable/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

const (
	preflightLoad = "LOAD"
	preflightSkip = "SKIP"
)

func (e *Endpoints) Post_RequestPreflightCheck(ctx *fiber.Ctx) error {
	type schema struct {
		URL string `json:"url" validate:"required,url"`
	}
	type response struct {
		Permission string `json:"permission"`
	}

	inputData := new(schema)
	if err := util.ParseAndValidateJSONBody(ctx, inputData); err != nil {
		return err
	}

	pageLoad, err := e.db.QueryPageLoadsByURL(inputData.URL)
	if err != nil {
		if errors.Is(err, db.ErrNoMatchingPageLoad) {
			goto respondLoad
		}
		return err
	}

	if pageLoad.NotLoadBefore.After(time.Now()) {
		return ctx.JSON(&response{preflightSkip})
	}

respondLoad:
	return ctx.JSON(&response{preflightLoad})
}

func (e *Endpoints) Post_DigestPageLoad(ctx *fiber.Ctx) error {
	type schema struct {
		URL           string    `json:"url" validate:"required,url"`
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		Content       string    `json:"content"`
		HTML          string    `json:"html" validate:"required"`
		NotLoadBefore int       `json:"notLoadBefore" validate:"gte=0"`
		OutboundLinks []string  `json:"outboundLinks" validate:"dive,url"`
		LoadedAt      time.Time `json:"loadedAt" validate:"required"`
	}

	crawlerID, err := getCrawlerIDHeader(ctx)
	if err != nil {
		return err
	}

	inputData := new(schema)
	if err := util.ParseAndValidateJSONBody(ctx, inputData); err != nil {
		return err
	}

	if inputData.NotLoadBefore == 0 {
		inputData.NotLoadBefore = 60
	}

	// putting this here will ensure that a worker is currently doing something
	// if not, an internal server error will be thrown
	// TODO: make this return a proper error instead of HTTP 500
	if err := e.db.UpdateTimeForJobByWorkerID(crawlerID, time.Now()); err != nil {
		return errors.WithStack(err)
	}

	pageLoadID, err := e.db.UpsertPageLoad(&db.PageLoad{
		URL:      inputData.URL,
		LoadedAt: inputData.LoadedAt,
		NotLoadBefore: util.Ptr(
			inputData.LoadedAt.Add(time.Duration(inputData.NotLoadBefore) * time.Minute),
		),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	pageID, err := e.db.UpsertPageInformation(&db.PageInformation{
		LoadID:                  pageLoadID,
		PageTitle:               util.PtrNilIfDefault(inputData.Title),
		PageMetaDescriptionText: util.PtrNilIfDefault(inputData.Description),
		PageContentText:         util.PtrNilIfDefault(inputData.Content),
		PageRawHTML:             inputData.HTML,
		RawHTMLSHA1:             sha1.Sum([]byte(inputData.HTML)),
		OutboundLinks:           inputData.OutboundLinks,
	})

	if err != nil {
		return errors.WithStack(err)
	}

	for _, outboundLink := range inputData.OutboundLinks {
		u, _ := url.Parse(outboundLink) // URL already validated
		_, err := e.db.AddDomainToQueue(u.Host, u.EscapedPath(), db.DefaultDomainQueuePriority)
		if err != nil && !(errors.Is(err, db.ErrDomainAlreadyQueued) || errors.Is(err, db.ErrDomainInBlocklist)) {
			return errors.WithStack(err)
		}
	}

	tokenMap := make(db.TokenMap)
	tokenMap.Add(search.TokeniseString(inputData.Content), db.IndexClassPageBody)
	tokenMap.Add(search.TokeniseString(inputData.Description), db.IndexClassPageDescription)
	tokenMap.Add(search.TokeniseString(inputData.Title), db.IndexClassPageTitle)

	if err := e.db.SearchIndexUpsert(&db.TokenSet{
		PageID: pageID,
		Tokens: tokenMap,
	}); err != nil {
		return errors.WithStack(err)
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}
