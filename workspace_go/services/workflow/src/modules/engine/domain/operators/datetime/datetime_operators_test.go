package datetime

import (
	"testing"
	"time"
)

// --- ParseDate tests ---

func TestParseDate_RFC3339(t *testing.T) {
	d, ok := ParseDate("2025-06-15T10:30:00Z", "UTC")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if d.Year() != 2025 || d.Month() != 6 || d.Day() != 15 {
		t.Fatalf("unexpected date: %v", d)
	}
}

func TestParseDate_DateOnly(t *testing.T) {
	d, ok := ParseDate("2025-06-15", "UTC")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if d.Year() != 2025 || d.Month() != 6 || d.Day() != 15 {
		t.Fatalf("unexpected date: %v", d)
	}
}

func TestParseDate_TimeValue(t *testing.T) {
	now := time.Now()
	d, ok := ParseDate(now, "UTC")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if d.Year() != now.Year() {
		t.Fatal("expected same year")
	}
}

func TestParseDate_UnixTimestamp(t *testing.T) {
	// 2025-01-01T00:00:00Z in Unix seconds
	ts := int64(1735689600)
	d, ok := ParseDate(ts, "UTC")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if d.Year() != 2025 || d.Month() != 1 || d.Day() != 1 {
		t.Fatalf("unexpected date from unix: %v", d)
	}
}

func TestParseDate_InvalidString(t *testing.T) {
	_, ok := ParseDate("not-a-date", "UTC")
	if ok {
		t.Fatal("expected ok=false for invalid date string")
	}
}

func TestParseDate_Timezone(t *testing.T) {
	d, ok := ParseDate("2025-06-15", "America/Sao_Paulo")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if d.Location().String() != "America/Sao_Paulo" {
		t.Fatalf("expected Sao Paulo timezone, got %s", d.Location())
	}
}

func TestParseDate_InvalidTimezone(t *testing.T) {
	// Invalid timezone should fall back to UTC
	d, ok := ParseDate("2025-06-15", "Invalid/Timezone")
	if !ok {
		t.Fatal("expected ok=true (fallback to UTC)")
	}
	if d.Location() != time.UTC {
		t.Fatalf("expected UTC fallback, got %s", d.Location())
	}
}

// --- ParseTime tests ---

func TestParseTime_HourMinuteSecond(t *testing.T) {
	d, ok := ParseTime("14:30:00", "UTC")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if d.Hour() != 14 || d.Minute() != 30 {
		t.Fatalf("unexpected time: %v", d)
	}
}

func TestParseTime_HourMinute(t *testing.T) {
	d, ok := ParseTime("14:30", "UTC")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if d.Hour() != 14 || d.Minute() != 30 {
		t.Fatalf("unexpected time: %v", d)
	}
}

func TestParseTime_Invalid(t *testing.T) {
	_, ok := ParseTime("not-a-time", "UTC")
	if ok {
		t.Fatal("expected ok=false for invalid time string")
	}
}

// --- DateOnly tests ---

func TestDateOnly(t *testing.T) {
	original := time.Date(2025, 6, 15, 14, 30, 45, 123, time.UTC)
	result := DateOnly(original)
	if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
		t.Fatalf("expected time zeroed, got %v", result)
	}
	if result.Year() != 2025 || result.Month() != 6 || result.Day() != 15 {
		t.Fatalf("expected same date, got %v", result)
	}
}

// --- CompareDates tests ---

func TestCompareDates(t *testing.T) {
	d1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)

	if CompareDates(d1, d2) != -1 {
		t.Fatal("expected d1 before d2")
	}
	if CompareDates(d2, d1) != 1 {
		t.Fatal("expected d2 after d1")
	}
	if CompareDates(d1, d1) != 0 {
		t.Fatal("expected equal")
	}
}

// --- CompareTimes tests ---

func TestCompareTimes(t *testing.T) {
	loc := time.UTC
	t1 := time.Date(2025, 1, 1, 8, 0, 0, 0, loc)
	t2 := time.Date(2025, 1, 1, 14, 30, 0, 0, loc)

	if CompareTimes(t1, t2) != -1 {
		t.Fatal("expected t1 before t2")
	}
	if CompareTimes(t2, t1) != 1 {
		t.Fatal("expected t2 after t1")
	}
	if CompareTimes(t1, t1) != 0 {
		t.Fatal("expected equal")
	}
}

// --- BeforeDateOperator tests ---

func TestBeforeDateOperator_Before(t *testing.T) {
	op := &BeforeDateOperator{}
	if op.Name() != "beforeDate" {
		t.Fatalf("expected name 'beforeDate', got %q", op.Name())
	}
	result, err := op.Evaluate("UTC", "2025-01-01", "2025-06-15")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 2025-01-01 before 2025-06-15")
	}
}

