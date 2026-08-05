package http

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WriteGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch st.Code() {

	case codes.InvalidArgument:
		http.Error(w, st.Message(), http.StatusBadRequest)

	case codes.NotFound:
		http.Error(w, st.Message(), http.StatusNotFound)

	case codes.AlreadyExists:
		http.Error(w, st.Message(), http.StatusConflict)

	case codes.PermissionDenied:
		http.Error(w, st.Message(), http.StatusForbidden)

	case codes.Unauthenticated:
		http.Error(w, st.Message(), http.StatusUnauthorized)

	case codes.FailedPrecondition:
		http.Error(w, st.Message(), http.StatusPreconditionFailed)

	default:
		http.Error(w, st.Message(), http.StatusInternalServerError)
	}
}
