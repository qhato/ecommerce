package commands

import (
	"context"
	"fmt"

	"github.com/qhato/ecommerce/internal/return/domain"
)

type ReturnCommandHandler struct {
	returnRepo domain.ReturnRepository
}

func NewReturnCommandHandler(returnRepo domain.ReturnRepository) *ReturnCommandHandler {
	return &ReturnCommandHandler{returnRepo: returnRepo}
}

func (h *ReturnCommandHandler) HandleCreateReturn(ctx context.Context, cmd CreateReturnCommand) (*domain.ReturnRequest, error) {
	returnReq, err := domain.NewReturnRequest(cmd.OrderID, fmt.Sprintf("%d", cmd.CustomerID), domain.ReturnReason(cmd.Reason))
	if err != nil {
		return nil, err
	}

	if err := h.returnRepo.Create(ctx, returnReq); err != nil {
		return nil, fmt.Errorf("failed to create return: %w", err)
	}

	return returnReq, nil
}

func (h *ReturnCommandHandler) HandleApproveReturn(ctx context.Context, cmd ApproveReturnCommand) (*domain.ReturnRequest, error) {
	returnReq, err := h.returnRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find return: %w", err)
	}
	if returnReq == nil {
		return nil, domain.ErrReturnNotFound
	}

	returnReq.Approve(1) // Default approver
	if err := h.returnRepo.Update(ctx, returnReq); err != nil {
		return nil, fmt.Errorf("failed to approve return: %w", err)
	}

	return returnReq, nil
}

func (h *ReturnCommandHandler) HandleRejectReturn(ctx context.Context, cmd RejectReturnCommand) (*domain.ReturnRequest, error) {
	returnReq, err := h.returnRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find return: %w", err)
	}
	if returnReq == nil {
		return nil, domain.ErrReturnNotFound
	}

	returnReq.Reject(cmd.Reason)
	if err := h.returnRepo.Update(ctx, returnReq); err != nil {
		return nil, fmt.Errorf("failed to reject return: %w", err)
	}

	return returnReq, nil
}

func (h *ReturnCommandHandler) HandleReceiveReturn(ctx context.Context, cmd ReceiveReturnCommand) (*domain.ReturnRequest, error) {
	returnReq, err := h.returnRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find return: %w", err)
	}
	if returnReq == nil {
		return nil, domain.ErrReturnNotFound
	}

	returnReq.MarkAsReceived()
	if err := h.returnRepo.Update(ctx, returnReq); err != nil {
		return nil, fmt.Errorf("failed to receive return: %w", err)
	}

	return returnReq, nil
}

func (h *ReturnCommandHandler) HandleInspectReturn(ctx context.Context, cmd InspectReturnCommand) (*domain.ReturnRequest, error) {
	returnReq, err := h.returnRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find return: %w", err)
	}
	if returnReq == nil {
		return nil, domain.ErrReturnNotFound
	}

	returnReq.MarkAsInspected()
	if err := h.returnRepo.Update(ctx, returnReq); err != nil {
		return nil, fmt.Errorf("failed to inspect return: %w", err)
	}

	return returnReq, nil
}

func (h *ReturnCommandHandler) HandleProcessRefund(ctx context.Context, cmd ProcessRefundCommand) (*domain.ReturnRequest, error) {
	returnReq, err := h.returnRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find return: %w", err)
	}
	if returnReq == nil {
		return nil, domain.ErrReturnNotFound
	}

	returnReq.MarkAsRefunded(cmd.RefundAmount, cmd.RefundMethod)
	if err := h.returnRepo.Update(ctx, returnReq); err != nil {
		return nil, fmt.Errorf("failed to process refund: %w", err)
	}

	return returnReq, nil
}

func (h *ReturnCommandHandler) HandleCancelReturn(ctx context.Context, cmd CancelReturnCommand) (*domain.ReturnRequest, error) {
	returnReq, err := h.returnRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find return: %w", err)
	}
	if returnReq == nil {
		return nil, domain.ErrReturnNotFound
	}

	returnReq.Cancel()
	if err := h.returnRepo.Update(ctx, returnReq); err != nil {
		return nil, fmt.Errorf("failed to cancel return: %w", err)
	}

	return returnReq, nil
}
