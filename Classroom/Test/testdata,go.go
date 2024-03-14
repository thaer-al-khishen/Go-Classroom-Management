package Test

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"webapptrials/Classroom/Shared"
)

func TestData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	Shared.SendApiResponse(w, http.StatusOK, "Test api functional!", "Success!", "")
}
