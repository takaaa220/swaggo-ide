package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}

// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  model.Account
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /accounts/{id} [get]
func hello() {
	if true {
		fmt.Println("Hello, World!")

		switch {
		case true:
			fmt.Println("Hello, World!")
		case false:
			fmt.Println("Hello, World!")
		}
	}
}
