terraform {
  backend "http" {
    address = "http://localhost:8080/foo"
    lock_address = "http://localhost:8080/foo"
    unlock_address = "http://localhost:8080/foo"
 }
}

resource "null_resource" "test" {}
