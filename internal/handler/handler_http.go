package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetGrpcEndpointHandler(ctx *fiber.Ctx) error {
	if h.grpcEndpoint == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "endpoint cannot be empty"})
	}
	return ctx.Status(fiber.StatusFound).JSON(fiber.Map{"endpoint": h.grpcEndpoint})
}

func (h *Handler) GetHandler(ctx *fiber.Ctx) error {
	reqBody := &Body{}
	if err := parseRequest(ctx, reqBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	val, err := h.s.Get(reqBody.Key)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	respBody := &Body{reqBody.Key, val}
	return ctx.Status(fiber.StatusFound).JSON(respBody)
}

func (h *Handler) PutHandler(ctx *fiber.Ctx) error {
	reqBody := &Body{}
	if err := parseRequest(ctx, reqBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.s.Add(reqBody.Key, reqBody.Value, true); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (h *Handler) UpdateHandler(ctx *fiber.Ctx) error {
	reqBody := &Body{}
	if err := parseRequest(ctx, reqBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.s.Update(reqBody.Key, reqBody.Value, true); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) RemoveHandler(ctx *fiber.Ctx) error {
	reqBody := &Body{}
	if err := parseRequest(ctx, reqBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.s.Remove(reqBody.Key, true); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func parseRequest(ctx *fiber.Ctx, b *Body) error {
	if err := ctx.BodyParser(b); err != nil {
		return fmt.Errorf("error parsing request body %w", err)
	}
	return nil
}
