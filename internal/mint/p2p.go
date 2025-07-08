package mint

import (
	"context"

	types "wetee.app/worker/internal/model"
)

// SendMessageToSecret Sends a message to a randomly selected node of type 1 within the context, while adding OrgId information
func (m *Minter) SendMessageToSecret(ctx context.Context, message *types.Message) error {
	return nil
}
