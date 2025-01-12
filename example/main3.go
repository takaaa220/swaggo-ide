package main

import "fmt"

// @Summary		Show an account
// @Description	"get string by ID"
// @Tags			accounts
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Account ID"
// @Success		200	{object}	model.Account
// @Failure		400	{object}	httputil.HTTPError
// @Failure		404	{object}	httputil.HTTPError
// @Failure		500	{object}	httputil.HTTPError
// @Router			/accounts/{id} [get]
func yeah2() {
	fmt.Println("Hello, World!")
}
