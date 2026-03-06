package meter

import (
	"fmt"

	"github.com/evcc-io/evcc/api"
	homewizard "github.com/evcc-io/evcc/meter/homewizard-v2"
	"github.com/go-viper/mapstructure/v2"
)

func init() {
	registry.Add("homewizard-v2", NewHomeWizardV2FromConfig)
}

// NewHomeWizardV2FromConfig creates a HomeWizard meter from configuration
func NewHomeWizardV2FromConfig(other map[string]any) (api.Meter, error) {
	cc := homewizard.Config{
		Timeout: homewizard.DefaultTimeout,
	}

	// Decode loosely: sub-type specific fields (e.g. phases) are handled by the individual constructors
	if err := mapstructure.WeakDecode(other, &cc); err != nil {
		return nil, err
	}

	if cc.Host == "" || cc.Token == "" {
		return nil, fmt.Errorf("missing host or token - run 'evcc token homewizard'")
	}

	if cc.Usage == "" {
		return nil, fmt.Errorf("missing required parameter 'usage' (must be one of: grid, pv, charge, battery)")
	}

	// Dispatch based on usage
	switch cc.Usage {
	case "grid":
		return homewizard.NewHomeWizardP1FromConfig(cc, other)
	case "pv", "charge":
		return homewizard.NewHomeWizardKWHFromConfig(cc, other)
	case "battery":
		return homewizard.NewHomeWizardBatteryFromConfig(cc, other)
	default:
		return nil, fmt.Errorf("invalid usage '%s': must be one of: grid, pv, charge, battery", cc.Usage)
	}
}
