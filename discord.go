// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lunny/tango"
)

type Options struct {
	WebhookID    string
	WebhookToken string
	WebhookURL   string
	Color        string
	Avatar       string
	Source       string
}

func (opt *Options) IsValid() bool {
	return opt.WebhookID != "" && opt.WebhookToken != ""
}

func prepareOptions(opts []Options) Options {
	var opt Options
	if len(opts) > 0 {
		opt = opts[0]
	}
	if !opt.IsValid() {
		return opt
	}
	opt.WebhookURL = fmt.Sprintf("https://discordapp.com/api/webhooks/%s/%s", opt.WebhookID, opt.WebhookToken)
	if opt.Source == "" {
		opt.Source = "Tango"
	}
	return opt
}

type (
	// EmbedFooterObject for Embed Footer Structure.
	EmbedFooterObject struct {
		Text string `json:"text"`
	}

	// EmbedAuthorObject for Embed Author Structure
	EmbedAuthorObject struct {
		Name    string `json:"name"`
		URL     string `json:"url"`
		IconURL string `json:"icon_url"`
	}

	// EmbedFieldObject for Embed Field Structure
	EmbedFieldObject struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	// EmbedObject is for Embed Structure
	EmbedObject struct {
		Title       string             `json:"title"`
		Description string             `json:"description"`
		URL         string             `json:"url"`
		Color       int                `json:"color"`
		Footer      EmbedFooterObject  `json:"footer"`
		Author      EmbedAuthorObject  `json:"author"`
		Fields      []EmbedFieldObject `json:"fields"`
	}
	PayLoad struct {
		Wait      bool          `json:"wait"`
		Content   string        `json:"content"`
		Username  string        `json:"username"`
		AvatarURL string        `json:"avatar_url"`
		TTS       bool          `json:"tts"`
		Embeds    []EmbedObject `json:"embeds"`
	}
)

// Discord send message to chatroom
func Discord(opts ...Options) tango.HandlerFunc {
	opt := prepareOptions(opts)
	if !opt.IsValid() {
		return func(ctx *tango.Context) {
			ctx.Next()
		}
	}

	return func(ctx *tango.Context) {
		ctx.Next()

		if ctx.ResponseWriter.Status()/100 == 5 {
			body, err := ctx.Body()
			if err != nil {
				ctx.Error(err)
				return
			}

			p := ctx.Req().URL.Path
			if len(ctx.Req().URL.RawQuery) > 0 {
				p = p + "?" + ctx.Req().URL.RawQuery
			}

			var payload = PayLoad{
				Embeds: []EmbedObject{
					{
						Title: fmt.Sprintf("%s %d %v", p, ctx.ResponseWriter.Status(), ctx.Result),
						Description: fmt.Sprintf("Request %s From %s failed: %s",
							ctx.Req().URL.String(), ctx.IP(), string(body)),
						URL:   ctx.Req().URL.String(),
						Color: 16724530,
						Author: EmbedAuthorObject{
							Name:    opt.Source,
							IconURL: ctx.Req().URL.Scheme + "://" + ctx.Req().URL.Host + "/favicon.ico",
						},
					},
				},
			}

			b := new(bytes.Buffer)
			if err := json.NewEncoder(b).Encode(payload); err != nil {
				ctx.Error(err)
				return
			}

			_, err = http.Post(opt.WebhookURL, "application/json; charset=utf-8", b)
			if err != nil {
				ctx.Error(err)
			}
		}
	}
}
