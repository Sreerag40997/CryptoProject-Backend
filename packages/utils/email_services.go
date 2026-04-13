package utils

func BuildOTPEmail(otp string) string {
	return "<h2>Your OTP</h2><b>" + otp + "</b><p>Valid for 5 min</p>"
}