func TestBeforeDateOperator_After(t *testing.T) {
	op := &BeforeDateOperator{}
	result, err := op.Evaluate("UTC", "2025-06-15", "2025-01-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 2025-06-15 NOT before 2025-01-01")
	}
}

func TestBeforeDateOperator_Equal(t *testing.T) {
	op := &BeforeDateOperator{}
	result, err := op.Evaluate("UTC", "2025-06-15", "2025-06-15")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected same date NOT before itself")
	}
}

func TestBeforeDateOperator_InvalidField(t *testing.T) {
	op := &BeforeDateOperator{}
	result, err := op.Evaluate("UTC", "not-a-date", "2025-06-15")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false for unparsable field")
	}
}

func TestBeforeDateOperator_IgnoresTime(t *testing.T) {
	op := &BeforeDateOperator{}
	// Same date, different times → should be equal (not before)
	result, err := op.Evaluate("UTC", "2025-06-15T23:59:59Z", "2025-06-15T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected same date to NOT be before (ignores time)")
	}
}

// --- AfterDateOperator tests ---

func TestAfterDateOperator_After(t *testing.T) {
	op := &AfterDateOperator{}
	if op.Name() != "afterDate" {
		t.Fatalf("expected name 'afterDate', got %q", op.Name())
	}
	result, err := op.Evaluate("UTC", "2025-06-15", "2025-01-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 2025-06-15 after 2025-01-01")
	}
}

func TestAfterDateOperator_Before(t *testing.T) {
	op := &AfterDateOperator{}
	result, err := op.Evaluate("UTC", "2025-01-01", "2025-06-15")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 2025-01-01 NOT after 2025-06-15")
	}
}

func TestAfterDateOperator_Equal(t *testing.T) {
	op := &AfterDateOperator{}
	result, err := op.Evaluate("UTC", "2025-06-15", "2025-06-15")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected same date NOT after itself")
	}
}

// --- BeforeTimeOperator tests ---

func TestBeforeTimeOperator_Before(t *testing.T) {
	op := &BeforeTimeOperator{}
	if op.Name() != "beforeTime" {
		t.Fatalf("expected name 'beforeTime', got %q", op.Name())
	}
	result, err := op.Evaluate("UTC", "08:00", "12:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 08:00 before 12:00")
	}
}

func TestBeforeTimeOperator_After(t *testing.T) {
	op := &BeforeTimeOperator{}
	result, err := op.Evaluate("UTC", "14:30", "12:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 14:30 NOT before 12:00")
	}
}

func TestBeforeTimeOperator_Equal(t *testing.T) {
	op := &BeforeTimeOperator{}
	result, err := op.Evaluate("UTC", "12:00", "12:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected same time NOT before itself")
	}
}

func TestBeforeTimeOperator_InvalidField(t *testing.T) {
	op := &BeforeTimeOperator{}
	result, err := op.Evaluate("UTC", "not-a-time", "12:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false for unparsable time")
	}
}

// --- AfterTimeOperator tests ---

func TestAfterTimeOperator_After(t *testing.T) {
	op := &AfterTimeOperator{}
	if op.Name() != "afterTime" {
		t.Fatalf("expected name 'afterTime', got %q", op.Name())
	}
	result, err := op.Evaluate("UTC", "14:30", "12:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 14:30 after 12:00")
	}
}

func TestAfterTimeOperator_Before(t *testing.T) {
	op := &AfterTimeOperator{}
	result, err := op.Evaluate("UTC", "08:00", "12:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 08:00 NOT after 12:00")
	}
}

func TestAfterTimeOperator_Equal(t *testing.T) {
	op := &AfterTimeOperator{}
	result, err := op.Evaluate("UTC", "12:00", "12:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected same time NOT after itself")
	}
}

// --- BetweenDateOperator tests ---

