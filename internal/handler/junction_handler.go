package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/srq/signalflux/internal/domain"
	"github.com/srq/signalflux/internal/handler/dto"
)

// internal/handler/junction_handler.go — Junction HTTP Handlers
//
// LAYER   : Handler / HTTP Adapter (outermost layer)
// IMPORTS : net/http, internal/domain, internal/handler/dto
//
// PURPOSE : Handle all /junctions routes. Parse HTTP request → call service →
//           build DTO response. ZERO business logic here — all rules live in service/.
//
// ─────────────────────────────────────────────────────────────────────────────
// FLOW ORDER (implement in this order):
//
// 1. Define JunctionHandler struct
//    type JunctionHandler struct {
//        service domain.JunctionService
//    }
//
// 2. NewJunctionHandler(svc domain.JunctionService) *JunctionHandler
//    → Returns &JunctionHandler{service: svc}
//
// 3. CreateJunction — POST /junctions
//    Steps:
//      a. decodeJSON(r, &req) → on error: respondError(w, 400, "malformed JSON")
//      b. req.Validate() → on error: respondError(w, 400, err.Error())
//      c. junction, err := h.service.Create(ctx, req.Type, req.Location)
//      d. On error: mapServiceError(w, err); return
//      e. respondJSON(w, http.StatusCreated, dto.JunctionToResponse(junction))
//
// 4. ListJunctions — GET /junctions
//    Steps:
//      a. Parse query params:
//           page  = strconv.Atoi(r.URL.Query().Get("page"))  → default 1 if missing/invalid
//           limit = strconv.Atoi(r.URL.Query().Get("limit")) → default 20 if missing/invalid
//      b. junctions, total, err := h.service.List(ctx, page, limit)
//      c. On error: mapServiceError(w, err); return
//      d. respondJSON(w, 200, dto.PaginatedJunctions(junctions, total, page, limit))
//
// 5. GetJunction — GET /junctions/{id}
//    Steps:
//      a. id, err := parseUUID(chi.URLParam(r, "id"))
//         → on error: respondError(w, 400, "invalid junction id"); return
//      b. junction, err := h.service.GetByID(ctx, id)
//      c. On error: mapServiceError(w, err); return
//      d. respondJSON(w, 200, dto.JunctionToResponse(junction))
//
// 6. DeleteJunction — DELETE /junctions/{id}
//    Steps:
//      a. id, err := parseUUID(chi.URLParam(r, "id"))
//         → on error: respondError(w, 400, "invalid junction id"); return
//      b. err = h.service.Delete(ctx, id)
//      c. On error: mapServiceError(w, err); return
//      d. w.WriteHeader(http.StatusNoContent)  ← 204 No Content (no body)
//
// ─────────────────────────────────────────────────────────────────────────────
// ROUTE REGISTRATION (done in main.go, documented here for reference):
//   POST   /junctions         → handler.CreateJunction
//   GET    /junctions         → handler.ListJunctions
//   GET    /junctions/{id}    → handler.GetJunction
//   DELETE /junctions/{id}    → handler.DeleteJunction

// TODO: Define JunctionHandler struct with service field
// TODO: Implement NewJunctionHandler constructor
// TODO: Implement CreateJunction — decode body, validate, call service.Create, respond 201
// TODO: Implement ListJunctions — parse page/limit query params, call service.List, respond 200 paginated
// TODO: Implement GetJunction — parse UUID path param, call service.GetByID, respond 200
// TODO: Implement DeleteJunction — parse UUID path param, call service.Delete, respond 204

type JunctionHandler struct {
	service domain.JunctionService
}

func NewJunctionHandler(js domain.JunctionService) *JunctionHandler {
	return &JunctionHandler{service: js}
}

func (jh JunctionHandler) CreateJunction(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.CreateJunctionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "malformed JSON"}`, http.StatusBadRequest)
		return
	}

	j := domain.Junction{
		Type:     req.Type,
		Location: req.Location,
	}

	junction, err := jh.service.Create(r.Context(), j)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(dto.JunctionToResponse(junction)); err != nil {
	}
}

func (jh JunctionHandler) ListJunctions(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	row, total, err := jh.service.List(r.Context(), 0, 0)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dto.PaginatedJunctions(row, total, 0, 0)); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

}

func (jh JunctionHandler) GetJunction(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	junction, err := jh.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(dto.JunctionToResponse(junction)); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

}

func (jh JunctionHandler) DeleteJunction(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = jh.service.Delete(r.Context(), id)

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
