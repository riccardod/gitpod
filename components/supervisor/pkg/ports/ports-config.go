// Copyright (c) 2020 TypeFox GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package ports

import (
	"context"

	"github.com/gitpod-io/gitpod/supervisor/pkg/gitpod"
)

// ConfigInterface provides access to port configurations
type ConfigInterface interface {
	// Observe provides channels triggered whenever the port configurations are changed
	Observe(ctx context.Context) (<-chan []*gitpod.PortsItems, <-chan error)
}

// ConfigService provides access to port configurations
type ConfigService struct {
	configService *gitpod.ConfigService
}

// NewConfigService creates a new instance of ConfigService
func NewConfigService(configService *gitpod.ConfigService) *ConfigService {
	return &ConfigService{
		configService: configService,
	}
}

// Observe provides channels triggered whenever the port configurations are changed
func (service *ConfigService) Observe(ctx context.Context) (<-chan []*gitpod.PortsItems, <-chan error) {
	portConfigs := make(chan []*gitpod.PortsItems)
	portConfigErrs := make(chan error)
	go func() {
		defer close(portConfigs)
		defer close(portConfigErrs)
		configs, errs := service.configService.Observe(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case config := <-configs:
				if config != nil {
					portConfigs <- config.Ports
				} else {
					portConfigs <- nil
				}
			case err := <-errs:
				portConfigErrs <- err
			}
		}
	}()
	return portConfigs, portConfigErrs
}
