package remote

import (
	"github.com/cozy/cozy-stack/pkg/remote"
	"github.com/cozy/cozy-stack/web/jsonapi"
	"github.com/cozy/cozy-stack/web/middlewares"
	"github.com/cozy/cozy-stack/web/permissions"
	"github.com/cozy/echo"
)

func remoteGet(c echo.Context) error {
	doctype := c.Param("doctype")
	if err := permissions.AllowWholeType(c, permissions.GET, doctype); err != nil {
		return wrapRemoteErr(err)
	}
	remote, err := remote.Find(doctype)
	if err != nil {
		return wrapRemoteErr(err)
	}
	if remote.Verb != "GET" {
		return jsonapi.MethodNotAllowed("GET")
	}
	instance := middlewares.GetInstance(c)
	err = remote.ProxyTo(doctype, instance, c.Response(), c.Request())
	if err != nil {
		return wrapRemoteErr(err)
	}
	return nil
}

func remotePost(c echo.Context) error {
	doctype := c.Param("doctype")
	if err := permissions.AllowWholeType(c, permissions.POST, doctype); err != nil {
		return wrapRemoteErr(err)
	}
	remote, err := remote.Find(doctype)
	if err != nil {
		return wrapRemoteErr(err)
	}
	if remote.Verb != "POST" {
		return jsonapi.MethodNotAllowed("POST")
	}
	instance := middlewares.GetInstance(c)
	err = remote.ProxyTo(doctype, instance, c.Response(), c.Request())
	if err != nil {
		return wrapRemoteErr(err)
	}
	return nil
}

// Routes set the routing for the remote service
func Routes(router *echo.Group) {
	router.GET("/:doctype", remoteGet)
	router.POST("/:doctype", remotePost)
	// TODO add tests
}

func wrapRemoteErr(err error) error {
	switch err {
	case remote.ErrNotFoundRemote:
		return jsonapi.NotFound(err)
	case remote.ErrInvalidRequest:
		return jsonapi.BadRequest(err)
	case remote.ErrRequestFailed:
		return jsonapi.BadGateway(err)
	case remote.ErrInvalidVariables:
		return jsonapi.BadRequest(err)
	}
	return err
}
