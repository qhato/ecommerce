package rules

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/qhato/ecommerce/internal/offer/domain"
	"github.com/shopspring/decimal"
)

// ExpressionEvaluator implements domain.RuleEvaluator
// This is a simplified expression evaluator for offer rules
// In production, consider using a more robust rule engine like:
// - github.com/hyperjumptech/grule-rule-engine
// - github.com/google/cel-go
type ExpressionEvaluator struct{}

// NewExpressionEvaluator creates a new ExpressionEvaluator
func NewExpressionEvaluator() domain.RuleEvaluator {
	return &ExpressionEvaluator{}
}

// Evaluate evaluates a rule expression against the given context
// Supports simple expressions like:
// - "item.CategoryID == '123'"
// - "item.Price > 100"
// - "order.OrderSubtotal >= 50"
// - "item.SKUID in ['SKU-1', 'SKU-2']"
// - "item.ProductID == 'PROD-123' and item.Quantity >= 2"
func (e *ExpressionEvaluator) Evaluate(ruleExpression string, context map[string]interface{}) (bool, error) {
	if ruleExpression == "" {
		return true, nil
	}

	// Handle logical operators (AND, OR)
	if strings.Contains(ruleExpression, " and ") || strings.Contains(ruleExpression, " AND ") {
		return e.evaluateAnd(ruleExpression, context)
	}
	if strings.Contains(ruleExpression, " or ") || strings.Contains(ruleExpression, " OR ") {
		return e.evaluateOr(ruleExpression, context)
	}

	// Single expression evaluation
	return e.evaluateSingleExpression(ruleExpression, context)
}

func (e *ExpressionEvaluator) evaluateAnd(expression string, context map[string]interface{}) (bool, error) {
	parts := splitByOperator(expression, " and ", " AND ")
	for _, part := range parts {
		result, err := e.Evaluate(strings.TrimSpace(part), context)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}
	return true, nil
}

func (e *ExpressionEvaluator) evaluateOr(expression string, context map[string]interface{}) (bool, error) {
	parts := splitByOperator(expression, " or ", " OR ")
	for _, part := range parts {
		result, err := e.Evaluate(strings.TrimSpace(part), context)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

func (e *ExpressionEvaluator) evaluateSingleExpression(expression string, context map[string]interface{}) (bool, error) {
	expression = strings.TrimSpace(expression)

	// Handle 'in' operator
	if strings.Contains(expression, " in ") {
		return e.evaluateInOperator(expression, context)
	}

	// Handle comparison operators
	for _, op := range []string{">=", "<=", "==", "!=", ">", "<"} {
		if strings.Contains(expression, op) {
			return e.evaluateComparison(expression, op, context)
		}
	}

	return false, fmt.Errorf("unsupported expression: %s", expression)
}

func (e *ExpressionEvaluator) evaluateInOperator(expression string, context map[string]interface{}) (bool, error) {
	parts := strings.Split(expression, " in ")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid 'in' expression: %s", expression)
	}

	leftSide := strings.TrimSpace(parts[0])
	rightSide := strings.TrimSpace(parts[1])

	// Get left side value
	leftValue, err := e.resolveValue(leftSide, context)
	if err != nil {
		return false, err
	}

	// Parse array on right side: ['value1', 'value2']
	arrayPattern := regexp.MustCompile(`\[(.*?)\]`)
	matches := arrayPattern.FindStringSubmatch(rightSide)
	if len(matches) < 2 {
		return false, fmt.Errorf("invalid array syntax: %s", rightSide)
	}

	arrayContent := matches[1]
	arrayValues := strings.Split(arrayContent, ",")

	leftStr := fmt.Sprintf("%v", leftValue)
	for _, val := range arrayValues {
		val = strings.TrimSpace(val)
		val = strings.Trim(val, "'\"")
		if leftStr == val {
			return true, nil
		}
	}

	return false, nil
}

func (e *ExpressionEvaluator) evaluateComparison(expression string, operator string, context map[string]interface{}) (bool, error) {
	parts := strings.Split(expression, operator)
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid comparison expression: %s", expression)
	}

	leftSide := strings.TrimSpace(parts[0])
	rightSide := strings.TrimSpace(parts[1])

	leftValue, err := e.resolveValue(leftSide, context)
	if err != nil {
		return false, err
	}

	rightValue, err := e.resolveValue(rightSide, context)
	if err != nil {
		return false, err
	}

	return e.compare(leftValue, rightValue, operator)
}

