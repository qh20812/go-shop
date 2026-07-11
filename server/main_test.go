package main

import "testing"

// Unit test đầu tiên của bạn — CI sẽ chạy file này ở MỌI lần push.
// Test hàm thuần (không cần DB) là cách bắt đầu dễ nhất.

func TestSanitizeFilename(t *testing.T) {
	got := sanitizeFilename("My Photo (1).PNG")
	want := "my-photo--1-.png"
	if got != want {
		t.Errorf("sanitizeFilename() = %q, muốn %q", got, want)
	}
}

func TestSanitizeFilenameChanPathTraversal(t *testing.T) {
	// Kẻ xấu đặt tên file "../../etc/passwd" hòng ghi đè file hệ thống
	got := sanitizeFilename("../../etc/passwd")
	want := "passwd" // filepath.Base đã cắt bỏ đường dẫn
	if got != want {
		t.Errorf("sanitizeFilename() = %q, muốn %q", got, want)
	}
}