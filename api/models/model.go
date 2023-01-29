package models

import "time"

type Request struct{
	URL							string 						`json:"url" binding:"required,url"`
	CustomShort			string						`json:"short"`
	Expiry					time.Duration			`json:"expiry"`
}

type Response struct{
	URL									string					`json:"url" binding:"required,url"`
	CustomShort					string					`json:"short"`
	Expiry							time.Duration		`json:"expiry"`
	XRateRemaining			int							`json:"rateLimit"`
	XRateLimitReset			time.Duration		`json:"rateLimitReset"`
}

type ResolveRequest struct{
	URL							string 						`uri:"url" binding:"required,min=6,max=6"`
}

