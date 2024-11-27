package main

import (
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	receipt1 := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
		},
		Total: "6.49",
	}
	receipt2 := receipt1 // Identical receipt

	uuid1 := generateUUID(receipt1)
	uuid2 := generateUUID(receipt2)

	if uuid1 != uuid2 {
		t.Errorf("Expected identical UUIDs for identical receipts, got %s and %s", uuid1, uuid2)
	}
}

func TestCalculatePoints(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "14:30",
		Items: []Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
		},
		Total: "18.74",
	}
	points := calculatePoints(receipt)

	expectedPoints := 6 + 5 + 3 + 6 + 10 // Customize based on rules.
	if points != expectedPoints {
		t.Errorf("Expected %d points, got %d", expectedPoints, points)
	}
}

func TestCalculatePoints_Example1(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "35.35",
	}
	points := calculatePoints(receipt)

	expectedPoints := 6 + 10 + 3 + 3 + 6 // Customize based on rules.
	if points != expectedPoints {
		t.Errorf("Expected %d points, got %d", expectedPoints, points)
	}
}

func TestCalculatePoints_Example2(t *testing.T) {
	receipt := Receipt{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Items: []Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
		Total: "9.00",
	}
	points := calculatePoints(receipt)

	expectedPoints := 50 + 25 + 14 + 10 + 10 // Customize based on rules.
	if points != expectedPoints {
		t.Errorf("Expected %d points, got %d", expectedPoints, points)
	}
}
