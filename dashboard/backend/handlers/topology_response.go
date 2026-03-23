package handlers

import (
	"fmt"
	"strings"

	routerconfig "github.com/vllm-project/semantic-router/src/semantic-router/pkg/config"
)

type topologySignalMapping struct {
	signalType string
	names      []string
	confidence float64
	reason     string
	addPath    bool
}

// convertRouterResponse converts Router API response to TestQueryResult.
func convertRouterResponse(req TestQueryRequest, routerResp *RouterIntentResponse, configPath string) *TestQueryResult {
	result := newTestQueryResult(req)

	appendMatchedSignals(result, routerResp)
	appendSignalGroupHighlights(result)
	applyRouterDecision(result, routerResp)
	applyRecommendedModel(result, routerResp.RecommendedModel)
	appendEvaluatedRulesFromConfig(result, configPath)

	return result
}

func newTestQueryResult(req TestQueryRequest) *TestQueryResult {
	return &TestQueryResult{
		Query:           req.Query,
		Mode:            req.Mode,
		MatchedSignals:  []MatchedSignal{},
		MatchedModels:   []string{},
		HighlightedPath: []string{"client"},
		IsAccurate:      true,
		EvaluatedRules:  []EvaluatedRule{},
	}
}

func appendMatchedSignals(result *TestQueryResult, routerResp *RouterIntentResponse) {
	if routerResp.MatchedSignals == nil {
		return
	}

	for _, mapping := range topologySignalMappings(routerResp) {
		addMatchedSignals(result, mapping)
	}
}

func topologySignalMappings(routerResp *RouterIntentResponse) []topologySignalMapping {
	return []topologySignalMapping{
		{signalType: "keyword", names: routerResp.MatchedSignals.Keywords, confidence: 1.0, reason: "Keyword rule matched", addPath: true},
		{signalType: "embedding", names: routerResp.MatchedSignals.Embeddings, confidence: 0.85, reason: "Embedding similarity matched", addPath: true},
		{signalType: "domain", names: routerResp.MatchedSignals.Domains, confidence: routerResp.Classification.Confidence, reason: "Domain classification matched", addPath: true},
		{signalType: "fact_check", names: routerResp.MatchedSignals.FactCheck, confidence: 0.9, reason: "Fact check signal matched"},
		{signalType: "preference", names: routerResp.MatchedSignals.Preferences, confidence: 1.0, reason: "User preference matched", addPath: true},
		{signalType: "user_feedback", names: routerResp.MatchedSignals.UserFeedback, confidence: 1.0, reason: "User feedback matched", addPath: true},
		{signalType: "language", names: routerResp.MatchedSignals.Language, confidence: 0.95, reason: "Language detected", addPath: true},
		{signalType: "context", names: routerResp.MatchedSignals.Context, confidence: 1.0, reason: "Context token count matched", addPath: true},
		{signalType: "complexity", names: routerResp.MatchedSignals.Complexity, confidence: 0.9, reason: "Complexity level matched", addPath: true},
		{signalType: "modality", names: routerResp.MatchedSignals.Modality, confidence: 1.0, reason: "Modality signal matched", addPath: true},
		{signalType: "authz", names: routerResp.MatchedSignals.Authz, confidence: 1.0, reason: "Authorization signal matched", addPath: true},
		{signalType: "jailbreak", names: routerResp.MatchedSignals.Jailbreak, confidence: 1.0, reason: "Jailbreak signal matched", addPath: true},
		{signalType: "pii", names: routerResp.MatchedSignals.PII, confidence: 1.0, reason: "PII signal matched", addPath: true},
		{signalType: "projection", names: routerResp.MatchedSignals.Projection, confidence: 1.0, reason: "Projection mapping matched", addPath: true},
	}
}

func addMatchedSignals(result *TestQueryResult, mapping topologySignalMapping) {
	for _, name := range mapping.names {
		result.MatchedSignals = append(result.MatchedSignals, MatchedSignal{
			Type:       mapping.signalType,
			Name:       name,
			Confidence: mapping.confidence,
			Reason:     mapping.reason,
		})
		if mapping.addPath {
			result.HighlightedPath = append(result.HighlightedPath, fmt.Sprintf("signal-%s-%s", mapping.signalType, name))
		}
	}
}