func TestBetweenDateOperator_InRange(t *testing.T) {
	op := &BetweenDateOperator{}
	if op.Name() != "betweenDate" {
		t.Fatalf("expected name 'betweenDate', got %q", op.Name())
	}
	result, err := op.Evaluate("UTC", "2025-01-15", []interface{}{"2025-01-01", "2025-01-31"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 2025-01-15 between [2025-01-01, 2025-01-31]")
	}
}

func TestBetweenDateOperator_AtStart(t *testing.T) {
	op := &BetweenDateOperator{}
	// Inclusive by default
	result, err := op.Evaluate("UTC", "2025-01-01", []interface{}{"2025-01-01", "2025-01-31"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected start date inclusive")
	}
}

func TestBetweenDateOperator_AtEnd(t *testing.T) {
	op := &BetweenDateOperator{}
	result, err := op.Evaluate("UTC", "2025-01-31", []interface{}{"2025-01-01", "2025-01-31"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected end date inclusive")
	}
}

func TestBetweenDateOperator_OutOfRange(t *testing.T) {
	op := &BetweenDateOperator{}
	result, err := op.Evaluate("UTC", "2025-02-15", []interface{}{"2025-01-01", "2025-01-31"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 2025-02-15 NOT between [2025-01-01, 2025-01-31]")
	}
}

func TestBetweenDateOperator_MapFormat(t *testing.T) {
	op := &BetweenDateOperator{}
	result, err := op.Evaluate("UTC", "2025-01-15",
		map[string]interface{}{"start": "2025-01-01", "end": "2025-01-31"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected date in range with map format")
	}
}

func TestBetweenDateOperator_Exclusive(t *testing.T) {
	op := &BetweenDateOperator{}
	result, err := op.EvaluateRange("UTC", "2025-01-01", "2025-01-01", "2025-01-31", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected start date NOT in range (exclusive)")
	}
}

func TestBetweenDateOperator_InvalidField(t *testing.T) {
	op := &BetweenDateOperator{}
	result, err := op.Evaluate("UTC", "not-a-date", []interface{}{"2025-01-01", "2025-01-31"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false for unparsable field")
	}
}

// --- BetweenTimeOperator tests ---

func TestBetweenTimeOperator_InRange(t *testing.T) {
	op := &BetweenTimeOperator{}
	if op.Name() != "betweenTime" {
		t.Fatalf("expected name 'betweenTime', got %q", op.Name())
	}
	result, err := op.Evaluate("UTC", "10:30", []interface{}{"08:00", "17:00"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected 10:30 between [08:00, 17:00]")
	}
}

func TestBetweenTimeOperator_AtStart(t *testing.T) {
	op := &BetweenTimeOperator{}
	result, err := op.Evaluate("UTC", "08:00", []interface{}{"08:00", "17:00"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected start time inclusive")
	}
}

func TestBetweenTimeOperator_AtEnd(t *testing.T) {
	op := &BetweenTimeOperator{}
	result, err := op.Evaluate("UTC", "17:00", []interface{}{"08:00", "17:00"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Fatal("expected end time inclusive")
	}
}

func TestBetweenTimeOperator_OutOfRange(t *testing.T) {
	op := &BetweenTimeOperator{}
	result, err := op.Evaluate("UTC", "20:00", []interface{}{"08:00", "17:00"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected 20:00 NOT between [08:00, 17:00]")
	}
}

func TestBetweenTimeOperator_Exclusive(t *testing.T) {
	op := &BetweenTimeOperator{}
	result, err := op.EvaluateRange("UTC", "08:00", "08:00", "17:00", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected start time NOT in range (exclusive)")
	}
}

func TestBetweenTimeOperator_InvalidField(t *testing.T) {
	op := &BetweenTimeOperator{}
	result, err := op.Evaluate("UTC", "not-a-time", []interface{}{"08:00", "17:00"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false for unparsable time")
	}
}

func TestBetweenTimeOperator_InvalidRange(t *testing.T) {
	op := &BetweenTimeOperator{}
	// Only one element → extractTimeRange fails
	result, err := op.Evaluate("UTC", "10:00", []interface{}{"08:00"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Fatal("expected false for invalid time range")
	}
}

// --- getLocation tests ---

func TestGetLocation_ValidTimezone(t *testing.T) {
	loc := getLocation("America/New_York")
	if loc.String() != "America/New_York" {
		t.Fatalf("expected 'America/New_York', got %q", loc.String())
	}
}

func TestGetLocation_EmptyTimezone(t *testing.T) {
	loc := getLocation("")
	if loc != time.UTC {
		t.Fatalf("expected UTC for empty timezone, got %q", loc.String())
	}
}

func TestGetLocation_InvalidTimezone(t *testing.T) {
	loc := getLocation("Invalid/Zone")
	if loc != time.UTC {
		t.Fatalf("expected UTC fallback for invalid timezone, got %q", loc.String())
	}
}

// --- Metadata tests ---

func TestOperatorMetadata(t *testing.T) {
	ops := []struct {
		name      string
		isBetween bool
		op        interface {
			Name() string
		}
	}{
		{"beforeDate", false, &BeforeDateOperator{}},
		{"afterDate", false, &AfterDateOperator{}},
		{"beforeTime", false, &BeforeTimeOperator{}},
		{"afterTime", false, &AfterTimeOperator{}},
		{"betweenDate", true, &BetweenDateOperator{}},
		{"betweenTime", true, &BetweenTimeOperator{}},
	}
	for _, tt := range ops {
		if tt.op.Name() != tt.name {
			t.Errorf("expected name %q, got %q", tt.name, tt.op.Name())
		}
	}
}
