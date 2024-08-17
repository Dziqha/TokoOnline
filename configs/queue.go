package configs

import "log"

// OrderQueue menambahkan orderId ke antrian order_queue secara asinkron
func OrderQueue(orderId string) error {
	go func() {
		if err := AddToQueue("order_queue", orderId); err != nil {
			log.Printf("Error in OrderQueue for orderId %s: %v", orderId, err)
		}
	}()
	return nil
}

// DeleteQueue menambahkan orderId ke antrian delete_queue secara asinkron
func DeleteQueue(orderId string) error {
	go func() {
		if err := AddToQueue("delete_queue", orderId); err != nil {
			log.Printf("Error in DeleteQueue for orderId %s: %v", orderId, err)
		}
	}()
	return nil
}

// CancelQueue menambahkan orderId ke antrian cancel_queue secara asinkron
func CancelQueue(orderId string) error {
	go func() {
		if err := AddToQueue("cancel_queue", orderId); err != nil {
			log.Printf("Error in CancelQueue for orderId %s: %v", orderId, err)
		}
	}()
	return nil
}

func CheckOutQueue(orderId string) error {
	go func() {
		if err := AddToQueue("checkout_queue", orderId); err != nil {
			log.Printf("Error in CheckOutQueue for orderId %s: %v", orderId, err)
		}
	}()
	return nil
}
