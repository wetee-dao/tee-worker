package mint

import (
	"context"

	"wetee.app/worker/internal/store"
)

// SendMessageToSecret Sends a message to a randomly selected node of type 1 within the context, while adding OrgId information
func (m *Minter) SendMessageToSecret(ctx context.Context, message *store.Message) error {
	return nil
}
