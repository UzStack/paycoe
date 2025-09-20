package main

import (
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// TopUpResult parse natijasi.
type TopUpResult struct {
	Type         string // "top_up"
	AmountRaw    string // "3.300.000,00"
	AmountInt    int64  // 3300000 (butun pul birligi sifatida)
	Currency     string // "UZS"
	AmountPretty string // qayta formatlashni xohlasa alohida qilishingiz mumkin
}

func toNumber(amount string, log *zap.Logger) (int64, bool) {
	normalized := strings.ReplaceAll(strings.ReplaceAll(amount, ".", ""), ",", ".")
	if res, err := strconv.ParseFloat(normalized, 64); err == nil {
		return int64(res), true
	}
	return 0, false
}

var (
	// To['’`]?ldirish so'zini izlash (case-insensitive).
	reTopUpWord = regexp.MustCompile(`(?i)To['’` + "`" + `]?ldirish`)
	// ➕ 3.300.000,00 UZS
	reAmount = regexp.MustCompile(`➕\s*([\d\.]+,\d{2})\s*([A-Z]{3})`)
)

// ParseTopUp matndan To'ldirish operatsiyasi bo'lsa ajratib qaytaradi.
func ParseTopUp(text string, log *zap.Logger) *TopUpResult {
	// Unicode no-break space va boshqalarni oddiy bo'shliqqa
	cleaned := strings.ReplaceAll(text, "\u202f", " ")
	cleaned = strings.ReplaceAll(cleaned, "\u00a0", " ")

	if !reTopUpWord.MatchString(cleaned) {
		return nil
	}
	m := reAmount.FindStringSubmatch(cleaned)
	if len(m) != 3 {
		return nil
	}
	amountStr := m[1]
	currency := m[2]

	val, ok := toNumber(amountStr, log)
	if !ok {
		// Agar bu yerga decimal qo'shmoqchi bo'lsangiz shopspring/decimal ni ulanishingiz mumkin.
	}

	return &TopUpResult{
		Type:         "top_up",
		AmountRaw:    amountStr,
		AmountInt:    val,
		Currency:     currency,
		AmountPretty: amountStr + " " + currency,
	}
}
