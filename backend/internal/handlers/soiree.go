package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"gueuledebois/backend/internal/services"
)

// maxPhotoSize : 10 Mo.
const maxPhotoSize = 10 << 20 // 10 Mo

// mimesAutorises : jpeg/png/webp uniquement (spec).
var mimesAutorises = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// SoireeHandler expose UC06 (créer), UC07 (modifier), UC08 (supprimer),
// UC10 (historique) et l'upload minimal de photos (enrichissement
// UC06/UC07, spec "photos minimales").
type SoireeHandler struct {
	soirees   *services.SoireeService
	uploadDir string
}

func NewSoireeHandler(soirees *services.SoireeService, uploadDir string) *SoireeHandler {
	return &SoireeHandler{soirees: soirees, uploadDir: uploadDir}
}

type soireeRequest struct {
	Titre       string    `json:"titre"`
	Date        time.Time `json:"date"`
	Lieu        string    `json:"lieu"`
	Description string    `json:"description"`
}

func (h *SoireeHandler) decode(r *http.Request) (services.CreateSoireeInput, error) {
	var req soireeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return services.CreateSoireeInput{}, services.ErrValidation
	}
	return services.CreateSoireeInput{
		Titre:       req.Titre,
		DateSoiree:  req.Date,
		Lieu:        req.Lieu,
		Description: req.Description,
	}, nil
}

// Create — POST /api/soirees (UC06, inclut UC16).
func (h *SoireeHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	in, err := h.decode(r)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}

	soiree, err := h.soirees.Create(r.Context(), user.ID, in)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, soiree)
}

// Update — PUT /api/soirees/{id} (UC07, inclut UC16).
func (h *SoireeHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	id, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	in, err := h.decode(r)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}

	soiree, err := h.soirees.Update(r.Context(), user.ID, id, in)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, soiree)
}

// Delete — DELETE /api/soirees/{id} (UC08, inclut UC16).
func (h *SoireeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	id, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	if err := h.soirees.Delete(r.Context(), user.ID, id); err != nil {
		mapAndWriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Get — GET /api/soirees/{id}.
func (h *SoireeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	soiree, err := h.soirees.Get(r.Context(), id)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	photos, _ := h.soirees.ListPhotos(r.Context(), id)
	writeJSON(w, http.StatusOK, map[string]any{"soiree": soiree, "photos": photos})
}

// ListMine — GET /api/soirees (UC10, historique de l'utilisateur connecté).
func (h *SoireeHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	list, err := h.soirees.ListByUser(r.Context(), user.ID)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// UploadPhoto — POST /api/soirees/{id}/photos (multipart/form-data, champ
// "photo"). Validation MIME (jpeg/png/webp) + taille max 5 Mo, fichier
// renommé en UUID, écrit tel quel sur disque (aucun redimensionnement,
// cf. spec "photos minimales"). UC06 précise : "photo invalide -> photo
// rejetée, création poursuivie sans elle" -> ici on répond une erreur
// dédiée sur cet endpoint séparé, au frontend de ne pas bloquer le flux.
func (h *SoireeHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	soireeID, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxPhotoSize+1<<20) // marge pour le multipart overhead
	if err := r.ParseMultipartForm(maxPhotoSize + 1<<20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_input", "Fichier trop volumineux ou requête multipart invalide.")
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_input", "Champ 'photo' manquant.")
		return
	}
	defer file.Close()

	if header.Size > maxPhotoSize {
		writeError(w, http.StatusBadRequest, "invalid_input", "Photo trop volumineuse (max 5 Mo).")
		return
	}

	// Détection du type MIME réel (pas seulement l'en-tête déclaré par le
	// client) : on lit les 512 premiers octets.
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	contentType := http.DetectContentType(buf[:n])
	ext, ok := mimesAutorises[contentType]
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_input", "Format de photo non supporté (jpeg/png/webp uniquement).")
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Erreur de lecture du fichier.")
		return
	}

	filename := uuid.NewString() + ext
	dstPath := filepath.Join(h.uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Impossible d'écrire le fichier.")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Impossible d'écrire le fichier.")
		return
	}

	photo, err := h.soirees.AddPhoto(r.Context(), user.ID, soireeID, "/uploads/"+filename)
	if err != nil {
		_ = os.Remove(dstPath) // on ne laisse pas un fichier orphelin si l'association échoue
		mapAndWriteError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, photo)
}

func idFromURL(r *http.Request, param string) (int64, error) {
	raw := chi.URLParam(r, param)
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("id invalide")
	}
	return id, nil
}
