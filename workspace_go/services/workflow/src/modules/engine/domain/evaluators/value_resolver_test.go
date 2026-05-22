package evaluators

import (
	"errors"
	"testing"

	defEntities "workflow/src/modules/definitions/domain/entities"
)

func fv(ftype defEntities.FieldValueType, value string) *defEntities.FieldValue {
	return &defEntities.FieldValue{Type: ftype, Value: value}
}

func fvWithNode(ftype defEntities.FieldValueType, value, nodeID string) *defEntities.FieldValue {
	return &defEntities.FieldValue{Type: ftype, Value: value, NodeID: nodeID}
}

func TestResolve_Literal(t *testing.T) {
	r := NewValueResolver()
	val, err := r.Resolve(fv(defEntities.FieldValueLiteral, "hello"), nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hello" {
		t.Fatalf("expected 'hello', got %v", val)
	}
}

func TestResolve_Event(t *testing.T) {
	r := NewValueResolver()
	event := map[string]interface{}{"temperature": 42.5}
	val, err := r.Resolve(fv(defEntities.FieldValueEvent, "temperature"), event, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 42.5 {
		t.Fatalf("expected 42.5, got %v", val)
	}
}

func TestResolve_Event_Nested(t *testing.T) {
	r := NewValueResolver()
	event := map[string]interface{}{
		"device": map[string]interface{}{
			"sensor": map[string]interface{}{
				"temp": 99,
			},
		},
	}
	val, err := r.Resolve(fv(defEntities.FieldValueEvent, "device.sensor.temp"), event, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 99 {
		t.Fatalf("expected 99, got %v", val)
	}
}

func TestResolve_State(t *testing.T) {
	r := NewValueResolver()
	state := map[string]interface{}{"count": 10}
	val, err := r.Resolve(fv(defEntities.FieldValueState, "count"), nil, state, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 10 {
		t.Fatalf("expected 10, got %v", val)
	}
}

func TestResolve_Variable(t *testing.T) {
	r := NewValueResolver()
	state := map[string]interface{}{"threshold": 50}
	val, err := r.Resolve(fv(defEntities.FieldValueVariable, "threshold"), nil, state, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 50 {
		t.Fatalf("expected 50, got %v", val)
	}
}

func TestResolve_Input(t *testing.T) {
	r := NewValueResolver()
	inputs := map[string]interface{}{"userId": "abc123"}
	val, err := r.Resolve(fv(defEntities.FieldValueInput, "input.userId"), nil, nil, nil, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "abc123" {
		t.Fatalf("expected 'abc123', got %v", val)
	}
}

func TestResolve_NodeOutput_FullMap(t *testing.T) {
	r := NewValueResolver()
	outputs := map[string]interface{}{
		"node1": map[string]interface{}{"result": "ok"},
	}
	val, err := r.Resolve(fvWithNode(defEntities.FieldValueNodeOutput, "", "node1"), nil, nil, outputs, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := val.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", val)
	}
	if m["result"] != "ok" {
		t.Fatalf("expected result=ok, got %v", m["result"])
	}
}

func TestResolve_NodeOutput_Field(t *testing.T) {
	r := NewValueResolver()
	outputs := map[string]interface{}{
		"node1": map[string]interface{}{"status": "done"},
	}
	val, err := r.Resolve(fvWithNode(defEntities.FieldValueNodeOutput, "status", "node1"), nil, nil, outputs, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "done" {
		t.Fatalf("expected 'done', got %v", val)
	}
}

func TestResolve_NodeOutput_MissingNodeID(t *testing.T) {
	r := NewValueResolver()
	_, err := r.Resolve(fvWithNode(defEntities.FieldValueNodeOutput, "x", ""), nil, nil, nil, nil)
	if !errors.Is(err, ErrInvalidNodeID) {
		t.Fatalf("expected ErrInvalidNodeID, got %v", err)
	}
}

func TestResolve_NilField(t *testing.T) {
	r := NewValueResolver()
	_, err := r.Resolve(nil, nil, nil, nil, nil)
	if !errors.Is(err, ErrInvalidFieldValue) {
		t.Fatalf("expected ErrInvalidFieldValue, got %v", err)
	}
}

func TestResolve_Event_FieldNotFound(t *testing.T) {
	r := NewValueResolver()
	event := map[string]interface{}{"a": 1}
	_, err := r.Resolve(fv(defEntities.FieldValueEvent, "missing"), event, nil, nil, nil)
	if !errors.Is(err, ErrFieldNotFound) {
		t.Fatalf("expected ErrFieldNotFound, got %v", err)
	}
}

func TestResolve_Event_NilSource(t *testing.T) {
	r := NewValueResolver()
	_, err := r.Resolve(fv(defEntities.FieldValueEvent, "field"), nil, nil, nil, nil)
	if !errors.Is(err, ErrInvalidSource) {
		t.Fatalf("expected ErrInvalidSource, got %v", err)
	}
}

func TestResolve_UnknownType(t *testing.T) {
	r := NewValueResolver()
	_, err := r.Resolve(fv("unknown_type", "val"), nil, nil, nil, nil)
	if !errors.Is(err, ErrInvalidSource) {
		t.Fatalf("expected ErrInvalidSource, got %v", err)
	}
}

/*
 * LITERAL TEMPLATE INTERPOLATION TESTS
 * Cover the {{namespace.path}} resolution path inside the literal branch.
 * Backwards compat: a literal without `{{` is byte-identical to today
 * (verified by the existing TestResolve_Literal above).
 */

func TestResolve_Literal_Templates_NoBraces_Verbatim(t *testing.T) {
	r := NewValueResolver()
	val, err := r.Resolve(fv(defEntities.FieldValueLiteral, "plain text"), nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "plain text" {
		t.Fatalf("expected 'plain text', got %v", val)
	}
}

func TestResolve_Literal_Templates_SingleNamespace_Event(t *testing.T) {
	r := NewValueResolver()
	event := map[string]interface{}{
		"user": map[string]interface{}{"name": "Ana"},
	}
	val, err := r.Resolve(fv(defEntities.FieldValueLiteral, "Hi {{event.user.name}}"), event, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "Hi Ana" {
		t.Fatalf("expected 'Hi Ana', got %v", val)
	}
}

func TestResolve_Literal_Templates_AllFourNamespaces(t *testing.T) {
	r := NewValueResolver()
	event := map[string]interface{}{"a": "A"}
	state := map[string]interface{}{"b": "B"}
	inputs := map[string]interface{}{"c": "C"}
	outputs := map[string]interface{}{
		"n1": map[string]interface{}{"d": "D"},
	}
	val, err := r.Resolve(
		fv(defEntities.FieldValueLiteral, "{{event.a}}-{{state.b}}-{{input.c}}-{{output.n1.d}}"),
		event, state, outputs, inputs,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "A-B-C-D" {
		t.Fatalf("expected 'A-B-C-D', got %v", val)
	}
}

func TestResolve_Literal_Templates_MissingPath_EmptyString(t *testing.T) {
	r := NewValueResolver()
	val, err := r.Resolve(
		fv(defEntities.FieldValueLiteral, "X{{event.missing}}Y"),
		map[string]interface{}{}, nil, nil, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "XY" {
		t.Fatalf("expected 'XY', got %v", val)
	}
}

func TestResolve_Literal_Templates_NumberScalar(t *testing.T) {
	r := NewValueResolver()
	state := map[string]interface{}{"count": 42}
	val, err := r.Resolve(fv(defEntities.FieldValueLiteral, "n={{state.count}}"), nil, state, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "n=42" {
		t.Fatalf("expected 'n=42', got %v", val)
	}
}

func TestResolve_Literal_Templates_ObjectAsJSON(t *testing.T) {
	r := NewValueResolver()
	inputs := map[string]interface{}{
		"profile": map[string]interface{}{"a": 1},
	}
	val, err := r.Resolve(fv(defEntities.FieldValueLiteral, "{{input.profile}}"), nil, nil, nil, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != `{"a":1}` {
		t.Fatalf(`expected '{"a":1}', got %v`, val)
	}
}

func TestResolve_Literal_Templates_MalformedNoClosingBraces(t *testing.T) {
	r := NewValueResolver()
	val, err := r.Resolve(fv(defEntities.FieldValueLiteral, "Hi {{ unclosed"), nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "Hi {{ unclosed" {
		t.Fatalf("expected 'Hi {{ unclosed', got %v", val)
	}
}

func TestResolve_Literal_Templates_ArrayIndex_Valid(t *testing.T) {
	r := NewValueResolver()
	inputs := map[string]interface{}{
		"recipients": []interface{}{"alice@x.com", "bob@x.com"},
	}
	val, err := r.Resolve(fv(defEntities.FieldValueLiteral, "first={{input.recipients.0}}, second={{input.recipients.1}}"), nil, nil, nil, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "first=alice@x.com, second=bob@x.com" {
		t.Fatalf("expected 'first=alice@x.com, second=bob@x.com', got %v", val)
	}
}

func TestResolve_Literal_Templates_ArrayIndex_OutOfRange(t *testing.T) {
	r := NewValueResolver()
	inputs := map[string]interface{}{
		"recipients": []interface{}{"alice@x.com"},
	}
	val, err := r.Resolve(fv(defEntities.FieldValueLiteral, "X{{input.recipients.5}}Y"), nil, nil, nil, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "XY" {
		t.Fatalf("expected 'XY' (empty for out-of-range index), got %v", val)
	}
}

func TestResolve_Literal_Templates_ArrayIndex_NestedMap(t *testing.T) {
	r := NewValueResolver()
	inputs := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"name": "first"},
			map[string]interface{}{"name": "second"},
		},
	}
	val, err := r.Resolve(fv(defEntities.FieldValueLiteral, "{{input.items.1.name}}"), nil, nil, nil, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "second" {
		t.Fatalf("expected 'second', got %v", val)
	}
}
