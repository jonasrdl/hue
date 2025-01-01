package hue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/grandcat/zeroconf"
)

// DiscoverBridgeIP discovers a Philips Hue Bridge using mDNS and returns its IP.
func DiscoverBridgeIP() (string, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return "", fmt.Errorf("failed to initialize resolver: %w", err)
	}

	results := make(chan *zeroconf.ServiceEntry)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		if err := resolver.Browse(ctx, "_hue._tcp", "local.", results); err != nil {
			close(results)
		}
	}()

	var entry *zeroconf.ServiceEntry
loop:
	for {
		select {
		case result, ok := <-results:
			if !ok {
				break loop
			}
			if result != nil && len(result.AddrIPv4) > 0 {
				entry = result
				break loop
			}
		case <-ctx.Done():
			return "", errors.New("discovery timed out")
		}
	}

	if entry == nil {
		return "", errors.New("no bridges found")
	}

	return entry.AddrIPv4[0].String(), nil
}
