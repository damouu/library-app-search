package health

import "net/http"

type DependencyChecker interface {
	Check() error
}

type Handler struct {
	checker DependencyChecker
}

func NewHandler(checker DependencyChecker) *Handler {
	return &Handler{
		checker: checker,
	}
}

func LiveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.checker.Check(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}
