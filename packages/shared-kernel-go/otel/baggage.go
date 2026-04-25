package otel

import (
	"context"

	"go.opentelemetry.io/otel/baggage"
)

func SetTenantBaggage(ctx context.Context, tenantID string) (context.Context, error) {
	return setBaggageValue(ctx, "tenant_id", tenantID)
}

func GetTenantBaggage(ctx context.Context) string {
	return getBaggageValue(ctx, "tenant_id")
}

func SetGSTINBaggage(ctx context.Context, gstin string) (context.Context, error) {
	return setBaggageValue(ctx, "gstin", gstin)
}

func GetGSTINBaggage(ctx context.Context) string {
	return getBaggageValue(ctx, "gstin")
}

func SetTANBaggage(ctx context.Context, tan string) (context.Context, error) {
	return setBaggageValue(ctx, "tan", tan)
}

func GetTANBaggage(ctx context.Context) string {
	return getBaggageValue(ctx, "tan")
}

func SetPANBaggage(ctx context.Context, pan string) (context.Context, error) {
	return setBaggageValue(ctx, "pan", pan)
}

func GetPANBaggage(ctx context.Context) string {
	return getBaggageValue(ctx, "pan")
}

func setBaggageValue(ctx context.Context, key, value string) (context.Context, error) {
	member, err := baggage.NewMember(key, value)
	if err != nil {
		return ctx, err
	}

	bag := baggage.FromContext(ctx)
	bag, err = bag.SetMember(member)
	if err != nil {
		return ctx, err
	}

	return baggage.ContextWithBaggage(ctx, bag), nil
}

func getBaggageValue(ctx context.Context, key string) string {
	bag := baggage.FromContext(ctx)
	return bag.Member(key).Value()
}