func (e *ExpressionEvaluator) resolveValue(path string, context map[string]interface{}) (interface{}, error) {
	path = strings.TrimSpace(path)

	// Handle string literals
	if strings.HasPrefix(path, "'") && strings.HasSuffix(path, "'") {
		return strings.Trim(path, "'"), nil
	}
	if strings.HasPrefix(path, "\"") && strings.HasSuffix(path, "\"") {
		return strings.Trim(path, "\""), nil
	}

	// Handle numeric literals
	if num, err := strconv.ParseFloat(path, 64); err == nil {
		return num, nil
	}

	// Handle boolean literals
	if path == "true" {
		return true, nil
	}
	if path == "false" {
		return false, nil
	}

	// Resolve from context (e.g., "item.CategoryID", "order.OrderSubtotal")
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid path: %s", path)
	}

	// Get root object from context
	rootKey := parts[0]
	obj, exists := context[rootKey]
	if !exists {
		return nil, fmt.Errorf("key not found in context: %s", rootKey)
	}

	// If no nested path, return the object
	if len(parts) == 1 {
		return obj, nil
	}

	// Navigate nested path
	return e.getNestedValue(obj, parts[1:])
}

func (e *ExpressionEvaluator) getNestedValue(obj interface{}, path []string) (interface{}, error) {
	if len(path) == 0 {
		return obj, nil
	}

	fieldName := path[0]

	// Handle OfferItem
	if item, ok := obj.(domain.OfferItem); ok {
		return e.getOfferItemField(item, fieldName)
	}

	// Handle OfferContext
	if ctx, ok := obj.(*domain.OfferContext); ok {
		return e.getOfferContextField(ctx, fieldName)
	}

	return nil, fmt.Errorf("unsupported object type for field access: %T", obj)
}

func (e *ExpressionEvaluator) getOfferItemField(item domain.OfferItem, fieldName string) (interface{}, error) {
	switch fieldName {
	case "ItemID":
		return item.ItemID, nil
	case "SKUID":
		return item.SKUID, nil
	case "CategoryID":
		if item.CategoryID != nil {
			return *item.CategoryID, nil
		}
		return nil, nil
	case "Price":
		price, _ := item.Price.Float64()
		return price, nil
	case "SalePrice":
		if item.SalePrice != nil {
			price, _ := item.SalePrice.Float64()
			return price, nil
		}
		return nil, nil
	case "Quantity":
		return item.Quantity, nil
	case "Subtotal":
		subtotal, _ := item.Subtotal.Float64()
		return subtotal, nil
	case "ProductID":
		if item.ProductID != nil {
			return *item.ProductID, nil
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown field: %s", fieldName)
	}
}

func (e *ExpressionEvaluator) getOfferContextField(ctx *domain.OfferContext, fieldName string) (interface{}, error) {
	switch fieldName {
	case "OrderTotal":
		total, _ := ctx.OrderTotal.Float64()
		return total, nil
	case "OrderSubtotal":
		subtotal, _ := ctx.OrderSubtotal.Float64()
		return subtotal, nil
	case "CustomerID":
		if ctx.CustomerID != nil {
			return *ctx.CustomerID, nil
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown field: %s", fieldName)
	}
}

func (e *ExpressionEvaluator) compare(left, right interface{}, operator string) (bool, error) {
	switch operator {
	case "==":
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right), nil
	case "!=":
		return fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right), nil
	case ">", "<", ">=", "<=":
		return e.compareNumeric(left, right, operator)
	default:
		return false, fmt.Errorf("unsupported operator: %s", operator)
	}
}

func (e *ExpressionEvaluator) compareNumeric(left, right interface{}, operator string) (bool, error) {
	leftNum, err := e.toFloat64(left)
	if err != nil {
		return false, fmt.Errorf("left side is not numeric: %w", err)
	}

	rightNum, err := e.toFloat64(right)
	if err != nil {
		return false, fmt.Errorf("right side is not numeric: %w", err)
	}

	leftDec := decimal.NewFromFloat(leftNum)
	rightDec := decimal.NewFromFloat(rightNum)

	switch operator {
	case ">":
		return leftDec.GreaterThan(rightDec), nil
	case "<":
		return leftDec.LessThan(rightDec), nil
	case ">=":
		return leftDec.GreaterThanOrEqual(rightDec), nil
	case "<=":
		return leftDec.LessThanOrEqual(rightDec), nil
	default:
		return false, fmt.Errorf("unsupported numeric operator: %s", operator)
	}
}

func (e *ExpressionEvaluator) toFloat64(val interface{}) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", val)
	}
}

func splitByOperator(expression string, operators ...string) []string {
	for _, op := range operators {
		if strings.Contains(expression, op) {
			parts := strings.Split(expression, op)
			result := make([]string, 0, len(parts))
			for _, part := range parts {
				result = append(result, strings.TrimSpace(part))
			}
			return result
		}
	}
	return []string{expression}
}
