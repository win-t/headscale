package hscontrol

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/juanfont/headscale/hscontrol/util"
	"github.com/juanfont/headscale/loghandler"
	"github.com/rs/zerolog/log"
	"tailscale.com/tailcfg"
)

// // NoiseRegistrationHandler handles the actual registration process of a machine.
func (ns *noiseServer) NoiseRegistrationHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	log.Trace().Caller().Msgf("Noise registration handler for client %s", req.RemoteAddr)
	if req.Method != http.MethodPost {
		http.Error(writer, "Wrong method", http.StatusMethodNotAllowed)

		return
	}

	log.Trace().
		Any("headers", req.Header).
		Msg("Headers")

	body, _ := io.ReadAll(req.Body)
	registerRequest := tailcfg.RegisterRequest{}
	if err := json.Unmarshal(body, &registerRequest); err != nil {
		log.Error().
			Caller().
			Err(err).
			Msg("Cannot parse RegisterRequest")
		machineRegistrations.WithLabelValues("unknown", "web", "error", "unknown").Inc()
		http.Error(writer, "Internal error", http.StatusInternalServerError)

		return
	}

	ns.nodeKey = registerRequest.NodeKey

	if l, _ := req.Context().Value(loghandler.LogHandlerCtxKey).(*loghandler.LogHandlerCtx); l != nil {
		l.Node = util.NodePublicKeyStripPrefix(registerRequest.NodeKey)
	}

	ns.headscale.handleRegisterCommon(writer, req, registerRequest, ns.conn.Peer(), true)
}
