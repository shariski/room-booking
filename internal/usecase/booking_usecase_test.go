package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shariski/room-booking/internal/model"
	"github.com/shariski/room-booking/internal/repository"
	"github.com/shariski/room-booking/internal/usecase"
)

var (
	db             *sql.DB
	bookingUsecase *usecase.BookingUsecase
	testRoomID     uuid.UUID
	testRoomID2    uuid.UUID
	testUserID     uuid.UUID
)

func TestMain(m *testing.M) {
	// setup: connect DB, run migrations, seed data
	var err error
	db, err = sql.Open("pgx", "postgres://bobobox:bobobox@localhost:5432/bobobox_test?sslmode=disable")
	if err != nil {
		panic(err)
	}

	bookingRepo := repository.NewBookingRepository(db)
	bookingUsecase = usecase.NewBookingUsecase(bookingRepo)

	code := m.Run()
	db.Close()
	os.Exit(code)
}

func setupTest(t *testing.T) {
	t.Helper()
	db.Exec("truncate bookings, rooms, users restart identity cascade")

	// seed user
	db.QueryRow(
		"INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		"Test User", "test@test.com", "hashedpassword",
	).Scan(&testUserID)

	// seed room
	db.QueryRow(
		"INSERT INTO rooms (name, type, description) VALUES ($1, $2, $3) RETURNING id",
		"Test Room", "single", "Test room description",
	).Scan(&testRoomID)

	db.QueryRow(
		"INSERT INTO rooms (name, type, description) VALUES ($1, $2, $3) RETURNING id",
		"Test Room 2", "double", "Test room double",
	).Scan(&testRoomID2)

}

