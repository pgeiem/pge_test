package parser

// func (r *TariffRule) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {
// func unmarshalTariffRule(rule *engine.TariffRule, data []byte) error {

// 	fmt.Println(">>>>>>>>>>>>>>>>>>>")

// 	temp := map[string]ast.Node{}

// 	err := yaml.Unmarshal(data, &temp)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse rule type: %w", err)
// 	}

// 	if len(temp) != 1 {
// 		return fmt.Errorf("invalid rule, only one rule definititon per type is allowed")
// 	}

// 	for ruletype, node := range temp {
// 		switch strings.ToLower(ruletype) {
// 		case "free", "banned":
// 			r := struct {
// 				Name                    string `yaml:"name"`
// 				engine.RecurrentSegment `yaml:",inline"`
// 			}{}
// 			if err := yaml.NodeToValue(node, &r, decoderOptions()...); err != nil {
// 				return fmt.Errorf("failed to parse %s rule: %w", ruletype, err)
// 			}
// 			rule.Name = r.Name
// 			rule.RecurrentSegment = r.RecurrentSegment
// 			rule.StartAmount = 0
// 			rule.EndAmount = 0
// 			if ruletype == "free" {
// 				rule.DurationType = engine.FreeDuration
// 			} else {
// 				rule.DurationType = engine.BannedDuration
// 			}
// 			rule.IsInterpolable = false

// 		case "linear":
// 			r := struct {
// 				Name                    string `yaml:"name"`
// 				engine.RecurrentSegment `yaml:",inline"`
// 				HourlyRate              engine.Amount `yaml:"hourlyrate"`
// 			}{}
// 			if err := yaml.NodeToValue(node, &r, decoderOptions()...); err != nil {
// 				return fmt.Errorf("failed to parse %s rule: %w", ruletype, err)
// 			}
// 			rule.Name = r.Name
// 			rule.RecurrentSegment = r.RecurrentSegment
// 			rule.StartAmount = 0
// 			rule.EndAmount = 0 //TODO set from hourlyrate
// 			rule.DurationType = engine.PayingDuration
// 			rule.IsInterpolable = true

// 		case "flatrate":
// 			r := struct {
// 				Name                    string `yaml:"name"`
// 				engine.RecurrentSegment `yaml:",inline"`
// 				Amount                  engine.Amount `yaml:"amount"`
// 			}{}
// 			if err := yaml.NodeToValue(node, &r, decoderOptions()...); err != nil {
// 				return fmt.Errorf("failed to parse %s rule: %w", ruletype, err)
// 			}
// 			rule.Name = r.Name
// 			rule.RecurrentSegment = r.RecurrentSegment
// 			rule.StartAmount = r.Amount
// 			rule.EndAmount = r.Amount
// 			rule.DurationType = engine.PayingDuration
// 			rule.IsInterpolable = false

// 		default:
// 			return fmt.Errorf("invalid rule type: %s", ruletype)
// 		}
// 		fmt.Println(">>>", rule)
// 	}

// 	/**r = make(QuotaInventory)
// 	for _, t := range temp {
// 		quota := Quota(nil)
// 		// TODO return an error if both DurationQuota and CounterQuota are set
// 		if t.DurationQuota != nil {
// 			quota = t.DurationQuota
// 		} else if t.CounterQuota != nil {
// 			quota = t.CounterQuota
// 		}
// 		if quota != nil {
// 			if quota.GetName() == "" {
// 				return fmt.Errorf("missing quota name")
// 			}
// 			(*r)[quota.GetName()] = quota
// 		}
// 	}*/
// 	return nil
// }
