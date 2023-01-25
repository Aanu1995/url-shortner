package models

import "time"

type Request struct{
	URL							string 						`json:"url" validate:"required,min=4"`
	CustomShort			string						`json:"short"`
	Expiry					time.Duration			`json:"expiry"`
}

type Response struct{
	URL									string					`json:"url" validate:"required,min=4"`
	CustomShort					string					`json:"short"`
	Expiry							time.Duration		`json:"expiry"`
	XRateRemaining			int							`json:"rateLimit"`
	XRateLimitReset			time.Duration		`json:"rateLimitReset"`
}