func TestCreateBooking_Success(t *testing.T) {
	setupTest(t)
	request := &model.CreateBookingRequest{
		RoomID:    testRoomID,
		UserID:    testUserID,
		StartDate: time.Now().AddDate(0, 0, 1),
		EndDate:   time.Now().AddDate(0, 0, 3),
	}

	result, err := bookingUsecase.Create(context.Background(), request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestCreateBooking_OverlappingDates(t *testing.T) {
	setupTest(t)
	request := &model.CreateBookingRequest{
		RoomID:    testRoomID,
		UserID:    testUserID,
		StartDate: time.Now().AddDate(0, 0, 10),
		EndDate:   time.Now().AddDate(0, 0, 15),
	}
	_, err := bookingUsecase.Create(context.Background(), request)
	if err != nil {
		t.Fatalf("first booking should succeed, got %v", err)
	}

	request2 := &model.CreateBookingRequest{
		RoomID:    testRoomID,
		UserID:    testUserID,
		StartDate: time.Now().AddDate(0, 0, 12),
		EndDate:   time.Now().AddDate(0, 0, 17),
	}
	_, err = bookingUsecase.Create(context.Background(), request2)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}

	var appErr *model.AppError
	if !errors.As(err, &appErr) || appErr.Status != http.StatusConflict {
		t.Fatalf("expected 409 conflict, got %v", err)
	}
}

func TestCreateBooking_SameDatesDifferentRooms(t *testing.T) {
	setupTest(t)
	request := &model.CreateBookingRequest{
		RoomID:    testRoomID,
		UserID:    testUserID,
		StartDate: time.Now().AddDate(0, 0, 10),
		EndDate:   time.Now().AddDate(0, 0, 15),
	}
	_, err := bookingUsecase.Create(context.Background(), request)
	if err != nil {
		t.Fatalf("first booking should succeed, got %v", err)
	}

	request2 := &model.CreateBookingRequest{
		RoomID:    testRoomID2,
		UserID:    testUserID,
		StartDate: time.Now().AddDate(0, 0, 10),
		EndDate:   time.Now().AddDate(0, 0, 15),
	}
	_, err = bookingUsecase.Create(context.Background(), request2)
	if err != nil {
		t.Fatalf("different room bookings should succeed, got %v", err)
	}
}

func TestCreateBooking_AdjacentDates(t *testing.T) {
	setupTest(t)
	req1 := &model.CreateBookingRequest{
		RoomID:    testRoomID,
		UserID:    testUserID,
		StartDate: time.Now().AddDate(0, 0, 20),
		EndDate:   time.Now().AddDate(0, 0, 25),
	}
	_, err := bookingUsecase.Create(context.Background(), req1)
	if err != nil {
		t.Fatalf("first booking should succeed, got %v", err)
	}

	req2 := &model.CreateBookingRequest{
		RoomID:    testRoomID,
		UserID:    testUserID,
		StartDate: time.Now().AddDate(0, 0, 25),
		EndDate:   time.Now().AddDate(0, 0, 30),
	}
	_, err = bookingUsecase.Create(context.Background(), req2)
	if err != nil {
		t.Fatalf("adjacent booking should succeed, got %v", err)
	}
}

func TestCreateBooking_StartDateAfterEndDate(t *testing.T) {
	setupTest(t)
	request := &model.CreateBookingRequest{
		RoomID:    testRoomID,
		UserID:    testUserID,
		StartDate: time.Now().AddDate(0, 0, 10),
		EndDate:   time.Now().AddDate(0, 0, 5),
	}

	_, err := bookingUsecase.Create(context.Background(), request)
	if err == nil {
		t.Fatal("expected bad request error, got nil")
	}
}

func TestDeleteBooking_Success(t *testing.T) {
	setupTest(t)
	request := &model.CreateBookingRequest{
		RoomID:    testRoomID,
		UserID:    testUserID,
		StartDate: time.Now().AddDate(0, 0, 40),
		EndDate:   time.Now().AddDate(0, 0, 45),
	}
	booking, _ := bookingUsecase.Create(context.Background(), request)

	deleteReq := &model.DeleteBookingRequest{
		ID:     booking.ID,
		UserID: testUserID,
	}
	result, err := bookingUsecase.Delete(context.Background(), deleteReq)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.DeletedAt == nil {
		t.Fatal("expected deleted_at to be set")
	}
}

func TestDeleteBooking_ThenRebook(t *testing.T) {
	setupTest(t)
	// create
	request := &model.CreateBookingRequest{
		RoomID:    testRoomID,
		UserID:    testUserID,
		StartDate: time.Now().AddDate(0, 0, 50),
		EndDate:   time.Now().AddDate(0, 0, 55),
	}
	booking, _ := bookingUsecase.Create(context.Background(), request)

	// delete
	deleteReq := &model.DeleteBookingRequest{
		ID:     booking.ID,
		UserID: testUserID,
	}
	bookingUsecase.Delete(context.Background(), deleteReq)

	_, err := bookingUsecase.Create(context.Background(), request)
	if err != nil {
		t.Fatalf("rebook after cancel should succeed, got %v", err)
	}
}

func TestDeleteBooking_NotFound(t *testing.T) {
	setupTest(t)
	deleteReq := &model.DeleteBookingRequest{
		ID:     uuid.New(),
		UserID: testUserID,
	}
	_, err := bookingUsecase.Delete(context.Background(), deleteReq)
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
}

func TestCreateBooking_ConcurrentSameRoom(t *testing.T) {
	setupTest(t)
	startDate := time.Now().AddDate(0, 0, 60)
	endDate := time.Now().AddDate(0, 0, 65)

	results := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func() {
			req := &model.CreateBookingRequest{
				RoomID:    testRoomID,
				UserID:    testUserID,
				StartDate: startDate,
				EndDate:   endDate,
			}
			_, err := bookingUsecase.Create(context.Background(), req)
			results <- err
		}()
	}

	successCount := 0
	for i := 0; i < 10; i++ {
		err := <-results
		if err == nil {
			successCount++
		}
	}

	if successCount != 1 {
		t.Fatalf("expected exactly 1 success, got %d", successCount)
	}
}
