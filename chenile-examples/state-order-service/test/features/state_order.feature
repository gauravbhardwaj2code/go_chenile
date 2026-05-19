Feature: State order service

  Scenario: Confirm order
    When I POST a REST request to URL "/state-orders" with payload
      """
      {
        "id": "order-1",
        "name": "Order 1"
      }
      """
    Then the http status code is 200
    And success is true
    And the REST response key "state" is "created"
    When I POST a REST request to URL "/state-orders/event" with payload
      """
      {
        "id": "order-1",
        "event": "confirm"
      }
      """
    Then the http status code is 200
    And success is true
    And the REST response key "state" is "confirmed"