func appendSignalGroupHighlights(result *TestQueryResult) {
	if len(result.MatchedSignals) == 0 {
		return
	}

	signalTypes := make(map[string]bool)
	for _, signal := range result.MatchedSignals {
		signalTypes[signal.Type] = true
	}
	for signalType := range signalTypes {
		result.HighlightedPath = append(result.HighlightedPath, fmt.Sprintf("signal-group-%s", signalType))
	}
}

func applyRouterDecision(result *TestQueryResult, routerResp *RouterIntentResponse) {
	if routerResp.DecisionResult != nil {
		result.MatchedDecision = routerResp.DecisionResult.DecisionName
		result.HighlightedPath = append(result.HighlightedPath, fmt.Sprintf("decision-%s", routerResp.DecisionResult.DecisionName))
		result.EvaluatedRules = append(result.EvaluatedRules, EvaluatedRule{
			DecisionName: routerResp.DecisionResult.DecisionName,
			Conditions:   routerResp.DecisionResult.MatchedRules,
			MatchedCount: len(routerResp.DecisionResult.MatchedRules),
			TotalCount:   len(routerResp.DecisionResult.MatchedRules),
			IsMatch:      true,
		})
		return
	}

	if routerResp.RoutingDecision == "" {
		return
	}

	result.MatchedDecision = routerResp.RoutingDecision
	result.HighlightedPath = append(result.HighlightedPath, fmt.Sprintf("decision-%s", routerResp.RoutingDecision))
	if isSystemFallbackDecision(routerResp.RoutingDecision) {
		result.IsFallbackDecision = true
		result.FallbackReason = getFallbackReason(routerResp.RoutingDecision)
		result.HighlightedPath = append(result.HighlightedPath, "fallback-decision")
	}
}

func applyRecommendedModel(result *TestQueryResult, recommendedModel string) {
	if recommendedModel == "" {
		return
	}

	result.MatchedModels = append(result.MatchedModels, recommendedModel)
	result.HighlightedPath = append(
		result.HighlightedPath,
		fmt.Sprintf("model-%s", normalizeModelName(recommendedModel)),
	)
}

func appendEvaluatedRulesFromConfig(result *TestQueryResult, configPath string) {
	parsedConfig, err := routerconfig.Parse(configPath)
	if err != nil || parsedConfig == nil {
		return
	}

	matchedSignalNames := buildMatchedSignalNameSet(result.MatchedSignals)
	for _, decision := range parsedConfig.IntelligentRouting.Decisions {
		if result.MatchedDecision != "" && decision.Name == result.MatchedDecision {
			continue
		}
		result.EvaluatedRules = append(result.EvaluatedRules, buildEvaluatedRule(decision, matchedSignalNames))
	}
}

func buildMatchedSignalNameSet(signals []MatchedSignal) map[string]bool {
	matchedSignalNames := make(map[string]bool, len(signals)*2)
	for _, signal := range signals {
		key := fmt.Sprintf("%s:%s", signal.Type, signal.Name)
		normalizedKey := fmt.Sprintf("%s:%s", signal.Type, normalizeSignalName(signal.Name))
		matchedSignalNames[key] = true
		matchedSignalNames[normalizedKey] = true
	}
	return matchedSignalNames
}

func buildEvaluatedRule(decision routerconfig.Decision, matchedSignalNames map[string]bool) EvaluatedRule {
	rule := EvaluatedRule{
		DecisionName: decision.Name,
		RuleOperator: strings.ToUpper(decision.Rules.Operator),
		IsMatch:      false,
		Priority:     decision.Priority,
	}
	if rule.RuleOperator == "" {
		rule.RuleOperator = "AND"
	}

	for _, condition := range decision.Rules.Conditions {
		conditionKey := fmt.Sprintf("%s:%s", condition.Type, condition.Name)
		normalizedConditionKey := fmt.Sprintf("%s:%s", condition.Type, normalizeSignalName(condition.Name))
		rule.Conditions = append(rule.Conditions, conditionKey)
		rule.TotalCount++
		if matchedSignalNames[conditionKey] || matchedSignalNames[normalizedConditionKey] {
			rule.MatchedCount++
		}
	}

	switch {
	case rule.TotalCount == 0:
		rule.IsMatch = true
	case rule.RuleOperator == "OR":
		rule.IsMatch = rule.MatchedCount > 0
	default:
		rule.IsMatch = rule.MatchedCount == rule.TotalCount
	}

	return rule
}
