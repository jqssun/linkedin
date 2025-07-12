package linkedingo

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *Client) RegisterAndroidPush(ctx context.Context, token string) error {
	payload := c.createPushRegistrationPayload(token)
	r := bytes.NewReader(payload)

	deviceId, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate device id: %w", err)
	}
	trackHeader := fmt.Sprintf(
		`{"osName":"Android OS","osVersion":"35","clientVersion":"4.1.1090","clientMinorVersion":199800,"appId":"com.linkedin.android","deviceId":"%s","storeId":"us_googleplay","isAdTrackingLimited":true,"mpName":"voyager-android","mpVersion":"2.168.97"}`, // TODO: check
		deviceId.String(),
	)

	_, err = c.newAuthedRequest(http.MethodPost, linkedInVoyagerNotificationsDashPushRegistration).
		WithQueryParam("action", "register").
		WithXLIHeaders().
		WithCSRF().
		WithContentType("application/x-protobuf2 ;symbol-table=voyager-21304").
		WithHeader("Accept", "application/vnd.linkedin.deduped+x-protobuf").
		WithHeader("X-LI-Track", trackHeader).
		WithBody(r).
		Do(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to register push notification: %w", err)
	}
	return nil
}

func (c *Client) createPushRegistrationPayload(token string) []byte {
	const (
		protobufStart              = byte(0x00)
		leadingOrdinal             = byte(0x14)
		arrayStart                 = byte(0x01)
		pushNotificationTokensKey  = "pushNotificationTokens"
		pushNotificationEnabledKey = "pushNotificationEnabled"
	)

	payload := []byte{
		protobufStart,
		2, // number of keys
		leadingOrdinal,
		byte(len(pushNotificationTokensKey)),
	}
	payload = append(payload, pushNotificationTokensKey...)
	payload = append(payload,
		arrayStart,
		1, // Array Length
		leadingOrdinal,
		byte(len(token)),
		0x01, // Unknown
	)
	payload = append(payload, token...)
	payload = append(payload,
		leadingOrdinal,
		byte(len(pushNotificationEnabledKey)),
	)
	payload = append(payload, pushNotificationEnabledKey...)
	payload = append(payload, 0x08)

	return payload
}
