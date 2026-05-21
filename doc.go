// Package biliapi provides a typed Go client for Bilibili Web REST APIs.
//
// Create a client with NewClient, then call one of the service fields such as
// Client.Video, Client.Search, Client.User, Client.Comment, Client.Danmaku,
// Client.Fav, Client.History, Client.Login, Client.ClientInfo, or Client.Misc.
//
// WBI-signed endpoints are signed automatically. Endpoints that require login
// use Credentials from WithCredentials or SetCredentials. Mutating endpoints
// also require the bili_jct CSRF token and return ErrMissingCSRF before sending
// a request when it is absent.
package biliapi
