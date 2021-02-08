package fortunecookie

import (
	"fmt"
	"net/http"
)

func init() {

}

func Fc_handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "fortune cookie")

}
